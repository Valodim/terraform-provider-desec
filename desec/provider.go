package desec

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dsc "github.com/nrdcg/desec"
)

type DesecConfig struct {
	cache  *DesecCache
	client *dsc.Client
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
			"limit_read": {
				// see https://desec.readthedocs.io/en/latest/rate-limits.html
				// We choose 8r/s here, because going for the full 10/s sometimes runs into limiter failures.
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     8,
				Description: "Limit of API requests per second, for read operations.",
			},
			"limit_write": {
				// see https://desec.readthedocs.io/en/latest/rate-limits.html
				// We choose 5r/s here, because going for 6/s sometimes runs into limiter failures.
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Limit of API requests per second, for write operations.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"desec_rrset":  resourceRRSet(),
			"desec_domain": resourceDomain(),
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
	o.LimitRead = d.Get("limit_read").(int)
	o.LimitWrite = d.Get("limit_write").(int)

	c := dsc.New(token, o)
	api_uri := d.Get("api_uri").(string)
	if api_uri != "" {
		c.BaseURL = api_uri
	}

	cache := NewDesecCache()
	return &DesecConfig{&cache, c}, nil
}

func isNotFoundError(err error) bool {
	apiError, ok := err.(*dsc.APIError)
	if !ok {
		return false
	}
	return apiError != nil && apiError.StatusCode == http.StatusNotFound
}
