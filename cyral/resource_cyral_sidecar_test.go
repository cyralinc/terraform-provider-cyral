package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarResourceName = "sidecar"
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

var cloudFormationSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "cft"),
	Labels:                   []string{"test1"},
	SidecarProperties:        NewSidecarProperties("cloudFormation", "foo", ""),
	UserEndpoint:             "some.cft.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var dockerSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "docker"),
	Labels:                   []string{"test2"},
	SidecarProperties:        NewSidecarProperties("docker", "bar", ""),
	UserEndpoint:             "some.docker.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var helmSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "helm"),
	Labels:                   []string{"test3"},
	SidecarProperties:        NewSidecarProperties("helm", "baz", ""),
	UserEndpoint:             "some.helm.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var tfSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "tf"),
	Labels:                   []string{"test4"},
	SidecarProperties:        NewSidecarProperties("terraform", "qux", ""),
	UserEndpoint:             "some.tf.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var singleContainerSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "singleContainer"),
	Labels:                   []string{"test5"},
	SidecarProperties:        NewSidecarProperties("singleContainer", "quxx", ""),
	UserEndpoint:             "some.singleContainer.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var linuxSidecarConfig = SidecarData{
	Name:                     accTestName(sidecarResourceName, "linux"),
	Labels:                   []string{"test6"},
	SidecarProperties:        NewSidecarProperties("linux", "empty", ""),
	UserEndpoint:             "some.linux.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var bypassNeverSidecarConfig = SidecarData{
	Name:              accTestName(sidecarResourceName, "bypassNeverSidecar"),
	SidecarProperties: NewSidecarProperties("terraform", "a", ""),
	ServiceConfigs: SidecarServiceConfigs{
		"dispatcher": map[string]string{
			"bypass": "never",
		},
	},
	UserEndpoint: "some.user.endpoint",
}

var bypassAlwaysSidecarConfig = SidecarData{
	Name:              accTestName(sidecarResourceName, "bypassAlwaysSidecar"),
	SidecarProperties: NewSidecarProperties("terraform", "b", ""),
	ServiceConfigs: SidecarServiceConfigs{
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
	testUpdateConfigLinux, testUpdateFuncLinux := setupSidecarTest(
		linuxSidecarConfig,
	)
	testUpdateConfigBypassNever, testUpdateFuncBypassNever := setupSidecarTest(bypassNeverSidecarConfig)
	testUpdateConfigBypassAlways, testUpdateFuncBypassAlways := setupSidecarTest(bypassAlwaysSidecarConfig)

	resource.ParallelTest(
		t, resource.TestCase{
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
					Config: testUpdateConfigLinux,
					Check:  testUpdateFuncLinux,
				},
				{
					Config: testUpdateConfigBypassNever,
					Check:  testUpdateFuncBypassNever,
				},
				{
					Config: testUpdateConfigBypassAlways,
					Check:  testUpdateFuncBypassAlways,
				},
				{
					ImportState:       true,
					ImportStateVerify: true,
					ResourceName:      "cyral_sidecar.test_sidecar",
				},
			},
		},
	)
}

func setupSidecarTest(sidecarData SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(sidecarData)

	var deploymentMethod, logIntegrationID string
	if properties := sidecarData.SidecarProperties; properties != nil {
		deploymentMethod = properties.DeploymentMethod
		logIntegrationID = properties.LogIntegrationID
	}
	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "name", sidecarData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "deployment_method", deploymentMethod),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "activity_log_integration_id", logIntegrationID),
	}

	if bypassMode := sidecarData.BypassMode(); bypassMode != "" {
		testFunctions = append(
			testFunctions,
			resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "bypass_mode", bypassMode),
		)
	}

	return configuration, resource.ComposeTestCheckFunc(testFunctions...)
}

func formatSidecarDataIntoConfig(sidecarData SidecarData) string {
	var certBundleConfig string
	if sidecarData.CertificateBundleSecrets != nil {
		certBundleConfig = fmt.Sprintf(
			`
		certificate_bundle_secrets {
			sidecar {
				secret_id = "%s"
				type = "%s"
				engine = "%s"
			}
		}`,
			sidecarData.CertificateBundleSecrets["sidecar"].SecretId,
			sidecarData.CertificateBundleSecrets["sidecar"].Type,
			sidecarData.CertificateBundleSecrets["sidecar"].Engine,
		)
	}

	var servicesConfig string
	if bypassMode := sidecarData.BypassMode(); bypassMode != "" {
		servicesConfig += fmt.Sprintf(
			`
		bypass_mode = "%s"`, bypassMode,
		)
	}

	var deploymentMethod, logIntegrationID string
	if properties := sidecarData.SidecarProperties; properties != nil {
		deploymentMethod = properties.DeploymentMethod
		logIntegrationID = properties.LogIntegrationID
	}

	config := fmt.Sprintf(
		`
	resource "cyral_sidecar" "test_sidecar" {
      		name = "%s"
	      	deployment_method = "%s"
	      	activity_log_integration_id = "%s"
		labels = %s
		user_endpoint = "%s"
		%s
		%s
      	}`, sidecarData.Name,
		deploymentMethod,
		logIntegrationID,
		listToStr(sidecarData.Labels),
		sidecarData.UserEndpoint,
		certBundleConfig,
		servicesConfig,
	)

	return config
}
