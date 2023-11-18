package provider

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository/datamap"
)

func packagesSchemas() []core.PackageSchema[any] {
	v := []core.PackageSchema[any]{
		datalabel.PackageSchema(),
		datamap.PackageSchema(),
	}
	return v
}
