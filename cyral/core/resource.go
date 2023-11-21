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
			tflog.Debug(ctx, fmt.Sprintf("Init %s - %s", operation.ResourceName, operation.Type))
			c := m.(*client.Client)

			var resourceData SchemaReader
			if operation.SchemaReaderFactory != nil {
				if resourceData = operation.SchemaReaderFactory(); resourceData != nil {
					tflog.Debug(ctx, fmt.Sprintf("Calling ReadFromSchema. Schema: %#v", d))
					if err := resourceData.ReadFromSchema(d); err != nil {
						return utils.CreateError(
							fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.ResourceName),
							err.Error(),
						)
					}
					tflog.Debug(ctx, fmt.Sprintf("Succesful call to ReadFromSchema. resourceData: %#v", resourceData))
				}
			}

			url := operation.URLFactory(d, c)

			body, err := c.DoRequest(url, operation.HttpMethod, resourceData)
			if operation.RequestErrorHandler != nil {
				err = operation.RequestErrorHandler.HandleError(d, c, err)
			}
			if err != nil {
				return utils.CreateError(
					fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.ResourceName),
					err.Error(),
				)
			}

			// If a `SchemaWriterFactory` implementation is NOT provided and this is a creation operation,
			// use the `defaultSchemaWriterFactory`, assuming the response is a JSON with an `id` field.
			/// TODO: Remove this feature after refactoring all resources to use the `DefaultContext`.
			var responseDataFunc SchemaWriterFactoryFunc
			if body != nil {
				if operation.SchemaWriterFactory == nil && operation.Type == operationtype.Create {
					responseDataFunc = defaultSchemaWriterFactory
					tflog.Debug(ctx, "NewResponseData function set to defaultSchemaWriterFactory.")
				} else {
					responseDataFunc = operation.SchemaWriterFactory
				}
			}
			if responseDataFunc != nil {
				if responseData := responseDataFunc(d); responseData != nil {
					tflog.Debug(ctx, fmt.Sprintf("NewResponseData function call performed. d: %#v", d))
					if err := json.Unmarshal(body, responseData); err != nil {
						return utils.CreateError("Unable to unmarshall JSON", err.Error())
					}
					tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", responseData))
					tflog.Debug(ctx, fmt.Sprintf("Calling WriteToSchema: responseData: %#v", responseData))
					if err := responseData.WriteToSchema(d); err != nil {
						return utils.CreateError(
							fmt.Sprintf("Unable to %s resource %s", operation.Type, operation.ResourceName),
							err.Error(),
						)
					}
					tflog.Debug(ctx, fmt.Sprintf("Succesful call to WriteToSchema. d: %#v", d))
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("End %s - %s", operation.ResourceName, operation.Type))
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
