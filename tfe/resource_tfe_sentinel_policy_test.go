package tfe

import (
	"fmt"
	"testing"

	tfe "github.com/HappyPathway/go-tfe"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTFESentinelPolicy_basic(t *testing.T) {
	policy := &tfe.Policy{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFESentinelPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTFESentinelPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESentinelPolicyExists(
						"tfe_sentinel_policy.foobar", policy),
					testAccCheckTFESentinelPolicyAttributes(policy),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "name", "policy-test"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "policy", "main = rule { true }"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "enforce_mode", "hard-mandatory"),
				),
			},
		},
	})
}

func TestAccTFESentinelPolicy_update(t *testing.T) {
	policy := &tfe.Policy{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFESentinelPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccTFESentinelPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESentinelPolicyExists(
						"tfe_sentinel_policy.foobar", policy),
					testAccCheckTFESentinelPolicyAttributes(policy),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "name", "policy-test"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "policy", "main = rule { true }"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "enforce_mode", "hard-mandatory"),
				),
			},

			resource.TestStep{
				Config: testAccTFESentinelPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESentinelPolicyExists(
						"tfe_sentinel_policy.foobar", policy),
					testAccCheckTFESentinelPolicyAttributesUpdated(policy),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "name", "policy-test"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "policy", "main = rule { false }"),
					resource.TestCheckResourceAttr(
						"tfe_sentinel_policy.foobar", "enforce_mode", "soft-mandatory"),
				),
			},
		},
	})
}

func testAccCheckTFESentinelPolicyExists(
	n string, policy *tfe.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tfeClient := testAccProvider.Meta().(*tfe.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		p, err := tfeClient.Policies.Read(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if p.ID != rs.Primary.ID {
			return fmt.Errorf("SentinelPolicy not found")
		}

		*policy = *p

		return nil
	}
}

func testAccCheckTFESentinelPolicyAttributes(
	policy *tfe.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if policy.Name != "policy-test" {
			return fmt.Errorf("Bad name: %s", policy.Name)
		}

		if policy.Enforce[0].Mode != "hard-mandatory" {
			return fmt.Errorf("Bad enforce mode: %s", policy.Enforce[0].Mode)
		}

		return nil
	}
}

func testAccCheckTFESentinelPolicyAttributesUpdated(
	policy *tfe.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if policy.Name != "policy-test" {
			return fmt.Errorf("Bad name: %s", policy.Name)
		}

		if policy.Enforce[0].Mode != "soft-mandatory" {
			return fmt.Errorf("Bad enforce mode: %s", policy.Enforce[0].Mode)
		}

		return nil
	}
}

func testAccCheckTFESentinelPolicyDestroy(s *terraform.State) error {
	tfeClient := testAccProvider.Meta().(*tfe.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_sentinel_policy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := tfeClient.Policies.Read(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Sentinel policy %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTFESentinelPolicy_basic = `
resource "tfe_organization" "foobar" {
  name = "terraform-test"
  email = "admin@company.com"
}

resource "tfe_sentinel_policy" "foobar" {
  name = "policy-test"
  organization = "${tfe_organization.foobar.id}"
  policy = "main = rule { true }"
  enforce_mode = "hard-mandatory"
}`

const testAccTFESentinelPolicy_update = `
resource "tfe_organization" "foobar" {
  name = "terraform-test"
  email = "admin@company.com"
}

resource "tfe_sentinel_policy" "foobar" {
  name = "policy-test"
  organization = "${tfe_organization.foobar.id}"
  policy = "main = rule { false }"
  enforce_mode = "soft-mandatory"
}`
