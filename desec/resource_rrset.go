package desec

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dsc "github.com/nrdcg/desec"
)

func resourceRRSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRRSetCreate,
		ReadContext:   resourceRRSetRead,
		UpdateContext: resourceRRSetUpdate,
		DeleteContext: resourceRRSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CAA", "CNAME", "TXT", "SRV", "LOC", "MX", "NS", "SPF", "CERT", "DNSKEY", "DS", "NAPTR", "SMIMEA", "SSHFP", "TLSA", "URI", "PTR"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(3600),
			},
		},
	}
}

func resourceRRSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dsc.Client)

	var diags diag.Diagnostics

	r := schemaToRRset(d)
	rrset, err := c.Records.Create(r)
	if err != nil {
		return diag.FromErr(err)
	}

	rrsetIntoSchema(rrset, d)
	return diags
}

func idFromNames(domainName, subName, recordType string) string {
	if subName == "" {
		subName = "@"
	}
	return fmt.Sprintf("%s/%s/%s", domainName, subName, recordType)
}

func namesFromId(id string) (string, string, string, error) {
	var domainName string
	var subName string
	var recordType string

	idAttr := strings.SplitN(id, "/", 3)
	if len(idAttr) != 3 {
		return "", "", "", fmt.Errorf("invalid id %q specified, should be in format \"domainName/subName/type\" for import", id)
	}

	domainName = idAttr[0]
	subName = idAttr[1]
	recordType = idAttr[2]

	if subName == "@" {
		subName = ""
	}

	return domainName, subName, recordType, nil
}

func resourceRRSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dsc.Client)

	var diags diag.Diagnostics

	domainName, subName, recordType, err := namesFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	r, err := c.Records.Get(domainName, subName, recordType)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	rrsetIntoSchema(r, d)

	return diags
}

func resourceRRSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dsc.Client)

	domainName, subName, recordType, err := namesFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	r := schemaToRRset(d)
	rrset, err := c.Records.Update(domainName, subName, recordType, r)
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	rrsetIntoSchema(rrset, d)
	return diags
}

func resourceRRSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dsc.Client)

	domainName, subName, recordType, err := namesFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.Records.Delete(domainName, subName, recordType)
	if err != nil && !isNotFoundError(err) {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func rrsetIntoSchema(r *dsc.RRSet, d *schema.ResourceData) {
	d.Set("created", r.Created.Format(time.RFC3339))
	d.Set("domain", r.Domain)
	d.Set("name", r.Name)
	d.Set("records", r.Records)
	d.Set("subname", r.SubName)
	d.Set("ttl", r.TTL)
	d.Set("type", r.Type)

	id := idFromNames(r.Domain, r.SubName, r.Type)
	d.SetId(id)
}

func schemaToRRset(d *schema.ResourceData) dsc.RRSet {
	r := dsc.RRSet{
		Domain:  d.Get("domain").(string),
		SubName: d.Get("subname").(string),
		Type:    d.Get("type").(string),
	}

	recs := d.Get("records").(*schema.Set)
	r.Records = make([]string, recs.Len())
	for i, rec := range recs.List() {
		r.Records[i] = rec.(string)
	}
	r.TTL = d.Get("ttl").(int)

	return r
}
