package logging_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/logging"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	integrationLogsResourceName              = "integration-log"
	integrationLogsFullTerraformResourceName = "cyral_integration_logging.logs_integration_test"
)

var ProviderFactories = map[string]func() (*schema.Provider, error){
	"cyral": func() (*schema.Provider, error) {
		return provider.Provider(), nil
	},
}

var initialLogsConfigCloudWatch logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsCloudWatchTest"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		CloudWatch: &logging.CloudWatchConfig{
			Region: "us-east-2",
			Group:  "group2",
			Stream: "abcd",
		},
	},
}

var initialLogsConfigDataDog logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Datadog"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Datadog: &logging.DataDogConfig{
			ApiKey: "TESTING_API",
		},
	},
}
var initialLogsConfigElk logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsElkComplete"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Elk: &logging.ElkConfig{
			EsURL:     "http://es.com",
			KibanaURL: "http://kibana.com",
			EsCredentials: &logging.EsCredentials{
				Username: "gabriel",
				Password: "123",
			},
		},
	},
}

var initialLogsConfigElkEmptyEsCredentials logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsElkEmptyEsCredentials"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Elk: &logging.ElkConfig{
			EsURL: "http://es.com",
		},
	},
}

var initialLogsConfigSplunk logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Splunk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Splunk: &logging.SplunkConfig{
			Hostname:    "www.splunk.com",
			HecPort:     "9529",
			AccessToken: "ACCESS",
			Index:       "65",
			UseTLS:      true,
		},
	},
}

var initialLogsConfigSumologic logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Sumologic"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		SumoLogic: &logging.SumoLogicConfig{
			Address: "https://www.hostname.com.br/path",
		},
	},
}

var initialLogsConfigFluentbit logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Fluentbit"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		FluentBit: &logging.FluentBitConfig{
			Config: `[OUTPUT]
Name         stdout
Match        *`,
		},
	},
}

var updatedLogsConfigCloudWatch logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsCloudWatchTest"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		CloudWatch: &logging.CloudWatchConfig{
			Region: "us-east-1",
			Group:  "group1",
			Stream: "abcd",
		},
	},
}

var updatedLogsConfigDataDog logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Datadog"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Datadog: &logging.DataDogConfig{
			ApiKey: "TESTING_API",
		},
	},
}

var updatedLogsConfigElk logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsElkComplete"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Elk: &logging.ElkConfig{
			EsURL:     "http://esupdate.com",
			KibanaURL: "http://kibanaupdate.com",
			EsCredentials: &logging.EsCredentials{
				Username: "gabriel-update",
				Password: "1234",
			},
		},
	},
}

var updatedLogsConfigElkEmptyEsCredentials logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "LogsElkEmptyEsCredentials"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Elk: &logging.ElkConfig{
			EsURL: "http://esupdate1.com",
		},
	},
}

var updatedLogsConfigSplunk logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Splunk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		Splunk: &logging.SplunkConfig{
			Hostname:    "www.splunk2.com",
			HecPort:     "8090",
			AccessToken: "ACCESS",
			Index:       "65",
			UseTLS:      true,
		},
	},
}

var updatedLogsConfigSumologic logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Sumologic"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		SumoLogic: &logging.SumoLogicConfig{
			Address: "https://www.hostnameupdated.com.br/path",
		},
	},
}

var updatedLogsConfigFluentbit logging.LoggingIntegration = logging.LoggingIntegration{
	Name:             utils.AccTestName(integrationLogsResourceName, "Fluentbit"),
	ReceiveAuditLogs: false,
	LoggingIntegrationConfig: logging.LoggingIntegrationConfig{
		FluentBit: &logging.FluentBitConfig{
			Config: `[OUTPUT]
Name         stdout
Match        *`,
		},
	},
}

func TestAccLogsIntegrationResourceCloudWatch(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigCloudWatch)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigCloudWatch)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceDataDog(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigDataDog)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigDataDog)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceElk(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigElk)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigElk)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceElkEmptyEsCredentials(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigElkEmptyEsCredentials)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigElkEmptyEsCredentials)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceSplunk(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigSplunk)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigSplunk)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceSumologic(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigSumologic)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigSumologic)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func TestAccLogsIntegrationResourceFluentbit(t *testing.T) {
	testConfig, testFunc := setupLogsTest(initialLogsConfigFluentbit)
	testUpdateConfig, testUpdateFunc := setupLogsTest(updatedLogsConfigFluentbit)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func setupLogsTest(integrationData logging.LoggingIntegration) (string, resource.TestCheckFunc) {
	configuration, err := formatLogsIntegrationDataIntoConfig(integrationData, "logs_integration_test")
	if err != nil {
		log.Fatalf("%v", err)
		return "", nil
	}

	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "name", integrationData.Name),
	}...)

	if integrationData.FluentBit == nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "receive_audit_logs", "true"),
		}...)
	}

	switch {
	case integrationData.CloudWatch != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "cloudwatch.0.region", integrationData.CloudWatch.Region),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "cloudwatch.0.group", integrationData.CloudWatch.Group),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "cloudwatch.0.stream", integrationData.CloudWatch.Stream),
		}...)
	case integrationData.Datadog != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "datadog.0.api_key", integrationData.Datadog.ApiKey),
		}...)
	case integrationData.Elk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "elk.0.es_url", integrationData.Elk.EsURL),
		}...)
		if integrationData.Elk.EsCredentials != nil {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "elk.0.es_credentials.0.password", integrationData.Elk.EsCredentials.Password),
				resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "elk.0.es_credentials.0.username", integrationData.Elk.EsCredentials.Username),
			}...)
		}
	case integrationData.Splunk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "splunk.0.hostname", integrationData.Splunk.Hostname),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "splunk.0.hec_port", integrationData.Splunk.HecPort),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "splunk.0.access_token", integrationData.Splunk.AccessToken),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "splunk.0.index", integrationData.Splunk.Index),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "splunk.0.use_tls", fmt.Sprint(integrationData.Splunk.UseTLS)),
		}...)
	case integrationData.SumoLogic != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "sumo_logic.0.address", integrationData.SumoLogic.Address),
		}...)

	case integrationData.FluentBit != nil:
		checkFuncs = append(
			checkFuncs,
			[]resource.TestCheckFunc{
				resource.TestCheckResourceAttrWith(
					integrationLogsFullTerraformResourceName,
					"fluent_bit.0.config",
					func(value string) error {
						// string must contain the config.
						// We don't check exact value as it may contain trailing characters
						if strings.Contains(value, integrationData.FluentBit.Config) {
							return nil
						}
						return fmt.Errorf("expected %v, got %v", integrationData.FluentBit.Config, value)
					},
				),
				resource.TestCheckResourceAttr(
					integrationLogsFullTerraformResourceName,
					"fluent_bit.0.skip_validate",
					fmt.Sprint(integrationData.FluentBit.SkipValidate),
				),
			}...,
		)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

// this function formats LoggingIntegration into string.
// this is also used in datasource tests
func formatLogsIntegrationDataIntoConfig(data logging.LoggingIntegration, resName string) (string, error) {
	var config string
	switch {
	case data.CloudWatch != nil:
		config = fmt.Sprintf(`
		cloudwatch {
			group = "%s"
			region = "%s"
			stream = "%s"
		}`, data.CloudWatch.Group, data.CloudWatch.Region, data.CloudWatch.Stream)
	case data.Datadog != nil:
		config = fmt.Sprintf(`
		datadog {
			api_key = "%s"
		}`, data.Datadog.ApiKey)
	case data.Elk != nil:
		if data.Elk.EsCredentials != nil {
			config = fmt.Sprintf(`
			elk {
				es_url = "%s"
				kibana_url = "%s"
				es_credentials {
					username = "%s"
					password = "%s"
				}
			}`, data.Elk.EsURL, data.Elk.KibanaURL, data.Elk.EsCredentials.Username, data.Elk.EsCredentials.Password)
		} else {
			config = fmt.Sprintf(`
			elk {
				es_url = "%s"
				kibana_url = "%s"
			}`, data.Elk.EsURL, data.Elk.KibanaURL)
		}
	case data.Splunk != nil:
		config = fmt.Sprintf(`
		splunk {
			hostname = "%s"
			hec_port = "%s"
			access_token = "%s"
			index = "%s"
			use_tls = %t
		}`, data.Splunk.Hostname, data.Splunk.HecPort, data.Splunk.AccessToken, data.Splunk.Index, data.Splunk.UseTLS)
	case data.SumoLogic != nil:
		config = fmt.Sprintf(`
		sumo_logic {
			address = "%s"
		}`, data.SumoLogic.Address)
	case data.FluentBit != nil:
		// fluentbit use INI format, so we need a proper way to handle this
		config = fmt.Sprintf(`
		fluent_bit {
			skip_validate = %t
			config = <<-EOF
%s
			EOF
		}`, data.FluentBit.SkipValidate, data.FluentBit.Config,
		)
	default:
		return "", fmt.Errorf("Error in parsing config in test, %v", data)
	}

	if data.FluentBit == nil {
		return fmt.Sprintf(`
		resource "cyral_integration_logging" "%s" {
			name = "%s"
			receive_audit_logs = %t
			%s
		}`, resName, data.Name, data.ReceiveAuditLogs, config), nil
	} else {
		return fmt.Sprintf(`
		resource "cyral_integration_logging" "%s" {
			name = "%s"
			%s
		}`, resName, data.Name, config), nil
	}

}
