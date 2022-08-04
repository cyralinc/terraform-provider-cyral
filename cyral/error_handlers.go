package cyral

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type DeleteIgnoreHttpNotFound struct {
	resName string
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
	log.Printf("[DEBUG] %s not found. Skipping deletion.", h.resName)
	return nil
}
