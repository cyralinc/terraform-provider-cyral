package network

import "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"

const (
	resourceName = "cyral_repository_network_access_policy"

	defaultNetworkAccessPolicyEnabled    = true
	defaultNetworkAccessRulesBlockAccess = false
)

func repositoryTypesNetworkShield() []string {
	return []string{
		repository.SQLServer,
		repository.Oracle,
	}
}
