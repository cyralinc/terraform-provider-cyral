package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
)

// ValidateRepoType checks if a given repository type is valid.
func ValidateRepoType(param string) error {
	// This code was copied here to remove dependency of CRUD,
	// but we should move the CRUD code to CRUD-API (or somewhere
	// else) in the future.
	validValues := map[string]bool{
		"bigquery":   true,
		"cassandra":  true,
		"dremio":     true,
		"galera":     true,
		"mariadb":    true,
		"mongodb":    true,
		"mysql":      true,
		"oracle":     true,
		"postgresql": true,
		"redshift":   true,
		"snowflake":  true,
		"s3":         true,
		"sqlserver":  true,
	}
	if validValues[param] == false {
		return fmt.Errorf("repo type must be one of %v", param)
	}
	return nil
}

// ValidateDeploymentMethod checks if a given deployment parameter value is supported.
func ValidateDeploymentMethod(param string) error {
	validValues := map[string]bool{
		"docker":         true,
		"cloudFormation": true,
		"terraform":      true,
		"helm":           true,
		"helm3":          true,
		"automated":      true,
		"custom":         true,
		"terraformGKE":   true,
	}
	if validValues[param] == false {
		return fmt.Errorf("deployment method must be one of %v", validValues)
	}
	return nil
}

// ValidateRolePermissions check if role resource has only a single permissions block
func ValidateRolePermissions(permissions []interface{}) error {
	if len(permissions) > 1 {
		return fmt.Errorf("only a single permissions block is allowed")
	}

	return nil
}

// ValidateAWSRegion checks if a given aws region value is valid.
func ValidateAWSRegion(param string) error {
	resolver := endpoints.DefaultResolver()
	partitions := resolver.(endpoints.EnumPartitions).Partitions()

	validValues := map[string]bool{}

	for _, p := range partitions {
		// Only regions that support EC2 are valid for deployment
		if _, ok := p.Services()["ec2"]; ok {
			for id := range p.Regions() {
				validValues[id] = true
			}
		}
	}

	if validValues[param] == false {
		keys := make([]string, 0, len(validValues))
		for k := range validValues {
			keys = append(keys, k)
		}
		return fmt.Errorf("AWS region must be one of %v", keys)
	}
	return nil
}
