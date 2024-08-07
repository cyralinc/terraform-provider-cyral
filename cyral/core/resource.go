package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type URLFactoryFunc = func(d *schema.ResourceData, c *client.Client) string
type SchemaReaderFactoryFunc = func() SchemaReader
type SchemaWriterFactoryFunc = func(d *schema.ResourceData) SchemaWriter

type RequestErrorHandler interface {
	HandleError(ctx context.Context, d *schema.ResourceData, c *client.Client, err error) error
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
	// Human-readable resource name that will be used in log messages
	ResourceName string
	Type         operationtype.OperationType
	// Resource type
	ResourceType        resourcetype.ResourceType
	HttpMethod          string
	URLFactory          URLFactoryFunc
	RequestErrorHandler RequestErrorHandler
	SchemaReaderFactory SchemaReaderFactoryFunc
	SchemaWriterFactory SchemaWriterFactoryFunc
}

func CRUDResources(operations []ResourceOperationConfig) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	return handleRequests(operations)
}

func CreateResource(createConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return handleRequests(
		[]ResourceOperationConfig{
			createConfig, readConfig,
		},
	)
}

func ReadResource(readConfig ResourceOperationConfig) schema.ReadContextFunc {
	return handleRequests(
		[]ResourceOperationConfig{
			readConfig,
		},
	)
}

func UpdateResource(updateConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return handleRequests(
		[]ResourceOperationConfig{
			updateConfig, readConfig,
		},
	)
}

func DeleteResource(deleteConfig ResourceOperationConfig) schema.DeleteContextFunc {
	return handleRequests(
		[]ResourceOperationConfig{
			deleteConfig,
		},
	)
}

func handleRequests(operations []ResourceOperationConfig) func(context.Context, *schema.ResourceData, any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		for _, operation := range operations {
			tflog.Debug(ctx, fmt.Sprintf("Init handleRequests to %s %s %s", operation.Type, operation.ResourceType, operation.ResourceName))
			c := m.(*client.Client)

			var resourceData SchemaReader
			if operation.SchemaReaderFactory != nil {
				tflog.Debug(ctx, "=> Calling SchemaReaderFactory")
				if resourceData = operation.SchemaReaderFactory(); resourceData != nil {
					tflog.Debug(ctx, fmt.Sprintf("=> Calling ReadFromSchema. Schema: %#v", d))
					if err := resourceData.ReadFromSchema(d); err != nil {
						tflog.Debug(ctx, fmt.Sprintf("End handleRequests to %s %s %s - Error: %s", operation.Type, operation.ResourceType, operation.ResourceName, err.Error()))
						return utils.CreateError(
							fmt.Sprintf("Unable to %s %s %s", operation.Type, operation.ResourceType, operation.ResourceName),
							err.Error(),
						)
					}
					tflog.Debug(ctx, fmt.Sprintf("=> Succesful call to ReadFromSchema. resourceData: %#v", resourceData))
				}
			}

			url := operation.URLFactory(d, c)

			body, err := c.DoRequest(ctx, url, operation.HttpMethod, resourceData)
			if err != nil && operation.RequestErrorHandler != nil {
				tflog.Debug(ctx, "=> Calling operation.RequestErrorHandler.HandleError")
				err = operation.RequestErrorHandler.HandleError(ctx, d, c, err)
			}
			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("End handleRequests to %s %s %s - Error: %s", operation.Type, operation.ResourceType, operation.ResourceName, err.Error()))
				return utils.CreateError(
					fmt.Sprintf("Unable to %s %s %s", operation.Type, operation.ResourceType, operation.ResourceName),
					err.Error(),
				)
			}

			if operation.SchemaWriterFactory == nil {
				tflog.Debug(ctx, "=> No SchemaWriterFactory found.")
			} else if body != nil {
				if responseData := operation.SchemaWriterFactory(d); responseData != nil {
					tflog.Debug(ctx, fmt.Sprintf("=> operation.SchemaWriterFactory function call performed. d: %#v", d))
					if err := json.Unmarshal(body, responseData); err != nil {
						tflog.Debug(ctx, fmt.Sprintf("End handleRequests to %s %s %s - Error: %s", operation.Type, operation.ResourceType, operation.ResourceName, err.Error()))
						return utils.CreateError("Unable to unmarshall JSON", err.Error())
					}
					tflog.Debug(ctx, fmt.Sprintf("=> Response body (unmarshalled): %#v", responseData))
					tflog.Debug(ctx, fmt.Sprintf("=> Calling WriteToSchema: responseData: %#v", responseData))
					if err := responseData.WriteToSchema(d); err != nil {
						tflog.Debug(ctx, fmt.Sprintf("End handleRequests to %s %s %s - Error: %s", operation.Type, operation.ResourceType, operation.ResourceName, err.Error()))
						return utils.CreateError(
							fmt.Sprintf("Unable to %s %s %s", operation.Type, operation.ResourceType, operation.ResourceName),
							err.Error(),
						)
					}
					tflog.Debug(ctx, fmt.Sprintf("=> Succesful call to WriteToSchema. d: %#v", d))
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("End handleRequests to %s %s %s - Success", operation.Type, operation.ResourceType, operation.ResourceName))
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
