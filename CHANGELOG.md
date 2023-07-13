## 4.5.1 (July 12, 2023)

Minimum recommended Control Plane version: `v4.8.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible arguments and resource are not used.

Arguments incompatible with Control Planes previous to `v4.8`: `cyral_sidecar.activity_log_integration_id` and `cyral_sidecar.diagnostic_log_integration_id`.
Resource incompatible with Control Planes previous to `v4.7`: `cyral_integration_logging`.
Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

## Bug fixes:

- **Fix conflicting log integration ID**: [#420](https://github.com/cyralinc/terraform-provider-cyral/pull/420).
- **ENG-12107: Add missing validation to idp saml resource**: [#421](https://github.com/cyralinc/terraform-provider-cyral/pull/421).

## 4.5.0 (June 30, 2023)

Minimum recommended Control Plane version: `v4.8.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible arguments and resource are not used.

Arguments incompatible with Control Planes previous to `v4.8`: `cyral_sidecar.activity_log_integration_id` and `cyral_sidecar.diagnostic_log_integration_id`.
Resource incompatible with Control Planes previous to `v4.7`: `cyral_integration_logging`.
Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

## Features:

- **Add activity and diagnostic log ID for sidecar**: [#413](https://github.com/cyralinc/terraform-provider-cyral/pull/413).

## 4.4.0 (June 16, 2023)

Minimum recommended Control Plane version: `v4.7.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument and resource are not used.

Resource incompatible with Control Planes previous to `v4.7`: `cyral_integration_logging`.
Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

## Features:

- **Add data source cyral_sidecar_listener**: [#410](https://github.com/cyralinc/terraform-provider-cyral/pull/410).

### Documentation:

- **Docs improvement and smart ports guide**: [#408](https://github.com/cyralinc/terraform-provider-cyral/pull/408).

## 4.3.1 (June 7, 2023)

Minimum recommended Control Plane version: `v4.7.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument and resource are not used.

Resource incompatible with Control Planes previous to `v4.7`: `cyral_integration_logging`.
Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Documentation:

- **Update mongodb example to latest cyral_idp_okta module**: [#405](https://github.com/cyralinc/terraform-provider-cyral/pull/405).

## 4.3.0 (June 6, 2023)

Minimum recommended Control Plane version: `v4.7.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument and resource are not used.

Resource incompatible with Control Planes previous to `v4.7`: `cyral_integration_logging`.
Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Deprecate:

- **Deprecate old resources**: [#399](https://github.com/cyralinc/terraform-provider-cyral/pull/399).
- **ENG-11753: Add deprecation notice**: [#397](https://github.com/cyralinc/terraform-provider-cyral/pull/397).

### Features:

- **ENG-11769: Update Generic SAML Draft to return SP Metadata**: [#394](https://github.com/cyralinc/terraform-provider-cyral/pull/394).
- **ENG-11747: Add resource and datasource for new "log management" integration**: [#395](https://github.com/cyralinc/terraform-provider-cyral/pull/395).
- **Add examples for using Integration logging Resources and Data Sources**: [#400](https://github.com/cyralinc/terraform-provider-cyral/pull/400).

### Documentation:

- **Update guides to remove CP ports and fix format**: [#403](https://github.com/cyralinc/terraform-provider-cyral/pull/403).
- **Update guides to provider 4.0**: [#402](https://github.com/cyralinc/terraform-provider-cyral/pull/402).

## 4.2.0 (April 27, 2023)

Minimum recommended Control Plane version: `v4.2.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument is not used.

Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Deprecate:

- **ENG-11470: Deprecation notices and minor docs fixes**: [#389](https://github.com/cyralinc/terraform-provider-cyral/pull/389).

### Features:

- **ENG-11522: Add sidecar log integration ID**: [#387](https://github.com/cyralinc/terraform-provider-cyral/pull/387).

## 4.1.2 (March 17, 2023)

Minimum recommended Control Plane version: `v4.2.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument is not used.

Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- **Update ADFS default value for IdP resource**: [#364](https://github.com/cyralinc/terraform-provider-cyral/pull/364).

### Documentation:

- **Update docs index with migration guide**: [#373](https://github.com/cyralinc/terraform-provider-cyral/pull/373).
- **Overall documentation improments**: [#375](https://github.com/cyralinc/terraform-provider-cyral/pull/375).

### Improvements:

- **Update SDK**: [#365](https://github.com/cyralinc/terraform-provider-cyral/pull/365).
- **Replace POST by PUT in access rule creation**: [#376](https://github.com/cyralinc/terraform-provider-cyral/pull/376).
- **Update Conf Auth resource to fix recreation issue**: [#378](https://github.com/cyralinc/terraform-provider-cyral/pull/378).

## 4.1.1 (February 21, 2023)

Minimum recommended Control Plane version: `v4.2.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument is not used.

Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Documentation:

- **Improve cyral_repository docs**: [#361](https://github.com/cyralinc/terraform-provider-cyral/pull/361).

## 4.1.0 (February 14, 2023)

Minimum recommended Control Plane version: `v4.2.0`. It is safe to use this provider with all `v4` control planes
as long as the incompatible argument is not used.

Argument incompatible with Control Planes previous to `v4.2`: `cyral_repository.mongodb_settings.srv_record_name`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Features:

- **Teach Terraform provider about MongoDB SRV Records**: [#336](https://github.com/cyralinc/terraform-provider-cyral/pull/336).

## 4.0.4 (February 9, 2023)

Minimum required Control Plane version: `v4.0.0`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- **Make mongodb_settings mandatory if cyral_repository.type = mongodb**: [#355](https://github.com/cyralinc/terraform-provider-cyral/pull/355).

### Documentation fixes:

- **Fix broken link and update docs**: [#354](https://github.com/cyralinc/terraform-provider-cyral/pull/354).

## 4.0.3 (February 8, 2023)

Minimum required Control Plane version: `v4.0.0`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- **Update migration scripts with identity map fix and update guides**: [#352](https://github.com/cyralinc/terraform-provider-cyral/pull/352).

## 4.0.2 (February 3, 2023)

Minimum required Control Plane version: `v4.0.0`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- **Update 4.0 migration guide to properly handle for_each resource definitions**: [#350](https://github.com/cyralinc/terraform-provider-cyral/pull/350).

## 4.0.1 (February 1, 2023)

Minimum required Control Plane version: `v4.0.0`.

See the list of incompatible resources with Control Planes `v3.x` and provider `v3.x` in the [`4.0.0`](#400-january-27-2023) release documentation.

### Documentation fixes:

- **Fix typo in version name**: [#348](https://github.com/cyralinc/terraform-provider-cyral/pull/348).
- **Improve documentation**: [#349](https://github.com/cyralinc/terraform-provider-cyral/pull/349).

### Other:

- **Update release workflow**: [#347](https://github.com/cyralinc/terraform-provider-cyral/pull/347).

## 4.0.0 (January 27, 2023)

Minimum required Control Plane version: `v4.0.0`.

Checkout the [v4 Migration Guide](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/4.0-migration-guide)
if you are upgrading from provider versions `v2` or `v3`. It can save you a lot of time.

Resources incompatible with Control Planes `v3.x`:
`cyral_repository` and `cyral_repository_binding`.

New resources:
`cyral_sidecar_listener` and `cyral_repository_access_gateway`.

Removed resource arguments:

- `cyral_repository.host` -- use `cyral_repository.repo_node.host` instead.
- `cyral_repository.port` -- use `cyral_repository.repo_node.port` instead.
- `cyral_repository.properties.mongodb_replica_set.max_nodes` -- this argument is no
  longer used and is inferred from the number of `repo_node` blocks declared in
  `cyral_repository`.
- `cyral_repository.properties.mongodb_replica_set.replica_set_id` -- use
  `cyral_repository.mongodb_settings.replica_set_name` instead.
- `cyral_repository_binding.listener_host` -- use `cyral_sidecar_listener.network_address.host` instead.
- `cyral_repository_binding.listener_port` -- use `cyral_sidecar_listener.network_address.port` instead.
- `cyral_repository_binding.sidecar_as_idp_access_gateway` -- use `cyral_repository_access_gateway` instead.

### Features:

- **New resource: Listeners API**: [#281](https://github.com/cyralinc/terraform-provider-cyral/pull/281).
- **Access Gateway Binding Resource**: [#331](https://github.com/cyralinc/terraform-provider-cyral/pull/331).

### Backwards compatibility breaks:

- **Port-multiplexing Cyral Repository Changes**: [#326](https://github.com/cyralinc/terraform-provider-cyral/pull/326).
- **Repository Binding Resource Changes**: [#329](https://github.com/cyralinc/terraform-provider-cyral/pull/329).

### Bug fixes:

- **Improved migration script runtime by using less pipes**: [#324](https://github.com/cyralinc/terraform-provider-cyral/pull/324).
- **clean up supported repos**: [#332](https://github.com/cyralinc/terraform-provider-cyral/pull/332).
- **Adding Force New Directives to Access Gateway resource**: [#340](https://github.com/cyralinc/terraform-provider-cyral/pull/340).

### Documentation fixes:

- **Remove references to deprecated fields mongodb_port_alloc_range_low/high**: [#339](https://github.com/cyralinc/terraform-provider-cyral/pull/339).

## 3.0.5 (February 8, 2023)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with Control Planes `v2.x` and Terraform `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- **Fix identity map issues on 3.0-migration script**: [#353](https://github.com/cyralinc/terraform-provider-cyral/pull/353).
- **Removed mentions of cassandra and bigquery as they are not supported as repo types**: [#332](https://github.com/cyralinc/terraform-provider-cyral/pull/332).

### Improvements:

- **Improved migration script runtime by using less pipes**: [#324](https://github.com/cyralinc/terraform-provider-cyral/pull/324).

## 3.0.4 (November 18, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with Control Planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- **Modified the script to call terraform show less**: [#320](https://github.com/cyralinc/terraform-provider-cyral/pull/320).

## 3.0.3 (November 14, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with Control Planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- **Deprecate & Update Cyral Permissions in for Cyral Roles**: [#313](https://github.com/cyralinc/terraform-provider-cyral/pull/313).
- **Fix sidecar resource nil properties issue**: [#317](https://github.com/cyralinc/terraform-provider-cyral/pull/317).

## 3.0.2 (November 7, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with Control Planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- **Remove rewrite_on_violation, add enable_dataset_rewrites**: [#308](https://github.com/cyralinc/terraform-provider-cyral/pull/308).
- **Add support for linux sidecar deployment**: [#311](https://github.com/cyralinc/terraform-provider-cyral/pull/311).

### Documentation fixes:

- **Fix typo edit in migration script**: [#307](https://github.com/cyralinc/terraform-provider-cyral/pull/307).

### Security fixes:

- **Bump github.com/stretchr/testify from 1.8.0 to 1.8.1**: [#306](https://github.com/cyralinc/terraform-provider-cyral/pull/306).

## 3.0.1 (October 18, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with Control Planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Features:

- **Add Cyral Terraform 3.0 Migration Guide**: [#303](https://github.com/cyralinc/terraform-provider-cyral/pull/303).

### Bug fixes:

- **Policy rule identities field should be omitted by default**: [#301](https://github.com/cyralinc/terraform-provider-cyral/pull/301).

## 3.0.0 (October 5, 2022)

Minimum required Control Plane version: `v3.0.0`.

Checkout the [v3 Migration Guide](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/3.0-migration-guide)
if you are upgrading from provider versions `v2`. It can save you a lot of time.

Resources incompatible with Control Planes `v2.x`:
`cyral_datamap` (removed, refer to `cyral_repository_datamap` instead),
`cyral_identity_map` (removed, use `cyral_repository_access_rules` instead),
`cyral_integration_okta` (removed, refer to `cyral_integration_idp_okta` instead),
`cyral_integration_sso_*` (renamed, refer to `cyral_integration_idp_*` instead),
`cyral_repository_identity_map` (removed, use `cyral_repository_access_rules` instead),
`cyral_repository_local_account` (removed, use `cyral_repository_user_account` instead).

Removed resource arguments:

- `cyral_policy.properties`

Removed provider arguments:

- `auth_provider`
- `auth0_audience`
- `auth0_domain`
- `auth0_client_id`
- `auth0_client_secret`

Renamed resource arguments:

- `cyral_policy_rule.id` -- now
  `cyral_policy_rule.policy_rule_id`. `cyral_policy_rule.id` now contains a
  composed ID in the format `{policy_id}/{policy_rule_id}`.

### Features:

- **Add resource AccessRules**: [#280](https://github.com/cyralinc/terraform-provider-cyral/pull/280).
- **User Account resource for Gatekeeper Project**: [#288](https://github.com/cyralinc/terraform-provider-cyral/pull/288).
- **Adapt guides to v3**: [#297](https://github.com/cyralinc/terraform-provider-cyral/pull/297).

### Backwards compatibility breaks:

- **Remove deprecated resources for major version 3**: [#291](https://github.com/cyralinc/terraform-provider-cyral/pull/291).
- **Remove resources for local_accounts and identity_maps**: [#293](https://github.com/cyralinc/terraform-provider-cyral/pull/293).

### Bug fixes:

- **Fix acceptance tests for network access policy**: [#289](https://github.com/cyralinc/terraform-provider-cyral/pull/289).
- **Fix race condition in user accounts test**: [#295](https://github.com/cyralinc/terraform-provider-cyral/pull/295).

## 2.11.1 (November 14, 2022)

Minimum required Control Plane version: `v2.35.0`.
Resource incompatible with Control Planes previous to `v2.35`: argument `enable_dataset_rewrites` from resource `cyral_repository_conf_analysis`.

### Bug fixes:

- **Deprecate 'View Sidecars and Repositories'**: [#315](https://github.com/cyralinc/terraform-provider-cyral/pull/315).

## 2.11.0 (November 7, 2022)

Minimum required Control Plane version: `v2.35.0`.
Resource incompatible with Control Planes previous to `v2.35`: argument `enable_dataset_rewrites` from resource `cyral_repository_conf_analysis`.

### Features:

- **Deprecate rewrite_on_violation, add enable_dataset_rewrites**: [#310](https://github.com/cyralinc/terraform-provider-cyral/pull/310).
- **Add support for linux sidecar deployment**: [#311](https://github.com/cyralinc/terraform-provider-cyral/pull/311).

## 2.10.2 (October 18, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with Control Planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Bug fixes:

- **Policy rule identities field should be omitted by default**: [#302](https://github.com/cyralinc/terraform-provider-cyral/pull/302).

## 2.10.1 (October 6, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with Control Planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Documentation:

- **Use version constraint ~> for guides (v2)**: [#298](https://github.com/cyralinc/terraform-provider-cyral/pull/298).

## 2.10.0 (October 5, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with Control Planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Features:

- **Add support for local account automatic approvals**: [#296](https://github.com/cyralinc/terraform-provider-cyral/pull/296).

## 2.9.0 (September 16, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with Control Planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Features:

- **Add support for Generic SAML IdP Integration**: [#244](https://github.com/cyralinc/terraform-provider-cyral/pull/244).
- **Add Duo MFA integration resource**: [#282](https://github.com/cyralinc/terraform-provider-cyral/pull/282).
- **Add support for repository network shield**: [#285](https://github.com/cyralinc/terraform-provider-cyral/pull/285).

### Bug fixes:

- **Fix ACC test changes after merging fix**: [#275](https://github.com/cyralinc/terraform-provider-cyral/pull/275).
- **Fix acc tests and improve sidecar test**: [#266](https://github.com/cyralinc/terraform-provider-cyral/pull/266).
- **Fix issues with repository resource properties**: [#271](https://github.com/cyralinc/terraform-provider-cyral/pull/27`2`).
- **Fix issue with repository resource properties**: [#272](https://github.com/cyralinc/terraform-provider-cyral/pull/272).

### Improvements:

- **Refactor acceptance tests**: [#267](https://github.com/cyralinc/terraform-provider-cyral/pull/267).
- **Explain setting max value of port range**: [#277](https://github.com/cyralinc/terraform-provider-cyral/pull/277).

## 2.8.0 (August 5, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with Control Planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Bug fixes:

- **Fix resource import functions**: [#260](https://github.com/cyralinc/terraform-provider-cyral/pull/260).
- **Fix data map attribute deletion**: [#263](https://github.com/cyralinc/terraform-provider-cyral/pull/263).

### Features:

- **Add parameter to enable data masking**: [#252](https://github.com/cyralinc/terraform-provider-cyral/pull/252).
- **Add support for dynamodb**: [#253](https://github.com/cyralinc/terraform-provider-cyral/pull/253).
- **Support advanced sidecar options**: [#255](https://github.com/cyralinc/terraform-provider-cyral/pull/255).
- **Add data source to retrieve data labels**: [#261](https://github.com/cyralinc/terraform-provider-cyral/pull/261).
- **Add single container type to sidecar resource**: [#265](https://github.com/cyralinc/terraform-provider-cyral/pull/265).

## 2.7.2 (August 15, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- **Fix issues with repository resource properties**: [#271](https://github.com/cyralinc/terraform-provider-cyral/pull/271).

## 2.7.1 (July 19, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- **Fix token expiration error for large configuration scripts**: [#247](https://github.com/cyralinc/terraform-provider-cyral/pull/247).

## 2.7.0 (July 1, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- **Fix memory referencing issue and refactor resources**: [#223](https://github.com/cyralinc/terraform-provider-cyral/pull/223).

### Features:

- **Add data source for cyral_integration_idp**: [#239](https://github.com/cyralinc/terraform-provider-cyral/pull/239).
- **Add resources for per-repository data map**: [#230](https://github.com/cyralinc/terraform-provider-cyral/pull/230).
- **Add Terraform guides for basic usages of Cyral features**: [#238](https://github.com/cyralinc/terraform-provider-cyral/pull/238).
- **Support MongoDB repository replica sets**: [#228](https://github.com/cyralinc/terraform-provider-cyral/pull/228).

## 2.6.2 (June 6, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- **Fix repository name for IdP Okta module**: [#220](https://github.com/cyralinc/terraform-provider-cyral/pull/220);

## 2.6.1 (May 31, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- **Fix broken links in documentation**: [#216](https://github.com/cyralinc/terraform-provider-cyral/pull/216);

## 2.6.0 (May 27, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with Control Planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

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
