package policy

import (
	"fmt"
	"time"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolicyMetadata struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
	Type        string    `json:"type"`
	Tags        []string  `json:"tags"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description"`
}

type PolicyListResponse struct {
	// Policies is a list of policy identifiers.
	Policies []string `json:"Policies,omitempty"`
}

type Policy struct {
	Meta *PolicyMetadata `json:"meta"`
	Data []string        `json:"data,omitempty"`
	Tags []string        `json:"tags,omitempty"`
}

func (r Policy) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("created", r.Meta.Created.String()); err != nil {
		return fmt.Errorf("error setting 'created' field: %w", err)
	}
	if err := d.Set("data", r.Data); err != nil {
		return fmt.Errorf("error setting 'data' field: %w", err)
	}
	if err := d.Set("data_label_tags", r.Tags); err != nil {
		return fmt.Errorf("error setting 'data_label_tags' field: %w", err)
	}
	if err := d.Set("description", r.Meta.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("enabled", r.Meta.Enabled); err != nil {
		return fmt.Errorf("error setting 'enabled' field: %w", err)
	}
	if err := d.Set("last_updated", r.Meta.LastUpdated.String()); err != nil {
		return fmt.Errorf("error setting 'last_updated' field: %w", err)
	}
	if err := d.Set("name", r.Meta.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if err := d.Set("type", r.Meta.Type); err != nil {
		return fmt.Errorf("error setting 'type' field: %w", err)
	}
	if err := d.Set("version", r.Meta.Version); err != nil {
		return fmt.Errorf("error setting 'version' field: %w", err)
	}
	// Once the `tags` field is removed, this conditional logic should also be
	// removed and only the `metadata_tags` should be set.
	_, isDeprecatedFieldSet := d.GetOk("tags")
	if isDeprecatedFieldSet {
		if err := d.Set("tags", r.Meta.Tags); err != nil {
			return fmt.Errorf("error setting 'tags' field: %w", err)
		}
	} else {
		if err := d.Set("metadata_tags", r.Meta.Tags); err != nil {
			return fmt.Errorf("error setting 'metadata_tags' field: %w", err)
		}

	}

	return nil
}

func (r *Policy) ReadFromSchema(d *schema.ResourceData) error {
	r.Data = utils.GetStrListFromSchemaField(d, "data")
	r.Tags = utils.GetStrListFromSchemaField(d, "data_label_tags")
	metadataTags := utils.GetStrListFromSchemaField(d, "metadata_tags")
	if len(metadataTags) == 0 {
		metadataTags = utils.GetStrListFromSchemaField(d, "tags")
	}
	r.Meta = &PolicyMetadata{
		Tags: metadataTags,
	}

	if v, ok := d.Get("name").(string); ok {
		r.Meta.Name = v
	}

	if v, ok := d.Get("version").(string); ok {
		r.Meta.Version = v
	}

	if v, ok := d.Get("type").(string); ok {
		r.Meta.Type = v
	}

	if v, ok := d.Get("enabled").(bool); ok {
		r.Meta.Enabled = v
	}

	if v, ok := d.Get("description").(string); ok {
		r.Meta.Description = v
	}

	return nil
}
