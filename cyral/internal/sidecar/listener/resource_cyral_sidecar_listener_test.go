package listener_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	listen "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/listener"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarListenerTestSidecarResourceName = "sidecar-for-listeners"
)

func sidecarListenerSidecarConfig() string {
	return utils.FormatBasicSidecarIntoConfig(
		utils.BasicSidecarResName,
		utils.AccTestName(sidecarListenerTestSidecarResourceName, "sidecar"),
		"docker", "",
	)
}
func TestSidecarListenerResource(t *testing.T) {
	testSteps := make([]resource.TestStep, 0, 10)
	testSteps = append(testSteps, settingsConflictsTest()...)
	testSteps = append(testSteps, updateTest()...)
	testSteps = append(testSteps, settingsTest()...)
	testSteps = append(testSteps, multipleListenersAndImportTest()...)

	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: provider.ProviderFactories,
			Steps:             testSteps,
		},
	)
}

func updateTest() []resource.TestStep {
	// Start with a bare bones mySQL sidecar listener
	onlyRequiredFields := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8000,
		},
	}
	// Change port.
	changePort := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 443,
		},
	}
	// Add host.
	addHost := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 443,
			Host: "https://mysql.test.com",
		},
	}
	// Add settings.
	addSettings := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 443,
			Host: "https://s3.test.com",
		},
		MySQLSettings: &listen.MySQLSettings{
			DbVersion: "3.4.0",
		},
	}

	return []resource.TestStep{
		setupSidecarListenerTestStep(
			"update_test",
			onlyRequiredFields,
		),
		setupSidecarListenerTestStep(
			"update_test",
			changePort,
		),
		setupSidecarListenerTestStep(
			"update_test",
			addHost,
		),
		setupSidecarListenerTestStep(
			"update_test",
			addSettings,
		),
	}
}

func settingsTest() []resource.TestStep {
	// MySQL listener with no character set defined.
	mySQLNoCharSet := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8000,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &listen.MySQLSettings{
			DbVersion: "2.1.0",
		},
	}
	// MySQL listener with character set defined.
	mySQLWithCharSet := listen.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8001,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &listen.MySQLSettings{
			DbVersion:    "2.1.0",
			CharacterSet: "ujis_japanese_ci",
		},
	}
	// S3 listeners with proxy mode.
	s3 := listen.SidecarListener{
		RepoTypes: []string{"s3"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8002,
			Host: "https://mysql.test.com",
		},
		S3Settings: &listen.S3Settings{
			ProxyMode: true,
		},
	}
	// DynamoDB with proxy mode (required for all DynamoDB repo types).
	dynamodb := listen.SidecarListener{
		RepoTypes: []string{"dynamodb"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8003,
			Host: "https://mysql.test.com",
		},
		DynamoDbSettings: &listen.DynamoDbSettings{
			ProxyMode: true,
		},
	}
	// SQL Server settings test step
	sqlServerSettings := listen.SidecarListener{
		RepoTypes: []string{"sqlserver"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 8004,
			Host: "https://sqlserver.test.com",
		},
		SQLServerSettings: &listen.SQLServerSettings{
			Version: "16.0.1000",
		},
	}

	return []resource.TestStep{
		setupSidecarListenerTestStep(
			"mySQL_no_charset",
			mySQLNoCharSet,
		),
		setupSidecarListenerTestStep(
			"mySQL_with_charset",
			mySQLWithCharSet,
		),
		setupSidecarListenerTestStep(
			"s3_with_proxy",
			s3,
		),
		setupSidecarListenerTestStep(
			"dynamo_db_with_proxy",
			dynamodb,
		),
		setupSidecarListenerTestStep(
			"sqlserver_settings",
			sqlServerSettings,
		),
	}
}

func multipleListenersAndImportTest() []resource.TestStep {
	// Define first listener resource.
	listener1ResName := "listener1-multiple-test"
	listener1 := listen.SidecarListener{
		RepoTypes: []string{"oracle"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 6003,
		},
	}
	// Define second listener resource.
	listener2ResName := "listener2-multiple-test"
	listener2 := listen.SidecarListener{
		RepoTypes: []string{"postgresql"},
		NetworkAddress: &listen.NetworkAddress{
			Port: 6018,
		},
	}

	// Setup config containing both listeners.
	listener1Config := SetupSidecarListenerConfig(
		listener1ResName, listener1,
	)
	listener2Config := SetupSidecarListenerConfig(
		listener2ResName, listener2,
	)
	multipleListenersConfig := sidecarListenerSidecarConfig() +
		listener1Config + listener2Config

	// Setup check for both listeners.
	listener1Check := setupSidecarListenerCheck(listener1ResName, listener1)
	listener2Check := setupSidecarListenerCheck(listener2ResName, listener2)
	multipleListenersCheck := resource.ComposeTestCheckFunc(
		listener1Check, listener2Check,
	)

	// Create multiple listeners test step.
	multipleListenersTest := resource.TestStep{
		Config: multipleListenersConfig,
		Check:  multipleListenersCheck,
	}

	// Create import test.
	resourceToImport := fmt.Sprintf("cyral_sidecar_listener.%s", listener1ResName)
	importTest := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      resourceToImport,
	}

	return []resource.TestStep{
		multipleListenersTest,
		importTest,
	}
}

// settingsConflictsTest generates a test matrix to ensure that conflicting settings
// for different repository types produce appropriate errors.
func settingsConflictsTest() []resource.TestStep {
	// List of repo types with conflicting settings
	repoTypes := []string{
		"mysql",
		"s3",
		"dynamodb",
		"sqlserver",
	}
	var testSteps []resource.TestStep
	// Generate test steps for every pair of conflicting repo types
	for i := 0; i < len(repoTypes); i++ {
		for j := i + 1; j < len(repoTypes); j++ {
			// Create a listener with all conflicting repo types
			// Downstream test code will cut at [0], but this is fine for what we are testing here
			listener := listen.SidecarListener{
				RepoTypes: repoTypes,
				NetworkAddress: &listen.NetworkAddress{
					Port: 8000,
					Host: "https://mysql.test.com",
				},
			}
			// Apply conflicting settings to the listener
			appendSetting(&listener, repoTypes[i])
			appendSetting(&listener, repoTypes[j])
			// Create a test step with the listener
			testSteps = append(testSteps, resource.TestStep{
				Config: sidecarListenerSidecarConfig() + SetupSidecarListenerConfig("settings_conflict", listener),
				ExpectError: regexp.MustCompile(
					".*conflicts with.*",
				),
			})
		}
	}
	return testSteps
}

// appendSetting applies settings for a given repository type to the provided listen.
// The listener's repository types are updated accordingly using default values.
func appendSetting(l *listen.SidecarListener, repoType string) {
	switch repoType {
	case "mysql":
		l.MySQLSettings = &listen.MySQLSettings{
			DbVersion: "5.7",
		}
	case "s3":
		l.S3Settings = &listen.S3Settings{
			ProxyMode: true,
		}
	case "dynamodb":
		l.DynamoDbSettings = &listen.DynamoDbSettings{
			ProxyMode: true,
		}
	case "sqlserver":
		l.SQLServerSettings = &listen.SQLServerSettings{
			Version: "16.0.1000",
		}
	}
}

func setupSidecarListenerTestStep(resName string, listener listen.SidecarListener) resource.TestStep {
	return resource.TestStep{
		Config: sidecarListenerSidecarConfig() +
			SetupSidecarListenerConfig(resName, listener),
		Check: setupSidecarListenerCheck(resName, listener),
	}
}

func setupSidecarListenerCheck(resourceName string, listener listen.SidecarListener) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_sidecar_listener.%s", resourceName)
	var checkFuncs []resource.TestCheckFunc

	// Required attributes
	checkFuncs = append(
		checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrPair(
				resFullName, utils.SidecarIDKey,
				fmt.Sprintf("cyral_sidecar.%s", utils.BasicSidecarResName), "id",
			),
			resource.TestCheckResourceAttr(
				resFullName,
				fmt.Sprintf("%s.0", listen.RepoTypesKey), listener.RepoTypes[0],
			),
		}...,
	)

	// Optional attributes
	if listener.NetworkAddress != nil {
		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				// Exactly one Network Address conf.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.#", listen.NetworkAddressKey),
					"1",
				),
				// Check host.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.NetworkAddressKey, utils.HostKey),
					listener.NetworkAddress.Host,
				),
				// Check port.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.NetworkAddressKey, utils.PortKey),
					strconv.Itoa(listener.NetworkAddress.Port),
				),
			}...,
		)
	}

	if listener.MySQLSettings != nil {
		expectedCharSet := "unspecified"
		if listener.MySQLSettings.CharacterSet != "" {
			expectedCharSet = listener.MySQLSettings.CharacterSet
		}

		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				// Exactly one mySQL Settings.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.#", listen.MySQLSettingsKey),
					"1",
				),
				// Check DB version.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.MySQLSettingsKey, listen.DbVersionKey),
					listener.MySQLSettings.DbVersion,
				),
				// Check character set.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.MySQLSettingsKey, listen.CharacterSetKey),
					expectedCharSet,
				),
			}...,
		)
	}

	if listener.S3Settings != nil {
		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				// Exactly one S3 Settings.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.#", listen.S3SettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.S3SettingsKey, listen.ProxyModeKey),
					strconv.FormatBool(listener.S3Settings.ProxyMode),
				),
			}...,
		)
	}

	if listener.DynamoDbSettings != nil {
		checkFuncs = append(
			checkFuncs, []resource.TestCheckFunc{
				// Exactly one S3 Settings.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.#", listen.DynamoDbSettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", listen.DynamoDbSettingsKey, listen.ProxyModeKey),
					strconv.FormatBool(listener.DynamoDbSettings.ProxyMode),
				),
			}...,
		)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func SetupSidecarListenerConfig(resourceName string, listener listen.SidecarListener) string {
	var config string
	var networkAddressStr string
	if listener.NetworkAddress != nil {
		var host, port string
		if listener.NetworkAddress.Host != "" {
			host = fmt.Sprintf(`host = "%s"`, listener.NetworkAddress.Host)
		}
		if listener.NetworkAddress.Port != 0 {
			port = fmt.Sprintf(`port = %d`, listener.NetworkAddress.Port)
		}
		networkAddressStr = fmt.Sprintf(
			`
			network_address {
				%s
				%s
			}`, host, port,
		)
	}

	var settings string

	if listener.MySQLSettings != nil {
		dbVersion, charSet := "null", "null"
		if listener.MySQLSettings.CharacterSet != "" {
			charSet = fmt.Sprintf(`"%s"`, listener.MySQLSettings.CharacterSet)
		}
		if listener.MySQLSettings.DbVersion != "" {
			dbVersion = fmt.Sprintf(`"%s"`, listener.MySQLSettings.DbVersion)
		}
		settings += fmt.Sprintf(
			`
		mysql_settings {
			db_version = %s
			character_set = %s
		}`, dbVersion, charSet,
		)
	}

	if listener.DynamoDbSettings != nil {
		settings += fmt.Sprintf(
			`
		dynamodb_settings {
			proxy_mode = %s
		}`, strconv.FormatBool(listener.DynamoDbSettings.ProxyMode),
		)
	}

	if listener.S3Settings != nil {
		settings += fmt.Sprintf(
			`
		s3_settings {
			proxy_mode = %s
		}`, strconv.FormatBool(listener.S3Settings.ProxyMode),
		)
	}
	if listener.SQLServerSettings != nil {
		settings += fmt.Sprintf(
			`
		sqlserver_settings {
			version = "%s"
		}`, listener.SQLServerSettings.Version,
		)
	}

	config += fmt.Sprintf(
		`
	resource "cyral_sidecar_listener" "%s" {
		sidecar_id = %s
		repo_types = ["%s"]
		%s
		%s
	}`, resourceName, utils.BasicSidecarID, listener.RepoTypes[0], networkAddressStr, settings,
	)
	return config
}
