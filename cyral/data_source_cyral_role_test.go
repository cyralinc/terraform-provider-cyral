package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TODO: currently, we just test that the configs are valid. We need to add ACC
// tests for a full scenario, containing roles, IdPs and user groups (see
// resources `cyral_role` and `cyral_role_sso_groups`. -aholmquist 2022-08-05
/*
func roleDataSourceTestUserGroupsAndRoleNames() ([]*UserGroup, []string) {
	return []*UserGroup{
		{
			Name:        "tf-provider-test-user-group-1",
			Description: "description-1",
		},
		{
			Name:        "tf-provider-test-user-group-2",
			Description: "description-2",
		},
	}, []string{
		"tf-provider-test-role-1",
		"tf-provider-test-role-2",
	}
}
*/

func TestAccRoleDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: roleDataSourceConfig(
					"main_test",
					"tf-provider-test-user-group-1",
					[]string{}),
			},
			{
				Config: roleDataSourceConfig(
					"main_test",
					"tf-provider-test-user-group-2",
					[]string{}),
			},
		},
	})
}

func roleDataSourceConfig(dsourceName, nameFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_role" "%s" {
		name = "%s"
		depends_on = [%s]
	}`, dsourceName, nameFilter, listToStr(dependsOn))
}
