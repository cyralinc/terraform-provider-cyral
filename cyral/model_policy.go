package cyral

import (
	"time"
)

type Policy struct {
	Meta *PolicyMetadata `json:"meta" yaml:"meta"`
	Data []string        `json:"data,omitempty" yaml:"data,omitempty,flow"`
}

type PolicyMetadata struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Version     string    `json:"version" yaml:"version"`
	Created     time.Time `json:"created" yaml:"created"`
	LastUpdated time.Time `json:"lastUpdated" yaml:"lastUpdated"`
	Type        string    `json:"type" yaml:"type"`
	Tags        []string  `json:"tags" yaml:"tags"`
	Enabled     bool      `json:"enabled" yaml:"enabled"`
	Description string    `json:"description" yaml:"description"`
}

type PolicyListResponse struct {
	// Policies is a list of policy identifiers.
	Policies []string `json:"Policies,omitempty"`
}
