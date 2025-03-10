package core

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
)

type IgnoreNotFoundByMessage struct {
	ResName        string
	MessageMatches string
	OperationType  operationtype.OperationType
}

func (h *IgnoreNotFoundByMessage) HandleError(
	ctx context.Context,
	r *schema.ResourceData,
	_ *client.Client,
	err error,
) error {
	tflog.Debug(ctx, "==> Init HandleError core.IgnoreNotFoundByMessage")

	// If this is a 404 already, then we don't need to actually match the error message

	if httpError, ok := err.(*client.HttpError); ok && httpError.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, "===> Ignoring regex matching as the status code is already 404")
		r.SetId("")
		return nil
	}

	matched, regexpError := regexp.MatchString(
		h.MessageMatches,
		err.Error(),
	)

	if regexpError != nil {
		return fmt.Errorf("regex failed to compile trying to match '%s' in '%w'. Error: %w",
			h.MessageMatches, err, regexpError)
	}

	if matched {
		tflog.Debug(ctx, fmt.Sprintf("===> %s not found. Skipping %s operation. Error: %v",
			h.ResName, h.OperationType, err))
		r.SetId("")
		tflog.Debug(ctx, "==> End HandleError core.IgnoreNotFoundByMessage - Success")
		return nil
	}

	tflog.Debug(ctx, "==> End HandleError core.IgnoreNotFoundByMessage - No match found, thus returning the original error")
	return err
}

type IgnoreHttpNotFound struct {
	ResName string
}

func (h *IgnoreHttpNotFound) HandleError(
	ctx context.Context,
	r *schema.ResourceData,
	_ *client.Client,
	err error,
) error {
	tflog.Debug(ctx, "==> Init HandleError core.IgnoreHttpNotFound")
	httpError, ok := err.(*client.HttpError)
	if !ok || httpError.StatusCode != http.StatusNotFound {
		tflog.Debug(ctx, "==> End HandleError core.IgnoreHttpNotFound - Did not find a 404, thus returning the original error")
		return err
	}
	r.SetId("")
	tflog.Debug(ctx, fmt.Sprintf(
		"==> End HandleError core.IgnoreHttpNotFound - %s not found. Marking resource for recreation or deletion.", h.ResName))
	return nil
}
