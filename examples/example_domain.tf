resource "desec_domain" "example" {
  name = "f01981a8-ad38-44e5-9505-c999c6e6c669.de"
}

// desec.example.  IN  A  127.0.0.1
resource "desec_rrset" "desec-example_A" {
  domain = desec_domain.example.name
  subname = ""
  type = "A"
  records = [ "127.0.0.1" ]
  ttl = 3600
}

// desec.example.hello.  IN  A  127.0.0.2
resource "desec_rrset" "desec-example_hello_A" {
  domain = desec_domain.example.name
  subname = "hello"
  type = "A"
  records = [ "127.0.0.1" ]
  ttl = 3600
}
