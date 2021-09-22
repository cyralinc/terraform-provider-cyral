package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRepositoryConfAnalysisResource(t *testing.T) {
	repoName := "tf-test-repository"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorRedact(repoName),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorAnnotationGroups(repoName),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config:      testAccRepoConfAnalysisConfig_ErrorLogGroups(repoName),
				ExpectError: regexp.MustCompile(`Error running pre-apply refresh: exit status 1`),
			},
			{
				Config: testAccRepoConfAnalysisConfig_DefaultValues(repoName),
				Check:  testAccRepoConfAnalysisCheck_DefaultValues(),
			},
			{
				Config: testAccRepoConfAnalysisConfig_Updated(repoName),
				Check:  testAccRepoConfAnalysisCheck_Updated(),
			},
		},
	})
}

func testAccRepoConfAnalysisConfig_ErrorRedact(repositoryName string) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo" {
		type = "postgresql"
		host = "some-hostname"
		port = "3067"
		name = "%s"
	}

	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = cyral_repository.test_repo.id
		redact = "some-invalid-value"
	}
	`, repositoryName)
}

func testAccRepoConfAnalysisConfig_ErrorAnnotationGroups(repositoryName string) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo" {
		type = "postgresql"
		host = "some-hostname"
		port = "3067"
		name = "%s"
	}

	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = cyral_repository.test_repo.id
		comment_annotation_groups = [
			"some-invalid-value"
		]
	}
	`, repositoryName)
}

func testAccRepoConfAnalysisConfig_ErrorLogGroups(repositoryName string) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo" {
		type = "postgresql"
		host = "some-hostname"
		port = "3067"
		name = "%s"
	}

	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = cyral_repository.test_repo.id
		log_groups = [
			"some-invalid-value"
		]
	}
	`, repositoryName)
}

func testAccRepoConfAnalysisConfig_DefaultValues(repositoryName string) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo" {
		type = "postgresql"
		host = "some-hostname"
		port = "3067"
		name = "%s"
	}

	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = cyral_repository.test_repo.id
	}
	`, repositoryName)
}

func testAccRepoConfAnalysisCheck_DefaultValues() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"redact", "all"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"tag_sensitive_data", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"ignore_identifier_case", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"analyze_where_clause", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"alert_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_pre_configured_alerts", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"block_on_violation", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_filter_analysis", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"rewrite_on_violation", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"comment_annotation_groups.#", "0"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"log_groups.#", "0"),
	)
}

func testAccRepoConfAnalysisConfig_Updated(repositoryName string) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repo" {
		type = "postgresql"
		host = "some-hostname"
		port = "3067"
		name = "%s"
	}

	resource "cyral_repository_conf_analysis" "test_conf_analysis" {
		repository_id = cyral_repository.test_repo.id
		redact = "all"
		tag_sensitive_data = false
		ignore_identifier_case = false
		analyze_where_clause = false
		alert_on_violation = true
		disable_pre_configured_alerts = false
		block_on_violation = true
		disable_filter_analysis = false
		rewrite_on_violation = true
		comment_annotation_groups = [
			"identity"
		]
		log_groups= [
			"everything",
			"sensitive & dql",
			"sensitive & dml",
			"sensitive & ddl"
		]
	}
	`, repositoryName)
}

func testAccRepoConfAnalysisCheck_Updated() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"redact", "all"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"tag_sensitive_data", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"ignore_identifier_case", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"analyze_where_clause", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"alert_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_pre_configured_alerts", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"block_on_violation", "true"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"disable_filter_analysis", "false"),
		resource.TestCheckResourceAttr("cyral_repository_conf_analysis.test_conf_analysis",
			"rewrite_on_violation", "true"),
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
