package role_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/role"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const (
	roleSSOGroupsResourceName = "role-sso-groups"

	testRoleSSOGroupsRoleResName            = "test_role"
	testRoleSSOGroupsRoleFullResName        = "cyral_role.test_role"
	testRoleSSOGroupsIntegrationResName     = "test_integration"
	testRoleSSOGroupsIntegrationFullResName = "cyral_integration_idp_okta.test_integration"
)

func roleSSOGroupsTestRoleName() string {
	return utils.AccTestName(roleSSOGroupsResourceName, "role")
}

func roleSSOGroupsTestRole() string {
	return fmt.Sprintf(`
	resource "cyral_role" "%s" {
		name="%s"
	}`, testRoleSSOGroupsRoleResName, roleSSOGroupsTestRoleName())
}

func roleSSOGroupsTestOktaIntegration() string {
	return utils.FormatBasicIntegrationIdPOktaIntoConfig(
		testRoleSSOGroupsIntegrationResName,
		utils.AccTestName(roleSSOGroupsResourceName, "integration"),
		utils.TestSingleSignOnURL,
	)
}

func TestAccRoleSSOGroupsResource(t *testing.T) {
	// This resource needs to exist when the last step executes, which is an
	// Import test.
	importTestResourceName := "cyral_role_sso_groups.test_role_sso_groups"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyRoleId(),
				ExpectError: regexp.MustCompile(`The argument "role_id" is required`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptySSOGroup(),
				ExpectError: regexp.MustCompile(`At least 1 "sso_group" blocks are required.`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyGroupName(),
				ExpectError: regexp.MustCompile(`The argument "group_name" is required`),
			},
			{
				Config:      testAccRoleSSOGroupsConfig_EmptyIdPID(),
				ExpectError: regexp.MustCompile(`The argument "idp_id" is required`),
			},
			{
				Config: testAccRoleSSOGroupsConfig_SingleSSOGroup(),
				Check:  testAccRoleSSOGroupsCheck_SingleSSOGroup(),
			},
			{
				Config: testAccRoleSSOGroupsConfig_MultipleSSOGroups(),
				Check:  testAccRoleSSOGroupsCheck_MultipleSSOGroups(),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      importTestResourceName,
			},
		},
	})
}

func TestRoleSSOGroupsResourceUpgradeV0(t *testing.T) {
	previousState := map[string]interface{}{
		"id":      "roleID/SSOGroups",
		"role_id": "roleID",
	}
	actualNewState, err := role.UpgradeRoleSSOGroupsV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":      "roleID",
		"role_id": "roleID",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func testAccRoleSSOGroupsConfig_EmptyRoleId() string {
	var config string
	config += roleSSOGroupsTestOktaIntegration()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		sso_group {
			group_name="Everyone"
			idp_id = %s.id
		}
	}
	`, testRoleSSOGroupsIntegrationFullResName)
	return config
}

func testAccRoleSSOGroupsConfig_EmptySSOGroup() string {
	var config string
	config += roleSSOGroupsTestRole()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id = %s.id
	}`, testRoleSSOGroupsRoleFullResName)
	return config
}

func testAccRoleSSOGroupsConfig_EmptyGroupName() string {
	var config string
	config += roleSSOGroupsTestOktaIntegration()
	config += roleSSOGroupsTestRole()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id = %s.id
		sso_group {
			idp_id = %s.id
		}
	}`, testRoleSSOGroupsRoleFullResName, testRoleSSOGroupsIntegrationFullResName)
	return config
}

func testAccRoleSSOGroupsConfig_EmptyIdPID() string {
	var config string
	config += roleSSOGroupsTestOktaIntegration()
	config += roleSSOGroupsTestRole()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id = %s.id
		sso_group {
			group_name="Everyone"
		}
	}`, testRoleSSOGroupsRoleFullResName)
	return config
}

func testAccRoleSSOGroupsConfig_SingleSSOGroup() string {
	var config string
	config += roleSSOGroupsTestOktaIntegration()
	config += roleSSOGroupsTestRole()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id = %s.id
		sso_group {
			group_name="Everyone"
			idp_id = %s.id
		}
	}
	`, testRoleSSOGroupsRoleFullResName, testRoleSSOGroupsIntegrationFullResName)
	return config
}

func testAccRoleSSOGroupsCheck_SingleSSOGroup() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "role_id",
			testRoleSSOGroupsRoleFullResName, "id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.#", "1"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.group_name", "Everyone"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_id",
			testRoleSSOGroupsIntegrationFullResName, "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_name",
			testRoleSSOGroupsIntegrationFullResName, "samlp.0.display_name"),
	)
}

func testAccRoleSSOGroupsConfig_MultipleSSOGroups() string {
	var config string
	config += roleSSOGroupsTestOktaIntegration()
	config += roleSSOGroupsTestRole()
	config += fmt.Sprintf(`
	resource "cyral_role_sso_groups" "test_role_sso_groups" {
		role_id = %s.id
		sso_group {
			group_name="Admin"
			idp_id = %s.id
		}
		sso_group {
			group_name="Dev"
			idp_id = %s.id
		}
	}`, testRoleSSOGroupsRoleFullResName, testRoleSSOGroupsIntegrationFullResName,
		testRoleSSOGroupsIntegrationFullResName)
	return config
}

func testAccRoleSSOGroupsCheck_MultipleSSOGroups() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "role_id",
			testRoleSSOGroupsRoleFullResName, "id"),
		resource.TestCheckResourceAttr("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.#", "2"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.0.id"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.*", map[string]string{"group_name": "Admin"}),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_id",
			testRoleSSOGroupsIntegrationFullResName, "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.0.idp_name",
			testRoleSSOGroupsIntegrationFullResName, "samlp.0.display_name"),
		resource.TestCheckResourceAttrSet("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.1.id"),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_role_sso_groups.test_role_sso_groups",
			"sso_group.*", map[string]string{"group_name": "Dev"}),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.1.idp_id",
			testRoleSSOGroupsIntegrationFullResName, "id"),
		resource.TestCheckResourceAttrPair(
			"cyral_role_sso_groups.test_role_sso_groups", "sso_group.1.idp_name",
			testRoleSSOGroupsIntegrationFullResName, "samlp.0.display_name"),
	)
}
