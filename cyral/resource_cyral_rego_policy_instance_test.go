package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRegoPolicyInstanceResource(t *testing.T) {
	testSteps := []resource.TestStep{}
	testSteps = append(testSteps, getRequiredArgumentTestSteps()...)
	testSteps = append(
		testSteps,
		[]resource.TestStep{
			{
				Config: testAccRegoPolicyInstanceConfig_OnlyRequiredArguments(regoPolicyInstanceOnlyRequiredArguments),
				Check:  testAccRegoPolicyInstanceCheck(regoPolicyInstanceOnlyRequiredArguments),
			},
			{
				Config: testAccRegoPolicyInstanceConfig_AllArguments(regoPolicyInstanceAllArguments),
				Check:  testAccRegoPolicyInstanceCheck(regoPolicyInstanceAllArguments),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration"},
				ResourceName:            "cyral_rego_policy_instance.policy_1",
			},
		}...,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps:             testSteps,
	})
}

func getRequiredArgumentTestSteps() []resource.TestStep {
	requiredArgumentsTestSteps := []resource.TestStep{}
	requiredArguments := []string{
		regoPolicyInstanceNameKey,
		regoPolicyInstanceCategoryKey,
		regoPolicyInstanceTemplateIDKey,
	}
	for _, argument := range requiredArguments {
		requiredArgumentsTestSteps = append(requiredArgumentsTestSteps, resource.TestStep{
			Config: testAccRegoPolicyInstanceConfig_EmptyConfig(),
			ExpectError: regexp.MustCompile(
				fmt.Sprintf(`The argument "%s" is required, but no definition was found.`, argument),
			),
		})
	}
	return requiredArgumentsTestSteps
}

type RegoPolicyInstanceTestParameters struct {
	policy            RegoPolicyInstancePayload
	policyCategory    string
	isUpdateOperation bool
}

var (
	regoPolicyInstanceOnlyRequiredArguments = RegoPolicyInstanceTestParameters{
		policy: RegoPolicyInstancePayload{
			RegoPolicyInstance: RegoPolicyInstance{
				Name:       "some-rate-limit-policy",
				TemplateID: "rate-limit",
				Parameters: "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false}",
			},
		},
		policyCategory: "SECURITY",
	}
	regoPolicyInstanceAllArguments = RegoPolicyInstanceTestParameters{
		policy: RegoPolicyInstancePayload{
			RegoPolicyInstance: RegoPolicyInstance{
				Name:        "some-rate-limit-policy",
				Description: "Some description.",
				TemplateID:  "rate-limit",
				Parameters:  "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false}",
				Enabled:     true,
				Scope: &RegoPolicyInstanceScope{
					RepoIDs: []string{"2U4prk5o6yi1rTvvXyImz8lgbgG"},
				},
				Tags: []string{"tag1", "tag2"},
			},
			Duration: "70s",
		},
		policyCategory:    "SECURITY",
		isUpdateOperation: true,
	}
)

func testAccRegoPolicyInstanceConfig_EmptyConfig() string {
	return `
	resource "cyral_rego_policy_instance" "policy_1" {
	}
	`
}

func testAccRegoPolicyInstanceConfig_OnlyRequiredArguments(
	parameters RegoPolicyInstanceTestParameters,
) string {
	return fmt.Sprintf(`
	resource "cyral_rego_policy_instance" "policy_1" {
		name = %q
		category = %q
		template_id = %q
		parameters = %q
	}
	`,
		parameters.policy.RegoPolicyInstance.Name,
		parameters.policyCategory,
		parameters.policy.RegoPolicyInstance.TemplateID,
		parameters.policy.RegoPolicyInstance.Parameters,
	)
}

func testAccRegoPolicyInstanceConfig_AllArguments(
	parameters RegoPolicyInstanceTestParameters,
) string {
	return fmt.Sprintf(`
	resource "cyral_rego_policy_instance" "policy_1" {
		name = %q
		category = %q
		description = %q
		template_id = %q
		parameters = %q
		enabled = %t
		scope {
			repo_ids = %s
		}
		tags = %s
		duration = %q

	}
	`,
		parameters.policy.RegoPolicyInstance.Name,
		parameters.policyCategory,
		parameters.policy.RegoPolicyInstance.Description,
		parameters.policy.RegoPolicyInstance.TemplateID,
		parameters.policy.RegoPolicyInstance.Parameters,
		parameters.policy.RegoPolicyInstance.Enabled,
		listToStr(parameters.policy.RegoPolicyInstance.Scope.RepoIDs),
		listToStr(parameters.policy.RegoPolicyInstance.Tags),
		parameters.policy.Duration,
	)
}

func testAccRegoPolicyInstanceCheck(
	parameters RegoPolicyInstanceTestParameters,
) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceResourceIDKey),
		resource.TestCheckResourceAttrSet("cyral_rego_policy_instance.policy_1",
			regoPolicyInstancePolicyIDKey),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceCategoryKey, parameters.policyCategory),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceNameKey, parameters.policy.RegoPolicyInstance.Name),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceDescriptionKey, parameters.policy.RegoPolicyInstance.Description),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceTemplateIDKey, parameters.policy.RegoPolicyInstance.TemplateID),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceParametersKey, parameters.policy.RegoPolicyInstance.Parameters),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceEnabledKey, fmt.Sprintf("%t", parameters.policy.RegoPolicyInstance.Enabled)),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", regoPolicyInstanceTagsKey),
			fmt.Sprintf("%d", len(parameters.policy.RegoPolicyInstance.Tags))),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", regoPolicyInstanceCreatedKey), "1"),
	}

	var durationTestCheckFunc resource.TestCheckFunc
	if parameters.policy.Duration != "" {
		durationTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceDurationKey, parameters.policy.Duration)
	} else {
		durationTestCheckFunc = resource.TestCheckNoResourceAttr("cyral_rego_policy_instance.policy_1",
			regoPolicyInstanceDurationKey)
	}
	testCheckFuncs = append(testCheckFuncs, durationTestCheckFunc)

	var lastUpdatedTestCheckFunc resource.TestCheckFunc
	if parameters.isUpdateOperation {
		lastUpdatedTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", regoPolicyInstanceLastUpdatedKey), "1")
	} else {
		lastUpdatedTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", regoPolicyInstanceLastUpdatedKey), "0")
	}
	testCheckFuncs = append(testCheckFuncs, lastUpdatedTestCheckFunc)

	var scopeTestCheckFunc resource.TestCheckFunc
	if parameters.policy.RegoPolicyInstance.Scope != nil {
		repoIDs := parameters.policy.RegoPolicyInstance.Scope.RepoIDs
		scopeTestCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
				fmt.Sprintf("%s.0.%s.#", regoPolicyInstanceScopeKey, regoPolicyInstanceRepoIDsKey),
				fmt.Sprintf("%d", len(repoIDs)),
			),
		)
	} else {
		scopeTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", regoPolicyInstanceScopeKey), "0")
	}
	testCheckFuncs = append(testCheckFuncs, scopeTestCheckFunc)

	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}
