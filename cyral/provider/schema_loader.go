package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/awsiam"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/mfaduo"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/confextension/pagerduty"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml"
	idpsaml_draft "github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/idpsaml/draft"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/logging"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/teams"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessgateway"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessrules"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/binding"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confanalysis"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/confauth"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/network"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/useraccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlcertificate"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/credentials"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar/listener"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		awsiam.PackageSchema(),
		accessgateway.PackageSchema(),
		accessrules.PackageSchema(),
		binding.PackageSchema(),
		confanalysis.PackageSchema(),
		confauth.PackageSchema(),
		credentials.PackageSchema(),
		datalabel.PackageSchema(),
		datamap.PackageSchema(),
		hcvault.PackageSchema(),
		idpsaml.PackageSchema(),
		idpsaml_draft.PackageSchema(),
		listener.PackageSchema(),
		logging.PackageSchema(),
		mfaduo.PackageSchema(),
		network.PackageSchema(),
		pagerduty.PackageSchema(),
		repository.PackageSchema(),
		samlcertificate.PackageSchema(),
		sidecar.PackageSchema(),
		slack.PackageSchema(),
		teams.PackageSchema(),
		tokensettings.PackageSchema(),
		useraccount.PackageSchema(),
	}
	return v
}
