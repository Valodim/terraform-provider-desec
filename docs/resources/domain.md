---
page_title: "domain Resource - terraform-provider-desec"
subcategory: ""
description: |-
  Provides a desec domain resource.
---

# Resource `desec_domain`

The record resource maps to the [domain API](https://desec.readthedocs.io/en/latest/dns/domains.html)
of [desec.io](https://desec.io).

## Example Usage

```terraform
resource "desec_domain" "desec-example" {
  name = "desec.example"
}
```

## Argument Reference

A domain is identified only by its `name`.

- `name` - (Required) The domain name.

## Attributes Reference

- `id` - The domain ID. The content of this field is identical to `name`.
- `created` - An RFC3339 timestamp of when the domain entry was created.
- `keys` - A list of DNSSEC domain keys.
- `minimum_ttl` - This domain's minimum TTL value.
- `published` - An RFC3339 timestamp of when the domain was last published.

## Import

Domains can be imported using the domain name as ID.

