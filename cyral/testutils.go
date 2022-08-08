package cyral

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	basicRepositoryID = "cyral_repository.test_repository.id"
)

// TODO: actually use these
var accTestPrefixIsSet = false
var accTestPrefix string

func getAccTestPrefix() {
	if !accTestPrefixIsSet {
		accTestPrefix = "tf-provider-" + strconv.Itoa(rand.Int()) + "-"
		accTestPrefixIsSet = true
	}
}

func formatAttributes(attributes []string) string {
	if len(attributes) == 0 {
		return ""
	}
	s := fmt.Sprintf(`"%s"`, attributes[0])
	if len(attributes) > 1 {
		for _, attribute := range attributes[1:] {
			s += fmt.Sprintf(`, "%s"`, attribute)
		}
	}
	return s
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

// TODO: add comments -aholmquist 2022-08-08
func formatBasicRepositoryIntoConfig(name, typ, host string, port int) string {
	return fmt.Sprintf(`
	resource "cyral_repository" "test_repository" {
		name = "%s"
		type = "%s"
		host = "%s"
		port = %d
	}`, name, typ, host, port)
}

// TODO: currently unused, remove if not used until PR is ready
func formatBasicSidecarIntoConfig(name, deploymentMethod string) string {
	return fmt.Sprintf(`
	resource "cyral_sidecar" "test_sidecar" {
		name = "%s"
		deployment_method = "%s"
	}`, name, deploymentMethod)
}

func formatBasicPolicyIntoConfig(name string, data []string) string {
	return fmt.Sprintf(`
	resource "cyral_policy" "test_policy" {
		name = "%s"
		data = [%s]
	}`, name, formatAttributes(data))
}
