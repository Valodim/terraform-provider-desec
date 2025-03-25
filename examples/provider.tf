terraform {
  required_providers {
    desec = {
      source = "Valodim/desec"
      version = "0.5.0"
    }
  }
}

provider "desec" {
  // loaded from DESEC_API_TOKEN environment
  // api_token = "abcdefghijklmn-opqrstuvwxyz1"
}
