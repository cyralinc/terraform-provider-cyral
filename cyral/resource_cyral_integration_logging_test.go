package cyral

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationLogsResourceName              = "integration-log"
	integrationLogsFullTerraformResourceName = "cyral_integration_logging.logs_integration_test"
)

var initialLogsConfigCloudWatch LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "LogsCloudWatchTest"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		CloudWatch: &CloudWatchConfig{
			Region:           "us-east-2",
			Group:            "group2",
			Stream:           "abcd",
			LogRetentionDays: 1,
		},
	},
}

var initialLogsConfigDataDog LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Datadog"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Datadog: &DataDogConfig{
			ApiKey: "TESTING_API",
		},
	},
}
var initialLogsConfigElk LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "LogsElk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Elk: &ElkConfig{
			EsURL:     "http://es.com",
			KibanaURL: "http://kibana.com",
			EsCredentials: EsCredentials{
				Username: "gabriel",
				Password: "123",
			},
		},
	},
}

var initialLogsConfigSplunk LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Splunk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Splunk: &SplunkConfig{
			Hostname:    "www.splunk.com",
			HecPort:     "9529",
			AccessToken: "ACCESS",
			Index:       "65",
			UseTLS:      true,
		},
	},
}

var initialLogsConfigSumologic LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Sumologic"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		SumoLogic: &SumoLogicConfig{
			Address: "http://www.hostname.com.br",
		},
	},
}

var initialLogsConfigFluentbit LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Fluentbit"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		FluentBit: &FluentBitConfig{
			Config: `[OUTPUT]
Name         stdout
Match        *`,
		},
	},
}

var updatedLogsConfigCloudWatch LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "LogsCloudWatchTest"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		CloudWatch: &CloudWatchConfig{
			Region:           "us-east-1",
			Group:            "group1",
			Stream:           "abcd",
			LogRetentionDays: 1,
		},
	},
}

var updatedLogsConfigDataDog LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Datadog"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Datadog: &DataDogConfig{
			ApiKey: "TESTING_API",
		},
	},
}

var updatedLogsConfigElk LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "LogsElk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Elk: &ElkConfig{
			EsURL:     "http://esupdate.com",
			KibanaURL: "http://kibanaupdate.com",
			EsCredentials: EsCredentials{
				Username: "gabriel-update",
				Password: "1234",
			},
		},
	},
}

var updatedLogsConfigSplunk LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Splunk"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		Splunk: &SplunkConfig{
			Hostname:    "www.splunk2.com",
			HecPort:     "8090",
			AccessToken: "ACCESS",
			Index:       "65",
			UseTLS:      true,
		},
	},
}

var updatedLogsConfigSumologic LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Sumologic"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		SumoLogic: &SumoLogicConfig{
			Address: "http://www.hostnameupdated.com.br",
		},
	},
}

var updatedLogsConfigFluentbit LoggingIntegration = LoggingIntegration{
	Name:             accTestName(integrationLogsResourceName, "Fluentbit"),
	ReceiveAuditLogs: true,
	LoggingIntegrationConfig: LoggingIntegrationConfig{
		FluentBit: &FluentBitConfig{
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
		ProviderFactories: providerFactories,
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
		ProviderFactories: providerFactories,
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
		ProviderFactories: providerFactories,
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
		ProviderFactories: providerFactories,
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
		ProviderFactories: providerFactories,
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
		ProviderFactories: providerFactories,
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

func setupLogsTest(integrationData LoggingIntegration) (string, resource.TestCheckFunc) {
	configuration, err := formatLogsIntegrationDataIntoConfig(integrationData, "logs_integration_test")
	if err != nil {
		log.Fatalf("%v", err)
		return "", nil
	}

	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "name", integrationData.Name),
		resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "receive_audit_logs", "true"),
	}...)

	switch {
	case integrationData.CloudWatch != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.cloud_watch.0.region", integrationData.CloudWatch.Region),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.cloud_watch.0.group", integrationData.CloudWatch.Group),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.cloud_watch.0.stream", integrationData.CloudWatch.Stream),
			resource.TestCheckResourceAttrSet(integrationLogsFullTerraformResourceName, "config.0.cloud_watch.0.log_retention_days"),
		}...)
	case integrationData.Datadog != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.datadog.0.api_key", integrationData.Datadog.ApiKey),
		}...)
	case integrationData.Elk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.elk.0.es_url", integrationData.Elk.EsURL),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.elk.0.kibana_url", integrationData.Elk.KibanaURL),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.elk.0.es_credentials.0.password", integrationData.Elk.EsCredentials.Password),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.elk.0.es_credentials.0.username", integrationData.Elk.EsCredentials.Username),
		}...)
	case integrationData.Splunk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.splunk.0.hostname", integrationData.Splunk.Hostname),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.splunk.0.hec_port", integrationData.Splunk.HecPort),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.splunk.0.access_token", integrationData.Splunk.AccessToken),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.splunk.0.index", integrationData.Splunk.Index),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.splunk.0.use_tls", fmt.Sprint(integrationData.Splunk.UseTLS)),
		}...)
	case integrationData.SumoLogic != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config.0.sumo_logic.0.address", integrationData.SumoLogic.Address),
		}...)

	case integrationData.FluentBit != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrWith(integrationLogsFullTerraformResourceName, "config.0.fluentbit.0.config", func(value string) error {

				// string must contain the config.
				// We don't check exact value as it may contain trailing characters
				if strings.Contains(value, integrationData.FluentBit.Config) {
					return nil
				}
				return fmt.Errorf("expected %v, got %v", integrationData.FluentBit.Config, value)
			}),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

// this function formats LoggingIntegration into string.
// this is also used in datasource tests
func formatLogsIntegrationDataIntoConfig(data LoggingIntegration, resName string) (string, error) {
	var config string
	switch {
	case data.CloudWatch != nil:
		config = fmt.Sprintf(`
		cloud_watch {
			group = "%s"
			region = "%s"
			stream = "%s"
			log_retention_days = %d
		}`, data.CloudWatch.Group, data.CloudWatch.Region, data.CloudWatch.Stream, data.CloudWatch.LogRetentionDays)
	case data.Datadog != nil:
		config = fmt.Sprintf(`
		datadog {
			api_key = "%s"
		}`, data.Datadog.ApiKey)
	case data.Elk != nil:
		config = fmt.Sprintf(`
		elk {
			es_url = "%s"
			kibana_url = "%s"
			es_credentials {
				username = "%s"
				password = "%s"
			}
		}`, data.Elk.EsURL, data.Elk.KibanaURL, data.Elk.EsCredentials.Username, data.Elk.EsCredentials.Password)
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
		fluentbit {
			config = <<-EOF
%s
			EOF
		}`, data.FluentBit.Config)
	default:
		return "", fmt.Errorf("Error in parsing config scheme in test, %v", data)
	}

	return fmt.Sprintf(`
	resource "cyral_integration_logging" "%s" {
		name = "%s"
		receive_audit_logs = %t
		config {
			%s
		}
	}`, resName, data.Name, data.ReceiveAuditLogs, config), nil
}
