package cyral

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// )

// const (
// 	repositoryBindingResourceName = "sidecar-for-bindings"

// )

// var initialConfig = BindingResource {
// }

// var updatedConfig = BindingResource{

// }

// func repositoryBindingTestSidecarConfig() string {
// 	return formatBasicSidecarIntoConfig(
// 		basicSidecarResName,
// 		accTestName(sidecarListenerTestSidecarResourceName, "sidecar"),
// 		"docker",
// 	)
// }

// func TestAccRepositoryBindingResource(t *testing.T) {
// 	resource.ParallelTest(t, resource.TestCase{
// 		ProviderFactories: providerFactories,
// 		Steps: []resource.TestStep{
// 		},
// 	})
// }

// func repositoryBindingTestStep(resName string, binding BindingResource) resource.TestStep {
// 	return resource.TestStep{
// 		Config: sidecarListenerSidecarConfig() +
// 			setupSidecarListenerConfig(resName, listener),
// 		Check: setupSidecarListenerCheck(resName, listener),
// 	}
// }
