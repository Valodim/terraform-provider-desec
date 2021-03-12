---
page_title: "Provider: desec.io"
subcategory: ""
description: |-
  Terraform provider for interacting with desec.io API.
---

# desec.io Provider

The [desec.io](https://desec.io) provider is used to DNS using the [desec.io API](https://desec.readthedocs.io/en/latest/index.html).

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "desec" {
  api_token = "abcdefghijklmn-opqrstuvwxyz1"
}
```

## Schema

- **api_token** (String) API token to authenticate to the service.
- **api_uri** (String, Optional) The API base URI to use. Defaults to `https://desec.io/api/v1/`
- **limit_read** (Integer, Optional) Maximum number of read API requests to send, per second. Defaults to 8.
- **limit_write** (Integer, Optional) Maximum number of write API requests to send, per second. Defaults to 5.

