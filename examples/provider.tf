terraform {
  required_providers {
    desec = {
      source = "Valodim/desec"
      version = "0.6.1"
    }
  }
}

provider "desec" {
  // loaded from DESEC_API_TOKEN environment
  // api_token = "abcdefghijklmn-opqrstuvwxyz1"

  // default, print tokens as warning
  // token_create_mode = "print"
  // store tokens in state (only until refreshed)
  // token_create_mode = "storeState"
}
