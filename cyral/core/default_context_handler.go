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
//  3. `BaseURLFactory` must be provided for resources. It will be used to
//     create the POST endpoint and others in case `GetPutDeleteURLFactory`
//     is not provided.
//  4. `GetPutDeleteURLFactory` must be provided for data sources.
//  5. If `GetPutDeleteURLFactory` is NOT provided (data sources or resources),
//     the endpoint to perform GET, PUT and DELETE calls are composed by the
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
	// BaseURLFactory provides the URL used for POSTs and that
	// will also be used to compose the ID URL for GET, PUT and
	// DELETE in case `GetPutDeleteURLFactory` is not provided.
	BaseURLFactory         URLFactoryFunc
	GetPutDeleteURLFactory URLFactoryFunc
}

func DefaultSchemaWriterFactory(d *schema.ResourceData) SchemaWriter {
	return &IDBasedResponse{}
}

func (dch DefaultContextHandler) defaultOperationHandler(
	operationType ot.OperationType,
	httpMethod string,
	schemaReaderFactory SchemaReaderFactoryFunc,
	schemaWriterFactory SchemaWriterFactoryFunc,
	requestErrorHandler RequestErrorHandler,
) ResourceOperationConfig {
	// POST = https://<CP>/<apiVersion>/<feature>
	// GET, PUT and DELETE = https://<CP>/<apiVersion>/<feature>/<id>
	endpoint := func(d *schema.ResourceData, c *client.Client) string {
		var url string
		if httpMethod == http.MethodPost {
			url = dch.BaseURLFactory(d, c)
		} else if dch.GetPutDeleteURLFactory != nil {
			url = dch.GetPutDeleteURLFactory(d, c)
		} else {
			url = fmt.Sprintf("%s/%s", dch.BaseURLFactory(d, c), d.Id())
		}
		tflog.Debug(context.Background(), fmt.Sprintf("Returning base URL for %s '%s' operation '%s' and httpMethod %s: %s",
			dch.ResourceType, dch.ResourceName, operationType, httpMethod, url))
		return url
	}

	result := ResourceOperationConfig{
		ResourceName:        dch.ResourceName,
		Type:                operationType,
		ResourceType:        dch.ResourceType,
		HttpMethod:          httpMethod,
		URLFactory:          endpoint,
		SchemaReaderFactory: schemaReaderFactory,
		SchemaWriterFactory: schemaWriterFactory,
		RequestErrorHandler: requestErrorHandler,
	}

	return result
}

func (dch DefaultContextHandler) CreateContext() schema.CreateContextFunc {
	return dch.CreateContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName}, nil)
}

func (dch DefaultContextHandler) CreateContextCustomErrorHandling(getErrorHandler RequestErrorHandler,
	postErrorHandler RequestErrorHandler) schema.CreateContextFunc {
	// By default, assumes that if no SchemaWriterFactoryPostMethod is provided,
	// the POST api will return an ID
	schemaWriterPost := DefaultSchemaWriterFactory
	if dch.SchemaWriterFactoryPostMethod != nil {
		schemaWriterPost = dch.SchemaWriterFactoryPostMethod
	}
	return CreateResource(
		dch.defaultOperationHandler(ot.Create, http.MethodPost, dch.SchemaReaderFactory, schemaWriterPost, postErrorHandler),
		dch.defaultOperationHandler(ot.Create, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod, getErrorHandler),
	)
}

func (dch DefaultContextHandler) ReadContext() schema.ReadContextFunc {
	return dch.ReadContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName})
}

func (dch DefaultContextHandler) ReadContextCustomErrorHandling(getErrorHandler RequestErrorHandler) schema.ReadContextFunc {
	return ReadResource(
		dch.defaultOperationHandler(ot.Read, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod, getErrorHandler),
	)
}

func (dch DefaultContextHandler) UpdateContext() schema.UpdateContextFunc {
	return dch.UpdateContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName}, nil)
}

func (dch DefaultContextHandler) UpdateContextCustomErrorHandling(getErrorHandler RequestErrorHandler,
	putErrorHandler RequestErrorHandler) schema.UpdateContextFunc {
	return UpdateResource(
		dch.defaultOperationHandler(ot.Update, http.MethodPut, dch.SchemaReaderFactory, nil, putErrorHandler),
		dch.defaultOperationHandler(ot.Update, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod, getErrorHandler),
	)
}

func (dch DefaultContextHandler) DeleteContext() schema.DeleteContextFunc {
	return dch.DeleteContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName})
}

func (dch DefaultContextHandler) DeleteContextCustomErrorHandling(deleteErrorHandler RequestErrorHandler) schema.DeleteContextFunc {
	return DeleteResource(dch.defaultOperationHandler(ot.Delete, http.MethodDelete, nil, nil, deleteErrorHandler))
}
