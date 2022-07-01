## 2.7.0 (July 1, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: parameter `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- **Fix memory referencing issue and refactor resources**: [#223](https://github.com/cyralinc/terraform-provider-cyral/pull/223).

### Features:

- **Add resources for per-repository data map**: [#230](https://github.com/cyralinc/terraform-provider-cyral/pull/230).
- **Add data source for cyral_integration_idp**: [#239](https://github.com/cyralinc/terraform-provider-cyral/pull/239).
- **Support MongoDB repository replica sets**: [#228](https://github.com/cyralinc/terraform-provider-cyral/pull/228).

## 2.6.2 (June 6, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: parameter `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- **Fix repository name for IdP Okta module**: [#220](https://github.com/cyralinc/terraform-provider-cyral/pull/220);

## 2.6.1 (May 31, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: parameter `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- **Fix broken links in documentation**: [#216](https://github.com/cyralinc/terraform-provider-cyral/pull/216);

## 2.6.0 (May 27, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: parameter `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- **Improve provider error message when auth fails**: [#210](https://github.com/cyralinc/terraform-provider-cyral/pull/210);
- **Fix issue with sidecar resource when configured with empty labels**: [#210](https://github.com/cyralinc/terraform-provider-cyral/pull/210);
- **Fix bug where cyral_sidecar_cft_template data source was crashing when configured with Splunk log integrations**: [#215](https://github.com/cyralinc/terraform-provider-cyral/pull/215);

### Features:

- **Add Kubernetes Secret support to cyral_repository_local_account resource**: [#193](https://github.com/cyralinc/terraform-provider-cyral/pull/193);
- **Add GCP Secret Manager support to cyral_repository_local_account resource**: [#206](https://github.com/cyralinc/terraform-provider-cyral/pull/206);
- **Add labels argument to cyral_repository resource**: [#192](https://github.com/cyralinc/terraform-provider-cyral/pull/192);
- **Add rate_limit argument to cyral_policy_rule resource**: [#194](https://github.com/cyralinc/terraform-provider-cyral/pull/194);
- **Add provider optional argument tls_skip_verify to configure TLS verification**: [#204](https://github.com/cyralinc/terraform-provider-cyral/pull/204);
- **Add certificate_bundle_secrets argument to sidecar resource to support sidecar certificate bundle secrets**: [#190](https://github.com/cyralinc/terraform-provider-cyral/pull/190);
- **Add cyral_sidecar_bound_ports data source to retrieve sidecar ports that are bound to repositories**: [#209](https://github.com/cyralinc/terraform-provider-cyral/pull/209);
- **Add cyral_sidecar_id data source to return sidecar ID given a sidecar name**: [#211](https://github.com/cyralinc/terraform-provider-cyral/pull/211);
- **Add cyral_sidecar_instance_ids data source to return sidecar instance IDs**: [#212](https://github.com/cyralinc/terraform-provider-cyral/pull/212);

### Deprecate:

- **Deprecate enviroment_variable argument of cyral_repository_local_account resource**: [#193](https://github.com/cyralinc/terraform-provider-cyral/pull/193);

### Documentation:

- **Update policy docs to inform how to use 'any' for rows argument**: [#201](https://github.com/cyralinc/terraform-provider-cyral/pull/201);
- **Update policy docs to inform how to use 'any' for data argument**: [#203](https://github.com/cyralinc/terraform-provider-cyral/pull/203);
- **Update docs to be automatically generated**: [#183](https://github.com/cyralinc/terraform-provider-cyral/pull/183);

## 2.5.2 (April 19, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with Control Planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- **Datamap resource is always suggesting an update on terraform plan**: [#187](https://github.com/cyralinc/terraform-provider-cyral/pull/187);

## 2.5.1 (March 31, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with Control Planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- **Update hadolint, gorelease and go version**: [#181](https://github.com/cyralinc/terraform-provider-cyral/pull/181);

## 2.5.0 (March 31, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with Control Planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- **Changed listener_port to port so the demo works**: [#170](https://github.com/cyralinc/terraform-provider-cyral/pull/170);
- **Fix PagerDuty resource due to confExtensions breaking changes**: [#173](https://github.com/cyralinc/terraform-provider-cyral/pull/173);

### Features:

- **Add Access Gateway support for repo binding resource**: [#180](https://github.com/cyralinc/terraform-provider-cyral/pull/180);

### Deprecate:

- **Deprecate properties field from policy resource**: [#178](https://github.com/cyralinc/terraform-provider-cyral/pull/178);

## 2.4.4 (March 3, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- **Added denodo to the validValues map**: [#169](https://github.com/cyralinc/terraform-provider-cyral/pull/169);

### Documentation:

- **Change cyral_repository to cyral_sidecar for the sidecar_id**: [#165](https://github.com/cyralinc/terraform-provider-cyral/pull/165);

### Others:

- **ci: adds pre-commit support**: [#162](https://github.com/cyralinc/terraform-provider-cyral/pull/162);

## 2.4.3 (February 15, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- **Modifying repository_binding resource is not removing old resource**: [#159](https://github.com/cyralinc/terraform-provider-cyral/pull/159);

## 2.4.2 (February 3, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Documentation:

- **Guide for setting up CP and deploy a sidecar**: [#156](https://github.com/cyralinc/terraform-provider-cyral/pull/156);

## 2.4.1 (January 5, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- **Terraform provider creates a new identity map instead of updating old one when state changes**: [#152](https://github.com/cyralinc/terraform-provider-cyral/pull/152);
- **Fixes issue about missing sidecar returning 500 error**: [#148](https://github.com/cyralinc/terraform-provider-cyral/pull/148);

## 2.4.0 (December 8, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Documentation:

- **Remove wrong 'okta' reference**: [#135](https://github.com/cyralinc/terraform-provider-cyral/pull/135);
- **Update examples to use for_each instead of count**: [#137](https://github.com/cyralinc/terraform-provider-cyral/pull/137);
- **Fix wrong flattened resources in docs**: [#142](https://github.com/cyralinc/terraform-provider-cyral/pull/142);

### Features:

- **Add CYRAL_TF_CONTROL_PLANE env var support**: [#139](https://github.com/cyralinc/terraform-provider-cyral/pull/139);

## 2.3.1 (November 19, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- **Fix image links**: [#133](https://github.com/cyralinc/terraform-provider-cyral/pull/133);

## 2.3.0 (November 19, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- **Fix documentation accordingly to Terraform standards**: [#131](https://github.com/cyralinc/terraform-provider-cyral/pull/131);

### Features:

- **Add resource to manage sso groups to roles**: [#106](https://github.com/cyralinc/terraform-provider-cyral/pull/106);
- **Add docs reference to Okta IdP module**: [#129](https://github.com/cyralinc/terraform-provider-cyral/pull/129);
- **Add sidecar user endpoint**: [#132](https://github.com/cyralinc/terraform-provider-cyral/pull/132);

## 2.2.1 (November 12, 2021)

Minimum required Control Plane version: `v2.25.0`.

### Bug fixes:

- **Fix IdP registration at CP-level**: [#128](https://github.com/cyralinc/terraform-provider-cyral/pull/128);

## 2.2.0 (November 4, 2021)

Minimum required Control Plane version: `v2.25.0`.

### Bug fixes:

- **Fix cyclic dependency issue in SAML certificate data source**: [#121](https://github.com/cyralinc/terraform-provider-cyral/pull/121);

### Deprecated resources:

- `cyral_integration_sso_*` renamed to `cyral_integration_idp_*`

## 2.1.1 (October 21, 2021)

Minimum required Control Plane version: `v2.24.0`.

### Bug fixes:

- **Remove unnecessary PreCheck from Terraform Provider Tests**: [#117](https://github.com/cyralinc/terraform-provider-cyral/pull/117);

## 2.1.0 (October 18, 2021)

Minimum required Control Plane version: `v2.24.0`.

### Features:

- **Generic SAML integration and new SSO resources**: [#115](https://github.com/cyralinc/terraform-provider-cyral/pull/115);

## 2.0.1 (September 30, 2021)

Minimum required Control Plane version: `v2.22.0`.

### Bug fixes:

- **Omitting access_duration in cyral_identity_map resulted in plan change on every plan**: [#111](https://github.com/cyralinc/terraform-provider-cyral/pull/111);

## 2.0.0 (September 24, 2021)

Minimum required Control Plane version: `v2.22.0`.

### Backwards compatibility breaks:

- **Resource cyral_sidecar**: changed parameters;
- **Data source cyral_sidecar_template**: data source replaced by `cyral_sidecar_cft_template` and template restricted to Cloudformation.

### Features:

- **Script to rotate service account secrets**: [#64](https://github.com/cyralinc/terraform-provider-cyral/pull/64);
- **Improve tooling**: [#92](https://github.com/cyralinc/terraform-provider-cyral/pull/92);
- **Resource Sidecar Credentials**: [#93](https://github.com/cyralinc/terraform-provider-cyral/pull/93);
- **Resource Repository Conf Analysis**: [#108](https://github.com/cyralinc/terraform-provider-cyral/pull/108);

## 1.2.2 (June 21, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Bug fixes:

- **Fix missing Helm3 support**: fix missing helm 3 support in data source `cyral_sidecar_template` ([#63](https://github.com/cyralinc/terraform-provider-cyral/pull/63));

---

## 1.2.1 (June 18, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Bug fixes:

- **Fix publishing issue**: fix issue publishing binaries ([#62](https://github.com/cyralinc/terraform-provider-cyral/pull/62));

---

## 1.2.0 (June 17, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Features:

- **Data Source SAML Certificate**: added new data source to retrieve certificate ([#60](https://github.com/cyralinc/terraform-provider-cyral/pull/60));
- **Docker Build**: added new feature in makefile to allow docker build without requiring local Go environment ([#61](https://github.com/cyralinc/terraform-provider-cyral/pull/61));
- **Resource Integration Hashicorp Vault**: added new resource to support Hashicorp Vault integration ([#58](https://github.com/cyralinc/terraform-provider-cyral/pull/58));
- **Resource Integration Pager Duty**: added new resource to support Pager Duty integration ([#56](https://github.com/cyralinc/terraform-provider-cyral/pull/56));

---

## 1.1.0 (April 23, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:

- **Resource Integration Slack Alerts**: added new resource to support Slack Alerts integration ([#43](https://github.com/cyralinc/terraform-provider-cyral/pull/43));
- **Resource Integration Microsoft Teams**: added new resource to support Microsoft Teams integration ([#44](https://github.com/cyralinc/terraform-provider-cyral/pull/44));
- **Resource Repository Local Account**: added new resource to support repositories local accounts ([#46](https://github.com/cyralinc/terraform-provider-cyral/pull/46));
- **Resource Integration Okta**: added new resource to support Okta integration ([#47](https://github.com/cyralinc/terraform-provider-cyral/pull/47));
- **Resource Repository Configuration Authentication**: added new resource to support repository configuration authentication ([#48](https://github.com/cyralinc/terraform-provider-cyral/pull/48));
- **Data Source Sidecar Template**: added new data source to support sidecar templates ([#50](https://github.com/cyralinc/terraform-provider-cyral/pull/50));
- **Resource Identity Map**: added new resource to support identity maps ([#51](https://github.com/cyralinc/terraform-provider-cyral/pull/51));

### Improvements:

- **Increase Test Coverage**: defined standards and increased the test coverage. ([#41](https://github.com/cyralinc/terraform-provider-cyral/pull/41));

---

## 1.0.0 (March 26, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:

- **Resource Sidecar**: added new resource to support sidecars ([#23](https://github.com/cyralinc/terraform-provider-cyral/pull/23));
- **Resource Datamap**: added new resource to support datamaps ([#24](https://github.com/cyralinc/terraform-provider-cyral/pull/24));
- **Resource Repository Binding**: added new resource to support binding repositories to sidecars ([#25](https://github.com/cyralinc/terraform-provider-cyral/pull/25));
- **Resource Policy**: added new resource to support policies ([#26](https://github.com/cyralinc/terraform-provider-cyral/pull/26));
- **Resource Policy Rule**: added new resource to support policy rules ([#28](https://github.com/cyralinc/terraform-provider-cyral/pull/28));
- **Resource: Datadog Integration**: added new resource to support Datadog integration ([#30](https://github.com/cyralinc/terraform-provider-cyral/pull/30));
- **Resource ELK Integration**: added new resource to support ELK integration ([#31](https://github.com/cyralinc/terraform-provider-cyral/pull/31));
- **Resource Splunk Integration**: added new resource to support Splunk integration ([#33](https://github.com/cyralinc/terraform-provider-cyral/pull/33));
- **Resource Sumo Logic Integration**: added new resource to support Sumo Logic integration ([#34](https://github.com/cyralinc/terraform-provider-cyral/pull/34));
- **Resource Logstash Integration**: added new resource to support Logstash integration ([#35](https://github.com/cyralinc/terraform-provider-cyral/pull/35));
- **Resource Looker Integration**: added new resource to support Looker integration ([#36](https://github.com/cyralinc/terraform-provider-cyral/pull/36)).

### Improvements:

- **Replace repo ID property**: replaced repository identification in the state file from `name` to `id`. ([#22](https://github.com/cyralinc/terraform-provider-cyral/pull/22));

---

## 0.2.0 (February 24, 2021)

### Features:

- **Terraform Plugin SDK v2**: migration to `Terraform Plugin SDK v2` and support for Terraform `v0.12` to `v0.14` ([#21](https://github.com/cyralinc/terraform-provider-cyral/pull/21)).

### Improvements:

- **API Port**: Removed `control_plane_api_port` and added this information to the existing `control_plane` parameter. ([#20](https://github.com/cyralinc/terraform-provider-cyral/pull/20));
- **Unit tests**: added unit tests for `config.go` ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8));

---

## 0.1.0 (May 15, 2020)

### Features:

- **Terraform import**: added support for `terraform import` statement for `repository` resource ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8)).

---

## 0.0.2 (May 2, 2020)

### Bug fixes:

- **Error handling**: fixed bug in error handling ([#5](https://github.com/cyralinc/terraform-provider-cyral/pull/5)).

---

## 0.0.1 (May 2, 2020)

### Features:

- **Provider with Repository**: draft `provider` and `repository` resource ([#1](https://github.com/cyralinc/terraform-provider-cyral/pull/1), ([#2](https://github.com/cyralinc/terraform-provider-cyral/pull/2))).
