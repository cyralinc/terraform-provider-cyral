package role_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: More tests -aholmquist 2022-08-29

func TestAccRoleDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: roleDataSourceConfig(
					"main_test",
					"",
					[]string{}),
			},
		},
	})
}

func roleDataSourceConfig(dsourceName, nameFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_role" "%s" {
		name = "%s"
		depends_on = %s
	}`, dsourceName, nameFilter, utils.ListToStr(dependsOn))
}
