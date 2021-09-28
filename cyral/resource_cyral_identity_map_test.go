package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialIdentityMapConfig IdentityMapResource = IdentityMapResource{
	IdentityType: "user",
	IdentityName: "tf-test-identity-map",
	AccessDuration: &AccessDuration{
		Days:    7,
		Hours:   10,
		Minutes: 30,
		Seconds: 10,
	},
}

var updatedIdentityMapConfig IdentityMapResource = IdentityMapResource{
	IdentityType: "user",
	IdentityName: "tf-test-identity-map",
	AccessDuration: &AccessDuration{
		Days:    0,
		Hours:   0,
		Minutes: 0,
		Seconds: 0,
	},
}

var identityMapConfigWithoutAccessDuration IdentityMapResource = IdentityMapResource{
	IdentityType: "user",
	IdentityName: "tf-test-identity-map",
}

func TestAccIdentityMapResource(t *testing.T) {
	testConfig, testFunc := setupIdentityMapTest(initialIdentityMapConfig)
	testUpdateConfig, testUpdateFunc := setupIdentityMapTest(updatedIdentityMapConfig)
	testWithoutAccessDurationConfig, testWithoutAccessDurationFunc :=
		setupIdentityMapTest(identityMapConfigWithoutAccessDuration)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				Config: testWithoutAccessDurationConfig,
				Check:  testWithoutAccessDurationFunc,
			},
		},
	})
}

func setupIdentityMapTest(integrationData IdentityMapResource) (string, resource.TestCheckFunc) {
	configuration := formatIdentityMapIntegrationDataIntoConfig(integrationData)
	var testFunction resource.TestCheckFunc
	if integrationData.AccessDuration != nil {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"identity_type", integrationData.IdentityType),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"identity_name", integrationData.IdentityName),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.#", "1"),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.0.days", fmt.Sprintf("%d", integrationData.AccessDuration.Days)),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.0.hours", fmt.Sprintf("%d", integrationData.AccessDuration.Hours)),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.0.minutes", fmt.Sprintf("%d", integrationData.AccessDuration.Minutes)),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.0.seconds", fmt.Sprintf("%d", integrationData.AccessDuration.Seconds)),
		)
	} else {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"identity_type", integrationData.IdentityType),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"identity_name", integrationData.IdentityName),
			resource.TestCheckResourceAttr("cyral_identity_map.tf_test_cyral_sidecar_template",
				"access_duration.#", "0"),
		)
	}
	return configuration, testFunction
}

func formatIdentityMapIntegrationDataIntoConfig(data IdentityMapResource) string {
	config := `
	resource "cyral_repository" "test_repo_repository" {
		type = "mongodb"
		host = "mongo.local"
		port = 3333
		name = "tf-repo-test"
	}
	
	resource "cyral_repository_local_account" "tf_test_repository_account" {
		repository_id = cyral_repository.test_repo_repository.id
		enviroment_variable {
			database_name = "tf_test_db_name"
			local_account = "tf_test_repo_account"
			variable_name = "CYRAL_DBSECRETS_TF_TEST_VARIABLE_NAME"
		}
	}
	
	`
	if data.AccessDuration != nil {
		config = fmt.Sprintf(`%s
	  resource "cyral_identity_map" "tf_test_cyral_sidecar_template" {
			repository_id               = cyral_repository.test_repo_repository.id
			repository_local_account_id = cyral_repository_local_account.tf_test_repository_account.id
			identity_type               = "%s"
			identity_name               = "%s"
			access_duration {
				days    = %d
				hours   = %d
				minutes = %d
				seconds = %d
			}
	  }`, config, data.IdentityType, data.IdentityName, data.AccessDuration.Days, data.AccessDuration.Hours, data.AccessDuration.Minutes, data.AccessDuration.Seconds)
	} else {
		config = fmt.Sprintf(`%s
	  resource "cyral_identity_map" "tf_test_cyral_sidecar_template" {
			repository_id               = cyral_repository.test_repo_repository.id
			repository_local_account_id = cyral_repository_local_account.tf_test_repository_account.id
			identity_type               = "%s"
			identity_name               = "%s"
	  }`, config, data.IdentityType, data.IdentityName)
	}
	return config
}
