package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarBoundPortsDataSourceName = "data-sidecar-bound-ports"
)

func dsourceSidecarBoundPortsSampleSidecarConfig() string {
	return formatBasicSidecarIntoConfig(
		basicSidecarResName,
		accTestName(sidecarBoundPortsDataSourceName, "sidecar"),
		"cloudFormation", "",
	)
}

func TestAccSidecarBoundPortsDataSource(t *testing.T) {
	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccSidecarBoundPortsConfig_EmptySidecarID(),
					ExpectError: regexp.MustCompile(`The argument "sidecar_id" is required`),
				},
				{
					Config: testAccSidecarBoundPortsConfig_NoBindings(),
					Check:  testAccSidecarBoundPortsCheck_NoBindings(),
				},
				{
					Config: testAccSidecarBoundPortsConfig_MultipleBindings(),
					Check:  testAccSidecarBoundPortsCheck_MultipleBindings(),
				},
			},
		},
	)
}

func testAccSidecarBoundPortsConfig_EmptySidecarID() string {
	return `
	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
	}
	`
}

func testAccSidecarBoundPortsConfig_NoBindings() string {
	var config string
	config += dsourceSidecarBoundPortsSampleSidecarConfig()
	config += fmt.Sprintf(
		`
	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
		sidecar_id = %s
	}`, basicSidecarID,
	)
	return config
}

func testAccSidecarBoundPortsCheck_NoBindings() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.#", "0",
		),
	)
}

func testAccSidecarBoundPortsConfig_MultipleBindings() string {
	var config string
	config += dsourceSidecarBoundPortsSampleSidecarConfig()

	// Repo 1
	config += formatBasicRepositoryIntoConfig(
		"repo_1",
		accTestName(sidecarBoundPortsDataSourceName, "repo1"),
		"mysql",
		"mysql.com",
		3306,
	)
	config += formatBasicSidecarListenerIntoConfig(
		"listener_1",
		basicSidecarID,
		"mysql",
		3306,
	)
	config += formatBasicRepositoryBindingIntoConfig(
		"binding_1",
		basicSidecarID,
		"cyral_repository.repo_1.id",
		"cyral_sidecar_listener.listener_1.listener_id",
	)
	// Repo 2
	config += formatBasicRepositoryIntoConfig(
		"repo_2",
		accTestName(sidecarBoundPortsDataSourceName, "repo2"),
		"mongodb",
		"mongo.com",
		27017,
	)
	config += formatBasicSidecarListenerIntoConfig(
		"listener_2",
		basicSidecarID,
		"mongodb",
		27017,
	)
	config += formatBasicRepositoryBindingIntoConfig(
		"binding_2",
		basicSidecarID,
		"cyral_repository.repo_2.id",
		"cyral_sidecar_listener.listener_2.listener_id",
	)
	// Repo 3
	config += formatBasicRepositoryIntoConfig(
		"repo_3",
		accTestName(sidecarBoundPortsDataSourceName, "repo3"),
		"oracle",
		"oracle.com",
		1234,
	)
	config += formatBasicSidecarListenerIntoConfig(
		"listener_3",
		basicSidecarID,
		"oracle",
		1234,
	)
	config += formatBasicRepositoryBindingIntoConfig(
		"binding_3",
		basicSidecarID,
		"cyral_repository.repo_3.id",
		"cyral_sidecar_listener.listener_3.listener_id",
	)
	// Repo 4
	config += formatBasicRepositoryIntoConfig(
		"repo_4",
		accTestName(sidecarBoundPortsDataSourceName, "repo4"),
		"s3",
		"s3.com",
		5678,
	)
	config += formatBasicSidecarListenerIntoConfig(
		"listener_4",
		basicSidecarID,
		"s3",
		5678,
	)
	config += formatBasicRepositoryBindingIntoConfig(
		"binding_4",
		basicSidecarID,
		"cyral_repository.repo_4.id",
		"cyral_sidecar_listener.listener_4.listener_id",
	)
	config += fmt.Sprintf(
		`
	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
		// depends_on is needed here so that we can retrieve the sidecar bound ports
		// only after the bindings are created. Otherwise, the data source would
		// retrieve the bound ports before the bindings are created, which in
		// this case would be zero ports.
		depends_on = [
			cyral_repository_binding.binding_1,
			cyral_repository_binding.binding_2,
			cyral_repository_binding.binding_3,
			cyral_repository_binding.binding_4
		]
		sidecar_id = %s
	}`, basicSidecarID,
	)

	return config
}

func testAccSidecarBoundPortsCheck_MultipleBindings() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.#", "4",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.0", "1234",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.1", "3306",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.2", "5678",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.3", "27017",
		),
	)
}
