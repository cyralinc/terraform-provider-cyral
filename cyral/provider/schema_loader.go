package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/teams"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessgateway"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/accessrules"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/useraccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlcertificate"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		accessgateway.PackageSchema(),
		accessrules.PackageSchema(),
		datalabel.PackageSchema(),
		datamap.PackageSchema(),
		hcvault.PackageSchema(),
		repository.PackageSchema(),
		samlcertificate.PackageSchema(),
		slack.PackageSchema(),
		teams.PackageSchema(),
		tokensettings.PackageSchema(),
		useraccount.PackageSchema(),
	}
	return v
}
