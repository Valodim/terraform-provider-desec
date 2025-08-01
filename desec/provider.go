package desec

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dsc "github.com/nrdcg/desec"
)

type TokenCreateMode int

const (
	TokenCreateModePrint TokenCreateMode = iota
	TokenCreateModeStoreState
)

type DesecConfig struct {
	cache           *DesecCache
	client          *dsc.Client
	tokenCreateMode TokenCreateMode
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DESEC_API_URI", ""),
			},
			"api_token": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("DESEC_API_TOKEN", ""),
				Description:  "The API token for operations.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("[0-9a-zA-Z_-]{28}"), "API key looks invalid"),
			},
			"retry_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The max number of retries when sending an API request.",
			},
			"token_create_mode": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"print", "storeState"}, false),
				Default:      "print",
				Optional:     true,
				Description:  "The way that tokens are exposed to the user during creation. If 'print', the token will be output as a warning message during creation. If 'storeState', the token will be stored in the state initially (this may be insecure!), but cleared when the state is next refreshed.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"desec_rrset":        resourceRRSet(),
			"desec_domain":       resourceDomain(),
			"desec_token":        resourceToken(),
			"desec_token_policy": resourceTokenPolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("api_token").(string)
	if token == "" {
		return nil, diag.Errorf("missing config field: api_token")
	}

	o := dsc.NewDefaultClientOptions()
	o.HTTPClient = cleanhttp.DefaultClient()
	o.HTTPClient.Transport = logging.NewTransport("Desec", o.HTTPClient.Transport)
	o.Logger = log.Default()

	retry_max, retry_max_set := d.GetOk("retry_max")
	if retry_max_set {
		o.RetryMax = retry_max.(int)
	}

	tokenCreateMode := TokenCreateModePrint
	if tcm, ok := d.GetOk("token_create_mode"); ok {
		switch tcm {
		case "storeState":
			tokenCreateMode = TokenCreateModeStoreState
		case "print":
			tokenCreateMode = TokenCreateModePrint
		default:
			return nil, diag.Errorf("invalid value for config field: token_create_mode")
		}
	}

	c := dsc.New(token, o)
	api_uri := d.Get("api_uri").(string)
	if api_uri != "" {
		c.BaseURL = api_uri
	}

	cache := NewDesecCache()
	return &DesecConfig{&cache, c, tokenCreateMode}, nil
}

func isNotFoundError(err error) bool {
	apiError, ok := err.(*dsc.APIError)
	if !ok {
		return false
	}
	return apiError != nil && apiError.StatusCode == http.StatusNotFound
}
