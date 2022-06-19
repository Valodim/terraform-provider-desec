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
- **retry_max** (Integer, Optional) The max number of retries when sending an API request. The default value is determined by the deSEC API client [implementation](https://github.com/nrdcg/desec).
