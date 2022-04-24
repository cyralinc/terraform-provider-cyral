package client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

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

func ValidateRepositoryConfAnalysisRedact() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"all",
		"none",
		"watched",
	}, false)
}

func ValidateRepositoryConfAnalysisCommentAnnotationGroups() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"identity",
		"client",
		"repo",
		"sidecar",
	}, false)
}

func ValidateRepositoryConfAnalysisLogSettings() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"everything",
		"dql",
		"dml",
		"ddl",
		"sensitive & dql",
		"sensitive & dml",
		"sensitive & ddl",
		"privileged",
		"port-scan",
		"auth-failure",
		"full-table-scan",
		"violations",
		"connections",
		"sensitive",
		"data-classification",
		"audit",
		"error",
		"new-connections",
		"closed-connections",
	}, false)
}
