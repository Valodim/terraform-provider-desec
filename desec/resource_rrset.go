package desec

import (
	"context"
	"fmt"
	"reflect"
	"sort"
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 178),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CAA", "CERT", "CNAME", "DNSKEY", "DS", "LOC", "MX", "NAPTR", "NS", "PTR", "SMIMEA", "SPF", "SRV", "SSHFP", "TLSA", "TXT", "URI"}, false),
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// records may be reordered by the server.
					// we cheat here, by checking the whole list as a sorted set of normalized values.
					// this is inefficient, but luckily the lists are always very small.
					o, n := d.GetChange("records")
					if (o == nil) != (n == nil) {
						return false
					}
					no := normalizeRecordSetInterface(o.(*schema.Set).List())
					nn := normalizeRecordSetInterface(n.(*schema.Set).List())
					return reflect.DeepEqual(no, nn)
				},
			},
			"ttl": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(60, 604800),
			},
		},
	}
}

func resourceRRSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	var diags diag.Diagnostics

	r := schemaToRRset(d)
	rrset, err := c.Records.Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	rrsetIntoSchema(rrset, d)
	return diags
}

func resourceRRSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	c := conf.client

	var diags diag.Diagnostics

	r, err := conf.cache.GetRRSetById(ctx, c, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if r == nil {
		d.SetId("")
	} else {
		rrsetIntoSchema(r, d)
	}

	return diags
}

func resourceRRSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	domainName, subName, recordType, err := namesFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	r := schemaToRRset(d)
	rrset, err := c.Records.Update(ctx, domainName, subName, recordType, r)
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
	conf := m.(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	domainName, subName, recordType, err := namesFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.Records.Delete(ctx, domainName, subName, recordType)
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
	d.Set("subname", r.SubName)
	d.Set("ttl", r.TTL)
	d.Set("type", r.Type)
	d.Set("records", normalizeRecordSet(r.Records))

	id := idFromNames(r.Domain, r.SubName, r.Type)
	d.SetId(id)
}

func schemaToRRset(d *schema.ResourceData) dsc.RRSet {
	rtype := d.Get("type").(string)
	r := dsc.RRSet{
		Domain:  d.Get("domain").(string),
		SubName: d.Get("subname").(string),
		Type:    rtype,
	}

	recs := d.Get("records").(*schema.Set)
	r.Records = make([]string, recs.Len())
	for i, rec := range recs.List() {
		if (rtype == "TXT" || rtype == "SPF") && rec.(string)[0] != '"' {
			r.Records[i] = fmt.Sprintf("\"%s\"", rec.(string))
		} else {
			r.Records[i] = rec.(string)
		}
	}
	r.TTL = d.Get("ttl").(int)

	return r
}

func normalizeRecordSetInterface(s []interface{}) []string {
	result := make([]string, len(s))
	for i, rec := range s {
		result[i] = normalizeLongRecord(rec.(string))
	}
	sort.Strings(result)
	return result
}

func normalizeRecordSet(s []string) []string {
	result := make([]string, len(s))
	for i, rec := range s {
		result[i] = normalizeLongRecord(rec)
	}
	sort.Strings(result)
	return result
}

func normalizeLongRecord(s string) string {
	if s == "" || s[0] != '"' {
		return s
	}
	s = strings.Trim(s, "\"")
	s = strings.Replace(s, "\" \"", "", -1)
	return s
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
