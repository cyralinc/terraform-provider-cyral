package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialRoleSSOGroupsConfig = RoleSSOGroupsData{
	Id:             "bb1ee321-0cd4-4471-b440-c18fd5fdff46",
	ConnectionName: "Okta",
	GroupName:      "group-test-1",
}

func TestAccRoleSSOGroupsResource(t *testing.T) {
	testConfig, testFunc := setupRoleSSOGroupsTest(initialRoleSSOGroupsConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
		},
	})
}

func setupRoleSSOGroupsTest(mapSSOGroupsData RoleSSOGroupsData) (string, resource.TestCheckFunc) {

	configuration := formatRoleSSOGroupsDataIntoConfig(mapSSOGroupsData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_map_sso_groups.test_map_sso_groups",
			"id", mapSSOGroupsData.Id),
		resource.TestCheckResourceAttr("cyral_map_sso_groups.test_map_sso_groups",
			"connectionName", mapSSOGroupsData.ConnectionName),
		resource.TestCheckResourceAttr("cyral_map_sso_groups.test_map_sso_groups",
			"groupName", mapSSOGroupsData.GroupName),
	)

	return configuration, testFunction
}

func formatRoleSSOGroupsDataIntoConfig(data RoleSSOGroupsData) string {
	return fmt.Sprintf(`
      resource "cyral_map_sso_groups" "test_map_sso_groups" {
      	id = "%s"
				connectionName = "%s"
				groupName = "%s"
      }`, data.Id, data.ConnectionName, data.GroupName)
}
