package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	deprecated_policy "github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated/policy"
	deprecated_policy_rule "github.com/cyralinc/terraform-provider-cyral/cyral/internal/deprecated/policy/rule"
	integration_awsiam "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/awsiam"
	integration_mfa_duo "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/mfaduo"
	integration_pager_duty "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/pagerduty"
	integration_hcvault "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	integration_idp_saml "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml"
	integration_idp_saml_draft "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml/draft"
	integration_logging "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/logging"
	integration_slack "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	integration_teams "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/teams"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/permission"
	policyv2 "github.com/cyralinc/terraform-provider-cyral/cyral/internal/policy/v2"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/regopolicy"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	repository_accessgateway "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessgateway"
	repository_accessrules "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessrules"
	repository_binding "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/binding"
	repository_confanalysis "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confanalysis"
	repository_confauth "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confauth"
	repository_datamap "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
	repository_network "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/network"
	repository_useraccount "github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/useraccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/role"
	role_ssogroups "github.com/cyralinc/terraform-provider-cyral/cyral/internal/role/ssogroups"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlcertificate"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/serviceaccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	sidecar_credentials "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/credentials"
	sidecar_health "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/health"
	sidecar_instance "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/instance"
	sidecar_instance_stats "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/instance/stats"
	sidecar_listener "github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/listener"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/systeminfo"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		datalabel.PackageSchema(),
		deprecated_policy.PackageSchema(),
		deprecated_policy_rule.PackageSchema(),
		integration_awsiam.PackageSchema(),
		integration_hcvault.PackageSchema(),
		integration_idp_saml.PackageSchema(),
		integration_idp_saml_draft.PackageSchema(),
		integration_logging.PackageSchema(),
		integration_mfa_duo.PackageSchema(),
		integration_pager_duty.PackageSchema(),
		integration_slack.PackageSchema(),
		integration_teams.PackageSchema(),
		permission.PackageSchema(),
		policyv2.PackageSchema(),
		regopolicy.PackageSchema(),
		repository.PackageSchema(),
		repository_accessgateway.PackageSchema(),
		repository_accessrules.PackageSchema(),
		repository_binding.PackageSchema(),
		repository_confanalysis.PackageSchema(),
		repository_confauth.PackageSchema(),
		repository_datamap.PackageSchema(),
		repository_network.PackageSchema(),
		repository_useraccount.PackageSchema(),
		role.PackageSchema(),
		role_ssogroups.PackageSchema(),
		samlcertificate.PackageSchema(),
		serviceaccount.PackageSchema(),
		sidecar.PackageSchema(),
		sidecar_credentials.PackageSchema(),
		sidecar_health.PackageSchema(),
		sidecar_listener.PackageSchema(),
		sidecar_instance.PackageSchema(),
		sidecar_instance_stats.PackageSchema(),
		systeminfo.PackageSchema(),
		tokensettings.PackageSchema(),
	}
	return v
}
