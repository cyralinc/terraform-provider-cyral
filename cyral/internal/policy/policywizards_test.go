package policyv2

import (
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPolicyWizardsDataSource(t *testing.T) {
	dsName := "data.cyral_policy_wizards.wizard_list"
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "cyral_policy_wizards" "wizard_list" {
}
`,
				Check: checkAllWizards(dsName),
			},
			{
				Config: `
data "cyral_policy_wizards" "wizard_list" {
  wizard_id = "data-firewall"
}
`,
				Check: checkOneWizard(dsName, "data-firewall"),
			},
			{
				Config: `
data "cyral_policy_wizards" "wizard_list" {
  wizard_id = "XXX"
}
`,
				Check: resource.TestCheckResourceAttr(dsName, "wizards.#", "0"),
			},
		},
	})
}

// checkAllWizards ensures that a few well known wizard ids are present in the
// datasource state. It does not attempt to make very exhaustive checks because
// wizard names, descriptions (and even the wizard list) is subject to change.
func checkAllWizards(dsName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckTypeSetElemNestedAttrs(
			dsName, "wizards.*",
			map[string]string{
				"id": "data-firewall",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			dsName, "wizards.*",
			map[string]string{
				"id": "data-masking",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs(
			dsName, "wizards.*",
			map[string]string{
				"id": "user-segmentation",
			},
		),
	)
}

// checkOneWizard ensures that the data source state contains only one wizard
// with the given id.
func checkOneWizard(dsName, id string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(dsName, "wizards.#", "1"),
		resource.TestCheckTypeSetElemNestedAttrs(
			dsName, "wizards.*",
			map[string]string{
				"id": id,
			},
		),
	)
}
