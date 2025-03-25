resource "desec_token" "example" {
	name = "test token"
	perm_create_domain = false
	perm_delete_domain = true
	perm_manage_tokens = false
}

resource "desec_token_policy" "example-default" {
	token_id = desec_token.example.id
	perm_write = false
}

resource "desec_token_policy" "example" {
	token_id = desec_token.example.id
	domain = desec_domain.example.name
	perm_write = true

	depends_on = [desec_token_policy.example-default]
}
