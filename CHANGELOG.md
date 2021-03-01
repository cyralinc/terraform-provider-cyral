## 0.2.0 (February 24, 2021)

### Features:
* **Terraform Plugin SDK v2**: migration to `Terraform Plugin SDK v2` and support for Terraform `v0.12` to `v0.14` ([#21](github.com/cyralinc/terraform-provider-cyral/pull/21)).

### Improvements:
* **API Port**: Removed `control_plane_api_port` and added this information to the existing `control_plane` parameter. ([#20](github.com/cyralinc/terraform-provider-cyral/pull/20));
* **Unit tests**: added unit tests for `config.go` ([#8](github.com/cyralinc/terraform-provider-cyral/pull/8));

## 0.1.0 (May 15, 2020)

### Features:
* **Terraform import**: added support for `terraform import` statement for `repository` resource ([#8](github.com/cyralinc/terraform-provider-cyral/pull/8)).

## 0.0.2 (May 2, 2020)

### Bug fixes:
* **Error handling**: fixed bug in error handling ([#5](github.com/cyralinc/terraform-provider-cyral/pull/5)).

## 0.0.1 (May 2, 2020)

### Features:
* **Provider with Repository**: draft `provider` and `repository` resource ([#1](github.com/cyralinc/terraform-provider-cyral/pull/1), ([#2](github.com/cyralinc/terraform-provider-cyral/pull/2))).
