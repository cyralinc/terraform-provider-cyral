package policyv2

const (
	resourceName    = "cyral_policy_v2"
	dataSourceName  = resourceName
	apiPathLocal    = "v2/policies/local"
	apiPathGlobal   = "v2/policies/global"
	apiPathApproval = "v2/policies/approval"
)

func getAPIPath(policyType string) string {
	switch policyType {
	case "POLICY_TYPE_LOCAL", "local":
		return apiPathLocal
	case "POLICY_TYPE_GLOBAL", "global":
		return apiPathGlobal
	case "POLICY_TYPE_APPROVAL", "approval":
		return apiPathApproval
	default:
		return ""
	}
}
