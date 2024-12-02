package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

type ResourceMethod func(context.Context, *client.Client, *schema.ResourceData) error

// GenericContextHandler can be used by resource implementations to ensure that
// the recommended best practices are consistently followed (e.g., handling of
// 404 errors, following a create/update with a get etc). The resource implementation
// needs to supply functions that implement the basic CRUD operations on the resource
// using gRPC or whatever else. Note that if REST APIs are used, it is recommended
// to use the DefaultContextHandler instead.
type GenericContextHandler struct {
	ResourceName string
	ResourceType resourcetype.ResourceType
	Create       ResourceMethod
	Read         ResourceMethod
	Update       ResourceMethod
	Delete       ResourceMethod
}

type method struct {
	method       ResourceMethod
	name         string
	errorHandler func(context.Context, *schema.ResourceData, error) error
}

func (gch *GenericContextHandler) handleResourceNotFoundError(
	ctx context.Context, rd *schema.ResourceData, err error,
) error {
	var isNotFoundError bool
	if status.Code(err) == codes.NotFound {
		isNotFoundError = true
	} else if httpError, ok := err.(*client.HttpError); ok &&
		httpError.StatusCode == http.StatusNotFound {
		isNotFoundError = true
	}
	if isNotFoundError {
		tflog.Debug(
			ctx,
			fmt.Sprintf(
				"==> Resource %s not found, marking for recreation or deletion.",
				gch.ResourceName,
			),
		)
		rd.SetId("")
		return nil
	}
	return err
}

// CreateContext is used to create a resource instance.
func (gch *GenericContextHandler) CreateContext(
	ctx context.Context,
	rd *schema.ResourceData,
	pd any,
) diag.Diagnostics {
	c := pd.(*client.Client)
	return gch.executeMethods(
		ctx, c, rd, []method{
			{
				method: gch.Create,
				name:   "create",
			},
			{
				method:       gch.Read,
				name:         "read",
				errorHandler: gch.handleResourceNotFoundError,
			},
		},
	)
}

// UpdateContext is used to update a resource instance.
func (gch *GenericContextHandler) UpdateContext(
	ctx context.Context,
	rd *schema.ResourceData,
	pd any,
) diag.Diagnostics {
	c := pd.(*client.Client)
	return gch.executeMethods(
		ctx, c, rd, []method{
			{
				method: gch.Update,
				name:   "update",
			},
			{
				method:       gch.Read,
				name:         "read",
				errorHandler: gch.handleResourceNotFoundError,
			},
		},
	)
}

// ReadContext is used to read a resource instance.
func (gch *GenericContextHandler) ReadContext(
	ctx context.Context,
	rd *schema.ResourceData,
	pd any,
) diag.Diagnostics {
	c := pd.(*client.Client)
	return gch.executeMethods(
		ctx, c, rd, []method{
			{
				method:       gch.Read,
				name:         "read",
				errorHandler: gch.handleResourceNotFoundError,
			},
		},
	)
}

// DeleteContext is used to delete a resource instance.
func (gch *GenericContextHandler) DeleteContext(
	ctx context.Context,
	rd *schema.ResourceData,
	pd any,
) diag.Diagnostics {
	c := pd.(*client.Client)
	return gch.executeMethods(
		ctx, c, rd, []method{
			{
				method:       gch.Delete,
				name:         "delete",
				errorHandler: gch.handleResourceNotFoundError,
			},
		},
	)
}

func (gch *GenericContextHandler) executeMethods(
	ctx context.Context, c *client.Client, rd *schema.ResourceData, methods []method,
) diag.Diagnostics {
	for _, m := range methods {
		tflog.Debug(ctx, fmt.Sprintf("resource %s: operation %s", gch.ResourceName, m.name))
		err := m.method(ctx, c, rd)
		if err != nil {
			tflog.Debug(
				ctx,
				fmt.Sprintf("resource %s: operation %s - error: %v", gch.ResourceName, m.name, err),
			)
			if m.errorHandler != nil {
				err = m.errorHandler(ctx, rd, err)
			}
		}
		if err != nil {
			return utils.CreateError(
				fmt.Sprintf("error in operation %s on resource %s", m.name, gch.ResourceName),
				err.Error(),
			)
		}
		tflog.Debug(
			ctx,
			fmt.Sprintf("resource %s: operation %s - success", gch.ResourceName, m.name),
		)
	}
	return nil
}
