package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	tprovACCPrefix = "tfprov-acc-"

	basicRepositoryResName        = "test_repository"
	basicRepositoryID             = "cyral_repository.test_repository.id"
	basicRepositoryBindingResName = "test_repository_binding"
	basicSidecarResName           = "test_sidecar"
	basicSidecarID                = "cyral_sidecar.test_sidecar.id"
	basicListenerResName          = "test_listener"
	basicListenerID               = "cyral_sidecar_listener.test_listener.listener_id"
	basicPolicyResName            = "test_policy"
	basicPolicyID                 = "cyral_policy.test_policy.id"

	testSingleSignOnURL = "https://some-test-sso-url.com"
)

// accTestName attempts to make resource names unique to a specific resource
// type, and avoid name clashes with other resources that exist in the testing
// control plane.
//
// Use this for every resource for which name clashes may occur.
//
// Example usage for cyral_datalabel resource:
//
//	accTestName("datalabel", "label1")
//
// Example usage for cyral_datalabel data source:
//
//	accTestName("data-datalabel", "label1")
//
// Note that doing it like above will prevent that the tests attempt to create a
// label called LABEL1 simultaneously, which would cause a failure.
//
// The convention is to use hyphen-separated words if possible, and prefix data
// sources with "data", to distinguish them and their counterpart resources.
func accTestName(resourceType, suffix string) string {
	return fmt.Sprintf("%s%s-%s", tprovACCPrefix, resourceType, suffix)
}

func hasAccTestPrefix(name string) bool {
	return strings.HasPrefix(name, tprovACCPrefix)
}

func importStateComposedIDFunc(
	resName string,
	idAtts []string,
	sep string,
) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		res, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resName)
		}
		var idParts []string
		for _, idAtt := range idAtts {
			idParts = append(idParts, res.Primary.Attributes[idAtt])
		}
		return marshalComposedID(idParts, sep), nil
	}
}

func formatBasicRepositoryIntoConfig(resName, repoName, typ, host string, port int) string {
	if typ == MongoDB {
		return fmt.Sprintf(`
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
		}`, resName, repoName, typ, host, port)
	} else {
		return fmt.Sprintf(`
		resource "cyral_repository" "%s" {
			name = "%s"
			type = "%s"
			repo_node {
				host = "%s"
				port = %d
			}
		}`, resName, repoName, typ, host, port)
	}
}

func formatBasicRepositoryBindingIntoConfig(resName, sidecarID, repositoryID, listenerID string) string {
	return fmt.Sprintf(`
	resource "cyral_repository_binding" "%s" {
		sidecar_id    = %s
		repository_id = %s
		listener_binding {
			listener_id = %s
			node_index = 0
		}
	}`, resName, sidecarID, repositoryID, listenerID)
}

func formatBasicSidecarListenerIntoConfig(resName, sidecarID, repoType string, listenerPort int) string {
	return fmt.Sprintf(`
	resource "cyral_sidecar_listener" "%s" {
		sidecar_id = %s
		repo_types = ["%s"]
		network_address {
			port = %d
		}
	}`, resName, sidecarID, repoType, listenerPort)
}

func formatBasicSidecarIntoConfig(resName, sidecarName, deploymentMethod string) string {
	return fmt.Sprintf(`
	resource "cyral_sidecar" "%s" {
		name              = "%s"
		deployment_method = "%s"
	}`, resName, sidecarName, deploymentMethod)
}

func formatBasicPolicyIntoConfig(name string, data []string) string {
	return fmt.Sprintf(`
	resource "cyral_policy" "%s" {
		name = "%s"
		data = %s
	}`, basicPolicyResName, name, listToStr(data))
}

func formatBasicIntegrationIdPOktaIntoConfig(resName, displayName, ssoURL string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_okta" "%s" {
		samlp {
			display_name = "%s"
			config {
				single_sign_on_service_url = "%s"
			}
		}
	}`, resName, displayName, ssoURL)
}

func formatBasicIntegrationIdPSAMLDraftIntoConfig(resName, displayName, idpType string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		idp_type = "%s"
	}`, resName, displayName, idpType)
}

func notZeroRegex() *regexp.Regexp {
	return regexp.MustCompile("[^0]|([0-9]{2,})")
}

// dsourceCheckTypeFilter is used by data source tests that accept type
// filters. When the data source test is run, there might be unexpected
// resources present in the control plane. To avoid test checks that fail
// non-deterministically, this function simply checks that all objects match the
// given type filter.
//
// Example usage:
//
// dsourceCheckTypeFilter(
//
//	"data.cyral_datalabel.test_datalabel",
//	"datalabel_list.%d.type",
//	"CUSTOM",
//
// ),
func dsourceCheckTypeFilter(
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
				return fmt.Errorf("Expected all objects in %s "+
					"to have type equal to type filter %q, "+
					"but got: %s", listKey, typeFilter, actualType)
			}
		}
		return nil
	}
}
