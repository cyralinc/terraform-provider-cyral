package internal_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	arn0 = "arn:aws:iam::421552459763:role/Role1"
	arn1 = "arn:aws:iam::421552459763:role/Role2"
)

func TestIntegrationAWSIAMAuthN(t *testing.T) {
	resourceName := utils.AccTestName("aws-iam-integration", "main_test")
	testCreate := getCreateStep(resourceName)
	testUpdate1 := getUpdateStepRemoveARN(resourceName)
	testUpdate2 := getUpdateStepAddARN(resourceName)
	testUpdate3 := getUpdateStepChangeName(resourceName)
	testUpdate4 := getUpdateStepChangeDescription(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			testCreate,
			testUpdate1,
			testUpdate2,
			testUpdate3,
			testUpdate4,
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      fmt.Sprintf("cyral_integration_aws_iam.%s", resourceName),
			},
		},
	})
}

func getCreateStep(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  role_arns = ["%s", "%s"]
}
	`

	config = fmt.Sprintf(config, resName, arn0, arn1)

	return resource.TestStep{
		Config: config,
		Check: AWSIAMCreateCheck(
			resName,
		),
	}
}

func AWSIAMCreateCheck(resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.0", arn0),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.1", arn1),
	)
}

func getUpdateStepRemoveARN(resName string) resource.TestStep {
	// We will delete one ARN
	config := `
resource "cyral_integration_aws_iam" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  role_arns = ["%s"]
}
	`

	config = fmt.Sprintf(config, resName, arn0)

	return resource.TestStep{
		Config: config,
		Check: AWSIAMUpdateCheck1(
			resName,
		),
	}
}

func AWSIAMUpdateCheck1(resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.0", arn0),
	)
}

func getUpdateStepAddARN(resName string) resource.TestStep {
	// Add the ARN back. Also, we'll add them in a different order.
	config := `
resource "cyral_integration_aws_iam" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  role_arns = ["%s", "%s"]
}
	`

	config = fmt.Sprintf(config, resName, arn1, arn0)

	return resource.TestStep{
		Config: config,
		Check: AWSIAMUpdateCheck2(
			resName,
		),
	}
}

func AWSIAMUpdateCheck2(resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.1", arn0),
	)
}

func getUpdateStepChangeName(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam" "%s" {
  name = "terraform_aws_integration_2"
  description = "whatever"
  role_arns = ["%s", "%s"]
}
	`

	config = fmt.Sprintf(config, resName, arn1, arn0)

	return resource.TestStep{
		Config: config,
		Check: AWSIAMUpdateCheck3(
			resName,
		),
	}
}

func AWSIAMUpdateCheck3(resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.1", arn0),
	)
}

func getUpdateStepChangeDescription(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam" "%s" {
  name = "terraform_aws_integration_2"
  description = "whatever whatever"
  role_arns = ["%s", "%s"]
}
	`

	config = fmt.Sprintf(config, resName, arn1, arn0)

	return resource.TestStep{
		Config: config,
		Check: AWSIAMUpdateCheck4(
			resName,
		),
	}
}

func AWSIAMUpdateCheck4(resName string) resource.TestCheckFunc {
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"role_arns.1", arn0),
	)
}
