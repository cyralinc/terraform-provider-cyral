package cyral

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarListenerResourceName    = "cyral_sidecar_listener"
	plainMySQLListenerResourceName = "plain_mysql"
	plainMySQLListener             = `resource "%s" "%s" {
  sidecar_id = %s
  tcp_listener {
    port = 3306
  }
  repo_types =["mysql"]
}
`
	multiplexedMySQLListenerResourceName = "multiplexed_mysql"
	multiplexedMySQLListener             = `resource "%s" "%s" {
    sidecar_id = %s
    tcp_listener {
		port = 3307
    }
    multiplexed = true
    mysql_settings {
		db_version = "5.7"
    }
    repo_types =["mysql"]
}
`
	s3ListenerResourceName = "s3"
	s3Listener             = `resource "%s" "%s" {
    sidecar_id = %s
    tcp_listener {
		port = 443
	}
    s3_settings {
		proxy_mode = true
    }
    repo_types =["s3"]
}
`
	mariaDbSocketListenerResourceName = "maria_db_socket"
	mariaDbSocketListener             = `resource "%s" "%s" {
    sidecar_id = %s
    unix_listener {
		file = "/var/run/mysqld/mysql.sock"
    }
    repo_types =["mariadb"]
}
`
)

func TestSidecarListenerResource(t *testing.T) {

	testConfig, testFunc := setupSidecarListenerTest()
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_sidecar_listener.plain_mysql",
			},
		},
	})

}

func setupSidecarListenerTest() (string, resource.TestCheckFunc) {
	configuration := createSidecarListenerConfig()
	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "id"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "listener_id"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "repo_types.#", "1"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "repo_types.0", "mysql"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "tcp_listener.#", "1"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "tcp_listener.0.port", "3306"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, plainMySQLListenerResourceName), "unix_listener.#", "0"),

		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("%s.%s", sidecarListenerResourceName, multiplexedMySQLListenerResourceName), "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("%s.%s", sidecarListenerResourceName, s3ListenerResourceName), "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "id"),
		resource.TestCheckResourceAttrSet(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "listener_id"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "repo_types.#", "1"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "repo_types.0", "mariadb"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "unix_listener.#", "1"),
		resource.TestCheckResourceAttr(fmt.Sprintf("%s.%s", sidecarListenerResourceName, mariaDbSocketListenerResourceName), "unix_listener.0.file", "/var/run/mysqld/mysql.sock"),
	)
	return configuration, testFunction

}

func createSidecarListenerConfig() string {
	var config string
	// generate random string of eight characters
	id := uuid.New()

	config += formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName(sidecarResourceName, id.String()),
		"docker",
	)
	config += "\n"
	config += fmt.Sprintf(plainMySQLListener, sidecarListenerResourceName, plainMySQLListenerResourceName, basicSidecarID)
	config += fmt.Sprintf(multiplexedMySQLListener, sidecarListenerResourceName, multiplexedMySQLListenerResourceName, basicSidecarID)
	config += fmt.Sprintf(s3Listener, sidecarListenerResourceName, s3ListenerResourceName, basicSidecarID)
	config += fmt.Sprintf(mariaDbSocketListener, sidecarListenerResourceName, mariaDbSocketListenerResourceName, basicSidecarID)
	return config
}
