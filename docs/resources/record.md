---
page_title: "record Resource - terraform-provider-desec"
subcategory: ""
description: |-
  Provides a desec RRSet resource.
---

# Resource `desec_rrset`

The record resource maps to the [RRSet API](https://desec.readthedocs.io/en/latest/dns/rrsets.html)
of [desec.io](https://desec.io).

## Example Usage

```terraform
resource "desec_rrset" "hello-a" {
  domain = "desec.example"
  subname = "hello"
  type = "A"
  records = [ "127.0.0.3" ]
  ttl = 3600
}
```

## Argument Reference

A record set is identified by `domain`, `subname`, and `type`.

- `domain` - (Required) The record's domain part.
- `subname` - (Required) The record's subdomain part. May be empty string to denote the zone apex.
- `type` - (Required) The record type. Such as A, AAAA, ...

Each record set contains `records` and `ttl`.

- `records` - (Required) The record content, as a set of strings.
- `ttl` - (Required) The TTL to set for the records, must be an integer.

## Import

RRSets can be imported using a composite ID formed of domain name, subdomain name, and type.

```
$ terraform import desec_rrset.hello-a desec.example/hello/A
```

where:

* `desec.example` - The domain name.
* `hello` - The subdomain name. Can be `@` to denote the zone apex (i.e. the domain name itself).
* `A` - The record type.
