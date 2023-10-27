package internal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
		internal.RegoPolicyInstanceNameKey,
		internal.RegoPolicyInstanceCategoryKey,
		internal.RegoPolicyInstanceTemplateIDKey,
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
	policy            internal.RegoPolicyInstancePayload
	policyCategory    string
	isUpdateOperation bool
}

var (
	regoPolicyInstanceOnlyRequiredArguments = RegoPolicyInstanceTestParameters{
		policy: internal.RegoPolicyInstancePayload{
			RegoPolicyInstance: internal.RegoPolicyInstance{
				Name:       "some-rate-limit-policy",
				TemplateID: "rate-limit",
				Parameters: "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false}",
			},
		},
		policyCategory: "SECURITY",
	}
	regoPolicyInstanceAllArguments = RegoPolicyInstanceTestParameters{
		policy: internal.RegoPolicyInstancePayload{
			RegoPolicyInstance: internal.RegoPolicyInstance{
				Name:        "some-rate-limit-policy",
				Description: "Some description.",
				TemplateID:  "rate-limit",
				Parameters:  "{\"rateLimit\":7,\"labels\":[\"EMAIL\"],\"alertSeverity\":\"high\",\"block\":false}",
				Enabled:     true,
				Scope: &internal.RegoPolicyInstanceScope{
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
		utils.ListToStr(parameters.policy.RegoPolicyInstance.Scope.RepoIDs),
		utils.ListToStr(parameters.policy.RegoPolicyInstance.Tags),
		parameters.policy.Duration,
	)
}

func testAccRegoPolicyInstanceCheck(
	parameters RegoPolicyInstanceTestParameters,
) resource.TestCheckFunc {
	testCheckFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceResourceIDKey),
		resource.TestCheckResourceAttrSet("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstancePolicyIDKey),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceCategoryKey, parameters.policyCategory),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceNameKey, parameters.policy.RegoPolicyInstance.Name),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceDescriptionKey, parameters.policy.RegoPolicyInstance.Description),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceTemplateIDKey, parameters.policy.RegoPolicyInstance.TemplateID),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceParametersKey, parameters.policy.RegoPolicyInstance.Parameters),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceEnabledKey, fmt.Sprintf("%t", parameters.policy.RegoPolicyInstance.Enabled)),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", internal.RegoPolicyInstanceTagsKey),
			fmt.Sprintf("%d", len(parameters.policy.RegoPolicyInstance.Tags))),
		resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", internal.RegoPolicyInstanceCreatedKey), "1"),
	}

	var durationTestCheckFunc resource.TestCheckFunc
	if parameters.policy.Duration != "" {
		durationTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceDurationKey, parameters.policy.Duration)
	} else {
		durationTestCheckFunc = resource.TestCheckNoResourceAttr("cyral_rego_policy_instance.policy_1",
			internal.RegoPolicyInstanceDurationKey)
	}
	testCheckFuncs = append(testCheckFuncs, durationTestCheckFunc)

	var lastUpdatedTestCheckFunc resource.TestCheckFunc
	if parameters.isUpdateOperation {
		lastUpdatedTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", internal.RegoPolicyInstanceLastUpdatedKey), "1")
	} else {
		lastUpdatedTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", internal.RegoPolicyInstanceLastUpdatedKey), "0")
	}
	testCheckFuncs = append(testCheckFuncs, lastUpdatedTestCheckFunc)

	var scopeTestCheckFunc resource.TestCheckFunc
	if parameters.policy.RegoPolicyInstance.Scope != nil {
		repoIDs := parameters.policy.RegoPolicyInstance.Scope.RepoIDs
		scopeTestCheckFunc = resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
				fmt.Sprintf("%s.0.%s.#", internal.RegoPolicyInstanceScopeKey, internal.RegoPolicyInstanceRepoIDsKey),
				fmt.Sprintf("%d", len(repoIDs)),
			),
		)
	} else {
		scopeTestCheckFunc = resource.TestCheckResourceAttr("cyral_rego_policy_instance.policy_1",
			fmt.Sprintf("%s.#", internal.RegoPolicyInstanceScopeKey), "0")
	}
	testCheckFuncs = append(testCheckFuncs, scopeTestCheckFunc)

	return resource.ComposeTestCheckFunc(
		testCheckFuncs...,
	)
}
