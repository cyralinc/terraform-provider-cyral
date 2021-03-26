## 1.0.0 (March 26, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:
* **Resource Sidecar**: added new resource to support sidecars ([#23](github.com/cyralinc/terraform-provider-cyral/pull/23));
* **Resource Datamap**: added new resource to support datamaps ([#24](github.com/cyralinc/terraform-provider-cyral/pull/24));
* **Resource Repository Binding**: added new resource to support binding repositories to sidecars ([#25](github.com/cyralinc/terraform-provider-cyral/pull/25));
* **Resource Policy**: added new resource to support policies ([#26](github.com/cyralinc/terraform-provider-cyral/pull/26));
* **Resource Policy Rule**: added new resource to support policy rules ([#28](github.com/cyralinc/terraform-provider-cyral/pull/28));
* **Resource: Datadog Integration**: added new resource to support Datadog integration ([#30](github.com/cyralinc/terraform-provider-cyral/pull/30));
* **Resource ELK Integration**: added new resource to support ELK integration ([#31](github.com/cyralinc/terraform-provider-cyral/pull/31));
* **Resource Splunk Integration**: added new resource to support Splunk integration ([#33](github.com/cyralinc/terraform-provider-cyral/pull/33));
* **Resource Sumo Logic Integration**: added new resource to support Sumo Logic integration ([#34](github.com/cyralinc/terraform-provider-cyral/pull/34));
* **Resource Logstash Integration**: added new resource to support Logstash integration ([#35](github.com/cyralinc/terraform-provider-cyral/pull/35));
* **Resource Looker Integration**: added new resource to support Looker integration ([#36](github.com/cyralinc/terraform-provider-cyral/pull/36)).


### Improvements:
* **Replace repo ID property**: replaced repository identification in the state file from `name` to `id`. ([#22](github.com/cyralinc/terraform-provider-cyral/pull/22));

--------------

## 0.2.0 (February 24, 2021)

### Features:
* **Terraform Plugin SDK v2**: migration to `Terraform Plugin SDK v2` and support for Terraform `v0.12` to `v0.14` ([#21](github.com/cyralinc/terraform-provider-cyral/pull/21)).

### Improvements:
* **API Port**: Removed `control_plane_api_port` and added this information to the existing `control_plane` parameter. ([#20](github.com/cyralinc/terraform-provider-cyral/pull/20));
* **Unit tests**: added unit tests for `config.go` ([#8](github.com/cyralinc/terraform-provider-cyral/pull/8));

--------------

## 0.1.0 (May 15, 2020)

### Features:
* **Terraform import**: added support for `terraform import` statement for `repository` resource ([#8](github.com/cyralinc/terraform-provider-cyral/pull/8)).

--------------

## 0.0.2 (May 2, 2020)

### Bug fixes:
* **Error handling**: fixed bug in error handling ([#5](github.com/cyralinc/terraform-provider-cyral/pull/5)).

--------------

## 0.0.1 (May 2, 2020)

### Features:
* **Provider with Repository**: draft `provider` and `repository` resource ([#1](github.com/cyralinc/terraform-provider-cyral/pull/1), ([#2](github.com/cyralinc/terraform-provider-cyral/pull/2))).
