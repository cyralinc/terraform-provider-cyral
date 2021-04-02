package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type URLCreatorFunc = func(d *schema.ResourceData, c *client.Client) string

type IModel interface {
	WriteResourceData(d *schema.ResourceData)
	ReadResourceData(d *schema.ResourceData)
}

type FunctionConfig struct {
	Name         string
	HttpMethod   string
	CreateURL    URLCreatorFunc
	ResourceData IModel
	ResponseData IModel
}

func CreateResource(createConfig FunctionConfig, readConfig FunctionConfig) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", createConfig.Name)
		c := m.(*client.Client)

		createConfig.ResourceData.ReadResourceData(d)

		url := createConfig.CreateURL(d, c)

		body, err := c.DoRequest(url, http.MethodPost, createConfig.ResourceData)
		if err != nil {
			return createError("Unable to create integration", fmt.Sprintf("%v", err))
		}

		if err := json.Unmarshal(body, &createConfig.ResponseData); err != nil {
			return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
		}
		log.Printf("[DEBUG] Response body (unmarshalled): %#v", createConfig.ResponseData)

		createConfig.ResponseData.WriteResourceData(d)

		log.Printf("[DEBUG] End %s", createConfig.Name)

		return ReadResource(readConfig)(ctx, d, m)
	}
}

func ReadResource(config FunctionConfig) schema.ReadContextFunc {
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

		config.ResponseData.WriteResourceData(d)

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}

func (config *FunctionConfig) Update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init %s", config.Name)
	c := m.(*client.Client)

	config.ResourceData.ReadResourceData(d)

	url := config.CreateURL(d, c)

	if _, err := c.DoRequest(url, config.HttpMethod, config.ResourceData); err != nil {
		return createError("Unable to update integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End %s", config.Name)

	readConfig := *config.ReadFunctionConfig

	return readConfig.Read(ctx, d, m)
}

func (config *FunctionConfig) Delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init %s", config.Name)
	c := m.(*client.Client)

	url := config.CreateURL(d, c)

	if _, err := c.DoRequest(url, config.HttpMethod, nil); err != nil {
		return createError("Unable to delete integration", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End %s", config.Name)

	return diag.Diagnostics{}
}
