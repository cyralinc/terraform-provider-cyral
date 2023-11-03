package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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

// Teaches a resource or data source how to read from the Terraform schema and
// store in the data structure defined for it.
type SchemaReader interface {
	ReadFromSchema(d *schema.ResourceData) error
}

// Teaches a resource or data source how to write to the Terraform schema from
// the data stored in the data structure defined for it.
type SchemaWriter interface {
	WriteToSchema(d *schema.ResourceData) error
}

type SchemaType string

const (
	DataSourceSchemaType = SchemaType("dataSource")
	ResourceSchemaType   = SchemaType("resource")
)

// The `SchemaDescriptor` describes the resource for a given schema.
type SchemaDescriptor struct {
	// Resource or data source name
	Name   string
	Type   SchemaType
	Schema func() *schema.Resource
}

// The `PackageSchema` is used to centralize the description of the existing
// schemas in a given package. It should be implemented in the `schema.go`
// file of a given package.
type PackageSchema interface {
	Name() string
	Schemas() []*SchemaDescriptor
}

type ResourceOperationConfig struct {
	Name       string
	HttpMethod string
	CreateURL  URLCreatorFunc
	RequestErrorHandler
	NewResourceData func() SchemaReader
	// TODO provide a default implementation returning the IDBasedResponse:
	// func(_ *schema.ResourceData) core.SchemaWriter { return &core.IDBasedResponse{} }
	NewResponseData func(d *schema.ResourceData) SchemaWriter
}

func CRUDResources(resourceOperations []ResourceOperation) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
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
) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		for _, operation := range resourceOperations {
			log.Printf("[DEBUG] Init %s", operation.Config.Name)
			c := m.(*client.Client)

			var resourceData SchemaReader
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
