package cyral

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccSidecarBoundPortsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
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
	})
}

func testAccSidecarBoundPortsConfig_EmptySidecarID() string {
	return `
	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
	}
	`
}

func testAccSidecarBoundPortsConfig_NoBindings() string {
	return `
	resource "cyral_sidecar" "sidecar_1" {
		name = "tf-provider-sidecar-1-bound-ports-sidecar"
		deployment_method = "cloudFormation"
	}

	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
		sidecar_id = cyral_sidecar.sidecar_1.id
	}
	`
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
	return `
	resource "cyral_sidecar" "sidecar_1" {
		name = "tf-provider-sidecar-1-bound-ports-sidecar"
		deployment_method = "cloudFormation"
	}

	resource "cyral_repository" "repo_1" {
		name = "tf-provider-repo-1-bound-ports-sidecar"
		type = "mysql"
		host = "mysql.com"
		port = 3306
	}

	resource "cyral_repository_binding" "repo_binding_1" {
		repository_id = cyral_repository.repo_1.id
		sidecar_id = cyral_sidecar.sidecar_1.id
		listener_port = 3306
		enabled = true
	}

	resource "cyral_repository" "repo_2" {
		name = "tf-provider-repo-2-bound-ports-sidecar"
		type = "mongodb"
		host = "mongodb.com"
		port = 27017
	}

	resource "cyral_repository_binding" "repo_binding_2" {
		repository_id = cyral_repository.repo_2.id
		sidecar_id = cyral_sidecar.sidecar_1.id
		listener_port = 27017
		enabled = true
	}

	data "cyral_sidecar_bound_ports" "sidecar_bound_ports_1" {
		// depends_on is needed here so that we can retrieve the sidecar bound ports
		// only after the bindings are created. Otherwise, the data source would
		// retrieve the bound ports before the bindings are created, which in
		// this case would be zero ports.
		depends_on = [
			cyral_repository_binding.repo_binding_1,
			cyral_repository_binding.repo_binding_2
		]
		sidecar_id = cyral_sidecar.sidecar_1.id
	}
	`
}

func testAccSidecarBoundPortsCheck_MultipleBindings() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.#", "2",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.0", "3306",
		),
		resource.TestCheckResourceAttr(
			"data.cyral_sidecar_bound_ports.sidecar_bound_ports_1",
			"bound_ports.1", "27017",
		),
	)
}

func TestGetBindingPorts_NoPorts(t *testing.T) {
	ports := getBindingPorts(BindingConfig{}, RepoData{})

	assert.Len(t, ports, 0)
}

func TestGetBindingPorts_SinglePort(t *testing.T) {
	binding := BindingConfig{
		Listener: &WrapperListener{
			Port: 1234,
		},
	}
	ports := getBindingPorts(binding, RepoData{})

	expectedPorts := []uint32{1234}

	assert.Equal(t, expectedPorts, ports)
}

func TestGetBindingPorts_MultiplePorts(t *testing.T) {
	binding := BindingConfig{
		Listener: &WrapperListener{
			Port: 1234,
		},
		TcpListeners: &TCPListeners{
			Listeners: []*TCPListener{
				{
					Port: 47017,
				},
				{
					Port: 37017,
				},
			},
		},
		AdditionalListeners: []*TCPListener{
			{
				Port: 457,
			},
			{
				Port: 443,
			},
		},
	}
	repo := RepoData{
		MaxAllowedListeners: 3,
	}
	ports := getBindingPorts(binding, repo)

	expectedPorts := []uint32{443, 457, 1234, 1235, 1236, 37017, 47017}

	assert.ElementsMatch(t, expectedPorts, ports)
}
