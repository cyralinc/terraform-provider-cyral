package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/src/client"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResourceOperation struct {
	Type   OperationType
	Config ResourceOperationConfig
}

type URLCreatorFunc = func(d *schema.ResourceData, c *client.Client) string

type RequestErrorHandler interface {
	HandleError(d *schema.ResourceData, c *client.Client, err error) error
}

// TODO Rename as `SchemaReader` and document properly.
// Teaches a resource or data source how to read from the Terraform schema and
// store in the data structure defined for it.
type ResourceData interface {
	ReadFromSchema(d *schema.ResourceData) error
}

// TODO Rename as `SchemaWriter` and document properly.
// Teaches a resource or data source how to write to the Terraform schema from
// the data stored in the data structure defined for it.
type ResponseData interface {
	WriteToSchema(d *schema.ResourceData) error
}

type SchemaType string

const (
	DataSourceSchema = SchemaType("dataSource")
	ResourceSchema   = SchemaType("resource")
)

// The `SchemaRegister` register was created to decouple the Provider schema from the
// resources and data sources definition. Instead of having the provider to depend on
// each of the resources and data sources, the `SchemaRegister` inverts the dependency
// and makes resources and data sources responsible for registering themselves to the
// provider, thus making the provider code unaware of particular implementations of
// each of its components.
//
// In order to use the `SchemaRegister`, implement an `init()` function that will
// be responsible for providing the resource or data source name, the schema and
// also the resource type. In order to better organize the code, it is
// recommended that this initialization is performed in the `datasource.go` and
// `resource.go` files.
type SchemaRegister struct {
	// Resource or data source name
	Name   string
	Type   SchemaType
	Schema func() *schema.Resource
}

type ResourceOperationConfig struct {
	Name       string
	HttpMethod string
	CreateURL  URLCreatorFunc
	RequestErrorHandler
	NewResourceData func() ResourceData
	// TODO provide a default implementation
	NewResponseData func(d *schema.ResourceData) ResponseData
}

func CRUDResources(resourceOperations []ResourceOperation) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return handleRequests(resourceOperations)
}

func CreateResource(createConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return handleRequests(
		[]ResourceOperation{
			{
				Type:   OperationTypeCreate,
				Config: createConfig,
			},
			{
				Type:   OperationTypeRead,
				Config: readConfig,
			},
		},
	)
}

func ReadResource(readConfig ResourceOperationConfig) schema.ReadContextFunc {
	return handleRequests(
		[]ResourceOperation{
			{
				Type:   OperationTypeRead,
				Config: readConfig,
			},
		},
	)
}

func UpdateResource(updateConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return handleRequests(
		[]ResourceOperation{
			{
				Type:   OperationTypeUpdate,
				Config: updateConfig,
			},
			{
				Type:   OperationTypeRead,
				Config: readConfig,
			},
		},
	)
}

func DeleteResource(deleteConfig ResourceOperationConfig) schema.DeleteContextFunc {
	return handleRequests(
		[]ResourceOperation{
			{
				Type:   OperationTypeDelete,
				Config: deleteConfig,
			},
		},
	)
}

func handleRequests(
	resourceOperations []ResourceOperation,
) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		for _, operation := range resourceOperations {
			log.Printf("[DEBUG] Init %s", operation.Config.Name)
			c := m.(*client.Client)

			var resourceData ResourceData
			if operation.Config.NewResourceData != nil {
				if resourceData = operation.Config.NewResourceData(); resourceData != nil {
					log.Printf("[DEBUG] Calling ReadFromSchema. Schema: %#v", d)
					if err := resourceData.ReadFromSchema(d); err != nil {
						return utils.CreateError(
							fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.Config.Name),
							err.Error(),
						)
					}
					log.Printf("[DEBUG] Succesful call to ReadFromSchema. resourceData: %#v", resourceData)
				}
			}

			url := operation.Config.CreateURL(d, c)

			body, err := c.DoRequest(url, operation.Config.HttpMethod, resourceData)
			if operation.Config.RequestErrorHandler != nil {
				err = operation.Config.RequestErrorHandler.HandleError(d, c, err)
			}
			if err != nil {
				return utils.CreateError(
					fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.Config.Name),
					err.Error(),
				)
			}

			if body != nil && operation.Config.NewResponseData != nil {
				if responseData := operation.Config.NewResponseData(d); responseData != nil {
					log.Printf("[DEBUG] NewResponseData function call performed. d: %#v", d)
					if err := json.Unmarshal(body, responseData); err != nil {
						return utils.CreateError("Unable to unmarshall JSON", err.Error())
					}
					log.Printf("[DEBUG] Response body (unmarshalled): %#v", responseData)
					log.Printf("[DEBUG] Calling WriteToSchema: responseData: %#v", responseData)
					if err := responseData.WriteToSchema(d); err != nil {
						return utils.CreateError(
							fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.Config.Name),
							err.Error(),
						)
					}
					log.Printf("[DEBUG] Succesful call to WriteToSchema. d: %#v", d)
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
