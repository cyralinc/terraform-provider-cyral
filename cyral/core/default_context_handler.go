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

// Implementation of a default context handler that can be used by all resources.
//
//  1. `SchemaWriterFactoryGetMethod“ must be provided.
//  2. In case `SchemaWriterFactoryPostMethod“ is not provided,
//     it will assume that a call to POST returns a JSON with
//     an `id` field, meaning it will use the
//     `IDBasedResponse` struct in such cases.
//  3. `BaseURLFactory` must be provided. It will be used to
//     create the POST endpoint and others in case `IdBasedURLFactory`
//     is not provided.
//  4. If `IdBasedURLFactory` is NOT provided, the endpoint to
//     perform GET, PUT and DELETE calls are composed by the
//     `BaseURLFactory` endpoint plus the ID specification as follows:
//     - POST:    https://<CP>/<apiVersion>/<featureName>
//     - GET:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - PUT:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - DELETE:  https://<CP>/<apiVersion>/<featureName>/<id>
type DefaultContextHandler struct {
	ResourceName        string
	ResourceType        rt.ResourceType
	SchemaReaderFactory SchemaReaderFactoryFunc
	// SchemaWriterFactoryGetMethod defines how the schema will be
	// written in GET operations.
	SchemaWriterFactoryGetMethod SchemaWriterFactoryFunc
	// SchemaWriterFactoryPostMethod defines how the schema will be
	// written in POST operations.
	SchemaWriterFactoryPostMethod SchemaWriterFactoryFunc
	// BaseURLFactory provides the base URL used for POSTs and that
	// will also be used to compose the ID URL in case the later
	// is not provided.
	BaseURLFactory    URLFactoryFunc
	IdBasedURLFactory URLFactoryFunc
}

func DefaultSchemaWriterFactory(d *schema.ResourceData) SchemaWriter {
	return &IDBasedResponse{}
}

func (dch DefaultContextHandler) defaultOperationHandler(
	operationType ot.OperationType,
	httpMethod string,
	schemaReaderFactory SchemaReaderFactoryFunc,
	schemaWriterFactory SchemaWriterFactoryFunc,
) ResourceOperationConfig {
	// POST = https://<CP>/<apiVersion>/<feature>
	// GET, PUT and DELETE = https://<CP>/<apiVersion>/<feature>/<id>
	endpoint := func(d *schema.ResourceData, c *client.Client) string {
		var url string
		if httpMethod == http.MethodPost {
			url = dch.BaseURLFactory(d, c)
		} else if dch.IdBasedURLFactory != nil {
			url = dch.IdBasedURLFactory(d, c)
		} else {
			url = fmt.Sprintf("%s/%s", dch.BaseURLFactory(d, c), d.Id())
		}
		tflog.Debug(context.Background(), fmt.Sprintf("Returning base URL for %s '%s' operation '%s' and httpMethod %s: %s",
			dch.ResourceType, dch.ResourceName, operationType, httpMethod, url))
		return url
	}

	var errorHandler RequestErrorHandler
	if httpMethod == http.MethodGet {
		errorHandler = &ReadIgnoreHttpNotFound{ResName: dch.ResourceName}
	} else if httpMethod == http.MethodDelete {
		errorHandler = &DeleteIgnoreHttpNotFound{ResName: dch.ResourceName}
	}
	result := ResourceOperationConfig{
		ResourceName:        dch.ResourceName,
		Type:                operationType,
		ResourceType:        dch.ResourceType,
		HttpMethod:          httpMethod,
		URLFactory:          endpoint,
		SchemaReaderFactory: schemaReaderFactory,
		SchemaWriterFactory: schemaWriterFactory,
		RequestErrorHandler: errorHandler,
	}

	return result
}

func (dch DefaultContextHandler) CreateContext() schema.CreateContextFunc {
	// By default, assumes that if no SchemaWriterFactoryPostMethod is provided,
	// the POST api will return an ID
	schemaWriterPost := DefaultSchemaWriterFactory
	if dch.SchemaWriterFactoryPostMethod != nil {
		schemaWriterPost = dch.SchemaWriterFactoryPostMethod
	}
	return CreateResource(
		dch.defaultOperationHandler(ot.Create, http.MethodPost, dch.SchemaReaderFactory, schemaWriterPost),
		dch.defaultOperationHandler(ot.Create, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod),
	)
}

func (dch DefaultContextHandler) ReadContext() schema.ReadContextFunc {
	return ReadResource(dch.ReadResourceOperationConfig())
}
func (dch DefaultContextHandler) ReadResourceOperationConfig() ResourceOperationConfig {
	return dch.defaultOperationHandler(ot.Read, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod)
}

func (dch DefaultContextHandler) UpdateContext() schema.UpdateContextFunc {
	return UpdateResource(
		dch.defaultOperationHandler(ot.Update, http.MethodPut, dch.SchemaReaderFactory, nil),
		dch.defaultOperationHandler(ot.Update, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod))
}

func (dch DefaultContextHandler) DeleteContext() schema.DeleteContextFunc {
	return DeleteResource(dch.defaultOperationHandler(ot.Delete, http.MethodDelete, nil, nil))
}
