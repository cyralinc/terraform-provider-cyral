package cyral

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

func defaultRequestErrorHandler(
	operationType string,
	logInfo string,
) RequestErrorHandler {
	switch operationType {
	case read:
		return &defaultReadErrorHandler{logInfo: logInfo}
	default:
		return &dummyRequestErrorHandler{}
	}
}

type defaultReadErrorHandler struct {
	logInfo string
}

func (h *defaultReadErrorHandler) HandleError(
	d *schema.ResourceData,
	c *client.Client,
	err error,
) error {
	httpError, ok := err.(*client.HttpError)
	if ok && httpError.Is404Status() {
		log.Printf("[DEBUG] Resource with ID %s not found. Resetting its ID. "+
			"Extra information: %s", d.Id(), h.logInfo)
		d.SetId("")
		return nil
	}
	return err
}

type dummyRequestErrorHandler struct{}

func (h *dummyRequestErrorHandler) HandleError(
	d *schema.ResourceData,
	c *client.Client,
	err error,
) error {
	return err
}

type DeleteIgnoreHttpNotFound struct {
	resName string
}

func (h *DeleteIgnoreHttpNotFound) HandleError(
	_ *schema.ResourceData,
	_ *client.Client,
	err error,
) error {
	httpError, ok := err.(*client.HttpError)
	if ok && httpError.Is404Status() {
		log.Printf("[DEBUG] %s not found. Skipping deletion.", h.resName)
		return nil
	}
	return err
}
