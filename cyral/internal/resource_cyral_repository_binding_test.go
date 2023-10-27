package internal_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repoBindingSidecarName = "sidecar-for-bindings-test"
	repoBindingRepoName    = "repo-for-bindings-test"
)

var initialConfig = internal.Binding{
	Enabled: true,
	ListenerBindings: []*internal.ListenerBinding{
		{
			NodeIndex: 0,
		},
	},
}

var updatedConfig = internal.Binding{
	Enabled: false,
	ListenerBindings: []*internal.ListenerBinding{
		{
			NodeIndex: 0,
		},
	},
}

func bindingRepoSidecarListenerConfig() string {
	config := utils.FormatBasicRepositoryIntoConfig(
		utils.BasicRepositoryResName,
		utils.AccTestName(repoBindingRepoName, "repo"),
		"mongodb",
		"mongo.local",
		27017,
	)
	config += utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName(repoBindingSidecarName, "sidecar"),
		"docker", "",
	)

	config += utils.FormatBasicSidecarListenerIntoConfig(
		utils.BasicListenerResName,
		utils.BasicSidecarID,
		"mongodb",
		27017,
	)
	return config
}

func TestAccRepositoryBindingResource(t *testing.T) {
	intialTest := repositoryBindingTestStep("binding", initialConfig)
	updateTest := repositoryBindingTestStep("binding", updatedConfig)
	resourceToImport := fmt.Sprintf("cyral_repository_binding.%s", "binding")
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      resourceToImport,
	}
	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				intialTest,
				updateTest,
				importTest,
			},
		},
	)
}

func repositoryBindingTestStep(resName string, binding internal.Binding) resource.TestStep {
	config := bindingRepoSidecarListenerConfig() +
		repoBindingConfig(resName, binding)
	return resource.TestStep{
		Config: config,
		Check:  repoBindingCheck(resName, binding),
	}
}

func repoBindingCheck(resName string, binding internal.Binding) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_binding.%s", resName)
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, internal.SidecarIDKey,
			fmt.Sprintf("cyral_sidecar.%s", utils.BasicSidecarResName), "id",
		),
		resource.TestCheckResourceAttrPair(
			resFullName, internal.RepositoryIDKey,
			fmt.Sprintf("cyral_repository.%s", utils.BasicRepositoryResName), "id",
		),
		resource.TestCheckResourceAttr(
			resFullName,
			internal.BindingEnabledKey, strconv.FormatBool(binding.Enabled),
		),
	}

	for i, binding := range binding.ListenerBindings {
		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttrPair(
					resFullName, fmt.Sprintf("%s.%d.%s", internal.ListenerBindingKey, i, internal.ListenerIDKey),
					fmt.Sprintf("cyral_sidecar_listener.%s", utils.BasicListenerResName),
					internal.ListenerIDKey,
				),
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.%d.%s", internal.ListenerBindingKey, i, internal.NodeIndexKey),
					strconv.Itoa(int(binding.NodeIndex)),
				),
			}...,
		)
	}
	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoBindingConfig(resName string, binding internal.Binding) string {
	config := fmt.Sprintf(
		`
	resource "cyral_repository_binding" "%s" {
		sidecar_id  = %s
		repository_id  = %s
		enabled = %s`,
		resName, utils.BasicSidecarID, utils.BasicRepositoryID,
		strconv.FormatBool(binding.Enabled),
	)

	for _, binding := range binding.ListenerBindings {
		config += fmt.Sprintf(
			`
		listener_binding {
			listener_id = %s
			node_index = %d
		}`, utils.BasicListenerID, binding.NodeIndex,
		)
	}
	config += `
	}`
	return config
}
