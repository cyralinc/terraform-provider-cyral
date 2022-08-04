package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	create = "create"
	read   = "read"
	update = "update"
	delete = "delete"
)

type URLCreatorFunc = func(d *schema.ResourceData, c *client.Client) string

type ResourceData interface {
	ReadFromSchema(d *schema.ResourceData) error
}

type ResponseData interface {
	WriteToSchema(d *schema.ResourceData) error
}

type ResourceOperationConfig struct {
	Name            string
	HttpMethod      string
	CreateURL       URLCreatorFunc
	NewResourceData func() ResourceData
	NewResponseData func(d *schema.ResourceData) ResponseData
}

func CreateResource(createConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return HandleRequest(create, createConfig, &readConfig)
}

func ReadResource(readConfig ResourceOperationConfig) schema.ReadContextFunc {
	return HandleRequest(read, readConfig, nil)
}

func UpdateResource(updateConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return HandleRequest(update, updateConfig, &readConfig)
}

func DeleteResource(deleteConfig ResourceOperationConfig) schema.DeleteContextFunc {
	return HandleRequest(delete, deleteConfig, nil)
}

func HandleRequest(
	operationType string,
	config ResourceOperationConfig,
	consecutiveConfig *ResourceOperationConfig,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		var resourceData ResourceData
		if config.NewResourceData != nil {
			if resourceData = config.NewResourceData(); resourceData != nil {
				if err := resourceData.ReadFromSchema(d); err != nil {
					return createError(
						fmt.Sprintf("Unable to %s resource", operationType),
						err.Error(),
					)
				}
			}
		}

		url := config.CreateURL(d, c)

		body, err := c.DoRequest(url, config.HttpMethod, resourceData)
		if err != nil {
			return createError(
				fmt.Sprintf("Unable to %s resource", operationType),
				err.Error(),
			)
		}

		if config.NewResponseData != nil {
			if responseData := config.NewResponseData(d); responseData != nil {
				if err := json.Unmarshal(body, responseData); err != nil {
					return createError("Unable to unmarshall JSON", err.Error())
				}
				log.Printf("[DEBUG] Response body (unmarshalled): %#v", responseData)

				if err := responseData.WriteToSchema(d); err != nil {
					return createError(
						fmt.Sprintf("Unable to %s resource", operationType),
						err.Error(),
					)
				}
			}
		}

		log.Printf("[DEBUG] End %s", config.Name)

		if consecutiveConfig != nil {
			return HandleRequest(read, *consecutiveConfig, nil)(ctx, d, m)
		}
		return diag.Diagnostics{}
	}
}

type IDBasedResponse struct {
	ID string `json:"id"`
}

func (response IDBasedResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(response.ID)
	return nil
}
