package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func getTestCBS() *CertificateBundleSecrets {
	cbs := make(CertificateBundleSecrets)
	cbs["sidecar"] = &CertificateBundleSecret{
		SecretId: "someSecret",
		Type:     "aws",
		Engine:   "someEngine",
	}
	return &cbs
}

var cloudFormationSidecarConfig *SidecarData = &SidecarData{
	Name:   "tf-provider-TestAccSidecarResource-cft",
	Labels: []string{"test1"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
	UserEndpoint:             "some.cft.user.endpoint",
	CertificateBundleSecrets: *getTestCBS(),
}

var dockerSidecarConfig *SidecarData = &SidecarData{
	Name:   "tf-provider-TestAccSidecarResource-docker",
	Labels: []string{"test2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "docker",
	},
	UserEndpoint:             "some.docker.user.endpoint",
	CertificateBundleSecrets: *getTestCBS(),
}

var helmSidecarConfig *SidecarData = &SidecarData{
	Name:   "tf-provider-TestAccSidecarResource-helm",
	Labels: []string{"test3"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "helm",
	},
	UserEndpoint:             "some.helm.user.endpoint",
	CertificateBundleSecrets: *getTestCBS(),
}

var tfSidecarConfig *SidecarData = &SidecarData{
	Name:   "tf-provider-TestAccSidecarResource-tf",
	Labels: []string{"test4"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "terraform",
	},
	UserEndpoint:             "some.tf.user.endpoint",
	CertificateBundleSecrets: *getTestCBS(),
}

func TestAccSidecarResource(t *testing.T) {
	testConfig, testFunc := setupSidecarTest(cloudFormationSidecarConfig)
	testUpdateConfigDocker, testUpdateFuncDocker := setupSidecarTest(dockerSidecarConfig)
	testUpdateConfigHelm, testUpdateFuncHelm := setupSidecarTest(helmSidecarConfig)
	testUpdateConfigTF, testUpdateFuncTF := setupSidecarTest(tfSidecarConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfigDocker,
				Check:  testUpdateFuncDocker,
			},
			{
				Config: testUpdateConfigHelm,
				Check:  testUpdateFuncHelm,
			},
			{
				Config: testUpdateConfigTF,
				Check:  testUpdateFuncTF,
			},
		},
	})
}

func setupSidecarTest(integrationData *SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "deployment_method", integrationData.SidecarProperty.DeploymentMethod),
	)

	return configuration, testFunction
}

func formatSidecarDataIntoConfig(data *SidecarData) string {
	return fmt.Sprintf(`
      resource "cyral_sidecar" "test_sidecar" {
      	name = "%s"
      	deployment_method = "%s"
		labels = ["%s"]
		user_endpoint = "%s"
		certificate_bundle_secrets {
			sidecar {
				secret_id = "%s"
				type = "%s"
				engine = "%s"
			}
		}
      }`, data.Name,
		data.SidecarProperty.DeploymentMethod,
		data.Labels[0],
		data.UserEndpoint,
		data.CertificateBundleSecrets["sidecar"].SecretId,
		data.CertificateBundleSecrets["sidecar"].Type,
		data.CertificateBundleSecrets["sidecar"].Engine)
}
