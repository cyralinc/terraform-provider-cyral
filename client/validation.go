package client

import "fmt"

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

// ValidateAWSRegion checks if a given aws region value is valid.
func ValidateAWSRegion(param string) error {
	validValues := map[string]bool{
		"us-east-2":      true,
		"us-east-1":      true,
		"us-west-1":      true,
		"us-west-2":      true,
		"af-south-1":     true,
		"ap-east-1":      true,
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-south-1":     true,
		"eu-west-3":      true,
		"eu-north-1":     true,
		"me-south-1":     true,
		"sa-east-1":      true,
	}
	if validValues[param] == false {
		return fmt.Errorf("AWS region must be one of %v", validValues)
	}
	return nil
}
