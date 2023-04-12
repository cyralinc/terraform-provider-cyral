package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarListenerTestSidecarResourceName = "sidecar-for-listeners"
)

func sidecarListenerSidecarConfig() string {
	return formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName(sidecarListenerTestSidecarResourceName, "sidecar"),
		"docker", "",
	)
}

func TestSidecarListenerResource(t *testing.T) {
	testSteps := make([]resource.TestStep, 0, 10)
	testSteps = append(testSteps, updateTest()...)
	testSteps = append(testSteps, settingsTest()...)
	testSteps = append(testSteps, multipleListenersAndImportTest()...)

	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps:             testSteps,
		},
	)
}

func updateTest() []resource.TestStep {
	// Start with a bare bones mySQL sidecar listener
	onlyRequiredFields := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 8000,
		},
	}
	// Change port.
	changePort := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 443,
		},
	}
	// Add host.
	addHost := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 443,
			Host: "https://mysql.test.com",
		},
	}
	// Add settings.
	addSettings := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 443,
			Host: "https://s3.test.com",
		},
		MySQLSettings: &MySQLSettings{
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
	mySQLNoCharSet := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 8000,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &MySQLSettings{
			DbVersion: "2.1.0",
		},
	}
	// MySQL listener with character set defined.
	mySQLWithCharSet := SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &NetworkAddress{
			Port: 8001,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &MySQLSettings{
			DbVersion:    "2.1.0",
			CharacterSet: "ujis_japanese_ci",
		},
	}
	// S3 listeners with proxy mode.
	s3 := SidecarListener{
		RepoTypes: []string{"s3"},
		NetworkAddress: &NetworkAddress{
			Port: 8002,
			Host: "https://mysql.test.com",
		},
		S3Settings: &S3Settings{
			ProxyMode: true,
		},
	}
	// DynamoDB with proxy mode (required for all DynamoDB repo types).
	dynamodb := SidecarListener{
		RepoTypes: []string{"dynamodb"},
		NetworkAddress: &NetworkAddress{
			Port: 8003,
			Host: "https://mysql.test.com",
		},
		DynamoDbSettings: &DynamoDbSettings{
			ProxyMode: true,
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
	}
}

func multipleListenersAndImportTest() []resource.TestStep {
	// Define first listener resource.
	listener1ResName := "listener1-multiple-test"
	listener1 := SidecarListener{
		RepoTypes: []string{"oracle"},
		NetworkAddress: &NetworkAddress{
			Port: 6003,
		},
	}
	// Define second listener resource.
	listener2ResName := "listener2-multiple-test"
	listener2 := SidecarListener{
		RepoTypes: []string{"postgresql"},
		NetworkAddress: &NetworkAddress{
			Port: 6018,
		},
	}

	// Setup config containing both listeners.
	listener1Config := setupSidecarListenerConfig(
		listener1ResName, listener1,
	)
	listener2Config := setupSidecarListenerConfig(
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

func setupSidecarListenerTestStep(resName string, listener SidecarListener) resource.TestStep {
	return resource.TestStep{
		Config: sidecarListenerSidecarConfig() +
			setupSidecarListenerConfig(resName, listener),
		Check: setupSidecarListenerCheck(resName, listener),
	}
}

func setupSidecarListenerCheck(resourceName string, listener SidecarListener) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_sidecar_listener.%s", resourceName)
	var checkFuncs []resource.TestCheckFunc

	// Required attributes
	checkFuncs = append(
		checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrPair(
				resFullName, SidecarIDKey,
				fmt.Sprintf("cyral_sidecar.%s", basicSidecarResName), "id",
			),
			resource.TestCheckResourceAttr(
				resFullName,
				fmt.Sprintf("%s.0", RepoTypesKey), listener.RepoTypes[0],
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
					fmt.Sprintf("%s.#", NetworkAddressKey),
					"1",
				),
				// Check host.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", NetworkAddressKey, HostKey),
					listener.NetworkAddress.Host,
				),
				// Check port.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", NetworkAddressKey, PortKey),
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
					fmt.Sprintf("%s.#", MySQLSettingsKey),
					"1",
				),
				// Check DB version.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", MySQLSettingsKey, DbVersionKey),
					listener.MySQLSettings.DbVersion,
				),
				// Check character set.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", MySQLSettingsKey, CharacterSetKey),
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
					fmt.Sprintf("%s.#", S3SettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", S3SettingsKey, ProxyModeKey),
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
					fmt.Sprintf("%s.#", DynamoDbSettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", DynamoDbSettingsKey, ProxyModeKey),
					strconv.FormatBool(listener.DynamoDbSettings.ProxyMode),
				),
			}...,
		)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func setupSidecarListenerConfig(resourceName string, listener SidecarListener) string {
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
	switch {
	case listener.MySQLSettings != nil:
		dbVersion, charSet := "null", "null"
		if listener.MySQLSettings.CharacterSet != "" {
			charSet = fmt.Sprintf(`"%s"`, listener.MySQLSettings.CharacterSet)
		}
		if listener.MySQLSettings.DbVersion != "" {
			dbVersion = fmt.Sprintf(`"%s"`, listener.MySQLSettings.DbVersion)
		}
		settings = fmt.Sprintf(
			`
		mysql_settings {
			db_version = %s
			character_set = %s
		}`, dbVersion, charSet,
		)
	case listener.DynamoDbSettings != nil:
		settings = fmt.Sprintf(
			`
		dynamodb_settings {
			proxy_mode = %s
		}`, strconv.FormatBool(listener.DynamoDbSettings.ProxyMode),
		)
	case listener.S3Settings != nil:
		settings = fmt.Sprintf(
			`
		s3_settings {
			proxy_mode = %s
		}`, strconv.FormatBool(listener.S3Settings.ProxyMode),
		)
	}

	config += fmt.Sprintf(
		`
	resource "cyral_sidecar_listener" "%s" {
		sidecar_id = %s
		repo_types = ["%s"]
		%s
		%s
	}`, resourceName, basicSidecarID, listener.RepoTypes[0], networkAddressStr, settings,
	)
	return config
}
