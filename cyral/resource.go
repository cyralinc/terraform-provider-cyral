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
	create = ResourceType("create")
	read   = ResourceType("read")
	update = ResourceType("update")
	delete = ResourceType("delete")
)

type ResourceType string

type ResourceConfig struct {
	Type            ResourceType
	OperationConfig ResourceOperationConfig
}

type URLCreatorFunc = func(d *schema.ResourceData, c *client.Client) string

type RequestErrorHandler interface {
	HandleError(d *schema.ResourceData, c *client.Client, err error) error
}

type ResourceData interface {
	ReadFromSchema(d *schema.ResourceData) error
}

type ResponseData interface {
	WriteToSchema(d *schema.ResourceData) error
}

type ResourceOperationConfig struct {
	Name       string
	HttpMethod string
	CreateURL  URLCreatorFunc
	RequestErrorHandler
	NewResourceData func() ResourceData
	NewResponseData func(d *schema.ResourceData) ResponseData
}

func CRUDResources(config []ResourceConfig) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return HandleRequest(config)
}

func CreateResource(createConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return HandleRequest(
		[]ResourceConfig{
			{
				Type:            create,
				OperationConfig: createConfig,
			},
			{
				Type:            read,
				OperationConfig: readConfig,
			},
		},
	)
}

func ReadResource(readConfig ResourceOperationConfig) schema.ReadContextFunc {
	return HandleRequest(
		[]ResourceConfig{
			{
				Type:            read,
				OperationConfig: readConfig,
			},
		},
	)
}

func UpdateResource(updateConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return HandleRequest(
		[]ResourceConfig{
			{
				Type:            update,
				OperationConfig: updateConfig,
			},
			{
				Type:            read,
				OperationConfig: readConfig,
			},
		},
	)
}

func DeleteResource(deleteConfig ResourceOperationConfig) schema.DeleteContextFunc {
	return HandleRequest(
		[]ResourceConfig{
			{
				Type:            delete,
				OperationConfig: deleteConfig,
			},
		},
	)
}

func HandleRequest(
	configSlice []ResourceConfig,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		for _, config := range configSlice {
			log.Printf("[DEBUG] Init %s", config.OperationConfig.Name)
			c := m.(*client.Client)

			var resourceData ResourceData
			if config.OperationConfig.NewResourceData != nil {
				if resourceData = config.OperationConfig.NewResourceData(); resourceData != nil {
					if err := resourceData.ReadFromSchema(d); err != nil {
						return createError(
							fmt.Sprintf("Unable to %s resource", config.Type),
							err.Error(),
						)
					}
				}
			}

			url := config.OperationConfig.CreateURL(d, c)

			body, err := c.DoRequest(url, config.OperationConfig.HttpMethod, resourceData)
			if config.OperationConfig.RequestErrorHandler != nil {
				err = config.OperationConfig.RequestErrorHandler.HandleError(d, c, err)
			}
			if err != nil {
				return createError(
					fmt.Sprintf("Unable to %s resource", config.Type),
					err.Error(),
				)
			}

			if body != nil && config.OperationConfig.NewResponseData != nil {
				if responseData := config.OperationConfig.NewResponseData(d); responseData != nil {
					if err := json.Unmarshal(body, responseData); err != nil {
						return createError("Unable to unmarshall JSON", err.Error())
					}
					log.Printf("[DEBUG] Response body (unmarshalled): %#v", responseData)

					if err := responseData.WriteToSchema(d); err != nil {
						return createError(
							fmt.Sprintf("Unable to %s resource", config.Type),
							err.Error(),
						)
					}
				}
			}

			log.Printf("[DEBUG] End %s", config.OperationConfig.Name)
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
