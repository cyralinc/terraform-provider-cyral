package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	ot "github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	rt "github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Implementation of a default context handler that can be used by all resources
// which API follows these principles:
//  1. The resource is backed by an ID coming from the API.
//  2. The creation is a POST that returns a JSON with an `id` field, meaning
//     it can be used with the `IDBasedResponse` struct.
//  3. The endpoint to perform GET, PUT and DELETE calls are composed by the
//     POST endpoint plus the ID specification like the following:
//     - POST:    https://<CP>/<apiVersion>/<featureName>
//     - GET:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - PUT:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - DELETE:  https://<CP>/<apiVersion>/<featureName>/<id>
type DefaultContextHandler struct {
	ResourceName        string
	ResourceType        rt.ResourceType
	SchemaReaderFactory SchemaReaderFactoryFunc
	SchemaWriterFactory SchemaWriterFactoryFunc
	BaseURLFactory      URLFactoryFunc
}

func defaultSchemaWriterFactory(d *schema.ResourceData) SchemaWriter {
	return &IDBasedResponse{}
}

func defaultOperationHandler(
	resourceName string,
	resourceType rt.ResourceType,
	operationType ot.OperationType,
	baseURLFactory URLFactoryFunc,
	httpMethod string,
	schemaReaderFactory SchemaReaderFactoryFunc,
	schemaWriterFactory SchemaWriterFactoryFunc,
) ResourceOperationConfig {
	// POST = https://<CP>/<apiVersion>/<feature>
	// GET, PUT and DELETE = https://<CP>/<apiVersion>/<feature>/<id>
	endpoint := func(d *schema.ResourceData, c *client.Client) string {
		url := baseURLFactory(d, c)
		if d.Id() != "" {
			url = fmt.Sprintf("%s/%s", baseURLFactory(d, c), d.Id())
		}
		tflog.Debug(context.Background(), fmt.Sprintf("Returning base URL for %s '%s' operation '%s' and httpMethod %s: %s",
			resourceType, resourceName, operationType, httpMethod, url))
		return url
	}

	var errorHandler RequestErrorHandler
	if httpMethod == http.MethodGet {
		errorHandler = &ReadIgnoreHttpNotFound{ResName: resourceName}
	} else if httpMethod == http.MethodDelete {
		errorHandler = &DeleteIgnoreHttpNotFound{ResName: resourceName}
	}
	result := ResourceOperationConfig{
		ResourceName:        resourceName,
		Type:                operationType,
		ResourceType:        resourceType,
		HttpMethod:          httpMethod,
		URLFactory:          endpoint,
		SchemaReaderFactory: schemaReaderFactory,
		SchemaWriterFactory: schemaWriterFactory,
		RequestErrorHandler: errorHandler,
	}

	return result
}

func (dch DefaultContextHandler) CreateContext() schema.CreateContextFunc {
	return CreateResource(
		defaultOperationHandler(dch.ResourceName, dch.ResourceType, ot.Create, dch.BaseURLFactory, http.MethodPost, dch.SchemaReaderFactory, nil),
		defaultOperationHandler(dch.ResourceName, dch.ResourceType, ot.Create, dch.BaseURLFactory, http.MethodGet, nil, dch.SchemaWriterFactory),
	)
}

func (dch DefaultContextHandler) ReadContext() schema.ReadContextFunc {
	return ReadResource(dch.ReadResourceOperationConfig())
}
func (dch DefaultContextHandler) ReadResourceOperationConfig() ResourceOperationConfig {
	return defaultOperationHandler(dch.ResourceName, dch.ResourceType, ot.Read, dch.BaseURLFactory, http.MethodGet, nil, dch.SchemaWriterFactory)
}

func (dch DefaultContextHandler) UpdateContext() schema.UpdateContextFunc {
	return UpdateResource(
		defaultOperationHandler(dch.ResourceName, dch.ResourceType, ot.Update, dch.BaseURLFactory, http.MethodPut, dch.SchemaReaderFactory, nil),
		defaultOperationHandler(dch.ResourceName, dch.ResourceType, ot.Update, dch.BaseURLFactory, http.MethodGet, nil, dch.SchemaWriterFactory))
}

func (dch DefaultContextHandler) DeleteContext() schema.DeleteContextFunc {
	return DeleteResource(defaultOperationHandler(
		dch.ResourceName, dch.ResourceType, ot.Delete, dch.BaseURLFactory, http.MethodDelete, nil, nil))
}
