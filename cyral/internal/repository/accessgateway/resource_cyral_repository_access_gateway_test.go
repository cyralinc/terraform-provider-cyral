package accessgateway_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	repoAccessGatewayResourceName = "access-gateway"
)

type accessGatewayTestConfig struct {
	sidecarResName  string
	bindingResName  string
	listenerResName string
	listenerPort    int
}

func (ag *accessGatewayTestConfig) sidecarID() string {
	return fmt.Sprintf("cyral_sidecar.%s.id", ag.sidecarResName)
}

func (ag *accessGatewayTestConfig) bindingID() string {
	return fmt.Sprintf("cyral_repository_binding.%s.binding_id", ag.bindingResName)
}

func (ag *accessGatewayTestConfig) listenerID() string {
	return fmt.Sprintf("cyral_sidecar_listener.%s.listener_id", ag.listenerResName)
}

func accessGatewayConfig(ag accessGatewayTestConfig) string {
	config := utils.FormatBasicRepositoryIntoConfig(
		utils.BasicRepositoryResName,
		utils.AccTestName(repoAccessGatewayResourceName, "repo"),
		"mongodb",
		"mongo.local",
		ag.listenerPort,
	)
	config += utils.FormatBasicSidecarIntoConfig(
		ag.sidecarResName,
		utils.AccTestName(repoAccessGatewayResourceName, ag.sidecarResName),
		"docker", "",
	)

	config += utils.FormatBasicSidecarListenerIntoConfig(
		ag.listenerResName,
		ag.sidecarID(),
		"mongodb",
		ag.listenerPort,
	)

	config += utils.FormatBasicRepositoryBindingIntoConfig(
		ag.bindingResName,
		ag.sidecarID(),
		utils.BasicRepositoryID,
		ag.listenerID(),
	)
	return config
}

func TestAccRepositoryAccessGatewayResource(t *testing.T) {
	accessGatewayResName := "test-access-gateway"
	initialConfig := accessGatewayTestConfig{
		sidecarResName:  "test-sidecar-1",
		bindingResName:  "test-binding-1",
		listenerResName: "test-listener-1",
		listenerPort:    27017,
	}
	updateSidecarConfig := accessGatewayTestConfig{
		sidecarResName:  "test-sidecar-2",
		bindingResName:  "test-binding-1",
		listenerResName: "test-listener-1",
		listenerPort:    27017,
	}
	updateBindingConfig := accessGatewayTestConfig{
		sidecarResName:  "test-sidecar-2",
		bindingResName:  "test-binding-2",
		listenerResName: "test-listener-2",
		listenerPort:    27018,
	}
	intialTest := repositoryAccessGatewayTestStep(accessGatewayResName, initialConfig)
	updateSidecarTest := repositoryAccessGatewayTestStep(accessGatewayResName, updateSidecarConfig)
	updateBindingTest := repositoryAccessGatewayTestStep(accessGatewayResName, updateBindingConfig)
	resourceToImport := fmt.Sprintf("cyral_repository_access_gateway.%s", accessGatewayResName)
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      resourceToImport,
	}
	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: provider.ProviderFactories,
			Steps: []resource.TestStep{
				intialTest,
				updateSidecarTest,
				updateBindingTest,
				importTest,
			},
		},
	)
}

func repositoryAccessGatewayTestStep(resName string, ag accessGatewayTestConfig) resource.TestStep {
	return resource.TestStep{
		Config: accessGatewayConfig(ag) +
			repoAccessGatewayConfig(resName, ag.sidecarID(), ag.bindingID()),
		Check: repoAccessGatewayCheck(resName, ag.sidecarResName, ag.bindingResName),
	}
}

func repoAccessGatewayCheck(resName, sidecarResName, bindingResName string) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_access_gateway.%s", resName)
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, utils.SidecarIDKey,
			fmt.Sprintf("cyral_sidecar.%s", sidecarResName), "id",
		),
		resource.TestCheckResourceAttrPair(
			resFullName, utils.RepositoryIDKey,
			fmt.Sprintf("cyral_repository.%s", utils.BasicRepositoryResName), "id",
		),
		resource.TestCheckResourceAttrPair(
			resFullName, utils.BindingIDKey,
			fmt.Sprintf("cyral_repository_binding.%s", bindingResName), "binding_id",
		),
	}
	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func repoAccessGatewayConfig(resName, sidecarID, bindingID string) string {
	return fmt.Sprintf(
		`
	resource "cyral_repository_access_gateway" "%s" {
		repository_id  = %s
		sidecar_id  = %s
		binding_id = %s
	}`, resName, utils.BasicRepositoryID, sidecarID, bindingID,
	)
}
