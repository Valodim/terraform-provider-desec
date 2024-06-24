terraform {
  required_providers {
    desec = {
      source = "Valodim/desec"
      version = "0.5.0"
    }
  }
}

provider "desec" {
  api_token = "abcdefghijklmn-opqrstuvwxyz1"
}

resource "desec_domain" "desec-example" {
  name = "desec.example"
}

// desec.example.  IN  A  127.0.0.1
resource "desec_rrset" "desec-example_A" {
  domain = desec_domain.desec-example.name
  subname = ""
  type = "A"
  records = [ "127.0.0.1" ]
  ttl = 3600
}

// desec.example.hello.  IN  A  127.0.0.2
resource "desec_rrset" "desec-example_hello_A" {
  domain = desec_domain.desec-example.name
  subname = "hello"
  type = "A"
  records = [ "127.0.0.1" ]
  ttl = 3600
}
