package cyral

import (
	"regexp"
)

const (
	sidecarDataSourceName = "data-sidecar"
)

func testSidecarDataSourceSidecars() []IdentifiedSidecarInfo {
	return []IdentifiedSidecarInfo{
		{
			Sidecar: SidecarData{
				Name:   accTestName(sidecarDataSourceName, "sidecar-1"),
				Labels: []string{"label-1", "label-2"},
				Properties: &SidecarProperties{
					DeploymentMethod: "cloudformation",
				},
				Services: SidecarServicesConfig{
					"dispatcher": map[string]string{
						"bypass": "never",
					},
				},
				UserEndpoint: "user-endpoint-1",
				CertificateBundleSecrets: CertificateBundleSecrets{
					"sidecar": &CertificateBundleSecret{
						Engine:   "engine-1",
						SecretId: "secret-id-1",
						Type:     "secret-type-1",
					},
				},
			},
		},
		{
			Sidecar: SidecarData{
				Name:   accTestName(sidecarDataSourceName, "sidecar-2"),
				Labels: []string{"label-3", "label-4"},
				Properties: &SidecarProperties{
					DeploymentMethod: "terraform",
				},
				Services: SidecarServicesConfig{
					"dispatcher": map[string]string{
						"bypass": "always",
					},
				},
				UserEndpoint: "user-endpoint-2",
				CertificateBundleSecrets: CertificateBundleSecrets{
					"sidecar": &CertificateBundleSecret{
						Engine:   "engine-2",
						SecretId: "secret-id-2",
						Type:     "secret-type-2",
					},
				},
			},
		},
	}
}

func TestAccSidecarDataSource(t *testing.T) {
	sidecars := testSidecarDataSourceSidecars()

	testCompleteSidecars := testSidecarDataSource(sidecars)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			testCompleteSidecar,
		},
	})
}

func testSidecarDataSource(dsourceName string, sidecars []IdentifiedSidecarInfo) resource.TestStep {
	var config string
	var checkFuncs []resource.TestFunc

	for _, sidecar := range sidecars {
		config += formatSidecarDataIntoConfig(&sidecar.Sidecar)
	}

	config += fmt.Sprintf(`
		data "cyral_sidecar" "%s" {
		}`, dsourceName)

	dsourceFullName := fmt.Sprintf("data.cyral_sidecar.%s", dsourceName)
	checkFuncs = append(checkFuncs,
		resource.TestMatchResourceAttr(dsourceFullName,
			"sidecar_list.#", notZeroRegex(),
		),
		resource.TestMatchResourceAttr(dsourceFullName,
			"sidecar_list.0.id", regexp.MustCompile(".+"),
		),
		resource.TestMatchResourceAttr(dsourceFullName,
			"sidecar_list.sidecar.id", notZeroRegex(),
		),
	)

	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checkFuncs...),
	}
}
