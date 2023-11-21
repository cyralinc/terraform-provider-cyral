package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/hcvault"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/slack"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/integration/teams"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/useraccount"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/samlcertificate"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		datalabel.PackageSchema(),
		datamap.PackageSchema(),
		hcvault.PackageSchema(),
		samlcertificate.PackageSchema(),
		slack.PackageSchema(),
		teams.PackageSchema(),
		useraccount.PackageSchema(),
	}
	return v
}
