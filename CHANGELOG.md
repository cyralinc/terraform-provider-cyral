## 4.13.1 (August 21, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.16`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.13.0`](#4130-august-21-2024)
release documentation.

### Documentation:

- Fix wrong doc parameters in repo conf auth ([#563](https://github.com/cyralinc/terraform-provider-cyral/pull/563))

## 4.13.0 (August 21, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.16`.

Arguments incompatible with control planes previous to `v4.16`:

- `cyral_repository.redshift_settings.aws_region`
- `cyral_repository.redshift_settings.cluster_identifier`
- `cyral_repository.redshift_settings.workgroup_name`

Resource incompatible with control planes previous to `v4.15`:

- `cyral_policy_v2`

Argument incompatible with control planes previous to `v4.14`:

- `cyral_repository_user_account.auth_scheme.azure_key_vault`

Argument incompatible with control planes previous to `v4.12`:

- `cyral_repository.mongodb_settings.flavor`

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- ENG-14351: Fix invalid TLS option check ([#561](https://github.com/cyralinc/terraform-provider-cyral/pull/561))

### Documentation:

- Minor doc fix ([#554](https://github.com/cyralinc/terraform-provider-cyral/pull/554))

### Features:

- ENG-14270: redshift iam ([#559](https://github.com/cyralinc/terraform-provider-cyral/pull/559))

### Improvements:

- Bump github.com/hashicorp/terraform-plugin-docs from 0.18.0 to 0.19.4 ([#543](https://github.com/cyralinc/terraform-provider-cyral/pull/543))
- Bump alpine from 3.20.1 to 3.20.2 ([#553](https://github.com/cyralinc/terraform-provider-cyral/pull/553))
- Bump golang.org/x/oauth2 from 0.21.0 to 0.22.0 ([#556](https://github.com/cyralinc/terraform-provider-cyral/pull/556))
- Bump hashicorp/terraform from 1.9.0 to 1.9.4 ([#557](https://github.com/cyralinc/terraform-provider-cyral/pull/557))
- Bump golang from 1.22.5-alpine3.20 to 1.22.6-alpine3.20 ([#558](https://github.com/cyralinc/terraform-provider-cyral/pull/558))
- Bump hashicorp/terraform from 1.9.4 to 1.9.5 ([#562](https://github.com/cyralinc/terraform-provider-cyral/pull/562))

## 4.12.0 (July 15, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.15`.

Resource incompatible with control planes previous to `v4.15`:

- `cyral_policy_v2`

Argument incompatible with control planes previous to `v4.14`:

- `cyral_repository_user_account.auth_scheme.azure_key_vault`

Argument incompatible with control planes previous to `v4.12`:

- `cyral_repository.mongodb_settings.flavor`

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Features:

- ENG-12949: Policy engine providers ([#547](https://github.com/cyralinc/terraform-provider-cyral/pull/547))

## 4.11.0 (June 6, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.14`.

Argument incompatible with control planes previous to `v4.14`:

- `cyral_repository_user_account.auth_scheme.azure_key_vault`

Argument incompatible with control planes previous to `v4.12`:

- `cyral_repository.mongodb_settings.flavor`

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Documentation:

- Improve examples for resources policy and policy_rules ([#545](https://github.com/cyralinc/terraform-provider-cyral/pull/545))

### Improvements:

- ENG-14083: Add Azure Key Vault user account auth scheme ([#542](https://github.com/cyralinc/terraform-provider-cyral/pull/542))
- Bump golang.org/x/oauth2 from 0.20.0 to 0.21.0 ([#544](https://github.com/cyralinc/terraform-provider-cyral/pull/544))

## 4.10.1 (May 31, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.10.0`](#4100-april-9-2024)
release documentation.

### Documentation:

- ENG-13977: Add deprecation message in existing ELK Terraform ([#540](https://github.com/cyralinc/terraform-provider-cyral/pull/540))
- Document access rule order better ([#541](https://github.com/cyralinc/terraform-provider-cyral/pull/541))

### Improvements:

- Change PASSPHRASE -> GPG_PASSPHRASE ([#528](https://github.com/cyralinc/terraform-provider-cyral/pull/528))
- Bump golang.org/x/net from 0.22.0 to 0.23.0 ([#530](https://github.com/cyralinc/terraform-provider-cyral/pull/530))
- Bump golang.org/x/oauth2 from 0.19.0 to 0.20.0 ([#534](https://github.com/cyralinc/terraform-provider-cyral/pull/534))
- Bump alpine from 3.19.1 to 3.20.0 ([#537](https://github.com/cyralinc/terraform-provider-cyral/pull/537))

## 4.10.0 (April 9, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12`.

Argument incompatible with control planes previous to `v4.12`:

- `cyral_repository.mongodb_settings.flavor`

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Documentation:

- ENG-13675: update documentation for identity_provider field ([#516](https://github.com/cyralinc/terraform-provider-cyral/pull/516))
- Fix docs formatting ([#520](https://github.com/cyralinc/terraform-provider-cyral/pull/520))

### Improvements:

- Standardize error handling and refactor old resources ([#521](https://github.com/cyralinc/terraform-provider-cyral/pull/521))
- Refactor remaining resources and data sources ([#522](https://github.com/cyralinc/terraform-provider-cyral/pull/522))
- Bump golang.org/x/oauth2 from 0.18.0 to 0.19.0 ([#523](https://github.com/cyralinc/terraform-provider-cyral/pull/523))
- Ignore regex analysis for 404 status code ([#524](https://github.com/cyralinc/terraform-provider-cyral/pull/524))

## 4.9.3 (March 14, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.9.0`](#490-january-31-2024)
release documentation.

### Documentation:

- Fix broken links and improve import documentation ([#519](https://github.com/cyralinc/terraform-provider-cyral/pull/519))

### Improvements:

- Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.32.0 to 2.33.0 ([#512](https://github.com/cyralinc/terraform-provider-cyral/pull/512))
- Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#514](https://github.com/cyralinc/terraform-provider-cyral/pull/514))
- Bump golang.org/x/oauth2 from 0.17.0 to 0.18.0 ([#515](https://github.com/cyralinc/terraform-provider-cyral/pull/515))
- Bump google.golang.org/protobuf from 1.32.0 to 1.33.0 ([#517](https://github.com/cyralinc/terraform-provider-cyral/pull/517))
- Bump hashicorp/terraform from 1.7.3 to 1.7.5 ([#518](https://github.com/cyralinc/terraform-provider-cyral/pull/518))

## 4.9.2 (February 28, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.9.0`](#490-january-31-2024)
release documentation.

### Documentation:

- Fix role name in RDS IAM required permissions ([#513](https://github.com/cyralinc/terraform-provider-cyral/pull/513))

## 4.9.1 (February 14, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.9.0`](#490-january-31-2024)
release documentation.

### Bug fixes:

- Add missing mongodb flavor configuration ([#509](https://github.com/cyralinc/terraform-provider-cyral/pull/509))

### Documentation:

- Rollback operations supported by datasetprotection policy ([#501](https://github.com/cyralinc/terraform-provider-cyral/pull/501))
- Add AWS RDS IAM auth guide ([#508](https://github.com/cyralinc/terraform-provider-cyral/pull/508))

### Improvements:

- Bump golang.org/x/oauth2 from 0.16.0 to 0.17.0 ([#507](https://github.com/cyralinc/terraform-provider-cyral/pull/507))

## 4.9.0 (January 31, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12`.

Argument incompatible with control planes previous to `v4.12`:

- `cyral_repository.mongodb_settings.flavor`

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- Update CFT template data source info and disable tests ([#500](https://github.com/cyralinc/terraform-provider-cyral/pull/500))

### Features:

- ENG-13251: Add optional MongoDBSettings 'flavor' field to repository resource ([#503](https://github.com/cyralinc/terraform-provider-cyral/pull/503))

### Improvements:

- Bump github.com/google/uuid from 1.5.0 to 1.6.0 ([#497](https://github.com/cyralinc/terraform-provider-cyral/pull/497))
- Bump github.com/hashicorp/terraform-plugin-docs from 0.16.0 to 0.18.0 ([#498](https://github.com/cyralinc/terraform-provider-cyral/pull/498))
- Bump hashicorp/terraform from 1.3.9 to 1.7.1 ([#499](https://github.com/cyralinc/terraform-provider-cyral/pull/499))
- Bump alpine from 3.18.5 to 3.19.1 ([#502](https://github.com/cyralinc/terraform-provider-cyral/pull/502))
- Bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.31.0 to 2.32.0 ([#504](https://github.com/cyralinc/terraform-provider-cyral/pull/504))

## 4.8.1 (January 18, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.8.0`](#480-january-3-2024)
release documentation.

### Documentation:

- Add repo-level policy guide and enhance cyral_rego_policy_instance docs ([#495](https://github.com/cyralinc/terraform-provider-cyral/pull/495))
- Improve descriptions and file organization for repo-level policy guide ([#496](https://github.com/cyralinc/terraform-provider-cyral/pull/496))

## 4.8.0 (January 3, 2024)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.12.0`.

Resource incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Data source incompatible with control planes previous to `v4.12`:

- `cyral_access_token_settings`

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Documentation:

- Fix missing reference in cyral.core docs ([#473](https://github.com/cyralinc/terraform-provider-cyral/pull/473))

## Features:

- ENG-12954, ENG-12955: Add datasource and resource fot token settings ([#479](https://github.com/cyralinc/terraform-provider-cyral/pull/479))

### Improvements:

- Create core package and refactor data label + datamap ([#439](https://github.com/cyralinc/terraform-provider-cyral/pull/439))
- Simplify core package ([#474](https://github.com/cyralinc/terraform-provider-cyral/pull/474))

## 4.7.2 (October 19, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.10.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.7.0`](#470-october-9-2023)
release documentation.

### Documentation:

- S3 guide ([#429](https://github.com/cyralinc/terraform-provider-cyral/pull/429))
- Update docs to avoid errors during destroy commands ([#469](https://github.com/cyralinc/terraform-provider-cyral/pull/469))

## 4.7.1 (October 11, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.10.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.7.0`](#470-october-9-2023)
release documentation.

### Bug fixes:

- Fix test so that it works running in parallel with other tests ([#466](https://github.com/cyralinc/terraform-provider-cyral/pull/466))

## 4.7.0 (October 9, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.10.0`.

Arguments incompatible with control planes previous to `v4.10`:

- `cyral_integration_logging.skip_validate`

Resource incompatible with control planes previous to `v4.10`:

- `cyral_integration_aws_iam`

Data sources incompatible with control planes previous to `v4.10`:

- `cyral_sidecar_health`
- `cyral_sidecar_instance_stats`
- `cyral_system_info`

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Features:

- ENG-12292: Add: SQL Server settings (version field) ([#438](https://github.com/cyralinc/terraform-provider-cyral/pull/438))
- Add `skip_validate` option to Fluent Bit logging integration resource ([#445](https://github.com/cyralinc/terraform-provider-cyral/pull/445))
- ENG-12558: add AuthType to the repo conf auth payload ([#450](https://github.com/cyralinc/terraform-provider-cyral/pull/450))
- ENG-12557: add CRUD operations for AWS IAM AuthN integration ([#451](https://github.com/cyralinc/terraform-provider-cyral/pull/451))
- ENG-12192: Add service account resource ([#453](https://github.com/cyralinc/terraform-provider-cyral/pull/453))
- ENG-12511: Add data source for systemInfo API ([#455](https://github.com/cyralinc/terraform-provider-cyral/pull/455))
- ENG-12678: Add data source for sidecar health API ([#456](https://github.com/cyralinc/terraform-provider-cyral/pull/456))
- ENG-12679: Add data source for sidecar instance API ([#457](https://github.com/cyralinc/terraform-provider-cyral/pull/457))
- ENG-12680: Add data source for sidecar instance stats API ([#459](https://github.com/cyralinc/terraform-provider-cyral/pull/459))
- ENG-12728: Add template parameters section to the rego_policy_instance docs ([#463](https://github.com/cyralinc/terraform-provider-cyral/pull/463))

### Bug fixes:

- Removing the `helm` and `cloudFormation` sidecar deployment types ([#452](https://github.com/cyralinc/terraform-provider-cyral/pull/452))
- Fix terraform tests that were failing in the E2E tests report ([#454](https://github.com/cyralinc/terraform-provider-cyral/pull/454))

## 4.6.0 (August 17, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Features:

- ENG-12334: Add data label tags field to the policy resource ([#433](https://github.com/cyralinc/terraform-provider-cyral/pull/433))
- ENG-12179: Add classification rule field to data label resource ([#436](https://github.com/cyralinc/terraform-provider-cyral/pull/436))
- ENG-12399, ENG-4406: Create rego policy instance resource and fix policy resource ([#440](https://github.com/cyralinc/terraform-provider-cyral/pull/440))

## 4.5.4 (August 3, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.5.0`](#450-june-30-2023)
release documentation.

### Bug fixes:

- Fix error handling in case of state file out of sync ([#430](https://github.com/cyralinc/terraform-provider-cyral/pull/430))

## 4.5.3 (July 27, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.5.0`](#450-june-30-2023)
release documentation.

### Documentation:

- Docs improvements ([#426](https://github.com/cyralinc/terraform-provider-cyral/pull/426))

## 4.5.2 (July 24, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.5.0`](#450-june-30-2023)
release documentation.

## Bug fixes:

- ENG-12193: Terrraform provider panics when ELK integration creds are not set ([#424](https://github.com/cyralinc/terraform-provider-cyral/pull/424))

### Documentation:

- Improve es_credentials description ([#423](https://github.com/cyralinc/terraform-provider-cyral/pull/423))

## 4.5.1 (July 12, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.5.0`](#450-june-30-2023)
release documentation.

## Bug fixes:

- Fix conflicting log integration ID ([#420](https://github.com/cyralinc/terraform-provider-cyral/pull/420))
- ENG-12107: Add missing validation to idp saml resource ([#421](https://github.com/cyralinc/terraform-provider-cyral/pull/421))

## 4.5.0 (June 30, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.8.0`.

Arguments incompatible with control planes previous to `v4.8`:

- `cyral_sidecar.activity_log_integration_id`
- `cyral_sidecar.diagnostic_log_integration_id`

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

## Features:

- Add activity and diagnostic log ID for sidecar ([#413](https://github.com/cyralinc/terraform-provider-cyral/pull/413))

## 4.4.0 (June 16, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.7.0`.

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

## Features:

- Add data source cyral_sidecar_listener ([#410](https://github.com/cyralinc/terraform-provider-cyral/pull/410))

### Documentation:

- Docs improvement and smart ports guide ([#408](https://github.com/cyralinc/terraform-provider-cyral/pull/408))

## 4.3.1 (June 7, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.7.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.3.0`](#430-june-6-2023)
release documentation.

### Documentation:

- Update mongodb example to latest cyral_idp_okta module ([#405](https://github.com/cyralinc/terraform-provider-cyral/pull/405))

## 4.3.0 (June 6, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.7.0`.

Resource incompatible with control planes previous to `v4.7`:

- `cyral_integration_logging`

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Deprecate:

- Deprecate old resources ([#399](https://github.com/cyralinc/terraform-provider-cyral/pull/399))
- ENG-11753: Add deprecation notice ([#397](https://github.com/cyralinc/terraform-provider-cyral/pull/397))

### Features:

- ENG-11769: Update Generic SAML Draft to return SP Metadata ([#394](https://github.com/cyralinc/terraform-provider-cyral/pull/394))
- ENG-11747: Add resource and datasource for new "log management" integration ([#395](https://github.com/cyralinc/terraform-provider-cyral/pull/395))
- Add examples for using Integration logging Resources and Data Sources ([#400](https://github.com/cyralinc/terraform-provider-cyral/pull/400))

### Documentation:

- Update guides to remove CP ports and fix format ([#403](https://github.com/cyralinc/terraform-provider-cyral/pull/403))
- Update guides to provider 4.0 ([#402](https://github.com/cyralinc/terraform-provider-cyral/pull/402))

## 4.2.0 (April 27, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.2.0`.

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Deprecate:

- ENG-11470: Deprecation notices and minor docs fixes ([#389](https://github.com/cyralinc/terraform-provider-cyral/pull/389))

### Features:

- ENG-11522: Add sidecar log integration ID ([#387](https://github.com/cyralinc/terraform-provider-cyral/pull/387))

## 4.1.2 (March 17, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.2.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.1.0`](#410-february-14-2023)
release documentation.

### Bug fixes:

- Update ADFS default value for IdP resource ([#364](https://github.com/cyralinc/terraform-provider-cyral/pull/364))
- Replace POST by PUT in access rule creation ([#376](https://github.com/cyralinc/terraform-provider-cyral/pull/376))
- Update Conf Auth resource to fix recreation issue ([#378](https://github.com/cyralinc/terraform-provider-cyral/pull/378))

### Documentation:

- Update docs index with migration guide ([#373](https://github.com/cyralinc/terraform-provider-cyral/pull/373))
- Overall documentation improments ([#375](https://github.com/cyralinc/terraform-provider-cyral/pull/375))

### Improvements:

- Update SDK ([#365](https://github.com/cyralinc/terraform-provider-cyral/pull/365))

## 4.1.1 (February 21, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.2.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.1.0`](#410-february-14-2023)
release documentation.

### Documentation:

- Improve cyral_repository docs ([#361](https://github.com/cyralinc/terraform-provider-cyral/pull/361))

## 4.1.0 (February 14, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.2.0`.

Argument incompatible with control planes previous to `v4.2`:

- `cyral_repository.mongodb_settings.srv_record_name`

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Features:

- Teach Terraform provider about MongoDB SRV Records ([#336](https://github.com/cyralinc/terraform-provider-cyral/pull/336))

## 4.0.4 (February 9, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.0.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.0.0`](#400-january-27-2023)
release documentation.

### Bug fixes:

- Make mongodb_settings mandatory if cyral_repository.type = mongodb ([#355](https://github.com/cyralinc/terraform-provider-cyral/pull/355))

### Documentation fixes:

- Fix broken link and update docs ([#354](https://github.com/cyralinc/terraform-provider-cyral/pull/354))

## 4.0.3 (February 8, 2023)

Minimum required Control Plane version: `v4.0.0`.

See the list of incompatible resources with control planes `v3.x` and provider `v3.x` in the [`v4.0.0`](#400-january-27-2023) release documentation.

### Bug fixes:

- Update migration scripts with identity map fix and update guides ([#352](https://github.com/cyralinc/terraform-provider-cyral/pull/352))

## 4.0.2 (February 3, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.0.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.0.0`](#400-january-27-2023)
release documentation.

### Bug fixes:

- Update 4.0 migration guide to properly handle for_each resource definitions ([#350](https://github.com/cyralinc/terraform-provider-cyral/pull/350))

## 4.0.1 (February 1, 2023)

It is safe to use this version with all `v4` control planes as long
as the new incompatible features are not used. These features require
a minimum version of the control plane and are detailed below.

The minimum control plane version required for full compatibility
with all the features in this release is `v4.0.0`.

See the list of incompatible attributes, data sources and resources
with previous `v4` control planes in the [`v4.0.0`](#400-january-27-2023)
release documentation.

### Documentation fixes:

- Fix typo in version name ([#348](https://github.com/cyralinc/terraform-provider-cyral/pull/348))
- Improve documentation ([#349](https://github.com/cyralinc/terraform-provider-cyral/pull/349))

### Other:

- Update release workflow ([#347](https://github.com/cyralinc/terraform-provider-cyral/pull/347))

## 4.0.0 (January 27, 2023)

Minimum required Control Plane version: `v4.0.0`.

Checkout the [v4 Migration Guide](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/4.0-migration-guide)
if you are upgrading from provider versions `v2` or `v3`. It can save you a lot of time.

Resources incompatible with control planes `v3.x`:

- `cyral_repository`
- `cyral_repository_binding`

New resources:

- `cyral_sidecar_listener`
- `cyral_repository_access_gateway`

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

- New resource: Listeners API ([#281](https://github.com/cyralinc/terraform-provider-cyral/pull/281))
- Access Gateway Binding Resource ([#331](https://github.com/cyralinc/terraform-provider-cyral/pull/331))

### Backwards compatibility breaks:

- Port-multiplexing Cyral Repository Changes ([#326](https://github.com/cyralinc/terraform-provider-cyral/pull/326))
- Repository Binding Resource Changes ([#329](https://github.com/cyralinc/terraform-provider-cyral/pull/329))

### Bug fixes:

- Improved migration script runtime by using less pipes ([#324](https://github.com/cyralinc/terraform-provider-cyral/pull/324))
- clean up supported repos ([#332](https://github.com/cyralinc/terraform-provider-cyral/pull/332))
- Adding Force New Directives to Access Gateway resource ([#340](https://github.com/cyralinc/terraform-provider-cyral/pull/340))

### Documentation fixes:

- Remove references to deprecated fields mongodb_port_alloc_range_low/high ([#339](https://github.com/cyralinc/terraform-provider-cyral/pull/339))

## 3.0.5 (February 8, 2023)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with control planes `v2.x` and Terraform `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- Fix identity map issues on 3.0-migration script ([#353](https://github.com/cyralinc/terraform-provider-cyral/pull/353))
- Removed mentions of cassandra and bigquery as they are not supported as repo types ([#332](https://github.com/cyralinc/terraform-provider-cyral/pull/332))

### Improvements:

- Improved migration script runtime by using less pipes ([#324](https://github.com/cyralinc/terraform-provider-cyral/pull/324))

## 3.0.4 (November 18, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with control planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- Modified the script to call terraform show less ([#320](https://github.com/cyralinc/terraform-provider-cyral/pull/320))

## 3.0.3 (November 14, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with control planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- Deprecate & Update Cyral Permissions in for Cyral Roles ([#313](https://github.com/cyralinc/terraform-provider-cyral/pull/313))
- Fix sidecar resource nil properties issue ([#317](https://github.com/cyralinc/terraform-provider-cyral/pull/317))

## 3.0.2 (November 7, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with control planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Bug fixes:

- Remove rewrite_on_violation, add enable_dataset_rewrites ([#308](https://github.com/cyralinc/terraform-provider-cyral/pull/308))
- Add support for linux sidecar deployment ([#311](https://github.com/cyralinc/terraform-provider-cyral/pull/311))

### Documentation fixes:

- Fix typo edit in migration script ([#307](https://github.com/cyralinc/terraform-provider-cyral/pull/307))

### Security fixes:

- Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#306](https://github.com/cyralinc/terraform-provider-cyral/pull/306))

## 3.0.1 (October 18, 2022)

Minimum required Control Plane version: `v3.0.0`.

See the list of incompatible resources with control planes `v2.x` and provider `v2.x` in the [`3.0.0`](#300-october-5-2022) release documentation.

### Features:

- Add Cyral Terraform 3.0 Migration Guide ([#303](https://github.com/cyralinc/terraform-provider-cyral/pull/303))

### Bug fixes:

- Policy rule identities field should be omitted by default ([#301](https://github.com/cyralinc/terraform-provider-cyral/pull/301))

## 3.0.0 (October 5, 2022)

Minimum required Control Plane version: `v3.0.0`.

Checkout the [v3 Migration Guide](https://registry.terraform.io/providers/cyralinc/cyral/latest/docs/guides/3.0-migration-guide)
if you are upgrading from provider versions `v2`. It can save you a lot of time.

Resources incompatible with control planes `v2.x`:
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

- Add resource AccessRules ([#280](https://github.com/cyralinc/terraform-provider-cyral/pull/280))
- User Account resource for Gatekeeper Project ([#288](https://github.com/cyralinc/terraform-provider-cyral/pull/288))
- Adapt guides to v3 ([#297](https://github.com/cyralinc/terraform-provider-cyral/pull/297))

### Backwards compatibility breaks:

- Remove deprecated resources for major version 3 ([#291](https://github.com/cyralinc/terraform-provider-cyral/pull/291))
- Remove resources for local_accounts and identity_maps ([#293](https://github.com/cyralinc/terraform-provider-cyral/pull/293))

### Bug fixes:

- Fix acceptance tests for network access policy ([#289](https://github.com/cyralinc/terraform-provider-cyral/pull/289))
- Fix race condition in user accounts test ([#295](https://github.com/cyralinc/terraform-provider-cyral/pull/295))

## 2.11.1 (November 14, 2022)

Minimum required Control Plane version: `v2.35.0`.
Resource incompatible with control planes previous to `v2.35`: argument `enable_dataset_rewrites` from resource `cyral_repository_conf_analysis`.

### Bug fixes:

- Deprecate 'View Sidecars and Repositories' ([#315](https://github.com/cyralinc/terraform-provider-cyral/pull/315))

## 2.11.0 (November 7, 2022)

Minimum required Control Plane version: `v2.35.0`.
Resource incompatible with control planes previous to `v2.35`: argument `enable_dataset_rewrites` from resource `cyral_repository_conf_analysis`.

### Features:

- Deprecate rewrite_on_violation, add enable_dataset_rewrites ([#310](https://github.com/cyralinc/terraform-provider-cyral/pull/310))
- Add support for linux sidecar deployment ([#311](https://github.com/cyralinc/terraform-provider-cyral/pull/311))

## 2.10.2 (October 18, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with control planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Bug fixes:

- Policy rule identities field should be omitted by default ([#302](https://github.com/cyralinc/terraform-provider-cyral/pull/302))

## 2.10.1 (October 6, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with control planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Documentation:

- Use version constraint ~> for guides (v2) ([#298](https://github.com/cyralinc/terraform-provider-cyral/pull/298))

## 2.10.0 (October 5, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with control planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Features:

- Add support for local account automatic approvals ([#296](https://github.com/cyralinc/terraform-provider-cyral/pull/296))

## 2.9.0 (September 16, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with control planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Features:

- Add support for Generic SAML IdP Integration ([#244](https://github.com/cyralinc/terraform-provider-cyral/pull/244))
- Add Duo MFA integration resource ([#282](https://github.com/cyralinc/terraform-provider-cyral/pull/282))
- Add support for repository network shield ([#285](https://github.com/cyralinc/terraform-provider-cyral/pull/285))

### Bug fixes:

- Fix ACC test changes after merging fix ([#275](https://github.com/cyralinc/terraform-provider-cyral/pull/275))
- Fix acc tests and improve sidecar test ([#266](https://github.com/cyralinc/terraform-provider-cyral/pull/266))
- Fix issues with repository resource properties ([#271](https://github.com/cyralinc/terraform-provider-cyral/pull/27`2`))
- Fix issue with repository resource properties ([#272](https://github.com/cyralinc/terraform-provider-cyral/pull/272))

### Improvements:

- Refactor acceptance tests ([#267](https://github.com/cyralinc/terraform-provider-cyral/pull/267))
- Explain setting max value of port range ([#277](https://github.com/cyralinc/terraform-provider-cyral/pull/277))

## 2.8.0 (August 5, 2022)

Minimum required Control Plane version: `v2.34.0`.
Resource incompatible with control planes between `v2.32` and `v2.34`: argument `bypass_mode` from resource `cyral_sidecar`.

### Bug fixes:

- Fix resource import functions ([#260](https://github.com/cyralinc/terraform-provider-cyral/pull/260))
- Fix data map attribute deletion ([#263](https://github.com/cyralinc/terraform-provider-cyral/pull/263))

### Features:

- Add parameter to enable data masking ([#252](https://github.com/cyralinc/terraform-provider-cyral/pull/252))
- Add support for dynamodb ([#253](https://github.com/cyralinc/terraform-provider-cyral/pull/253))
- Support advanced sidecar options ([#255](https://github.com/cyralinc/terraform-provider-cyral/pull/255))
- Add data source to retrieve data labels ([#261](https://github.com/cyralinc/terraform-provider-cyral/pull/261))
- Add single container type to sidecar resource ([#265](https://github.com/cyralinc/terraform-provider-cyral/pull/265))

## 2.7.2 (August 15, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- Fix issues with repository resource properties ([#271](https://github.com/cyralinc/terraform-provider-cyral/pull/271))

## 2.7.1 (July 19, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- Fix token expiration error for large configuration scripts ([#247](https://github.com/cyralinc/terraform-provider-cyral/pull/247))

## 2.7.0 (July 1, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- Fix memory referencing issue and refactor resources ([#223](https://github.com/cyralinc/terraform-provider-cyral/pull/223))

### Features:

- Add data source for cyral_integration_idp ([#239](https://github.com/cyralinc/terraform-provider-cyral/pull/239))
- Add resources for per-repository data map ([#230](https://github.com/cyralinc/terraform-provider-cyral/pull/230))
- Add Terraform guides for basic usages of Cyral features ([#238](https://github.com/cyralinc/terraform-provider-cyral/pull/238))
- Support MongoDB repository replica sets ([#228](https://github.com/cyralinc/terraform-provider-cyral/pull/228))

## 2.6.2 (June 6, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- Fix repository name for IdP Okta module ([#220](https://github.com/cyralinc/terraform-provider-cyral/pull/220))

## 2.6.1 (May 31, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Documentation:

- Fix broken links in documentation ([#216](https://github.com/cyralinc/terraform-provider-cyral/pull/216))

## 2.6.0 (May 27, 2022)

Minimum required Control Plane version: `v2.32.0`.
Resource incompatible with control planes between `v2.29` and `v2.31`: argument `certificate_bundle_secrets` from resource `cyral_sidecar`.

### Bug fixes:

- Improve provider error message when auth fails ([#210](https://github.com/cyralinc/terraform-provider-cyral/pull/210))
- Fix issue with sidecar resource when configured with empty labels ([#210](https://github.com/cyralinc/terraform-provider-cyral/pull/210))
- Fix bug where cyral_sidecar_cft_template data source was crashing when configured with Splunk log integrations ([#215](https://github.com/cyralinc/terraform-provider-cyral/pull/215))

### Features:

- Add Kubernetes Secret support to cyral_repository_local_account resource ([#193](https://github.com/cyralinc/terraform-provider-cyral/pull/193))
- Add GCP Secret Manager support to cyral_repository_local_account resource ([#206](https://github.com/cyralinc/terraform-provider-cyral/pull/206))
- Add labels argument to cyral_repository resource ([#192](https://github.com/cyralinc/terraform-provider-cyral/pull/192))
- Add rate_limit argument to cyral_policy_rule resource ([#194](https://github.com/cyralinc/terraform-provider-cyral/pull/194))
- Add provider optional argument tls_skip_verify to configure TLS verification ([#204](https://github.com/cyralinc/terraform-provider-cyral/pull/204))
- Add certificate_bundle_secrets argument to sidecar resource to support sidecar certificate bundle secrets ([#190](https://github.com/cyralinc/terraform-provider-cyral/pull/190))
- Add cyral_sidecar_bound_ports data source to retrieve sidecar ports that are bound to repositories ([#209](https://github.com/cyralinc/terraform-provider-cyral/pull/209))
- Add cyral_sidecar_id data source to return sidecar ID given a sidecar name ([#211](https://github.com/cyralinc/terraform-provider-cyral/pull/211))
- Add cyral_sidecar_instance_ids data source to return sidecar instance IDs ([#212](https://github.com/cyralinc/terraform-provider-cyral/pull/212))

### Deprecate:

- Deprecate enviroment_variable argument of cyral_repository_local_account resource ([#193](https://github.com/cyralinc/terraform-provider-cyral/pull/193))

### Documentation:

- Update policy docs to inform how to use 'any' for rows argument ([#201](https://github.com/cyralinc/terraform-provider-cyral/pull/201))
- Update policy docs to inform how to use 'any' for data argument ([#203](https://github.com/cyralinc/terraform-provider-cyral/pull/203))
- Update docs to be automatically generated ([#183](https://github.com/cyralinc/terraform-provider-cyral/pull/183))

## 2.5.2 (April 19, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with control planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- Datamap resource is always suggesting an update on terraform plan ([#187](https://github.com/cyralinc/terraform-provider-cyral/pull/187))

## 2.5.1 (March 31, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with control planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- Update hadolint, gorelease and go version ([#181](https://github.com/cyralinc/terraform-provider-cyral/pull/181))

## 2.5.0 (March 31, 2022)

Minimum required Control Plane version: `v2.29.0`.
Resource incompatible with control planes between `v2.25` and `v2.28`: `cyral_integration_pager_duty`.

### Bug fixes:

- Changed listener_port to port so the demo works ([#170](https://github.com/cyralinc/terraform-provider-cyral/pull/170))
- Fix PagerDuty resource due to confExtensions breaking changes ([#173](https://github.com/cyralinc/terraform-provider-cyral/pull/173))

### Features:

- Add Access Gateway support for repo binding resource ([#180](https://github.com/cyralinc/terraform-provider-cyral/pull/180))

### Deprecate:

- Deprecate properties field from policy resource ([#178](https://github.com/cyralinc/terraform-provider-cyral/pull/178))

## 2.4.4 (March 3, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- Added denodo to the validValues map ([#169](https://github.com/cyralinc/terraform-provider-cyral/pull/169))

### Documentation:

- Change cyral_repository to cyral_sidecar for the sidecar_id ([#165](https://github.com/cyralinc/terraform-provider-cyral/pull/165))

### Others:

- ci: adds pre-commit support ([#162](https://github.com/cyralinc/terraform-provider-cyral/pull/162))

## 2.4.3 (February 15, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- Modifying repository_binding resource is not removing old resource ([#159](https://github.com/cyralinc/terraform-provider-cyral/pull/159))

## 2.4.2 (February 3, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Documentation:

- Guide for setting up CP and deploy a sidecar ([#156](https://github.com/cyralinc/terraform-provider-cyral/pull/156))

## 2.4.1 (January 5, 2022)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- Terraform provider creates a new identity map instead of updating old one when state changes ([#152](https://github.com/cyralinc/terraform-provider-cyral/pull/152))
- Fixes issue about missing sidecar returning 500 error ([#148](https://github.com/cyralinc/terraform-provider-cyral/pull/148))

## 2.4.0 (December 8, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Documentation:

- Remove wrong 'okta' reference ([#135](https://github.com/cyralinc/terraform-provider-cyral/pull/135))
- Update examples to use for_each instead of count ([#137](https://github.com/cyralinc/terraform-provider-cyral/pull/137))
- Fix wrong flattened resources in docs ([#142](https://github.com/cyralinc/terraform-provider-cyral/pull/142))

### Features:

- Add CYRAL_TF_CONTROL_PLANE env var support ([#139](https://github.com/cyralinc/terraform-provider-cyral/pull/139))

## 2.3.1 (November 19, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- Fix image links ([#133](https://github.com/cyralinc/terraform-provider-cyral/pull/133))

## 2.3.0 (November 19, 2021)

Minimum required Control Plane version: `v2.25.3`.

### Bug fixes:

- Fix documentation accordingly to Terraform standards ([#131](https://github.com/cyralinc/terraform-provider-cyral/pull/131))

### Features:

- Add resource to manage sso groups to roles ([#106](https://github.com/cyralinc/terraform-provider-cyral/pull/106))
- Add docs reference to Okta IdP module ([#129](https://github.com/cyralinc/terraform-provider-cyral/pull/129))
- Add sidecar user endpoint ([#132](https://github.com/cyralinc/terraform-provider-cyral/pull/132))

## 2.2.1 (November 12, 2021)

Minimum required Control Plane version: `v2.25.0`.

### Bug fixes:

- Fix IdP registration at CP-level ([#128](https://github.com/cyralinc/terraform-provider-cyral/pull/128))

## 2.2.0 (November 4, 2021)

Minimum required Control Plane version: `v2.25.0`.

### Bug fixes:

- Fix cyclic dependency issue in SAML certificate data source ([#121](https://github.com/cyralinc/terraform-provider-cyral/pull/121))

### Deprecated resources:

- `cyral_integration_sso_*` renamed to `cyral_integration_idp_*`

## 2.1.1 (October 21, 2021)

Minimum required Control Plane version: `v2.24.0`.

### Bug fixes:

- Remove unnecessary PreCheck from Terraform Provider Tests ([#117](https://github.com/cyralinc/terraform-provider-cyral/pull/117))

## 2.1.0 (October 18, 2021)

Minimum required Control Plane version: `v2.24.0`.

### Features:

- Generic SAML integration and new SSO resources ([#115](https://github.com/cyralinc/terraform-provider-cyral/pull/115))

## 2.0.1 (September 30, 2021)

Minimum required Control Plane version: `v2.22.0`.

### Bug fixes:

- Omitting access_duration in cyral_identity_map resulted in plan change on every plan ([#111](https://github.com/cyralinc/terraform-provider-cyral/pull/111))

## 2.0.0 (September 24, 2021)

Minimum required Control Plane version: `v2.22.0`.

### Backwards compatibility breaks:

- Resource cyral_sidecar: changed parameters;
- Data source cyral_sidecar_template: data source replaced by `cyral_sidecar_cft_template` and template restricted to Cloudformation.

### Features:

- Script to rotate service account secrets ([#64](https://github.com/cyralinc/terraform-provider-cyral/pull/64))
- Improve tooling ([#92](https://github.com/cyralinc/terraform-provider-cyral/pull/92))
- Resource Sidecar Credentials ([#93](https://github.com/cyralinc/terraform-provider-cyral/pull/93))
- Resource Repository Conf Analysis ([#108](https://github.com/cyralinc/terraform-provider-cyral/pull/108))

## 1.2.2 (June 21, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Bug fixes:

- Fix missing Helm3 support: fix missing helm 3 support in data source `cyral_sidecar_template` ([#63](https://github.com/cyralinc/terraform-provider-cyral/pull/63))

---

## 1.2.1 (June 18, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Bug fixes:

- Fix publishing issue: fix issue publishing binaries ([#62](https://github.com/cyralinc/terraform-provider-cyral/pull/62))

---

## 1.2.0 (June 17, 2021)

Minimum required Control Plane version: `v2.19.0`.

### Features:

- Data Source SAML Certificate: added new data source to retrieve certificate ([#60](https://github.com/cyralinc/terraform-provider-cyral/pull/60))
- Docker Build: added new feature in makefile to allow docker build without requiring local Go environment ([#61](https://github.com/cyralinc/terraform-provider-cyral/pull/61))
- Resource Integration Hashicorp Vault: added new resource to support Hashicorp Vault integration ([#58](https://github.com/cyralinc/terraform-provider-cyral/pull/58))
- Resource Integration Pager Duty: added new resource to support Pager Duty integration ([#56](https://github.com/cyralinc/terraform-provider-cyral/pull/56))

---

## 1.1.0 (April 23, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:

- Resource Integration Slack Alerts: added new resource to support Slack Alerts integration ([#43](https://github.com/cyralinc/terraform-provider-cyral/pull/43))
- Resource Integration Microsoft Teams: added new resource to support Microsoft Teams integration ([#44](https://github.com/cyralinc/terraform-provider-cyral/pull/44))
- Resource Repository Local Account: added new resource to support repositories local accounts ([#46](https://github.com/cyralinc/terraform-provider-cyral/pull/46))
- Resource Integration Okta: added new resource to support Okta integration ([#47](https://github.com/cyralinc/terraform-provider-cyral/pull/47))
- Resource Repository Configuration Authentication: added new resource to support repository configuration authentication ([#48](https://github.com/cyralinc/terraform-provider-cyral/pull/48))
- Data Source Sidecar Template: added new data source to support sidecar templates ([#50](https://github.com/cyralinc/terraform-provider-cyral/pull/50))
- Resource Identity Map: added new resource to support identity maps ([#51](https://github.com/cyralinc/terraform-provider-cyral/pull/51))

### Improvements:

- Increase Test Coverage: defined standards and increased the test coverage. ([#41](https://github.com/cyralinc/terraform-provider-cyral/pull/41))

---

## 1.0.0 (March 26, 2021)

Minimum required Control Plane version: `v2.17.0`.

### Features:

- Resource Sidecar: added new resource to support sidecars ([#23](https://github.com/cyralinc/terraform-provider-cyral/pull/23))
- Resource Datamap: added new resource to support datamaps ([#24](https://github.com/cyralinc/terraform-provider-cyral/pull/24))
- Resource Repository Binding: added new resource to support binding repositories to sidecars ([#25](https://github.com/cyralinc/terraform-provider-cyral/pull/25))
- Resource Policy: added new resource to support policies ([#26](https://github.com/cyralinc/terraform-provider-cyral/pull/26))
- Resource Policy Rule: added new resource to support policy rules ([#28](https://github.com/cyralinc/terraform-provider-cyral/pull/28))
- Resource: Datadog Integration: added new resource to support Datadog integration ([#30](https://github.com/cyralinc/terraform-provider-cyral/pull/30))
- Resource ELK Integration: added new resource to support ELK integration ([#31](https://github.com/cyralinc/terraform-provider-cyral/pull/31))
- Resource Splunk Integration: added new resource to support Splunk integration ([#33](https://github.com/cyralinc/terraform-provider-cyral/pull/33))
- Resource Sumo Logic Integration: added new resource to support Sumo Logic integration ([#34](https://github.com/cyralinc/terraform-provider-cyral/pull/34))
- Resource Logstash Integration: added new resource to support Logstash integration ([#35](https://github.com/cyralinc/terraform-provider-cyral/pull/35))
- Resource Looker Integration: added new resource to support Looker integration ([#36](https://github.com/cyralinc/terraform-provider-cyral/pull/36)).

### Improvements:

- Replace repo ID property: replaced repository identification in the state file from `name` to `id`. ([#22](https://github.com/cyralinc/terraform-provider-cyral/pull/22))

---

## 0.2.0 (February 24, 2021)

### Features:

- Terraform Plugin SDK v2: migration to `Terraform Plugin SDK v2` and support for Terraform `v0.12` to `v0.14` ([#21](https://github.com/cyralinc/terraform-provider-cyral/pull/21)).

### Improvements:

- API Port: Removed `control_plane_api_port` and added this information to the existing `control_plane` parameter. ([#20](https://github.com/cyralinc/terraform-provider-cyral/pull/20))
- Unit tests: added unit tests for `config.go` ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8))

---

## 0.1.0 (May 15, 2020)

### Features:

- Terraform import: added support for `terraform import` statement for `repository` resource ([#8](https://github.com/cyralinc/terraform-provider-cyral/pull/8)).

---

## 0.0.2 (May 2, 2020)

### Bug fixes:

- Error handling: fixed bug in error handling ([#5](https://github.com/cyralinc/terraform-provider-cyral/pull/5)).

---

## 0.0.1 (May 2, 2020)

### Features:

- Provider with Repository: draft `provider` and `repository` resource ([#1](https://github.com/cyralinc/terraform-provider-cyral/pull/1), ([#2](https://github.com/cyralinc/terraform-provider-cyral/pull/2))).
