package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/tokensettings"
)

func packagesSchemas() []core.PackageSchema {
	v := []core.PackageSchema{
		datalabel.PackageSchema(),
		datamap.PackageSchema(),
		tokensettings.PackageSchema(),
	}
	return v
}
