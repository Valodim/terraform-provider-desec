---
page_title: "token policy Resource - terraform-provider-desec"
subcategory: ""
description: |-
  Provides a desec token policy resource.
---

# Resource `desec_token_policy`

The token resource maps to the [token policy API](https://desec.readthedocs.io/en/latest/auth/tokens.html#token-scoping-policies)
of [desec.io](https://desec.io).

## Example Usage

```terraform
# default policy for the token
resource "desec_token_policy" "example-token-default" {
	token_id = desec_token.example.id
	perm_write = false
}

resource "desec_token_policy" "example-token-example-domain-write" {
	token_id = desec_token.example.id
	domain = desec_domain.example.name
	perm_write = true

	# desec will reject creating policies if no default policy exists for a token
	# without this line, the creation request may happen before the default policy and fail
  depends_on = [desec_token_policy.example-token-default]
}
```

## Argument Reference

A token is identified by its `id`.

## Attributes Reference

- `id` - (Read-Only) The id of the token policy. Note this policy is scoped to the `token_id`.
- `token_id` - (Required, Read-Only) The id of the token this policy applies to.
- `perm_write` - (Required) Indicates write permission for the RRset specified by (domain, subname, type) when using the general RRset management or dynDNS interface. Defaults to false.

- `domain` - Domain name to which the policy applies. Empty string (= null) for the default policy.
- `subname` - Subname to which the policy applies. Empty string (= null) for the default policy.
- `type` - Record type to which the policy applies. Empty string (= null) for the default policy.

## Import

Token policies can be imported by combining the token and policy ids: `$TOKEN_ID/$POLICY_ID`

