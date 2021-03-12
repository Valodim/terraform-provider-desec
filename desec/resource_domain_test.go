package desec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDesecDomainBasic(t *testing.T) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
		return
	}
	domainName := fmt.Sprintf("%s.example", uuid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDesecDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDesecDomainConfigBasic(domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("desec_domain.desec-example", "name", domainName),
					testAccCheckDesecDomainExists,
				),
			},
		},
	})
}

func testAccCheckDesecDomainDestroy(s *terraform.State) error {
	conf := testAccProvider.Meta().(*DesecConfig)
	conf.cache.Clear()
	c := conf.client

	domains, err := c.Domains.GetAll()
	if err != nil {
		return err
	}
	for _, domain := range domains {
		return fmt.Errorf("Domain still exists:\n%#v", domain.Name)
	}

	return nil
}

func testAccCheckDesecDomainConfigBasic(domainName string) string {
	return fmt.Sprintf(`
	resource "desec_domain" "desec-example" {
		name = "%s"
	}
	`, domainName)
}

func testAccCheckDesecDomainExists(s *terraform.State) error {
	c := testAccProvider.Meta().(*DesecConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "desec_domain" {
			continue
		}

		_, err := c.Domains.Get(rs.Primary.ID)
		if err != nil {
			return nil
		}
	}

	return nil
}
