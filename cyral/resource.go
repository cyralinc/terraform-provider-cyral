package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type URLCreatorFunc = func(d *schema.ResourceData, c *client.Client) string

type SchemaDataHandler interface {
	WriteToSchema(d *schema.ResourceData)
	ReadFromSchema(d *schema.ResourceData)
}

type ResourceOperationConfig struct {
	Name         string
	HttpMethod   string
	CreateURL    URLCreatorFunc
	ResourceData SchemaDataHandler
	ResponseData SchemaDataHandler
}

type IDBasedResponse struct {
	ID string `json:"id"`
}

func (response IDBasedResponse) WriteToSchema(d *schema.ResourceData) {
	d.SetId(response.ID)
}

func (response *IDBasedResponse) ReadFromSchema(d *schema.ResourceData) {
	response.ID = d.Id()
}

func CreateResource(createConfig ResourceOperationConfig, readConfig ResourceOperationConfig) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", createConfig.Name)
		c := m.(*client.Client)

		createConfig.ResourceData.ReadFromSchema(d)

		url := createConfig.CreateURL(d, c)

		body, err := c.DoRequest(url, createConfig.HttpMethod, createConfig.ResourceData)
		if err != nil {
			return createError("Unable to create integration", fmt.Sprintf("%v", err))
		}

		if err := json.Unmarshal(body, &createConfig.ResponseData); err != nil {
			return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
		}
		log.Printf("[DEBUG] Response body (unmarshalled): %#v", createConfig.ResponseData)

		createConfig.ResponseData.WriteToSchema(d)

		log.Printf("[DEBUG] End %s", createConfig.Name)

		return ReadResource(readConfig)(ctx, d, m)
	}
}

func ReadResource(config ResourceOperationConfig) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		url := config.CreateURL(d, c)

		body, err := c.DoRequest(url, config.HttpMethod, nil)
		if err != nil {
			return createError(fmt.Sprintf("Unable to read integration. IntegrationID: %s",
				d.Id()), fmt.Sprintf("%v", err))
		}

		if err := json.Unmarshal(body, &config.ResponseData); err != nil {
			return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
		}
		log.Printf("[DEBUG] Response body (unmarshalled): %#v", config.ResponseData)

		config.ResponseData.WriteToSchema(d)

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}

func UpdateResource(updateConfig ResourceOperationConfig, readConfig ResourceOperationConfig) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", updateConfig.Name)
		c := m.(*client.Client)

		updateConfig.ResourceData.ReadFromSchema(d)

		url := updateConfig.CreateURL(d, c)

		if _, err := c.DoRequest(url, updateConfig.HttpMethod, updateConfig.ResourceData); err != nil {
			return createError("Unable to update integration", fmt.Sprintf("%v", err))
		}

		log.Printf("[DEBUG] End %s", updateConfig.Name)

		return ReadResource(readConfig)(ctx, d, m)
	}
}

func DeleteResource(config ResourceOperationConfig) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		url := config.CreateURL(d, c)

		if _, err := c.DoRequest(url, config.HttpMethod, nil); err != nil {
			return createError("Unable to delete integration", fmt.Sprintf("%v", err))
		}

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}
