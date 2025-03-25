package desec

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dsc "github.com/nrdcg/desec"
)

func resourceToken() *schema.Resource {
	return &schema.Resource{
		// Creation is not supported, since it will create a secret that can't be shown
		CreateContext: resourceTokenCreate,
		ReadContext:   resourceTokenRead,
		UpdateContext: resourceTokenUpdate,
		DeleteContext: resourceTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			/*
			   {
			       "id": "3a6b94b5-d20e-40bd-a7cc-521f5c79fab3",
			       "created": "2018-09-06T09:08:43.762697Z",
			       "last_used": null,
			       "owner": "youremailaddress@example.com"",
			       "user_override": null,
			       "max_age": "365 00:00:00",
			       "max_unused_period": null,
			       "name": "my new token",
			       "perm_create_domain": false,
			       "perm_delete_domain": false,
			       "perm_manage_tokens": false,
			       "allowed_subnets": [
			           "0.0.0.0/0",
			           "::/0"
			       ],
			       "auto_policy": false,
			       "token": "4pnk7u-NHvrEkFzrhFDRTjGFyX_S"
			   }
			*/
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Typically only available right after creation, emptied on next refresh
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 178),
			},
			"perm_create_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"perm_delete_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"perm_manage_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allowed_subnets": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	t := schemaToToken(d)
	token, err := c.Tokens.Create(ctx, t.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	// TODO unify in create call
	token, err = c.Tokens.Update(ctx, token.ID, &t)
	if err != nil {
		return diag.FromErr(err)
	}

	tokenIntoSchema(token, d)
	return diags
}

func resourceTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	t, err := c.Tokens.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if t == nil {
		d.SetId("")
	} else {
		tokenIntoSchema(t, d)
	}

	return diags
}

func resourceTokenUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	t := schemaToToken(d)
	token, err := c.Tokens.Update(ctx, d.Id(), &t)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	tokenIntoSchema(token, d)

	return diags
}

func resourceTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	err := c.Tokens.Delete(ctx, d.Id())
	if err != nil && !isNotFoundError(err) {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func tokenIntoSchema(r *dsc.Token, d *schema.ResourceData) {
	d.SetId(r.ID)
	d.Set("created", (*r.Created).Format(time.RFC3339))
	d.Set("name", r.Name)
	d.Set("owner", r.Owner)
	d.Set("perm_create_domain", r.PermCreateDomain)
	d.Set("perm_delete_domain", r.PermDeleteDomain)
	d.Set("perm_manage_tokens", r.PermManageTokens)
	d.Set("auto_policy", r.AutoPolicy)
	d.Set("token", r.Value)
	if r.AllowedSubnets != nil {
		d.Set("allowed_subnets", r.AllowedSubnets)
	}
}

func schemaToToken(d *schema.ResourceData) dsc.Token {
	result := dsc.Token{
		Owner:            d.Get("owner").(string),
		Name:             d.Get("name").(string),
		PermCreateDomain: d.Get("perm_create_domain").(bool),
		PermDeleteDomain: d.Get("perm_delete_domain").(bool),
		PermManageTokens: d.Get("perm_manage_tokens").(bool),
		AutoPolicy:       d.Get("auto_policy").(bool),
	}
	result.AllowedSubnets = []string{}
	for _, as := range d.Get("allowed_subnets").([]any) {
		result.AllowedSubnets = append(result.AllowedSubnets, as.(string))
	}
	return result
}
