#!/usr/bin/env bash

# Set the color variable
red='\033[0;31m'
# Clear the color after that
clear='\033[0m'

if [ ${BASH_VERSION:0:1} \< 4 ]
then
    echo "Bash version 4 or higher is required by this script."
    echo "Please install the latest bash version and ensure your"
    echo "BASH_VERSION environmental variable is set correctly."
    exit
fi

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

user_account_import_args=()
declare -A access_rules_resource_id_to_import_args=()

user_account_resource_defs=()
declare -A access_rules_resource_id_to_defs=()

user_account_resource_addresses=()
declare -A access_rules_resouce_id_to_address=()

local_accounts_to_delete=()
identity_maps_to_delete=()

empty_access_duration="$(printf '%s' '[]')"

# Create array of all resources in the terraform state
tf_state=($(terraform state list | grep "cyral_repository_local_account\|cyral_repository_identity_map"))
tf_json=$(terraform show -json | jq ".values.root_module.resources[]")

# Find all cyral_repository_identity_maps and cyral_repository_local_accounts
for resource_address in ${tf_state[@]}; do
  if [[ $resource_address == cyral_repository_identity_map.* ]]
  then
    # We will need to delete this identity map from the .tf file, store its name
    identity_maps_to_delete+=($resource_address)
    # Escape the double quotes so we find it using jq
    escaped_resource_address=$(sed -e 's/\"/\\"/g'<<<$resource_address)
    # Get repo ID, local account ID, identity_type and access_duration for the identity map.
    values_arr=($(jq -r "select(.address == \"$escaped_resource_address\") | .values.repository_id, .values.repository_local_account_id, .values.identity_type, .values.access_duration"<<<$tf_json))
    if [[ ${values_arr[3]} != $empty_access_duration ]] && [[ ${values_arr[2]} == "user" ]]; then
        # Identity map was migrated to be an approval, which is not managed through terraform-- do nothing.
        continue
    fi
    # Remove [] and \" from the resource as they are not supported and substitute [ for _
    resource_address=$(sed -e 's/[]]//g;s/[[]/_/g;s/\\"//g'<<<$escaped_resource_address)
    # Construct name of the access rule that will be imported.
    access_rules_resouce_name=${resource_address##"cyral_repository_identity_map."}
    access_rules_resouce_address=cyral_repository_access_rules.${access_rules_resouce_name}
    # Construct import ID for the access rule that was migrated from this identity map.
    resource_id="${values_arr[0]}/${values_arr[1]}"
    # Store import name and ID as a key value pair
    import_args="${access_rules_resouce_address} ${resource_id}"
    access_rules_resource_id_to_import_args[${resource_id}]="${import_args}"
    # Save name of the new access rule, so that it can be added to the .tf file
    access_rules_resource_id_to_defs[${resource_id}]="resource \"cyral_repository_access_rules\" \"${access_rules_resouce_name}\" {}"
    access_rules_resouce_id_to_address[${resource_id}]="${access_rules_resouce_address}"
  elif [[ $resource_address == cyral_repository_local_account.* ]]
  then
    # We will need to delete this local account from the .tf file, store its name
    local_accounts_to_delete+=($resource_address)
    # Escape the double quotes so we find it using jq
    escaped_resource_address=$(sed -e 's/\"/\\"/g'<<<$resource_address)
    # Get local account ID for the local account.
    values_arr=($(jq -r "select(.address == \"$escaped_resource_address\") | .values.repository_id, .values.id"<<<$tf_json))
    # Construct import ID for the user account that was migrated from this local account.
    resource_id="${values_arr[0]}/${values_arr[1]}"
    # Remove [] and \" from the resource as they are not supported and substitute [ for _
    resource_address=$(sed -e 's/[]]//g;s/[[]/_/g;s/\\"//g'<<<$escaped_resource_address)
    # Construct name of the user account that will be imported.
    user_account_resource_name=${resource_address##"cyral_repository_local_account."}
    user_account_resource_address=cyral_repository_user_account.${user_account_resource_name}
    # Save name of the migrated user account, so that it can be added to the .tf file
    user_account_resource_defs+=("resource \"cyral_repository_user_account\" \"${user_account_resource_name}\" {}")
    # Store import name and ID as a key value pair
    import_args="${user_account_resource_address} ${resource_id}"
    user_account_import_args+=("${import_args}")
    user_account_resource_addresses+=("${user_account_resource_address}")
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
    printf '%s\n\n' "${access_rules_resource_id_to_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo; echo; echo;
echo "Now its time to upgrade your Cyral Terraform Provider to version 3!"
echo
echo -e "Before we proceed, you will need to do the following:
    1.  Open your Terraform .tf configuration file.
    2.  Change the version number of the cyral provider in the required_providers
        section of your .tf configuration file to '~>3.0'. It should look like this:
            cyral = {
                source  = \"cyralinc/cyral\"
                version = \"~>3.0\"
            }
    3. Ensure that new empty resource definitions were added to the end of
        your .tf file. The definitions will look like this:
            User Account
            resource \"cyral_repository_user_account\" \"<resource_name>\" {}

            Access Rules
            resource \"cyral_repository_access_rules\" \"<resource_name>\" {}
    ************************************************
    *                                              *
    *        ${red}IMPORTANT STEP PLEASE DONT SKIP${clear}       *
    *                                              *
    ************************************************
    4.  Find all references to cyral_repository_identity_map and
        cyral_repository_local_account in your .tf file and remove the
        entire resource definition for each one."
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
printf '%s\n' "${user_account_resource_addresses[@]}"
echo
echo "Importing the following cyral_repository_access_rules into your Terraform state:"
printf '%s\n' "${access_rules_resouce_id_to_address[@]}"
echo

for user_account_id in "${user_account_import_args[@]}";do
    terraform import $user_account_id
done
for access_rule_id in "${access_rules_resource_id_to_import_args[@]}";do
    terraform import $access_rule_id
done

local_accounts="[$(IFS=" "; echo "${local_accounts_to_delete[*]}")]"
terraform state rm ${local_accounts:1:${#local_accounts}-2}

identity_maps="[$(IFS=" "; echo "${identity_maps_to_delete[*]}")]"
terraform state rm ${identity_maps:1:${#identity_maps}-2}

for user_account in ${user_account_resource_addresses[@]};do
    terraform state show -no-color $user_account | grep -v "   user_account_id" | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
done
for access_rule in ${access_rules_resouce_id_to_address[@]};do
    terraform state show -no-color $access_rule | grep -v "   id " >> cyral_migration_repository_access_rules_and_user_accounts.txt
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
echo "It is finally time to remove the empty resources from your .tf files."
echo "Please perform the following actions: "
echo
echo "  1.  Remove the empty resource definitions for the "
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
