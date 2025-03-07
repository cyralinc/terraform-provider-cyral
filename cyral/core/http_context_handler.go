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

// HTTPContextHandler facilitates easy resource and data source implementation
// for any resource that is accessed and modified using an HTTP/REST API.
//
//  1. `SchemaWriterFactoryGetMethod“ must be provided.
//  2. In case `SchemaWriterFactoryPostMethod“ is not provided,
//     it will assume that a call to POST returns a JSON with
//     an `id` field, meaning it will use the
//     `IDBasedResponse` struct in such cases.
//  3. `BaseURLFactory` must be provided for resources. It will be used to
//     create the POST endpoint and others in case `ReadUpdateDeleteURLFactory`
//     is not provided.
//  4. `ReadUpdateDeleteURLFactory` must be provided for data sources.
//  5. If `ReadUpdateDeleteURLFactory` is NOT provided (data sources or resources),
//     the endpoints to perform GET, PUT, PATCH and DELETE calls are composed by the
//     `BaseURLFactory` endpoint plus the ID specification as follows:
//     - POST:    https://<CP>/<apiVersion>/<featureName>
//     - GET:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - PUT:     https://<CP>/<apiVersion>/<featureName>/<id>
//     - PATCH:   https://<CP>/<apiVersion>/<featureName>/<id>
//     - DELETE:  https://<CP>/<apiVersion>/<featureName>/<id>
type HTTPContextHandler struct {
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
	// will also be used to compose the ID URL for GET, PUT/PATCH and
	// DELETE in case `ReadUpdateDeleteURLFactory` is not provided.
	BaseURLFactory             URLFactoryFunc
	ReadUpdateDeleteURLFactory URLFactoryFunc

	// Http method for update operations. If not provided, assumes http.MethodPut
	UpdateMethod string
}

func DefaultSchemaWriterFactory(d *schema.ResourceData) SchemaWriter {
	return &IDBasedResponse{}
}

func (dch HTTPContextHandler) defaultOperationHandler(
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
		} else if dch.ReadUpdateDeleteURLFactory != nil {
			url = dch.ReadUpdateDeleteURLFactory(d, c)
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

func (dch HTTPContextHandler) CreateContext() schema.CreateContextFunc {
	return dch.CreateContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName}, nil)
}

func (dch HTTPContextHandler) CreateContextCustomErrorHandling(getErrorHandler RequestErrorHandler,
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

func (dch HTTPContextHandler) ReadContext() schema.ReadContextFunc {
	return dch.ReadContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName})
}

func (dch HTTPContextHandler) ReadContextCustomErrorHandling(getErrorHandler RequestErrorHandler) schema.ReadContextFunc {
	return ReadResource(
		dch.defaultOperationHandler(ot.Read, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod, getErrorHandler),
	)
}

func (dch HTTPContextHandler) UpdateContext() schema.UpdateContextFunc {
	return dch.UpdateContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName}, nil)
}

func (dch HTTPContextHandler) UpdateContextCustomErrorHandling(getErrorHandler RequestErrorHandler,
	putErrorHandler RequestErrorHandler) schema.UpdateContextFunc {
	updateMethod := http.MethodPut
	if dch.UpdateMethod != "" {
		updateMethod = dch.UpdateMethod
	}
	return UpdateResource(
		dch.defaultOperationHandler(ot.Update, updateMethod, dch.SchemaReaderFactory, nil, putErrorHandler),
		dch.defaultOperationHandler(ot.Update, http.MethodGet, nil, dch.SchemaWriterFactoryGetMethod, getErrorHandler),
	)
}

func (dch HTTPContextHandler) DeleteContext() schema.DeleteContextFunc {
	return dch.DeleteContextCustomErrorHandling(&IgnoreHttpNotFound{ResName: dch.ResourceName})
}

func (dch HTTPContextHandler) DeleteContextCustomErrorHandling(deleteErrorHandler RequestErrorHandler) schema.DeleteContextFunc {
	return DeleteResource(dch.defaultOperationHandler(ot.Delete, http.MethodDelete, nil, nil, deleteErrorHandler))
}
