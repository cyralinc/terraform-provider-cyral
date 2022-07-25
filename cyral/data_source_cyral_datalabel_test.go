package cyral

import (
	"fmt"
	"testing"
	//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func datalabelDataSourceTestDataLabels() []DataLabel {
	return []DataLabel{
		// TODO
	}
}

func TestAccDatalabelDataSource(t *testing.T) {
	// TODO

	// resource.Test(t, resource.TestCase{
	// 	ProviderFactories: providerFactories,
	// 	Steps: []resource.TestStep{
	// 	},
	// })
}

func datalabelDataSourceConfig(nameFilter, typeFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_datalabel" "test_datalabel" {
		depends_on = [%s]
		name = "%s"
		type = "%s"
	}`, formatAttributes(dependsOn), nameFilter, typeFilter)
}
