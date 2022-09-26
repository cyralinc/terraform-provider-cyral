package cyral

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarListenerResourceName = "cyral_sidecar_listener"
)

func TestUnmarshal(t *testing.T) {
	// create variable of type ReadidecarListenersAPIResponse
	// unmarshal the json into the variable
	var sidecarListenerResponse ReadSidecarListenersAPIResponse
	jsonString := `{"listenerConfig":{"id":"2FIv17glqx4adFJiXrPAQDdAKeH", "tcpListener":{"host":"", "port":3306}, "repoTypes":["mysql"], "multiplexed":false}}`
	err := json.Unmarshal([]byte(jsonString), &sidecarListenerResponse)
	if err != nil {
		t.Fatal(err)
	}

}
func TestSidecarListenerResource(t *testing.T) {

	testConfig, testFunc := setupSidecarListenerTest(cloudFormationSidecarConfig)
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

func setupSidecarListenerTest(sidecarData *SidecarData) (string, resource.TestCheckFunc) {
	// create a test config string from the sidecarData
	// create a test function that will check the sidecarData
	// return both
	configuration := createSidecarListenerConfig()
	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			fmt.Sprintf("%s.plain_mysql", sidecarListenerResourceName), "sidecar_id",
			"cyral_sidecar.test_sidecar", "id"),
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
	config += fmt.Sprintf(`
resource "cyral_sidecar_listener" "plain_mysql" {
  sidecar_id = %s
  tcp_listener {
    port = 3306
  }
  repo_types =["mysql"]
}
`, basicSidecarID)
	return config
}
