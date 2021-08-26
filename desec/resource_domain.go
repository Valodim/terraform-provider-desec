package desec

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dsc "github.com/nrdcg/desec"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"minimum_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"published": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dnskey": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ds": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"flags": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"keytype": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	domainName := d.Get("name").(string)
	domain, err := c.Domains.Create(ctx, domainName)
	if err != nil {
		return diag.FromErr(err)
	}

	domainIntoData(domain, d)
	return nil
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	domain, err := c.Domains.Get(ctx, d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	domainIntoData(domain, d)
	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	err := c.Domains.Delete(ctx, d.Id())
	if err != nil && !isNotFoundError(err) {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func domainIntoData(domain *dsc.Domain, d *schema.ResourceData) {
	d.Set("created", domain.Created.Format(time.RFC3339))
	d.Set("name", domain.Name)
	d.Set("minimum_ttl", domain.MinimumTTL)
	if domain.Published != nil {
		d.Set("published", domain.Published.Format(time.RFC3339))
	}
	keys := make([]*map[string]interface{}, len(domain.Keys))
	for i, k := range domain.Keys {
		key := make(map[string]interface{})
		key["dnskey"] = k.DNSKey
		key["ds"] = k.DS
		key["flags"] = k.Flags
		key["keytype"] = k.KeyType
		keys[i] = &key
	}
	d.Set("keys", keys)
	d.SetId(domain.Name)
}
