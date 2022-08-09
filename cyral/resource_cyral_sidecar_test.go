package cyral

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func getTestCBS() CertificateBundleSecrets {
	cbs := make(CertificateBundleSecrets)
	cbs["sidecar"] = &CertificateBundleSecret{
		SecretId: "someSecret",
		Type:     "aws",
		Engine:   "someEngine",
	}
	return cbs
}

var cloudFormationSidecarConfig *SidecarData = &SidecarData{
	Name:                     "tf-provider-TestAccSidecarResource-cft",
	Labels:                   []string{"test1"},
	SidecarProperty:          NewSidecarProperty("cloudFormation"),
	UserEndpoint:             "some.cft.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var dockerSidecarConfig *SidecarData = &SidecarData{
	Name:                     "tf-provider-TestAccSidecarResource-docker",
	Labels:                   []string{"test2"},
	SidecarProperty:          NewSidecarProperty("docker"),
	UserEndpoint:             "some.docker.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var helmSidecarConfig *SidecarData = &SidecarData{
	Name:                     "tf-provider-TestAccSidecarResource-helm",
	Labels:                   []string{"test3"},
	SidecarProperty:          NewSidecarProperty("helm"),
	UserEndpoint:             "some.helm.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var tfSidecarConfig *SidecarData = &SidecarData{
	Name:                     "tf-provider-TestAccSidecarResource-tf",
	Labels:                   []string{"test4"},
	SidecarProperty:          NewSidecarProperty("terraform"),
	UserEndpoint:             "some.tf.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var singleContainerSidecarConfig *SidecarData = &SidecarData{
	Name:                     "tf-provider-TestAccSidecarResource-singleContainer",
	Labels:                   []string{"test5"},
	SidecarProperty:          NewSidecarProperty("singleContainer"),
	UserEndpoint:             "some.singleContainer.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var failoverSidecarConfig *SidecarData = &SidecarData{
	Name:            "tf-provider-TestAccSidecarResource-failoverSidecar",
	SidecarProperty: NewSidecarProperty("terraform"),
	ServicesConfig: SidecarServicesConfig{
		"dispatcher": map[string]string{
			"bypass": "always",
		},
	},
	UserEndpoint: "some.user.endpoint",
}

var passthroughSidecarConfig *SidecarData = &SidecarData{
	Name:            "tf-provider-TestAccSidecarResource-passthroughSidecar",
	SidecarProperty: NewSidecarProperty("terraform"),
	ServicesConfig: SidecarServicesConfig{
		"dispatcher": map[string]string{
			"bypass": "always",
		},
	},
	UserEndpoint: "some.user.endpoint",
}

func TestAccSidecarResource(t *testing.T) {
	testConfig, testFunc := setupSidecarTest(cloudFormationSidecarConfig)
	testUpdateConfigDocker, testUpdateFuncDocker := setupSidecarTest(dockerSidecarConfig)
	testUpdateConfigHelm, testUpdateFuncHelm := setupSidecarTest(helmSidecarConfig)
	testUpdateConfigTF, testUpdateFuncTF := setupSidecarTest(tfSidecarConfig)
	testUpdateConfigSingleContainer, testUpdateFuncSingleContainer := setupSidecarTest(
		singleContainerSidecarConfig,
	)
	testUpdateConfigPassthrough, testUpdateFuncPassthrough := setupSidecarTest(passthroughSidecarConfig)
	testUpdateConfigFailover, testUpdateFuncFailover := setupSidecarTest(failoverSidecarConfig)

	resource.ParallelTest(t, resource.TestCase{
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
			{
				Config: testUpdateConfigSingleContainer,
				Check:  testUpdateFuncSingleContainer,
			},
			{
				Config: testUpdateConfigPassthrough,
				Check:  testUpdateFuncPassthrough,
			},
			{
				Config: testUpdateConfigFailover,
				Check:  testUpdateFuncFailover,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_sidecar.test_sidecar",
			},
		},
	})
}

func setupSidecarTest(sidecarData *SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(sidecarData)

	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "name", sidecarData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "deployment_method", sidecarData.SidecarProperty.DeploymentMethod),
	}

	if bypassMode := sidecarData.BypassMode(); bypassMode != "" {
		testFunctions = append(testFunctions,
			resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "bypass_mode", bypassMode))
	}

	return configuration, resource.ComposeTestCheckFunc(testFunctions...)
}

func formatSidecarDataIntoConfig(sidecarData *SidecarData) string {
	var certBundleConfig string
	if sidecarData.CertificateBundleSecrets != nil {
		certBundleConfig = fmt.Sprintf(`
		certificate_bundle_secrets {
			sidecar {
				secret_id = "%s"
				type = "%s"
				engine = "%s"
			}
		}`,
			sidecarData.CertificateBundleSecrets["sidecar"].SecretId,
			sidecarData.CertificateBundleSecrets["sidecar"].Type,
			sidecarData.CertificateBundleSecrets["sidecar"].Engine)
	}

	var servicesConfig string
	if bypassMode := sidecarData.BypassMode(); bypassMode != "" {
		servicesConfig += fmt.Sprintf(`
		bypass_mode = "%s"`, bypassMode)
	}

	config := fmt.Sprintf(`
	resource "cyral_sidecar" "test_sidecar" {
      		name = "%s"
	      	deployment_method = "%s"
		labels = [%s]
		user_endpoint = "%s"
		%s
		%s
      	}`, sidecarData.Name,
		sidecarData.SidecarProperty.DeploymentMethod,
		listToStr(sidecarData.Labels),
		sidecarData.UserEndpoint,
		certBundleConfig,
		servicesConfig)

	log.Printf("[DEBUG] Config:%s", config)

	return config
}
