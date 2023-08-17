package cyral

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	repositoryConfAnalysisResourceName = "repository-conf-analysis"
)

func repositoryConfAnalysisSampleRepositoryConfig() string {
	return utils.FormatBasicRepositoryIntoConfig(
		BasicRepositoryResName,
		utils.AccTestName(repositoryConfAnalysisResourceName, "repository"),
		"postgresql",
		"some-hostname",
		3067,
	)
}

func TestAccRepositoryConfAnalysisResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorRedact(),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorAnnotationGroups(),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorLogGroups(),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config: testAccRepoConfAnalysisConfig_DefaultValues(),
				Check:  testAccRepoConfAnalysisCheck_DefaultValues(),
			},
			{
				Config: testAccRepoConfAnalysisConfig_Updated(),
				Check:  testAccRepoConfAnalysisCheck_Updated(),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_repository_conf_analysis.test_conf_analysis",
			},
		},
	})
}

func TestRepositoryConfAnalysisResourceUpgradeV0(t *testing.T) {
	previousState := map[string]interface{}{
		"id":            "repositoryID/ConfAnalysis",
		"repository_id": "repositoryID",
	}
	actualNewState, err := upgradeRepositoryConfAnalysisV0(context.Background(),
		previousState, nil)
	require.NoError(t, err)
	expectedNewState := map[string]interface{}{
		"id":            "repositoryID",
		"repository_id": "repositoryID",
	}
	require.Equal(t, expectedNewState, actualNewState)
}

func testAccRepoConfAnalysisConfig_ErrorRedact() string {
	var config string
	config += repositoryConfAnalysisSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = %s
		redact = "some-invalid-value"
	}`, BasicRepositoryID)
	return config
}

func testAccRepoConfAnalysisConfig_ErrorAnnotationGroups() string {
	var config string
	config += repositoryConfAnalysisSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = %s
		comment_annotation_groups = [
			"some-invalid-value"
		]
	}`, BasicRepositoryID)
	return config
}

func testAccRepoConfAnalysisConfig_ErrorLogGroups() string {
	var config string
	config += repositoryConfAnalysisSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = %s
		log_groups = [
			"some-invalid-value"
		]
	}`, BasicRepositoryID)
	return config
}

func testAccRepoConfAnalysisConfig_DefaultValues() string {
	var config string
	config += repositoryConfAnalysisSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = %s
	}`, BasicRepositoryID)
	return config
}

func testAccRepoConfAnalysisCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"alert_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"block_on_violation", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"comment_annotation_groups.#", "0"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_filter_analysis", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_pre_configured_alerts", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"enable_data_masking", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.#", "0"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"redact", "all"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"enable_dataset_rewrites", "false"),
	)
}

func testAccRepoConfAnalysisConfig_Updated() string {
	var config string
	config += repositoryConfAnalysisSampleRepositoryConfig()
	config += fmt.Sprintf(`
	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = %s
		redact = "all"
		alert_on_violation = true
		disable_pre_configured_alerts = false
		block_on_violation = true
		disable_filter_analysis = false
		enable_dataset_rewrites = true
		enable_data_masking = true
		comment_annotation_groups = [
			"identity"
		]
		log_groups= [
			"everything",
			"sensitive & dql",
			"sensitive & dml",
			"sensitive & ddl"
		]
	}`, BasicRepositoryID)
	return config
}

func testAccRepoConfAnalysisCheck_Updated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"redact", "all"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"alert_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_pre_configured_alerts", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"block_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_filter_analysis", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"enable_data_masking", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"enable_dataset_rewrites", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"comment_annotation_groups.#", "1"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.#", "4"),
		resource.TestCheckTypeSetElemAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.*", "everything"),
		resource.TestCheckTypeSetElemAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.*", "sensitive & dql"),
		resource.TestCheckTypeSetElemAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.*", "sensitive & dml"),
		resource.TestCheckTypeSetElemAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.*", "sensitive & ddl"),
	)
}
