---
page_title: "token Resource - terraform-provider-desec"
subcategory: ""
description: |-
  Provides a desec token resource.
---

# Resource `desec_token`

The token resource maps to the [token API](https://desec.readthedocs.io/en/latest/auth/tokens.html)
of [desec.io](https://desec.io).

## Example Usage

```terraform
resource "desec_token" "example" {
  # all fields are optional
	name = "example token"
	perm_create_domain = false
	perm_delete_domain = true
	perm_manage_tokens = false
}
```

## Argument Reference

A token is identified by its `id`, which is assigned server-side.

## Attributes Reference

- `id` - The token ID, assigned server-side.

- `created` - (Read-Only) An RFC3339 timestamp of when the token entry was created.
- `owner` - (Read-Only) The owner account email address who created this token.
- `token` - (Read-Only) The actual token for import. SEE NOTE ON TOKEN ATTRIBUTE BELOW

- `allowed_subnets` - Exhaustive list of IP addresses or subnets clients must connect from in order to successfully authenticate with the token. Defaults to no restriction.
- `auto_policy` - When using this token to create a domain, automatically configure a permissive scoping policy for it.
- `name` - Token name. It is meant for user reference only and carries no operational meaning.
- `perm_create_domain` - Permission to create a new domain.
- `perm_delete_domain` - Permission to delete a domain.
- `perm_manage_tokens` - Permission to manage tokens (this one and also all others).

### NOTE ON TOKEN ATTRIBUTE

The `token` attribute is returned only once by the server, on creation of the token resource, but
never afterwards. That means when the resource is initially created, the `token` value becomes part
of the terraform state, and can be viewed with `terraform show`. It will be emptied when the state
is next refreshed.

## Import

Tokens can be imported by their token id. This is the recommended way to create tokens, since the
secret will otherwise briefly end up in the terraform state (see note above).

