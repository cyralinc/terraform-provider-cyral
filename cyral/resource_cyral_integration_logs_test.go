package cyral

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationLogsResourceName              = "integration-log"
	integrationLogsFullTerraformResourceName = "cyral_integration_logs.logs_integration"
)

var initialLogsConfigElk IntegrationLogConfig = IntegrationLogConfig{
	Name:             accTestName(integrationLogsResourceName, "LogsElk"),
	ReceiveAuditLogs: true,
	IntegrationConfigScheme: IntegrationConfigScheme{
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

var initialLogsConfigCloudWatch IntegrationLogConfig = IntegrationLogConfig{
	Name:             accTestName(integrationLogsResourceName, "LogsCloudWatch"),
	ReceiveAuditLogs: true,
	IntegrationConfigScheme: IntegrationConfigScheme{
		CloudWatch: &CloudWatchConfig{
			Region: "us-east-1",
			Group:  "group1",
			Stream: "abc",
		},
	},
}

var updatedLogsConfigElk IntegrationLogConfig = IntegrationLogConfig{
	Name:             accTestName(integrationLogsResourceName, "LogsElk"),
	ReceiveAuditLogs: true,
	IntegrationConfigScheme: IntegrationConfigScheme{
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

var updatedLogsConfigCloudWatch IntegrationLogConfig = IntegrationLogConfig{
	Name:             accTestName(integrationLogsResourceName, "LogsCloudWatch"),
	ReceiveAuditLogs: true,
	IntegrationConfigScheme: IntegrationConfigScheme{
		CloudWatch: &CloudWatchConfig{
			Region: "us-east-1",
			Group:  "group1",
			Stream: "abc",
		},
	},
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

func TestAccLogsIntegrationResourceCloudWatch(t *testing.T) {
	testConfig2, testFunc2 := setupLogsTest(initialLogsConfigCloudWatch)
	testUpdateConfig2, testUpdateFunc2 := setupLogsTest(updatedLogsConfigCloudWatch)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig2,
				Check:  testFunc2,
			},
			{
				Config: testUpdateConfig2,
				Check:  testUpdateFunc2,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      integrationLogsFullTerraformResourceName,
			},
		},
	})
}

func setupLogsTest(integrationData IntegrationLogConfig) (string, resource.TestCheckFunc) {
	configuration, err := formatLogsIntegrationDataIntoConfig(integrationData)
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
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.cloud_watch.0.region", integrationData.CloudWatch.Region),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.cloud_watch.0.group", integrationData.CloudWatch.Group),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.cloud_watch.0.stream", integrationData.CloudWatch.Stream),
		}...)
	case integrationData.Datadog != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.datadog.0.api_key", integrationData.Datadog.ApiKey),
		}...)
	case integrationData.Elk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.elk.0.es_url", integrationData.Elk.EsURL),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.elk.0.kibana_url", integrationData.Elk.KibanaURL),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.elk.0.es_credentials.0.password", integrationData.Elk.EsCredentials.Password),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.elk.0.es_credentials.0.username", integrationData.Elk.EsCredentials.Username),
		}...)
	case integrationData.Splunk != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.splunk.0.host", integrationData.Splunk.Host),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.splunk.0.hec_port", integrationData.Splunk.HecPort),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.splunk.0.access_token", integrationData.Splunk.AccessToken),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.splunk.0.index", integrationData.Splunk.Index),
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.splunk.0.use_tls", boolAsString(integrationData.Splunk.UseTLS)),
		}...)
	case integrationData.SumoLogic != nil:
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(integrationLogsFullTerraformResourceName, "config_scheme.0.sumo_logic.0.address", integrationData.SumoLogic.Address),
		}...)
	}

	testFunction := resource.ComposeTestCheckFunc(checkFuncs...)

	return configuration, testFunction
}

// this function formats IntegrationLogConfig into string.
// this is also used in datasource tests
func formatLogsIntegrationDataIntoConfig(data IntegrationLogConfig) (string, error) {
	var config string
	switch {
	case data.CloudWatch != nil:
		config = fmt.Sprintf(`
		cloud_watch {
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
			host = "%s"
			hec_port = "%s"
			access_token = "%s"
			index = "%s"
			use_tls = "%s"
		}`, data.Splunk.Host, data.Splunk.HecPort, data.Splunk.AccessToken, data.Splunk.Index, boolAsString(data.Splunk.UseTLS))
	case data.SumoLogic != nil:
		config = fmt.Sprintf(`
		sumo_logic {
			address = "%s"
		}`, data.SumoLogic.Address)
	default:
		return "", fmt.Errorf("Error in parsing config scheme in test, %v", data)
	}

	return fmt.Sprintf(`
	resource "cyral_integration_logs" "logs_integration" {
		name = "%s"
		receive_audit_logs = "%s"
		config_scheme {
			%s
		}
	}`, data.Name, boolAsString(data.ReceiveAuditLogs), config), nil
}
