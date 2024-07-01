package policyv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ChangeInfo represents information about changes to the policy
type ChangeInfo struct {
	Actor     string `json:"actor,omitempty"`
	ActorType string `json:"actorType,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// PolicyV2 represents the top-level policy structure
type PolicyV2 struct {
	Policy Policy `json:"policy,omitempty"`
}

type Scope struct {
	RepoIds []string `json:"repoIds,omitempty"`
}

// Policy represents the policy details
type Policy struct {
	ID          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Enabled     bool       `json:"enabled,omitempty"`
	Scope       *Scope     `json:"scope,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	ValidFrom   string     `json:"validFrom,omitempty"`
	ValidUntil  string     `json:"validUntil,omitempty"`
	Document    string     `json:"document,omitempty"`
	LastUpdated ChangeInfo `json:"lastUpdated,omitempty"`
	Created     ChangeInfo `json:"created,omitempty"`
	Enforced    bool       `json:"enforced,omitempty"`
	Type        string     `json:"type,omitempty"`
}

// WriteToSchema writes the policy data to the schema
func (r PolicyV2) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("id", r.Policy.ID); err != nil {
		return fmt.Errorf("error setting 'id' field: %w", err)
	}
	if err := d.Set("name", r.Policy.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if err := d.Set("description", r.Policy.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("enabled", r.Policy.Enabled); err != nil {
		return fmt.Errorf("error setting 'enabled' field: %w", err)
	}
	if err := d.Set("tags", r.Policy.Tags); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}
	if err := d.Set("valid_from", r.Policy.ValidFrom); err != nil {
		return fmt.Errorf("error setting 'valid_from' field: %w", err)
	}
	if err := d.Set("valid_until", r.Policy.ValidUntil); err != nil {
		return fmt.Errorf("error setting 'valid_until' field: %w", err)
	}

	if err := d.Set("document", r.Policy.Document); err != nil {
		return fmt.Errorf("error setting 'document' field: %w", err)
	}

	if err := d.Set("last_updated", map[string]interface{}{
		"actor":      r.Policy.LastUpdated.Actor,
		"actor_type": r.Policy.LastUpdated.ActorType,
		"timestamp":  r.Policy.LastUpdated.Timestamp,
	}); err != nil {
		return fmt.Errorf("error setting 'last_updated' field: %w", err)
	}
	if err := d.Set("created", map[string]interface{}{
		"actor":      r.Policy.Created.Actor,
		"actor_type": r.Policy.Created.ActorType,
		"timestamp":  r.Policy.Created.Timestamp,
	}); err != nil {
		return fmt.Errorf("error setting 'created' field: %w", err)
	}
	if err := d.Set("enforced", r.Policy.Enforced); err != nil {
		return fmt.Errorf("error setting 'enforced' field: %w", err)
	}
	if r.Policy.Type != "" {
		if err := d.Set("type", r.Policy.Type); err != nil {
			return fmt.Errorf("error setting 'type' field: %w", err)
		}
	}
	if err := d.Set("scope", flattenScope(r.Policy.Scope)); err != nil {
		return fmt.Errorf("error setting 'scope' field: %w", err)
	}
	d.SetId(r.Policy.ID)
	return nil
}

// ReadFromSchema reads the policy data from the schema
func (r *PolicyV2) ReadFromSchema(d *schema.ResourceData) error {
	r.Policy.ID = d.Get("id").(string)
	r.Policy.Name = d.Get("name").(string)
	r.Policy.Description = d.Get("description").(string)
	r.Policy.Enabled = d.Get("enabled").(bool)
	r.Policy.Tags = expandStringList(d.Get("tags").([]interface{}))
	r.Policy.ValidFrom = d.Get("valid_from").(string)
	r.Policy.ValidUntil = d.Get("valid_until").(string)
	r.Policy.Document = d.Get("document").(string)
	r.Policy.Enforced = d.Get("enforced").(bool)
	r.Policy.Type = d.Get("type").(string)
	if v, ok := d.GetOk("scope"); ok {
		r.Policy.Scope = scopeFromInterface(v.([]interface{}))
	}
	return nil
}

// expandStringList converts a list of interface{} to a list of strings
func expandStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

// flattenScope converts the Scope struct to a list of maps
func flattenScope(scope *Scope) []map[string]interface{} {
	if scope == nil {
		return nil
	}
	scopeMap := []map[string]interface{}{
		{
			"repo_ids": scope.RepoIds,
		},
	}
	return scopeMap
}

// scopeFromInterface converts the map to a Scope struct
func scopeFromInterface(s []interface{}) *Scope {
	if len(s) == 0 || s[0] == nil {
		return nil
	}
	m := s[0].(map[string]interface{})
	scope := Scope{
		RepoIds: expandStringList(m["repo_ids"].([]interface{})),
	}
	return &scope
}
