package internal_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	TestMetadataURL = "https://some-test-metadata-url.com"
	Base64Doc       = "PG1kOkVudGl0eURlc2NyaXB0b3IgeG1sbnM6bWQ9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDptZXRhZGF0YSIgZW50aXR5SUQ9Imh0dHA6Ly93d3cub2t0YS5jb20vZXhrMXNhZm84a3JFQXA4aTA1ZDciPgo8bWQ6SURQU1NPRGVzY3JpcHRvciBXYW50QXV0aG5SZXF1ZXN0c1NpZ25lZD0iZmFsc2UiIHByb3RvY29sU3VwcG9ydEVudW1lcmF0aW9uPSJ1cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoyLjA6cHJvdG9jb2wiPgo8bWQ6S2V5RGVzY3JpcHRvciB1c2U9InNpZ25pbmciPgo8ZHM6S2V5SW5mbyB4bWxuczpkcz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC8wOS94bWxkc2lnIyI+CjxkczpYNTA5RGF0YT4KPGRzOlg1MDlDZXJ0aWZpY2F0ZT5NSUlEcURDQ0FwQ2dBd0lCQWdJR0FYcWJuaktnTUEwR0NTcUdTSWIzRFFFQkN3VUFNSUdVTVFzd0NRWURWUVFHRXdKVlV6RVRNQkVHIEExVUVDQXdLUTJGc2FXWnZjbTVwWVRFV01CUUdBMVVFQnd3TlUyRnVJRVp5WVc1amFYTmpiekVOTUFzR0ExVUVDZ3dFVDJ0MFlURVUgTUJJR0ExVUVDd3dMVTFOUFVISnZkbWxrWlhJeEZUQVRCZ05WQkFNTURHUmxkaTAyTVRJMk5UY3lOekVjTUJvR0NTcUdTSWIzRFFFSiBBUllOYVc1bWIwQnZhM1JoTG1OdmJUQWVGdzB5TVRBM01USXhOalEyTlRSYUZ3MHpNVEEzTVRJeE5qUTNOVE5hTUlHVU1Rc3dDUVlEIFZRUUdFd0pWVXpFVE1CRUdBMVVFQ0F3S1EyRnNhV1p2Y201cFlURVdNQlFHQTFVRUJ3d05VMkZ1SUVaeVlXNWphWE5qYnpFTk1Bc0cgQTFVRUNnd0VUMnQwWVRFVU1CSUdBMVVFQ3d3TFUxTlBVSEp2ZG1sa1pYSXhGVEFUQmdOVkJBTU1ER1JsZGkwMk1USTJOVGN5TnpFYyBNQm9HQ1NxR1NJYjNEUUVKQVJZTmFXNW1iMEJ2YTNSaExtTnZiVENDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DIGdnRUJBS1Q1WDQvMEZqL3p6MXloTUZoY2VNVlh1RUZVWU1JYllpNUtaNE1OcTd0SjYwQjh6SkMvVm42YjdWRHNORk9ibGE2a05URXUgNEFxelAxWnNBQ2FWSmVKSjJZSlAzWWFEOHplN3dkMTQ2dmdaZnRHSmg1a2xUekZ2YWdUck11MHloVjBaTk5CUE1rYXFXT1V0MnhYRSAyVnZnalFrVTZ1Y0JHRjFwMy84eXhqVFlHaDd1eS83ZXIzaHZqK0s5cGNlUkMyQW1NaHZGaGVvbnk3MVdTVUIxZHZTY0kzNUZ4TzkxIDhjVDRSbkRrYmhpOVhXNXpSaXBzL0VjS2JqSEI1VzVpL2l2UGltbjdjLy8xUXBQV2FjdUFtNVRCNDBpRjA4ODFVSkd4OERXTmo5SnogUEU4T0h4YkJlMjJrVC93RUp2aEZwbnpFK0QzM2l5N01TdmhENE1uNHZva0NBd0VBQVRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQSBsWUF0cmZKU2FjODk3RnJFZW5ZNHZxaFU3ajFVRlJQMXZzdXRiN2dXdWUwSmZOSU0yK3VORWJLM09wcDFCb05qTzkwcUJsQVVOSlhuIGFkQ0JBd000NHR2WTZOZjhtRGx4c3YwdXFrOEU1RWZDVERGemh5Vlg1bDhPVVZnNzdlYzZORThISk53SXdWMXhMenpIaUgxOGQ3ZVUgbzhZem41a0VaeHhLNzhWUnBOOXE3UkwzamJOMUtWRDU4cXpiRWRucWFXVU43Z3EreXI0K2hKenpKZnZ3TDdOd1ZpaVRBRDliZTAvVSAzRDA0SnRRa0Z0blZseTB6YUI5TElPWE9iL3FpRUl4WElNMU1OVVpLc1hNQkU0WHVGdnVMUzV6SFZyd1JIZVovVXhZNEFMS1NuQnRDIGZnandxd21WbCt0Z0liSjBvRmRBVmdEbDVZSVRFUzBKTHY5WEl3PT08L2RzOlg1MDlDZXJ0aWZpY2F0ZT4KPC9kczpYNTA5RGF0YT4KPC9kczpLZXlJbmZvPgo8L21kOktleURlc2NyaXB0b3I+CjxtZDpOYW1lSURGb3JtYXQ+dXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6MS4xOm5hbWVpZC1mb3JtYXQ6dW5zcGVjaWZpZWQ8L21kOk5hbWVJREZvcm1hdD4KPG1kOk5hbWVJREZvcm1hdD51cm46b2FzaXM6bmFtZXM6dGM6U0FNTDoxLjE6bmFtZWlkLWZvcm1hdDplbWFpbEFkZHJlc3M8L21kOk5hbWVJREZvcm1hdD4KPG1kOlNpbmdsZVNpZ25PblNlcnZpY2UgQmluZGluZz0idXJuOm9hc2lzOm5hbWVzOnRjOlNBTUw6Mi4wOmJpbmRpbmdzOkhUVFAtUE9TVCIgTG9jYXRpb249Imh0dHBzOi8vZGV2LTYxMjY1NzI3Lm9rdGEuY29tL2FwcC9kZXYtNjEyNjU3MjdfdGVzdF8xL2V4azFzYWZvOGtyRUFwOGkwNWQ3L3Nzby9zYW1sIi8+CjxtZDpTaW5nbGVTaWduT25TZXJ2aWNlIEJpbmRpbmc9InVybjpvYXNpczpuYW1lczp0YzpTQU1MOjIuMDpiaW5kaW5nczpIVFRQLVJlZGlyZWN0IiBMb2NhdGlvbj0iaHR0cHM6Ly9kZXYtNjEyNjU3Mjcub2t0YS5jb20vYXBwL2Rldi02MTI2NTcyN190ZXN0XzEvZXhrMXNhZm84a3JFQXA4aTA1ZDcvc3NvL3NhbWwiLz4KPC9tZDpJRFBTU09EZXNjcmlwdG9yPgo8L21kOkVudGl0eURlc2NyaXB0b3I+"
)

func TestAccSAMLConfigurationDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccSAMLConfigurationConfig_EmptyMetadata(),
				ExpectError: regexp.MustCompile(`Error: Invalid combination of arguments`),
			},
			{
				Config:      testAccSAMLConfigurationConfig_BothMetadataType(),
				ExpectError: regexp.MustCompile(`Error: Invalid combination of arguments`),
			},
			{
				Config: testAccSAMLConfigurationConfig_MetadataURL(),
				Check:  testAccSAMLConfigurationCheck_MetadataURL(),
			},
			{
				Config: testAccSAMLConfigurationConfig_MetadataBase64Doc(),
				Check:  testAccSAMLConfigurationCheck_MetadataBase64Doc(),
			},
		},
	})
}

func testAccSAMLConfigurationConfig_EmptyMetadata() string {
	return `
	data "cyral_saml_configuration" "test_saml_configuration" {
	}
	`
}

func testAccSAMLConfigurationConfig_BothMetadataType() string {
	return fmt.Sprintf(`
	data "cyral_saml_configuration" "test_saml_configuration" {
		saml_metadata_url = "%s"
		base_64_saml_metadata_document = "%s"
	}
	`, TestMetadataURL, Base64Doc)
}

func testAccSAMLConfigurationConfig_MetadataURL() string {
	return fmt.Sprintf(`
	data "cyral_saml_configuration" "test_saml_configuration" {
		saml_metadata_url = "%s"
	}
	`, TestMetadataURL)
}

func testAccSAMLConfigurationCheck_MetadataURL() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.cyral_saml_configuration.test_saml_configuration",
			"saml_metadata_url", TestMetadataURL),
		resource.TestCheckNoResourceAttr("data.cyral_saml_configuration.test_saml_configuration",
			"base_64_saml_metadata_document"),
	)
}

func testAccSAMLConfigurationConfig_MetadataBase64Doc() string {
	return fmt.Sprintf(`
	data "cyral_saml_configuration" "test_saml_configuration" {
		base_64_saml_metadata_document = "%s"
	}
	`, Base64Doc)
}

func testAccSAMLConfigurationCheck_MetadataBase64Doc() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.cyral_saml_configuration.test_saml_configuration",
			"base_64_saml_metadata_document", Base64Doc),
		resource.TestCheckNoResourceAttr("data.cyral_saml_configuration.test_saml_configuration",
			"samlMetadataURL"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_using_jwks_url"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"name_id_policy_format"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"hide_on_login_page"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"back_channel_supported"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_post_binding_response"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_post_binding_authn_request"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_post_binding_logout"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_want_authn_requests_signed"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_want_assertions_signed"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"want_assertions_encrypted"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_force_authentication"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"disable_validate_signature"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"single_sign_on_service_url"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"signing_certificate"),
		resource.TestCheckResourceAttrSet("data.cyral_saml_configuration.test_saml_configuration",
			"allowed_clock_skew"),
	)
}
