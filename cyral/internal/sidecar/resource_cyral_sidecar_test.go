package sidecar_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func getTestCBS() sidecar.CertificateBundleSecrets {
	cbs := make(sidecar.CertificateBundleSecrets)
	cbs["sidecar"] = &sidecar.CertificateBundleSecret{
		SecretId: "someSecret",
		Type:     "aws",
		Engine:   "someEngine",
	}
	return cbs
}

var cloudFormationSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "cft"),
	Labels:                   []string{"test1"},
	SidecarProperties:        sidecar.NewSidecarProperties("cft-ec2", "foo", ""),
	UserEndpoint:             "some.cft.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var dockerSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "docker"),
	Labels:                   []string{"test2"},
	SidecarProperties:        sidecar.NewSidecarProperties("docker", "bar", ""),
	UserEndpoint:             "some.docker.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var helmSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "helm3"),
	Labels:                   []string{"test3"},
	SidecarProperties:        sidecar.NewSidecarProperties("helm3", "baz", ""),
	UserEndpoint:             "some.helm3.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var tfSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "tf"),
	Labels:                   []string{"test4"},
	SidecarProperties:        sidecar.NewSidecarProperties("terraform", "qux", ""),
	UserEndpoint:             "some.tf.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var singleContainerSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "singleContainer"),
	Labels:                   []string{"test5"},
	SidecarProperties:        sidecar.NewSidecarProperties("singleContainer", "quxx", ""),
	UserEndpoint:             "some.singleContainer.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var linuxSidecarConfig = sidecar.SidecarData{
	Name:                     utils.AccTestName(utils.SidecarResourceName, "linux"),
	Labels:                   []string{"test6"},
	SidecarProperties:        sidecar.NewSidecarProperties("linux", "empty", ""),
	UserEndpoint:             "some.linux.user.endpoint",
	CertificateBundleSecrets: getTestCBS(),
}

var bypassNeverSidecarConfig = sidecar.SidecarData{
	Name:              utils.AccTestName(utils.SidecarResourceName, "bypassNeverSidecar"),
	SidecarProperties: sidecar.NewSidecarProperties("terraform", "a", ""),
	ServicesConfig: sidecar.SidecarServicesConfig{
		"dispatcher": map[string]string{
			"bypass": "never",
		},
	},
	UserEndpoint: "some.user.endpoint",
}

var bypassAlwaysSidecarConfig = sidecar.SidecarData{
	Name:              utils.AccTestName(utils.SidecarResourceName, "bypassAlwaysSidecar"),
	SidecarProperties: sidecar.NewSidecarProperties("terraform", "b", ""),
	ServicesConfig: sidecar.SidecarServicesConfig{
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
			ProviderFactories: provider.ProviderFactories,
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

func setupSidecarTest(sidecarData sidecar.SidecarData) (string, resource.TestCheckFunc) {
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

func formatSidecarDataIntoConfig(sidecarData sidecar.SidecarData) string {
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
		utils.ListToStr(sidecarData.Labels),
		sidecarData.UserEndpoint,
		certBundleConfig,
		servicesConfig,
	)

	return config
}
