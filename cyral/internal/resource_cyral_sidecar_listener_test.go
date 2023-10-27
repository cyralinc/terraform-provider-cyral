package internal_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal"
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
	onlyRequiredFields := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 8000,
		},
	}
	// Change port.
	changePort := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 443,
		},
	}
	// Add host.
	addHost := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 443,
			Host: "https://mysql.test.com",
		},
	}
	// Add settings.
	addSettings := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 443,
			Host: "https://s3.test.com",
		},
		MySQLSettings: &internal.MySQLSettings{
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
	mySQLNoCharSet := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 8000,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &internal.MySQLSettings{
			DbVersion: "2.1.0",
		},
	}
	// MySQL listener with character set defined.
	mySQLWithCharSet := internal.SidecarListener{
		RepoTypes: []string{"mysql"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 8001,
			Host: "https://mysql.test.com",
		},
		MySQLSettings: &internal.MySQLSettings{
			DbVersion:    "2.1.0",
			CharacterSet: "ujis_japanese_ci",
		},
	}
	// S3 listeners with proxy mode.
	s3 := internal.SidecarListener{
		RepoTypes: []string{"s3"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 8002,
			Host: "https://mysql.test.com",
		},
		S3Settings: &internal.S3Settings{
			ProxyMode: true,
		},
	}
	// DynamoDB with proxy mode (required for all DynamoDB repo types).
	dynamodb := internal.SidecarListener{
		RepoTypes: []string{"dynamodb"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 8003,
			Host: "https://mysql.test.com",
		},
		DynamoDbSettings: &internal.DynamoDbSettings{
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
	listener1 := internal.SidecarListener{
		RepoTypes: []string{"oracle"},
		NetworkAddress: &internal.NetworkAddress{
			Port: 6003,
		},
	}
	// Define second listener resource.
	listener2ResName := "listener2-multiple-test"
	listener2 := internal.SidecarListener{
		RepoTypes: []string{"postgresql"},
		NetworkAddress: &internal.NetworkAddress{
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

func setupSidecarListenerTestStep(resName string, listener internal.SidecarListener) resource.TestStep {
	return resource.TestStep{
		Config: sidecarListenerSidecarConfig() +
			SetupSidecarListenerConfig(resName, listener),
		Check: setupSidecarListenerCheck(resName, listener),
	}
}

func setupSidecarListenerCheck(resourceName string, listener internal.SidecarListener) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_sidecar_listener.%s", resourceName)
	var checkFuncs []resource.TestCheckFunc

	// Required attributes
	checkFuncs = append(
		checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttrPair(
				resFullName, internal.SidecarIDKey,
				fmt.Sprintf("cyral_sidecar.%s", utils.BasicSidecarResName), "id",
			),
			resource.TestCheckResourceAttr(
				resFullName,
				fmt.Sprintf("%s.0", internal.RepoTypesKey), listener.RepoTypes[0],
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
					fmt.Sprintf("%s.#", internal.NetworkAddressKey),
					"1",
				),
				// Check host.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.NetworkAddressKey, internal.HostKey),
					listener.NetworkAddress.Host,
				),
				// Check port.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.NetworkAddressKey, internal.PortKey),
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
					fmt.Sprintf("%s.#", internal.MySQLSettingsKey),
					"1",
				),
				// Check DB version.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.MySQLSettingsKey, internal.DbVersionKey),
					listener.MySQLSettings.DbVersion,
				),
				// Check character set.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.MySQLSettingsKey, internal.CharacterSetKey),
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
					fmt.Sprintf("%s.#", internal.S3SettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.S3SettingsKey, internal.ProxyModeKey),
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
					fmt.Sprintf("%s.#", internal.DynamoDbSettingsKey),
					"1",
				),
				// Check proxy mode.
				resource.TestCheckResourceAttr(
					resFullName,
					fmt.Sprintf("%s.0.%s", internal.DynamoDbSettingsKey, internal.ProxyModeKey),
					strconv.FormatBool(listener.DynamoDbSettings.ProxyMode),
				),
			}...,
		)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func SetupSidecarListenerConfig(resourceName string, listener internal.SidecarListener) string {
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
	}`, resourceName, utils.BasicSidecarID, listener.RepoTypes[0], networkAddressStr, settings,
	)
	return config
}
