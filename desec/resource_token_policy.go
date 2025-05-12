package desec

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dsc "github.com/nrdcg/desec"
)

/* Implementation notes:
 *  - The internal ID is tokenID + policyID, since the policyID on its own isn't a unique identifier
 *  - The nullable fields are represented as empty strings internally
 */
func resourceTokenPolicy() *schema.Resource {
	return &schema.Resource{
		// Creation is not supported, since it will create a secret that can't be shown
		CreateContext: resourceTokenPolicyCreate,
		ReadContext:   resourceTokenPolicyRead,
		UpdateContext: resourceTokenPolicyUpdate,
		DeleteContext: resourceTokenPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTokenPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			/*
				{
						"id": "7aed3f71-bc81-4f7e-90ae-8f0df0d1c211",
						"domain": "example.com",
						"subname": null,
						"type": null,
						"perm_write": true
				}
			*/
			"token_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"subname": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "",
				Optional: true,
			},
			"perm_write": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceTokenPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	tokenId := d.Get("token_id").(string)
	tokenPolicy, err := c.TokenPolicies.Create(ctx, tokenId, schemaToTokenPolicy(d))
	if err != nil {
		return diag.FromErr(err)
	}

	tokenPolicyIntoSchema(tokenId, tokenPolicy, d)
	return diags
}

func resourceTokenPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	policy, err := c.TokenPolicies.GetOne(ctx, d.Get("token_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if policy == nil {
		d.SetId("")
	} else {
		tokenPolicyIntoSchema(d.Get("token_id").(string), policy, d)
	}

	return diags
}

func resourceTokenPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	t := schemaToTokenPolicy(d)
	tokenPolicy, err := c.TokenPolicies.Update(ctx, d.Get("token_id").(string), d.Id(), t)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	tokenPolicyIntoSchema(d.Get("token_id").(string), tokenPolicy, d)

	return diags
}

func resourceTokenPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	err := c.TokenPolicies.Delete(ctx, d.Get("token_id").(string), d.Id())
	if err != nil && !isNotFoundError(err) {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceTokenPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	pieces := strings.Split(d.Id(), "/")
	if len(pieces) != 2 {
		return nil, fmt.Errorf("invalid id: %s", d.Id())
	}

	d.Set("token_id", pieces[0])
	d.SetId(pieces[1])
	return []*schema.ResourceData{d}, nil
}

func tokenPolicyIntoSchema(tokenId string, r *dsc.TokenPolicy, d *schema.ResourceData) {
	d.SetId(r.ID)
	d.Set("token_id", tokenId)
	if r.Domain != nil {
		d.Set("domain", r.Domain)
	} else {
		d.Set("domain", "")
	}
	if r.SubName != nil {
		d.Set("subname", r.SubName)
	} else {
		d.Set("subname", "")
	}
	if r.Type != nil {
		d.Set("type", r.Type)
	} else {
		d.Set("type", "")
	}
	d.Set("perm_write", r.WritePermission)
}

func schemaToTokenPolicy(d *schema.ResourceData) dsc.TokenPolicy {
	result := dsc.TokenPolicy{
		WritePermission: d.Get("perm_write").(bool),
	}
	domain := d.Get("domain").(string)
	if domain != "" {
		result.Domain = &domain
	}
	subname := d.Get("subname").(string)
	if subname != "" {
		result.SubName = &subname
	}
	typ := d.Get("type").(string)
	if typ != "" {
		result.Type = &typ
	}
	return result
}

