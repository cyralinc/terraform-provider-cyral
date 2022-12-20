#!/bin/bash

# Set the color variable
red='\033[0;31m'
# Clear the color after that
clear='\033[0m'

if ! command -v terraform &> /dev/null
then
    echo "The Terraform CLI must be installed for this script to run."
    echo "Instructions for installation can be found here:"
    echo "https://learn.hashicorp.com/tutorials/terraform/install-cli"
    exit
fi

if ! command -v jq &> /dev/null
then
    echo "The tool jq must be installed for this script to run."
    echo "Instructions for installation can be found here:"
    echo "https://stedolan.github.io/jq/download/"
    exit
fi

echo "Welcome to Cyral's Terraform Provider Version 4.0 Migration script!"
echo
echo \
"This script will create new resource definitions for the
cyral_repository and cyral_repository_binding resource that
will be migrated. Additionally, cyral_sidecar listeners
will be created, which are now required to bind sidecars to
repositories."
echo
echo "Please set CYRAL_TF_FILE_PATH equal to the file path of your .tf file."
echo
read -p "Are you ready to continue? [N/y] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi
echo
echo "Searching for cyral_repository and cyral_repository_binding resources to migrate..."
echo

terraform state pull > terraform.tfstate.cyral.migration.backup
cp ${CYRAL_TF_FILE_PATH} cyral_terraform_migration_backup_configuration.txt

repo_import_args=()
binding_import_args=()
listener_import_args=()

repo_resource_defs=()
binding_resource_defs=()
listener_resource_defs=()

repo_resource_names=()
binding_resource_names=()
listener_resource_names=()

repos_to_delete=()
bindings_to_delete=()

access_gateway_resources=()

user_account_import_ids=()
access_rule_import_ids=()

user_account_resource_defs=()
access_rule_resource_defs=()

user_account_resource_names=()
access_rule_resource_names=()

local_accounts_to_delete=()
identity_maps_to_delete=()

empty_access_duration="$(printf '%s' '[]')"

# Create array of repo and binding resource to migrate
resources_to_migrate=($(terraform state list | grep "cyral_repository\|cyral_repository_binding\cyral_repository_local_account\|cyral_repository_identity_map"))
# Store terraform state JSON representation
tf_state_json=$(terraform show -json | jq ".values.root_module.resources[]")

# Find all cyral_repository and cyral_repository_local_accounts
for resource_address in ${resources_to_migrate[@]}; do
  if [[ $resource_address == cyral_repository.* ]]
  then
    # We will need to delete this repo from the .tf file, so we
    # store the full resource address.
    repos_to_delete+=($resource_address)
    # Get repo ID. This will be used to import the updated repo.
    repo_id=($(jq -r "select(.address == \"$resource_address\") | .values.id"<<<$tf_state_json))
    # Remove [] from the resource address and substitute [ for _
    # E.g. cyral_repository.repo[0] -> cyral_repository.repo_0
    resource_address=$(sed -e 's/[]]//g;s/[[]/_/g'<<<$resource_address)
    # Save an empty resource definition for a new repo, so that
    # it can be added to the .tf file.
    repo_resource_defs+=("resource \"cyral_repository\" \"${resource_address##"cyral_repository."}\" {}")
    # Store import resource address and repo ID, which will be the argument to terraform import.
    repo_import_args+=("${resource_address} ${repo_id}")
    repo_resource_names+=("${resource_address}")
  elif [[ $resource_address == cyral_repository_binding.* ]]
  then
    # We will need to delete this sidecar binding from the .tf file, store its name
    bindings_to_delete+=($resource_address)
    # Get ids required to import binding resource.
    id_values_arr=($(jq -r "select(.address == \"$resource_address\") | .values.sidecar_id, .values.repository_id, .values.sidecar_as_idp_access_gateway, .values.id"<<<$tf_state_json))
    # Construct import ID for the repository binding that was migrated in CP.
    import_id="${id_values_arr[0]}/${id_values_arr[1]}"
    if [[id_values_arr[2]]]
        access_gateway_resource_name=${resource_address##"cyral_repository_binding."}-access-gateway
        access_gateway_resources+=("resource \"cyral_repository_access_gateway\" \"${access_gateway_resource_name}\" {
            repository_id = ${id_values_arr[1]}
            sidecar_id = ${id_values_arr[0]}
            binding_id = ${id_values_arr[3]}
        }")
    fi
    # Remove [] from the resource address and substitute [ for _
    binding_resource_address=$(sed -e 's/[]]//g;s/[[]/_/g'<<<$resource_address)
    # Save empty resource definition for the binding, so that it can be added to the .tf file
    binding_resource_defs+=("resource \"cyral_repository_binding\" \"${resource_address##"cyral_repository_binding."}\" {}")
    # Store import argument and name for binding
    binding_import_args+=("${resource_address} ${import_id}")
    binding_resource_names+=("${resource_address}")
  fi
    if [[ $resource_address == cyral_repository_identity_map.* ]]
  then
    # We will need to delete this identity map from the .tf file, store its name
    identity_maps_to_delete+=($resource_address)
    # Get repo ID, local account ID, identity_type and access_duration for the identity map.
    values_arr=($(jq -r "select(.address == \"$resource_address\") | .values.repository_id, .values.repository_local_account_id, .values.identity_type, .values.access_duration"<<<$tf_json))
    if [[ ${values_arr[3]} != $empty_access_duration ]] && [[ ${values_arr[2]} == "user" ]]; then
        # Identity map was migrated to be an approval, which is not managed through terraform-- do nothing.
        continue
    fi
    # Construct import ID for the access rule that was migrated from this identity map.
    import_id="${values_arr[0]}/${values_arr[1]}"
    # Remove [] from the resource as they are not supported and substitute [ for _
    resource_address=$(sed -e 's/[]]//g;s/[[]/_/g'<<<$resource_address)
    # Construct name of the access rule that will be imported.
    import_name=cyral_repository_access_rules.${resource_address##"cyral_repository_identity_map."}
    # Save name of the new access rule, so that it can be added to the .tf file
    access_rule_resource_defs+=("resource \"cyral_repository_access_rules\" \"${resource_address##"cyral_repository_identity_map."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    access_rule_import_ids+=("${import_kv_pair}")
    access_rule_resource_names+=("${import_name}")
  elif [[ $resource_address == cyral_repository_local_account.* ]]
  then
    # We will need to delete this local account from the .tf file, store its name
    local_accounts_to_delete+=($resource_address)
    # Get local account ID for the local account.
    values_arr=($(jq -r "select(.address == \"$resource_address\") | .values.repository_id, .values.id"<<<$tf_json))
    # Construct import ID for the user account that was migrated from this local account.
    import_id="${values_arr[0]}/${values_arr[1]}"
    # Remove [] from the resource as they are not supported and substitute [ for _
    resource_address=$(sed -e 's/[]]//g;s/[[]/_/g'<<<$resource_address)
    # Construct name of the user account that will be imported.
    import_name=cyral_repository_user_account.${resource_address##"cyral_repository_local_account."}
    # Save name of the migrated user account, so that it can be added to the .tf file
    user_account_resource_defs+=("resource \"cyral_repository_user_account\" \"${resource_address##"cyral_repository_local_account."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    user_account_import_ids+=("${import_kv_pair}")
    user_account_resource_names+=("${import_name}")
  fi
done

echo "Found ${#repos_to_delete[@]} cyral_repository resources to migrate."
echo "Found ${#bindings_to_delete[@]} cyral_repository_bindings resources to migrate."
echo "Found ${#local_accounts_to_delete[@]} cyral_repository_local_accounts to migrate to cyral_repository_user_accounts."
echo "Found ${#identity_maps_to_delete[@]} cyral_repository_identity_maps to migrate to cyral_repository_access_rules."
echo
echo "The following file path was provided for CYRAL_TF_FILE_PATH: ${CYRAL_TF_FILE_PATH}"
read -p "Would you like this script to append these lines to the .tf file ${CYRAL_TF_FILE_PATH}? [N/y] " -n 1 -r
if [[  $REPLY =~ ^[Yy]$ ]]
then
    printf '\n\n' >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${binding_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${repo_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${user_account_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${access_rule_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo; echo; echo;
echo "Now its time to upgrade your Cyral Terraform Provider to version 4!"
echo
echo -e "Before we proceed, you will need to do the following:
    1.  Open your Terraform .tf configuration file.
    2.  Change the version number of the cyral provider in the required_providers
        section of your .tf configuration file to '~>4.0'. It should look like this:
            cyral = {
                source  = \"cyralinc/cyral\"
                version = \"~>4.0\"
            }
    3. Ensure that new empty resource definitions were added to the end of
        your .tf file. The definitions will look like this:
            Repositories
            resource \"cyral_repository\" \"<resource_name>\" {}

            Respository Bindings
            resource \"cyral_respository_binding\" \"<resource_name>\" {}

            User Account
            resource \"cyral_repository_user_account\" \"<resource_name>\" {}

            Access Rules
            resource \"cyral_repository_access_rules\" \"<resource_name>\" {}
    ************************************************
    *                                              *
    *        ${red}IMPORTANT STEP PLEASE DONT SKIP${clear}       *
    *                                              *
    ************************************************
    4.  Find all references to cyral_repository, cyral_repository_binding,
        cyral_repository_identity_map and cyral_repository_local_account
        resources in your .tf file and remove the entire resource definition
        for each one."
echo
read -p "Are you ready to upgrade Terraform? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

terraform init -upgrade

echo
echo "Removing the following cyral_repository resources from your tf state:"
printf '%s\n' "${repos_to_delete[@]}"
echo
echo "Removing the following cyral_sidecar_binding resources from your tf state:"
printf '%s\n' "${bindings_to_delete[@]}"
echo

for repo in ${repos_to_delete[@]};do
  terraform state rm $repo
done

for binding in ${bindings_to_delete[@]};do
  terraform state rm $binding
done

echo
echo "Removing the following cyral_repository_local_accounts:"
printf '%s\n' "${local_accounts_to_delete[@]}"
echo
echo "Removing the following cyral_repository_identity_maps:"
printf '%s\n' "${identity_maps_to_delete[@]}"
echo

for local_account in ${local_accounts_to_delete[@]};do
    terraform state rm $local_account
done

for identity_map in ${identity_maps_to_delete[@]};do
    terraform state rm $identity_map
done

echo
echo "Importing the following cyral_repository into your Terraform state:"
printf '%s\n' "${repo_resource_names[@]}"
echo
echo "Importing the following cyral_repository_access_rules into your Terraform state:"
printf '%s\n' "${binding_resource_names[@]}"
echo

for repo in "${repo_import_args[@]}";do
    terraform import $repo
done
for binding in "${binding_import_args[@]}";do
    terraform import $binding
done

echo
echo "Importing the following cyral_repository_user_accounts into your Terraform state:"
printf '%s\n' "${user_account_resource_names[@]}"
echo
echo "Importing the following cyral_repository_access_rules into your Terraform state:"
printf '%s\n' "${access_rule_resource_names[@]}"
echo

for user_account_id in ${user_account_import_ids[@]};do
    terraform import $user_account_id
done
for access_rule_id in ${access_rule_import_ids[@]};do
    terraform import $access_rule_id
done

# Only after we have imported the migrated bindings, we have access to the
# listener_ids that were created during CP migration. As a result, we need
# to repeat the entire process, except this time we are only importing
# listeners. Since listeners are an entirely new resource type, there is
# no need to remove them from the terraform state.

# Create array of migrated bindings
migrated_bindings=($(terraform state list | grep "cyral_repository_binding"))
# Store terraform state JSON representation
tf_state_json=$(terraform show -json | jq ".values.root_module.resources[]")

for binding in ${migrated_bindings[@]}; do
  # Get ids required to import listner resource.
  id_values_arr=($(jq -r "select(.address == \"$binding\") | .values.sidecar_id, .values.listener_binding[0].listener_id"<<<$tf_state_json))
  # Construct import ID for the listener that was created during CP migration.
  import_id="${id_values_arr[0]}/${id_values_arr[1]}"
  # Save empty resource definition for the listener, so that it can be added to the .tf file
  listener_resource_name=${binding##"cyral_repository_binding."}_listener
  listener_resource_defs+=("resource \"cyral_sidecar_listener\" \"${listener_resource_name}\" {}")
  # Store import argument and name for listener
  listener_import_args+=("${listener_resource_name} ${import_id}")
  listener_resource_names+=("${listener_resource_name}")
done

echo "Found ${#listener_resource_names[@]} cyral_sidecar_listener resources to import."
echo
echo "The following file path was provided for CYRAL_TF_FILE_PATH: ${CYRAL_TF_FILE_PATH}"
read -p "Would you like this script to append these lines to the .tf file ${CYRAL_TF_FILE_PATH}? [N/y] " -n 1 -r
if [[  $REPLY =~ ^[Yy]$ ]]
then
    printf '\n\n' >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${listener_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo "Importing the following cyral_sidecar_listeners into your Terraform state:"
printf '%s\n' "${listener_resource_names[@]}"
echo

for listener in "${listener_import_args[@]}";do
    terraform import $listener
done

# Once resources have been imported, we need to replace the empty resource definitions
# (which are required to exist prior to import), with resource definitions containing
# the correct values. We do this by printing the resource definition from the terraform
# state, then remove any computed values. Finally, we append the new resource definitions
# to a seperate .tf file, containing all of the resources that were migrated.
for binding in ${binding_resource_names[@]};do
    terraform state show -no-color $binding | grep -v "   binding_id" | grep -v "   id " >> cyral_migration_repositories_bindings_listeners.txt
done
for repo in ${repo_resource_names[@]};do
    terraform state show -no-color $repo | grep -v "   id " >> cyral_migration_repositories_bindings_listeners.txt
done
for listener in ${listener_resource_names[@]};do
    terraform state show -no-color $listener | grep -v "   listener_id " >> cyral_migration_repositories_bindings_listeners.txt
done
for access_gateway in ${access_gateway_resources[@]};do
    echo access_gateway >> cyral_migration_repositories_bindings_listeners.txt
done

for user_account in ${user_account_resource_names[@]};do
    terraform state show -no-color $user_account | grep -v "   user_account_id" | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
done
for access_rule in ${access_rule_resource_names[@]};do
    terraform state show -no-color $access_rule | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
done

mv cyral_migration_repository_access_rules_and_user_accounts.txt cyral_migration_repository_access_rules_and_user_accounts.tf
terraform fmt

mv cyral_migration_repositories_bindings_listeners.txt cyral_migration_repositories_bindings_listeners.tf
terraform fmt

echo; echo; echo;
echo "Now that the Terraform state is up-to-date, let's clean up your .tf file."
echo "This script created a new .tf file containing the resources definitions"
echo "for the cyral_repository_access_rules and cyral_repository_user_accounts"
echo "that were migrated into your Terraform state. The new .tf file is called:"
echo
echo "cyral_migration_repository_access_rules_and_user_accounts.tf"
echo
echo
echo "It is finally time to remove the empty resources from your .tf files."
echo "Please perform the following actions: "
echo
echo "  1.  Remove the empty resource definitions for the "
echo "      cyral_repository_access_rules and cyral_repository_user_accounts"
echo "      that were added to the end of your .tf file, which is named:"
echo "      ${CYRAL_TF_FILE_PATH}"
echo "  2.  Remove the empty resource definitions for the "
echo "      cyral_repository_access_rules and cyral_repository_user_accounts"
echo "      that were added to the end of your .tf file, which is named:"
echo "      ${CYRAL_TF_FILE_PATH}"
echo
echo "When you are done, run the following command:"
echo "terraform plan"
echo
echo
echo "If migration was successful, you should see the following message:"
echo "No changes. Your infrastructure matches the configuration."
echo
echo "If migration was not successful, you have the option to revert to your"
echo "previous Terraform state and try again."
echo
read -p "Would you revert to your previous state and try the migration again? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo "Your previous .tf file was copied before the migration was ran. It is called "
echo "cyral_terraform_migration_backup_configuration.txt"
echo
echo "Please perform the following action before proceding:"
echo "  1.  Replace the contents of your .tf file ${CYRAL_TF_FILE_PATH} "
echo "      with the contents of cyral_terraform_migration_backup_configuration.txt. "
echo "  2.  Delete the following files that were created by the script: "
echo "      - cyral_terraform_migration_backup_configuration.txt"
echo "      - cyral_migration_repository_access_rules_and_user_accounts.tf"
echo "      - cyral_migration_repositories_bindings_listeners.tf"
echo
read -p "Are you ready to revert to your pre-migration Terraform state? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

# Downgrade Terraform
terraform init -upgrade

# Replace current state with backup state.
mv 'terraform.tfstate.cyral.migration.backup' 'terraform.tfstate'

# Revert to old state.
terraform state push 'terraform.tfstate'
