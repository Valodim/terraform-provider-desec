package desec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDesecRRsetBasic(t *testing.T) {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
		return
	}
	domainName := fmt.Sprintf("%s.example", uuid)
	name := fmt.Sprintf("test.%s.", domainName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDesecDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDesecRRsetConfigBasic(domainName, "[ \"127.0.0.1\" ]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("desec_rrset.desec-example_hello_A", "name", name),
					resource.TestCheckResourceAttr("desec_rrset.desec-example_hello_A", "records.#", "1"),
					testAccCheckDesecRRsetExists,
				),
			},
			{
				Config: testAccCheckDesecRRsetConfigBasic(domainName, "[ \"127.0.0.1\", \"127.0.0.2\" ]"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("desec_rrset.desec-example_hello_A", "records.#", "2"),
					testAccCheckDesecRRsetExists,
				),
			},
		},
	})
}

func testAccCheckDesecRRsetConfigBasic(domainName, ips string) string {
	return fmt.Sprintf(`
	resource "desec_domain" "desec-example" {
		name = "%s"
	}
	resource "desec_rrset" "desec-example_hello_A" {
		domain = desec_domain.desec-example.name
		subname = "test"
		type = "A"
		records = %s
		ttl = 3600
	}
	resource "desec_rrset" "desec-example_three_TXT" {
		domain = desec_domain.desec-example.name
		subname = "three"
		type = "TXT"
		records = ["one", "two", "three"]
		ttl = 3600
	}
	resource "desec_rrset" "desec-example_quotes_TXT" {
		domain = desec_domain.desec-example.name
		subname = "quotes"
		type = "TXT"
		records = ["\"txt content with quotes\""]
		ttl = 3600
	}
	resource "desec_rrset" "desec-example_long_TXT" {
		domain = desec_domain.desec-example.name
		subname = "long"
		type = "TXT"
		records = ["txt content very long XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"]
		ttl = 3600
	}
	`, domainName, ips)
}

func testAccCheckDesecRRsetExists(s *terraform.State) error {
	c := testAccProvider.Meta().(*DesecConfig).client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "desec_rrset" {
			continue
		}

		domainName, subName, recordType, err := namesFromId(rs.Primary.ID)
		_, err = c.Records.Get(domainName, subName, recordType)
		if err != nil {
			return nil
		}
	}

	return nil
}
