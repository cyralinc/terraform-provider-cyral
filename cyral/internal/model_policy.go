package internal

import (
	"time"
)

type Policy struct {
	Meta *PolicyMetadata `json:"meta"`
	Data []string        `json:"data,omitempty"`
	Tags []string        `json:"tags,omitempty"`
}

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
