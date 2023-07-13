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
	create = OperationType("create")
	read   = OperationType("read")
	update = OperationType("update")
	delete = OperationType("delete")
)

type OperationType string

type ResourceOperation struct {
	Type   OperationType
	Config ResourceOperationConfig
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

func CRUDResources(config []ResourceOperation) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return HandleRequests(config)
}

func CreateResource(createConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return HandleRequests(
		[]ResourceOperation{
			{
				Type:   create,
				Config: createConfig,
			},
			{
				Type:   read,
				Config: readConfig,
			},
		},
	)
}

func ReadResource(readConfig ResourceOperationConfig) schema.ReadContextFunc {
	return HandleRequests(
		[]ResourceOperation{
			{
				Type:   read,
				Config: readConfig,
			},
		},
	)
}

func UpdateResource(updateConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return HandleRequests(
		[]ResourceOperation{
			{
				Type:   update,
				Config: updateConfig,
			},
			{
				Type:   read,
				Config: readConfig,
			},
		},
	)
}

func DeleteResource(deleteConfig ResourceOperationConfig) schema.DeleteContextFunc {
	return HandleRequests(
		[]ResourceOperation{
			{
				Type:   delete,
				Config: deleteConfig,
			},
		},
	)
}

func HandleRequests(
	resourceOperations []ResourceOperation,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		for _, operation := range resourceOperations {
			log.Printf("[DEBUG] Init %s", operation.Config.Name)
			c := m.(*client.Client)

			var resourceData ResourceData
			if operation.Config.NewResourceData != nil {
				if resourceData = operation.Config.NewResourceData(); resourceData != nil {
					if err := resourceData.ReadFromSchema(d); err != nil {
						return createError(
							fmt.Sprintf("Unable to %s resource", operation.Type),
							err.Error(),
						)
					}
				}
			}

			url := operation.Config.CreateURL(d, c)

			body, err := c.DoRequest(url, operation.Config.HttpMethod, resourceData)
			if operation.Config.RequestErrorHandler != nil {
				err = operation.Config.RequestErrorHandler.HandleError(d, c, err)
			}
			if err != nil {
				return createError(
					fmt.Sprintf("Unable to %s resource", operation.Type),
					err.Error(),
				)
			}

			if body != nil && operation.Config.NewResponseData != nil {
				if responseData := operation.Config.NewResponseData(d); responseData != nil {
					if err := json.Unmarshal(body, responseData); err != nil {
						return createError("Unable to unmarshall JSON", err.Error())
					}
					log.Printf("[DEBUG] Response body (unmarshalled): %#v", responseData)

					if err := responseData.WriteToSchema(d); err != nil {
						return createError(
							fmt.Sprintf("Unable to %s resource", operation.Type),
							err.Error(),
						)
					}
				}
			}

			log.Printf("[DEBUG] End %s", operation.Config.Name)
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
