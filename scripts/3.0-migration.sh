#!/bin/bash

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

echo "Welcome to Cyral's Terraform Provider Version 3 Migration script!"
echo
echo "This script will create new resource definitions for the"
echo "cyral_repository_user_accounts and cyral_repository_access_rules that will "
echo "be migrated."
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
echo "Searching for cyral_repository_identity_map and cyral_repository_local_account resources to migrate..."
echo

terraform state pull > terraform.tfstate.cyral.migration.backup
cp ${CYRAL_TF_FILE_PATH} cyral_terraform_migration_backup_configuration.txt

user_account_import_ids=()
access_rule_import_ids=()

user_account_resource_defs=()
access_rule_resource_defs=()

user_account_resource_names=()
access_rule_resource_names=()

local_accounts_to_delete=()
identity_maps_to_delete=()

empty_access_duration="$(printf '%s' '[]')"

# Create array of all resources in the terraform state
IFS=$'\n' read -r -d '' -a tf_state < <( terraform state list && printf '\0' )

# Find all cyral_repository_identity_maps and cyral_repository_local_accounts
for resource in ${tf_state[@]}; do
  if [[ $resource == cyral_repository_identity_map.* ]]
  then
    # We will need to delete this identity map from the .tf file, store its name
    identity_maps_to_delete+=($resource)
    # Get repo ID and local account ID for the identity map.
    repo_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_id" | sed 's/"//g' )
    local_account_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_local_account_id" | sed 's/"//g')
    identity_type=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.identity_type" | sed 's/"//g')
    access_duration=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.access_duration")
    if [[ $access_duration != $empty_access_duration ]] && [[ $identity_type == "user" ]]; then
        # Identity map was migrated to be an approval, which is not managed through terraform-- do nothing.
        continue
    fi
    # Construct import ID for the access rule that was migrated from this identity map.
    import_id="$repo_id/$local_account_id"
    # Construct name of the access rule that will be imported.
    import_name=cyral_repository_access_rules.${resource##"cyral_repository_identity_map."}
    # Save name of the new access rule, so that it can be added to the .tf file
    access_rule_resource_defs+=("resource \"cyral_repository_access_rules\" \"${resource##"cyral_repository_identity_map."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    access_rule_import_ids+=("${import_kv_pair}")
    access_rule_resource_names+=("${import_name}")
  elif [[ $resource == cyral_repository_local_account.* ]]
  then
    # We will need to delete this local account from the .tf file, store its name
    local_accounts_to_delete+=($resource)
    # Get repo ID and local account ID for the local account.
    repo_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_id" | sed 's/"//g' )
    local_account_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.id" | sed 's/"//g')
    # Construct import ID for the user account that was migrated from this local account.
    import_id="$repo_id/$local_account_id"
    # Construct name of the user account that will be imported.
    import_name=cyral_repository_user_account.${resource##"cyral_repository_local_account."}
    # Save name of the migrated user account, so that it can be added to the .tf file
    user_account_resource_defs+=("resource \"cyral_repository_user_account\" \"${resource##"cyral_repository_local_account."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    user_account_import_ids+=("${import_kv_pair}")
    user_account_resource_names+=("${import_name}")
  fi
done

echo "Found ${#local_accounts_to_delete[@]} cyral_repository_local_accounts to migrate to cyral_repository_user_accounts."
echo "Found ${#identity_maps_to_delete[@]} cyral_repository_identity_maps to migrate to cyral_repository_access_rules."
echo
echo "The following file path was provided for CYRAL_TF_FILE_PATH: ${CYRAL_TF_FILE_PATH}"
read -p "Would you like this script to append these lines to the .tf file ${CYRAL_TF_FILE_PATH}? [N/y] " -n 1 -r
if [[  $REPLY =~ ^[Yy]$ ]]
then
    printf '\n\n' >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${user_account_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${access_rule_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo; echo; echo;
echo "Now its time to upgrade your Cyral Terraform Provider to version 3!"
echo
echo "Before we proceed, you will need to do the following:
    1.  Open your Terraform .tf configuration file.
    2.  Change the version number of the cyral provider in the required_providers
        section of your .tf configuration file to 3.0.0. It should look like this:
            cyral = {
                source  = \"cyralinc/cyral\"
                version = \"3.0.0\"
            }
    3. Ensure that new empty resource definitions were added to the end of
        your .tf file. The definitions will look like this:
            User Account
            resource \"cyral_repository_user_account\" \"<resource_name>\" {}

            Access Rules
            resource \"cyral_repository_access_rules\" \"<resource_name>\" {}"
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
echo "Importing the following cyral_repository_user_accounts into your Terraform state:"
printf '%s\n' "${user_account_resource_names[@]}"
echo
echo "Importing the following cyral_repository_access_rules into your Terraform state:"
printf '%s\n' "${access_rule_resource_names[@]}"
echo

for ((i = 0; i < ${#user_account_import_ids[@]}; i++));do
    terraform import ${user_account_import_ids[$i]}
done
for ((i = 0; i < ${#access_rule_import_ids[@]}; i++));do
    terraform import ${access_rule_import_ids[$i]}
done

echo
echo "Removing the following cyral_repository_local_accounts:"
printf '%s\n' "${local_accounts_to_delete[@]}"
echo
echo "Removing the following cyral_repository_identity_maps:"
printf '%s\n' "${identity_maps_to_delete[@]}"
echo

for ((i = 0; i < ${#local_accounts_to_delete[@]}; i++));do
    terraform state rm ${local_accounts_to_delete[$i]}
done

for ((i = 0; i < ${#identity_maps_to_delete[@]}; i++));do
    terraform state rm ${identity_maps_to_delete[$i]}
done

for ((i = 0; i < ${#user_account_resource_names[@]}; i++));do
    terraform state show -no-color ${user_account_resource_names[$i]} | grep -v "   user_account_id" | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
done
for ((i = 0; i < ${#access_rule_resource_names[@]}; i++));do
    terraform state show -no-color ${access_rule_resource_names[$i]} | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
done

mv cyral_migration_repository_access_rules_and_user_accounts.txt cyral_migration_repository_access_rules_and_user_accounts.tf
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
echo "It is finally time to remove the old resources from your .tf files."
echo "Please perform the following actions: "
echo
echo "  1.  Find all references to cyral_repository_identity_maps and"
echo "      cyral_repository_local_accounts in your .tf file and remove the"
echo "      entire resource defintion for each one."
echo "  2.  Remove the empty resource definitions for the "
echo "      cyral_repository_access_rules and cyral_repository_user_accounts"
echo "      that were added to the end of your .tf file, which is named:"
echo "      ${CYRAL_TF_FILE_PATH}"
echo
echo
echo "When you are done, run the following command:"
echo "terraform plan"
echo
echo
echo "If migration was successful, you should see the following message:"
echo "No changes. Your infrastructure matches the configuration."
echo
echo "If migration was not successful, you have to option to revert to your"
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
