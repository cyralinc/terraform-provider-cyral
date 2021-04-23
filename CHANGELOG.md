## 1.1.0 (April 23, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:
* **Resource Integration Slack Alerts**: added new resource to support Slack Alerts integration ([#43](https://github.com/cyralinc/terraform-provider-cyral/pull/43));
* **Resource Integration Microsoft Teams**: added new resource to support Microsoft Teams integration ([#44](https://github.com/cyralinc/terraform-provider-cyral/pull/44));
* **Resource Repository Local Account**: added new resource to support repositories local accounts ([#46](https://github.com/cyralinc/terraform-provider-cyral/pull/46));
* **Resource Integration Okta**: added new resource to support Okta integration ([#47](https://github.com/cyralinc/terraform-provider-cyral/pull/47));
* **Resource Repository Configuration Authentication**: added new resource to support repository configuration authentication ([#48](https://github.com/cyralinc/terraform-provider-cyral/pull/48));
* **Data Source Sidecar Template**: added new data source to support sidecar templates ([#50](https://github.com/cyralinc/terraform-provider-cyral/pull/50));
* **Resource Identity Map**: added new resource to support identity maps ([#51](https://github.com/cyralinc/terraform-provider-cyral/pull/51));

### Improvements:
* **Increase Test Coverage**: defined standards and increased the test coverage. ([#41](https://github.com/cyralinc/terraform-provider-cyral/pull/41));

--------------

## 1.0.0 (March 26, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:
* **Resource Sidecar**: added new resource to support sidecars ([#23](https://github.com/cyralinc/terraform-provider-cyral/pull/23));
* **Resource Datamap**: added new resource to support datamaps ([#24](https://github.com/cyralinc/terraform-provider-cyral/pull/24));
* **Resource Repository Binding**: added new resource to support binding repositories to sidecars ([#25](https://github.com/cyralinc/terraform-provider-cyral/pull/25));
* **Resource Policy**: added new resource to support policies ([#26](https://github.com/cyralinc/terraform-provider-cyral/pull/26));
* **Resource Policy Rule**: added new resource to support policy rules ([#28](https://github.com/cyralinc/terraform-provider-cyral/pull/28));
* **Resource: Datadog Integration**: added new resource to support Datadog integration ([#30](https://github.com/cyralinc/terraform-provider-cyral/pull/30));
* **Resource ELK Integration**: added new resource to support ELK integration ([#31](https://github.com/cyralinc/terraform-provider-cyral/pull/31));
* **Resource Splunk Integration**: added new resource to support Splunk integration ([#33](https://github.com/cyralinc/terraform-provider-cyral/pull/33));
* **Resource Sumo Logic Integration**: added new resource to support Sumo Logic integration ([#34](https://github.com/cyralinc/terraform-provider-cyral/pull/34));
* **Resource Logstash Integration**: added new resource to support Logstash integration ([#35](https://github.com/cyralinc/terraform-provider-cyral/pull/35));
* **Resource Looker Integration**: added new resource to support Looker integration ([#36](https://github.com/cyralinc/terraform-provider-cyral/pull/36)).


### Improvements:
* **Replace repo ID property**: replaced repository identification in the state file from `name` to `id`. ([#22](https://github.com/cyralinc/terraform-provider-cyral/pull/22));

--------------

## 0.2.0 (February 24, 2021)

### Features:
* **Terraform Plugin SDK v2**: migration to `Terraform Plugin SDK v2` and support for Terraform `v0.12` to `v0.14` ([#21](https://github.com/cyralinc/terraform-provider-cyral/pull/21)).

### Improvements:
* **API Port**: Removed `control_plane_api_port` and added this information to the existing `control_plane` parameter. ([#20](https://github.com/cyralinc/terraform-provider-cyral/pull/20));
* **Unit tests**: added unit tests for `config.go` ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8));

--------------

## 0.1.0 (May 15, 2020)

### Features:
* **Terraform import**: added support for `terraform import` statement for `repository` resource ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8)).

--------------

## 0.0.2 (May 2, 2020)

### Bug fixes:
* **Error handling**: fixed bug in error handling ([#5](https://github.com/cyralinc/terraform-provider-cyral/pull/5)).

--------------

## 0.0.1 (May 2, 2020)

### Features:
* **Provider with Repository**: draft `provider` and `repository` resource ([#1](https://github.com/cyralinc/terraform-provider-cyral/pull/1), ([#2](https://github.com/cyralinc/terraform-provider-cyral/pull/2))).
