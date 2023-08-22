package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repoBindingSidecarName = "sidecar-for-bindings-test"
	repoBindingRepoName    = "repo-for-bindings-test"
)

var initialConfig = Binding{
	Enabled: true,
	ListenerBindings: []*ListenerBinding{
		{
			NodeIndex: 0,
		},
	},
}

var updatedConfig = Binding{
	Enabled: false,
	ListenerBindings: []*ListenerBinding{
		{
			NodeIndex: 0,
		},
	},
}

func bindingRepoSidecarListenerConfig() string {
	config := utils.FormatBasicRepositoryIntoConfig(
		BasicRepositoryResName,
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
		BasicListenerResName,
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

func repositoryBindingTestStep(resName string, binding Binding) resource.TestStep {
	config := bindingRepoSidecarListenerConfig() +
		repoBindingConfig(resName, binding)
	return resource.TestStep{
		Config: config,
		Check:  repoBindingCheck(resName, binding),
	}
}

func repoBindingCheck(resName string, binding Binding) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_binding.%s", resName)
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, SidecarIDKey,
			fmt.Sprintf("cyral_sidecar.%s", utils.BasicSidecarResName), "id",
		),
		resource.TestCheckResourceAttrPair(
			resFullName, RepositoryIDKey,
			fmt.Sprintf("cyral_repository.%s", BasicRepositoryResName), "id",
		),
		resource.TestCheckResourceAttr(
			resFullName,
			BindingEnabledKey, strconv.FormatBool(binding.Enabled),
		),
	}

	for i, binding := range binding.ListenerBindings {
		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttrPair(
					resFullName, fmt.Sprintf("%s.%d.%s", ListenerBindingKey, i, ListenerIDKey),
					fmt.Sprintf("cyral_sidecar_listener.%s", BasicListenerResName),
					ListenerIDKey,
				),
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.%d.%s", ListenerBindingKey, i, NodeIndexKey),
					strconv.Itoa(int(binding.NodeIndex)),
				),
			}...,
		)
	}
	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoBindingConfig(resName string, binding Binding) string {
	config := fmt.Sprintf(
		`
	resource "cyral_repository_binding" "%s" {
		sidecar_id  = %s
		repository_id  = %s
		enabled = %s`,
		resName, utils.BasicSidecarID, BasicRepositoryID,
		strconv.FormatBool(binding.Enabled),
	)

	for _, binding := range binding.ListenerBindings {
		config += fmt.Sprintf(
			`
		listener_binding {
			listener_id = %s
			node_index = %d
		}`, BasicListenerID, binding.NodeIndex,
		)
	}
	config += `
	}`
	return config
}
