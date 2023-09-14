package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	arn0 = "arn:aws:iam::421552459763:role/Role1"
	arn1 = "arn:aws:iam::421552459763:role/Role2"
)

func TestIntegrationAWSIAMAuthN(t *testing.T) {
	testCreate := getCreateStep("main_test")
	testUpdate1 := getUpdateStepRemoveARN("main_test")
	testUpdate2 := getUpdateStepAddARN("main_test")
	testUpdate3 := getUpdateStepChangeName("main_test")
	testUpdate4 := getUpdateStepChangeDescription("main_test")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			testCreate,
			testUpdate1,
			testUpdate2,
			testUpdate3,
			testUpdate4,
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_integration_aws_iam_authn.main_test",
			},
		},
	})
}

func getCreateStep(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam_authn" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  arns = ["%s", "%s"]
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
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam_authn.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.0", arn0),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.1", arn1),
	)
}

func getUpdateStepRemoveARN(resName string) resource.TestStep {
	// We will delete one ARN
	config := `
resource "cyral_integration_aws_iam_authn" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  arns = ["%s"]
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
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam_authn.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.#", "1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.0", arn0),
	)
}

func getUpdateStepAddARN(resName string) resource.TestStep {
	// Add the ARN back. Also, we'll add them in a different order.
	config := `
resource "cyral_integration_aws_iam_authn" "%s" {
  name = "terraform_aws_integration_1"
  description = "whatever"
  arns = ["%s", "%s"]
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
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam_authn.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_1"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.1", arn0),
	)
}

func getUpdateStepChangeName(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam_authn" "%s" {
  name = "terraform_aws_integration_2"
  description = "whatever"
  arns = ["%s", "%s"]
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
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam_authn.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.1", arn0),
	)
}

func getUpdateStepChangeDescription(resName string) resource.TestStep {
	config := `
resource "cyral_integration_aws_iam_authn" "%s" {
  name = "terraform_aws_integration_2"
  description = "whatever whatever"
  arns = ["%s", "%s"]
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
	resourceFullName := fmt.Sprintf("cyral_integration_aws_iam_authn.%s", resName)
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", "terraform_aws_integration_2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", "whatever whatever"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.#", "2"),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.0", arn1),
		resource.TestCheckResourceAttr(resourceFullName,
			"arns.1", arn0),
	)
}
