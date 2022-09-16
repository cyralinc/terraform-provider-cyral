package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repositoryIdentityMapResourceName = "repository-identity-map"
)

var initialIdentityMapConfig RepositoryIdentityMapResource = RepositoryIdentityMapResource{
	IdentityType: "user",
	IdentityName: accTestName(repositoryIdentityMapResourceName, "identity"),
	AccessDuration: &AccessDuration{
		Days:    7,
		Hours:   10,
		Minutes: 30,
		Seconds: 10,
	},
}

var updatedIdentityMapConfig RepositoryIdentityMapResource = RepositoryIdentityMapResource{
	IdentityType: "user",
	IdentityName: accTestName(repositoryIdentityMapResourceName, "identity"),
	AccessDuration: &AccessDuration{
		Days:    0,
		Hours:   0,
		Minutes: 0,
		Seconds: 0,
	},
}

var identityMapConfigWithoutAccessDuration RepositoryIdentityMapResource = RepositoryIdentityMapResource{
	IdentityType: "user",
	IdentityName: accTestName(repositoryIdentityMapResourceName, "identity"),
}

func TestAccRepositoryIdentityMapResource(t *testing.T) {
	testConfig, testFunc :=
		setupRepositoryIdentityMapTest(initialIdentityMapConfig)
	testUpdateConfig, testUpdateFunc :=
		setupRepositoryIdentityMapTest(updatedIdentityMapConfig)
	testWithoutAccessDurationConfig, testWithoutAccessDurationFunc :=
		setupRepositoryIdentityMapTest(identityMapConfigWithoutAccessDuration)

	importStateResName := "cyral_repository_identity_map.test_identity_map"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importStateComposedIDFunc(
					importStateResName,
					[]string{
						"repository_id",
						"identity_type",
						"identity_name",
						"repository_local_account_id",
					},
					"/",
				),
				ResourceName: importStateResName,
			},
		},
	})
}

func setupRepositoryIdentityMapTest(integrationData RepositoryIdentityMapResource) (string, resource.TestCheckFunc) {
	var configuration string
	configuration += formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		"tfprov-test-repository-identity-map-repository",
		"mongodb",
		"mongo.local",
		3333,
	)
	configuration += formatBasicRepositoryLocalAccountIntoConfig_Cyral(
		basicRepositoryLocalAccountResName,
		basicRepositoryID,
		"tfprov-test-repository-identity-map-locaccount",
		"some-password",
	)
	configuration += formatRepositoryIdentityMapDataIntoConfig(
		integrationData, basicRepositoryID, basicRepositoryLocalAccountID)

	var testFunction resource.TestCheckFunc
	if integrationData.AccessDuration != nil {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"identity_type", integrationData.IdentityType),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"identity_name", integrationData.IdentityName),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.#", "1"),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.0.days", fmt.Sprintf("%d", integrationData.AccessDuration.Days)),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.0.hours", fmt.Sprintf("%d", integrationData.AccessDuration.Hours)),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.0.minutes", fmt.Sprintf("%d", integrationData.AccessDuration.Minutes)),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.0.seconds", fmt.Sprintf("%d", integrationData.AccessDuration.Seconds)),
		)
	} else {
		testFunction = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"identity_type", integrationData.IdentityType),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"identity_name", integrationData.IdentityName),
			resource.TestCheckResourceAttr("cyral_repository_identity_map.test_identity_map",
				"access_duration.#", "0"),
		)
	}
	return configuration, testFunction
}

func formatRepositoryIdentityMapDataIntoConfig(
	data RepositoryIdentityMapResource,
	repositoryID, repositoryLocalAccID string,
) string {
	var config string
	if data.AccessDuration != nil {
		config = fmt.Sprintf(`
	resource "cyral_repository_identity_map" "test_identity_map" {
		repository_id               = %s
		repository_local_account_id = %s
		identity_type               = "%s"
		identity_name               = "%s"
		access_duration {
			days    = %d
			hours   = %d
			minutes = %d
			seconds = %d
		}
	}`, repositoryID, repositoryLocalAccID, data.IdentityType, data.IdentityName,
			data.AccessDuration.Days, data.AccessDuration.Hours,
			data.AccessDuration.Minutes, data.AccessDuration.Seconds)
	} else {
		config = fmt.Sprintf(`
	resource "cyral_repository_identity_map" "test_identity_map" {
		repository_id               = %s
		repository_local_account_id = %s
		identity_type               = "%s"
		identity_name               = "%s"
	}`, repositoryID, repositoryLocalAccID, data.IdentityType, data.IdentityName)
	}
	return config
}
