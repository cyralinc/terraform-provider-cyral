package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	TFProvACCPrefix = "tfprov-acc-"

	BasicRepositoryResName        = "test_repository"
	BasicRepositoryID             = "cyral_repository.test_repository.id"
	BasicRepositoryBindingResName = "test_repository_binding"
	BasicSidecarResName           = "test_sidecar"
	BasicSidecarID                = "cyral_sidecar.test_sidecar.id"
	BasicListenerResName          = "test_listener"
	BasicListenerID               = "cyral_sidecar_listener.test_listener.listener_id"
	BasicPolicyResName            = "test_policy"
	BasicPolicyID                 = "cyral_policy.test_policy.id"

	IntegrationIdPResourceName = "integration-idp"
	PolicyResourceName         = "policy"
	RepositoryResourceName     = "repository"
	RoleResourceName           = "role"
	SidecarResourceName        = "sidecar"

	TestSingleSignOnURL = "https://some-test-sso-url.com"
)

// AccTestName attempts to make resource names unique to a specific resource
// type, and avoid name clashes with other resources that exist in the testing
// control plane.
//
// Use this for every resource for which name clashes may occur.
//
// Example usage for cyral_datalabel resource:
//
//	AccTestName("datalabel", "label1")
//
// Example usage for cyral_datalabel data source:
//
//	AccTestName("data-datalabel", "label1")
//
// Note that doing it like above will prevent that the tests attempt to create a
// label called LABEL1 simultaneously, which would cause a failure.
//
// The convention is to use hyphen-separated words if possible, and prefix data
// sources with "data", to distinguish them and their counterpart resources.
func AccTestName(resourceType, suffix string) string {
	return fmt.Sprintf("%s%s-%s", TFProvACCPrefix, resourceType, suffix)
}

func HasAccTestPrefix(name string) bool {
	return strings.HasPrefix(name, TFProvACCPrefix)
}

func FormatBasicRepositoryIntoConfig(resName, repoName, typ, host string, port int) string {
	if typ == "mongodb" {
		return fmt.Sprintf(
			`
		resource "cyral_repository" "%s" {
			name = "%s"
			type = "%s"
			repo_node {
				host = "%s"
				port = %d
			}
			mongodb_settings {
				server_type = "standalone"
			}
		}`, resName, repoName, typ, host, port,
		)
	} else {
		return fmt.Sprintf(
			`
		resource "cyral_repository" "%s" {
			name = "%s"
			type = "%s"
			repo_node {
				host = "%s"
				port = %d
			}
		}`, resName, repoName, typ, host, port,
		)
	}
}

func FormatBasicRepositoryBindingIntoConfig(resName, sidecarID, repositoryID, listenerID string) string {
	return fmt.Sprintf(
		`
	resource "cyral_repository_binding" "%s" {
		sidecar_id    = %s
		repository_id = %s
		listener_binding {
			listener_id = %s
			node_index = 0
		}
	}`, resName, sidecarID, repositoryID, listenerID,
	)
}

func FormatBasicSidecarListenerIntoConfig(resName, sidecarID, repoType string, listenerPort int) string {
	return fmt.Sprintf(
		`
	resource "cyral_sidecar_listener" "%s" {
		sidecar_id = %s
		repo_types = ["%s"]
		network_address {
			port = %d
		}
	}`, resName, sidecarID, repoType, listenerPort,
	)
}

func FormatBasicSidecarIntoConfig(resName, sidecarName, deploymentMethod, logIntegrationID string) string {
	return fmt.Sprintf(
		`
	resource "cyral_sidecar" "%s" {
		name               = "%s"
		deployment_method  = "%s"
		activity_log_integration_id = "%s"
	}`, resName, sidecarName, deploymentMethod, logIntegrationID,
	)
}

func FormatBasicPolicyIntoConfig(name string, data []string) string {
	return fmt.Sprintf(
		`
	resource "cyral_policy" "%s" {
		name = "%s"
		data = %s
	}`, BasicPolicyResName, name, ListToStr(data),
	)
}

func FormatBasicIntegrationIdPOktaIntoConfig(resName, displayName, ssoURL string) string {
	return fmt.Sprintf(
		`
	resource "cyral_integration_idp_okta" "%s" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}`, resName, displayName, ssoURL,
	)
}

func FormatBasicIntegrationIdPSAMLDraftIntoConfig(resName, displayName, idpType string) string {
	return fmt.Sprintf(
		`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		idp_type = "%s"
	}`, resName, displayName, idpType,
	)
}

func FormatBasicDataSourcePermissionIntoConfig(resourceName string) string {
	return fmt.Sprintf(
		`
	data "cyral_permission" "%s" {}
	`, resourceName,
	)
}

func NotZeroRegex() *regexp.Regexp {
	return regexp.MustCompile("[^0]|([0-9]{2,})")
}

// DSourceCheckTypeFilter is used by data source tests that accept type
// filters. When the data source test is run, there might be unexpected
// resources present in the control plane. To avoid test checks that fail
// non-deterministically, this function simply checks that all objects match the
// given type filter.
//
// Example usage:
//
// DSourceCheckTypeFilter(
//
//	"data.cyral_datalabel.test_datalabel",
//	"datalabel_list.%d.type",
//	"CUSTOM",
//
// ),
func DSourceCheckTypeFilter(
	dsourceFullName, typeTemplate, typeFilter string,
) func(s *terraform.State) error {
	listKey := strings.Split(typeTemplate, ".")[0]
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dsourceFullName]
		if !ok {
			return fmt.Errorf("Not found: %s", dsourceFullName)
		}
		numObjects, err := strconv.Atoi(ds.Primary.Attributes[listKey+".#"])
		if err != nil {
			return err
		}
		for i := 0; i < numObjects; i++ {
			typeLocation := fmt.Sprintf(typeTemplate, i)
			actualType := ds.Primary.Attributes[typeLocation]
			if actualType != typeFilter {
				return fmt.Errorf(
					"Expected all objects in %s "+
						"to have type equal to type filter %q, "+
						"but got: %s", listKey, typeFilter, actualType,
				)
			}
		}
		return nil
	}
}

func DatalabelConfigResourceFullName(resName string) string {
	return fmt.Sprintf("cyral_datalabel.%s", resName)
}

func FormatDataLabelIntoConfig(resName, dataLabelName, dataLabelDescription,
	ruleType, ruleCode, ruleStatus string, dataLabelTags []string) string {
	var classificationRuleConfig string
	if ruleType != "" {
		classificationRuleConfig = fmt.Sprintf(`
 		classification_rule {
 			rule_type = "%s"
 			rule_code = "%s"
 			rule_status = "%s"
 		}`,
			ruleType,
			ruleCode,
			ruleStatus,
		)
	}
	return fmt.Sprintf(`
	resource "cyral_datalabel" "%s" {
		name  = "%s"
		description = "%s"
		tags = %s
		%s
	}`,
		resName,
		dataLabelName,
		dataLabelDescription,
		ListToStr(dataLabelTags),
		classificationRuleConfig,
	)
}

func FormatELKIntegrationDataIntoConfig(name, kibanaURL, esURL string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_elk" "elk_integration" {
		name = "%s"
		kibana_url = "%s"
		es_url = "%s"
	}`, name, kibanaURL, esURL)
}

func FormatDatadogIntegrationDataIntoConfig(name, apiKey string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_datadog" "datadog_integration" {
		name = "%s"
		api_key = "%s"
	}`, name, apiKey)
}

// FormatPolicyIntoConfig formats a policy map into a Terraform configuration string.
//
// resName is the resource name to be used in the Terraform configuration.
// policy is a map containing the policy data.
func FormatPolicyIntoConfig(resName string, policy map[string]interface{}) string {
	var config strings.Builder

	config.WriteString(fmt.Sprintf(`
resource "cyral_policy_v2" "%s" {
`, resName))

	if name, ok := policy["name"]; ok {
		config.WriteString(fmt.Sprintf("  name        = \"%s\"\n", name))
	}
	if description, ok := policy["description"]; ok {
		config.WriteString(fmt.Sprintf("  description = \"%s\"\n", description))
	}
	if enabled, ok := policy["enabled"]; ok {
		config.WriteString(fmt.Sprintf("  enabled     = %v\n", enabled))
	}
	if tags, ok := policy["tags"]; ok {
		config.WriteString(fmt.Sprintf("  tags        = %s\n", formatTags(tags)))
	}
	if scope, ok := policy["scope"]; ok {
		config.WriteString(fmt.Sprintf("  scope       %s\n", formatScope(scope)))
	}
	if validFrom, ok := policy["valid_from"]; ok {
		config.WriteString(fmt.Sprintf("  valid_from  = \"%s\"\n", validFrom))
	}
	if validUntil, ok := policy["valid_until"]; ok {
		config.WriteString(fmt.Sprintf("  valid_until = \"%s\"\n", validUntil))
	}
	if document, ok := policy["document"]; ok {
		escapedDocument := strings.ReplaceAll(document.(string), "${", "$${")
		config.WriteString(fmt.Sprintf("  document    = chomp(\n<<EOT\n%s\nEOT\n)\n", escapedDocument))
	}
	if enforced, ok := policy["enforced"]; ok {
		config.WriteString(fmt.Sprintf("  enforced    = %v\n", enforced))
	}
	if policyType, ok := policy["type"]; ok {
		config.WriteString(fmt.Sprintf("  type        = \"%s\"\n", policyType))
	}

	config.WriteString("}\n")

	return config.String()
}

// formatScope formats the scope data into a string suitable for Terraform configuration.
//
// scope is an interface that should be a map with string keys and slice of strings as values.
func formatScope(scope interface{}) string {
	if scope == nil {
		return "= null"
	}

	scopeMap := scope.(map[string][]string)
	if len(scopeMap) == 0 {
		return "= null"
	}

	var formattedScope strings.Builder
	formattedScope.WriteString("{\n")
	for _, value := range scopeMap {
		formattedScope.WriteString(fmt.Sprintf(`  %s = %s`, "repo_ids", formatTags(value)))
	}
	formattedScope.WriteString("\n}")

	return formattedScope.String()
}

// formatTags formats a slice of tags into a string suitable for Terraform configuration.
//
// tags is an interface that should be a slice of strings.
func formatTags(tags interface{}) string {
	if tags == nil {
		return "null"
	}

	tagList := tags.([]string)
	if len(tagList) == 0 {
		return "[]"
	}

	var formattedTags strings.Builder
	formattedTags.WriteString("[")
	for i, tag := range tagList {
		if i > 0 {
			formattedTags.WriteString(", ")
		}
		formattedTags.WriteString(fmt.Sprintf(`"%s"`, tag))
	}
	formattedTags.WriteString("]")

	return formattedTags.String()
}
