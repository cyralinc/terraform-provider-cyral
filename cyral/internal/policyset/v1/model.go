package policysetv1

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

// ChangeInfo represents information about changes to the policy set
type ChangeInfo struct {
	Actor     string `json:"actor,omitempty"`
	ActorType string `json:"actorType,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// ToMap converts ChangeInfo to a map
func (c ChangeInfo) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"actor":      c.Actor,
		"actor_type": c.ActorType,
		"timestamp":  c.Timestamp,
	}
}

// PolicySetPolicy represents a policy in the policy set
type PolicySetPolicy struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ToMap converts PolicySetPolicy to a map
func (p PolicySetPolicy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type": p.Type,
		"id":   p.ID,
	}
}

// Scope represents the scope of the policy set
type Scope struct {
	RepoIds []string `json:"repoIds,omitempty"`
}

// ToMap converts Scope to a list of maps
func (s *Scope) ToMap() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"repo_ids": s.RepoIds,
		},
	}
}

// PolicySet represents the policy set details
type PolicySet struct {
	ID               string            `json:"id,omitempty"`
	WizardID         string            `json:"wizardId,omitempty"`
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Scope            *Scope            `json:"scope,omitempty"`
	WizardParameters string            `json:"wizardParameters,omitempty"`
	Enabled          bool              `json:"enabled,omitempty"`
	Policies         []PolicySetPolicy `json:"policies,omitempty"`
	LastUpdated      ChangeInfo        `json:"lastUpdated,omitempty"`
	Created          ChangeInfo        `json:"created,omitempty"`
}

// WriteToSchema writes the policy set data to the schema
func (r *PolicySet) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("id", r.ID); err != nil {
		return fmt.Errorf("error setting 'id' field: %w", err)
	}
	if err := d.Set("wizard_id", r.WizardID); err != nil {
		return fmt.Errorf("error setting 'wizard_id' field: %w", err)
	}
	if err := d.Set("name", r.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if err := d.Set("description", r.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}
	if err := d.Set("tags", r.Tags); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}
	if err := d.Set("wizard_parameters", r.WizardParameters); err != nil {
		return fmt.Errorf("error setting 'wizard_parameters' field: %w", err)
	}
	if err := d.Set("enabled", r.Enabled); err != nil {
		return fmt.Errorf("error setting 'enabled' field: %w", err)
	}
	if err := d.Set("policies", policiesToMaps(r.Policies)); err != nil {
		return fmt.Errorf("error setting 'policies' field: %w", err)
	}
	if err := d.Set("last_updated", r.LastUpdated.ToMap()); err != nil {
		return fmt.Errorf("error setting 'last_updated' field: %w", err)
	}
	if err := d.Set("created", r.Created.ToMap()); err != nil {
		return fmt.Errorf("error setting 'created' field: %w", err)
	}
	if r.Scope != nil {
		if err := d.Set("scope", r.Scope.ToMap()); err != nil {
			return fmt.Errorf("error setting 'scope' field: %w", err)
		}
	}
	d.SetId(r.ID)
	return nil
}

func policiesToMaps(policies []PolicySetPolicy) []map[string]interface{} {
	var result []map[string]interface{}
	for _, policy := range policies {
		result = append(result, policy.ToMap())
	}
	return result
}

// ReadFromSchema reads the policy set data from the schema
func (r *PolicySet) ReadFromSchema(d *schema.ResourceData) error {
	r.ID = d.Get("id").(string)
	r.WizardID = d.Get("wizard_id").(string)
	r.Name = d.Get("name").(string)
	r.Description = d.Get("description").(string)
	r.Tags = utils.ConvertFromInterfaceList[string](d.Get("tags").([]interface{}))
	r.WizardParameters = d.Get("wizard_parameters").(string)
	r.Enabled = d.Get("enabled").(bool)
	if v, ok := d.GetOk("scope"); ok {
		r.Scope = scopeFromInterface(v.([]interface{}))
	}
	return nil
}

// scopeFromInterface converts the map to a Scope struct
func scopeFromInterface(s []interface{}) *Scope {
	if len(s) == 0 || s[0] == nil {
		// return an empty scope (ie a scope with a repo ids array of length 0)
		return &Scope{
			RepoIds: []string{},
		}
	}
	m := s[0].(map[string]interface{})
	scope := Scope{
		RepoIds: utils.ConvertFromInterfaceList[string](m["repo_ids"].([]interface{})),
	}
	return &scope
}
