package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	basicRepositoryResName             = "test_repository"
	basicRepositoryID                  = "cyral_repository.test_repository.id"
	basicRepositoryBindingResName      = "test_repository_binding"
	basicRepositoryLocalAccountResName = "test_repository_local_account"
	basicRepositoryLocalAccountID      = "cyral_repository_local_account.test_repository_local_account.id"
	basicSidecarResName                = "test_sidecar"
	basicSidecarID                     = "cyral_sidecar.test_sidecar.id"
	basicPolicyResName                 = "test_policy"
	basicPolicyID                      = "cyral_policy.test_policy.id"
)

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
	return fmt.Sprintf(`
	resource "cyral_repository" "%s" {
		name = "%s"
		type = "%s"
		host = "%s"
		port = %d
	}`, resName, repoName, typ, host, port)
}

func formatBasicRepositoryBindingIntoConfig(resName, sidecarID, repositoryID string, listenerPort int) string {
	return fmt.Sprintf(`
	resource "cyral_repository_binding" "%s" {
		sidecar_id    = %s
		repository_id = %s
		listener_port = %d
	}`, resName, sidecarID, repositoryID, listenerPort)
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
		data = [%s]
	}`, basicPolicyResName, name, listToStr(data))
}

func formatBasicRepositoryLocalAccountIntoConfig_Cyral(
	repositoryID, localAccount, password string,
) string {
	return fmt.Sprintf(`
	resource "cyral_repository_local_account" "%s" {
		repository_id = %s
		cyral_storage {
			local_account = "%s"
			password      = "%s"
		}
	}`, basicRepositoryLocalAccountResName, repositoryID, localAccount, password)
}

func notZeroRegex() *regexp.Regexp {
	// TODO: fix this regex -aholmquist 2022-08-09
	return regexp.MustCompile("^[0-9]*[^0]$")
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
// 	"data.cyral_datalabel.test_datalabel",
// 	"datalabel_list.%d.type",
// 	"CUSTOM",
// ),
//
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
