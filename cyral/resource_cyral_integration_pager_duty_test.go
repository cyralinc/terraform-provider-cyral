package cyral

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialPagerDutyIntegrationConfig PagerDutyIntegration = PagerDutyIntegration{

	ID:       "unitTest",
	Name:     "unitTest",
	APIToken: "unitTest",
}

var updatedPagerDutyIntegrationConfig PagerDutyIntegration = PagerDutyIntegration{

	ID:       "unitTest-updated",
	Name:     "unitTest-updated",
	APIToken: "unitTest-updated",
}

func TestAccPagerDutyIntegrationResource(t *testing.T) {
	testConfig, testFunc := setupPagerDutyIntegrationTest(initialPagerDutyIntegrationConfig)
	testUpdateConfig, testUpdateFunc := setupPagerDutyIntegrationTest(updatedPagerDutyIntegrationConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
		},
	})
}

func setupPagerDutyIntegrationTest(integrationData PagerDutyIntegration) (string, resource.TestCheckFunc) {
	configuration := formatPagerDutyIntegrationDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(

		resource.TestCheckResourceAttr("cyral_pager_duty_integration.pager_duty_integration", "id", integrationData.ID),
		resource.TestCheckResourceAttr("cyral_pager_duty_integration.pager_duty_integration", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_pager_duty_integration.pager_duty_integration", "api_token", integrationData.APIToken),
	)

	return configuration, testFunction
}

func formatPagerDutyIntegrationDataIntoConfig(data PagerDutyIntegration) string {
	return fmt.Sprintf(`
	resource "cyral_integration_pager_duty" "pager_duty_integration" {
	id = %s
	name = %s
	api_token = %s
	}`, data.ID, data.Name, data.APIToken)
}

func TestPagerDutyIntegration_MarshalJSON(t *testing.T) {
	type fields struct {
		ID       string
		Name     string
		APIToken string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "basic marshalling",
			fields: fields{
				ID:       "123456",
				Name:     "pager-duty-integration",
				APIToken: "this-is-a-token",
			},
			want: func() []byte {
				marshalled, _ := json.Marshal(struct {
					Name        string `json:"name"`
					Category    string `json:"category"`
					BuiltInType string `json:"builtInType"`
					Parameters  string `json:"parameters"`
				}{
					Name:        "pager-duty-integration",
					Category:    "builtin",
					BuiltInType: "pagerduty",
					Parameters:  "{\"apiToken\":\"this-is-a-token\"}",
				})
				return marshalled
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := PagerDutyIntegration{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				APIToken: tt.fields.APIToken,
			}
			got, err := json.Marshal(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("PagerDutyIntegration.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PagerDutyIntegration.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestPagerDutyIntegration_UnmarshalJSON(t *testing.T) {
	type fields struct {
		ID       string
		Name     string
		APIToken string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    PagerDutyIntegration
	}{
		{
			name: "basic testing",
			fields: fields{
				ID: "this-is-an-id",
			},
			args: args{
				func() []byte {
					marshalled, _ := json.Marshal(struct {
						Name           string `json:"name"`
						Category       string `json:"category"`
						BuiltInType    string `json:"builtInType"`
						Parameters     string `json:"parameters"`
						AuthTemplateID string `json:"authorizationTemplateID"`
					}{
						Name:           "pager-duty-integration",
						Category:       "builtin",
						BuiltInType:    "pagerduty",
						Parameters:     "{\"apiToken\":\"this-is-a-token\"}",
						AuthTemplateID: "uniquelygeneratedIDforpolicytemplate",
					})
					return marshalled
				}(),
			},
			want: PagerDutyIntegration{
				Name:     "pager-duty-integration",
				APIToken: "this-is-a-token",
				ID:       "this-is-an-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := PagerDutyIntegration{
				ID:       tt.fields.ID,
				Name:     tt.fields.Name,
				APIToken: tt.fields.APIToken,
			}
			if err := json.Unmarshal(tt.args.b, &data); (err != nil) != tt.wantErr {
				t.Errorf("PagerDutyIntegration.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, data) {
				t.Errorf("PagerDutyIntegration.UnmarshalJSON() got = %v, want = %v", data, tt.want)
			}
		})
	}
}
