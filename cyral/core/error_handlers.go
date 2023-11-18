package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
)

type DeleteIgnoreHttpNotFound struct {
	ResName string
}

func (h *DeleteIgnoreHttpNotFound) HandleError(
	_ *schema.ResourceData,
	_ *client.Client,
	err error,
) error {
	httpError, ok := err.(*client.HttpError)
	if !ok || httpError.StatusCode != http.StatusNotFound {
		return err
	}
	tflog.Debug(context.Background(), fmt.Sprintf("%s not found. Skipping deletion.", h.ResName))
	return nil
}

type ReadIgnoreHttpNotFound struct {
	ResName string
}

func (h *ReadIgnoreHttpNotFound) HandleError(
	r *schema.ResourceData,
	_ *client.Client,
	err error,
) error {
	httpError, ok := err.(*client.HttpError)
	if !ok || httpError.StatusCode != http.StatusNotFound {
		return err
	}
	r.SetId("")
	tflog.Debug(context.Background(), fmt.Sprintf("%s not found. Marking resource for recreation.", h.ResName))
	return nil
}
