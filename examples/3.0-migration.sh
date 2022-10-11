#!/bin/bash

if ! command -v terraform &> /dev/null
then
    echo "The Terraform cli must be installed for this script to run."
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

echo "Welcome to Cyral's Terraform 3.0 Migration script!"
echo
echo "You have the option to allow the script to create new resource definitions for "
echo "the cyral_repository_user_accounts and cyral_repository_access_rules that will "
echo "be migrated. Alternatively, you can create the new resource definitions manually."
echo
echo "If you would like this script to create the new resources definitions, please set"
echo "CYRAL_TF_FILE_PATH equal to the file path of your .tf file."
echo
read -p "Are you ready to continue? [N/y] " -n 1 -r
echo 
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi
echo "Searching for cyral_repository_identity_map and cyral_repository_local_account resources to migrate..."
echo

user_account_import_ids=()
access_rule_import_ids=()

user_accounts_to_add=()
access_rules_to_add=()

local_accounts_to_delete=()
identity_maps_to_delete=()

empty_access_duration="$(printf '%s' '[]')"

# Create array of all resources in the terraform state
IFS=$'\n' read -r -d '' -a tf_state < <( terraform state list && printf '\0' )

# Find all cyral_repository_identity_maps and cyral_repository_local_accounts
for resource in ${tf_state[@]}; do
  if [[ $resource == cyral_repository_identity_map.* ]] 
  then
    # Get repo ID and local account ID for the identity map.
    repo_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_id" | sed 's/"//g' )
    local_account_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_local_account_id" | sed 's/"//g')
    identity_type=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.identity_type" | sed 's/"//g')
    access_duration=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.access_duration")
    if [[ $access_duration != $empty_access_duration ]] && [[ $identity_type == "user" ]]; then
        # Identity map was migrated to be an approval, which is not managed through terraform-- do nothing. 
        continue
    fi
    # We will need to delete this identity map from the .tf file, store its name
    identity_maps_to_delete+=($resource)
    # Construct import ID for the access rule that was migrated from this identity map. 
    import_id="$repo_id/$local_account_id"
    # Construct name of the access rule that will be imported. 
    import_name=cyral_repository_access_rule.${resource##"cyral_repository_identity_map."}
    # Save name of the new access rule, so that it can be added to the .tf file
    access_rules_to_add+=("resource \"cyral_repository_access_rule\" \"${resource##"cyral_repository_identity_map."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    access_rule_import_ids+=("${import_kv_pair}")
  elif [[ $resource == cyral_repository_local_account.* ]]
  then
    # Get repo ID and local account ID for the local account.
    repo_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.repository_id" | sed 's/"//g' )
    local_account_id=$(terraform show -json | jq ".values.root_module.resources[] | select(.address == \"$resource\") | .values.id" | sed 's/"//g')
    # We will need to delete this local account from the .tf file, store its name
    local_accounts_to_delete+=($resource)
    # Construct import ID for the user account that was migrated from this local account. 
    import_id="$repo_id/$local_account_id"
    # Construct name of the user account that will be imported. 
    import_name=cyral_repository_user_account.${resource##"cyral_repository_local_account."}
    # Save name of the migrated user account, so that it can be added to the .tf file
    user_accounts_to_add+=("resource \"cyral_repository_user_account\" \"${resource##"cyral_repository_local_account."}\" {}")
    # Store import name and ID as a key value pair
    import_kv_pair="${import_name} ${import_id}"
    user_account_import_ids+=("${import_kv_pair}")
  fi
done

echo "Found ${#local_accounts_to_delete[@]} cyral_repository_local_accounts to migrate to cyral_repository_user_accounts."
echo "Found ${#identity_maps_to_delete[@]} cyral_repository_identity_maps to migrate to cyral_repository_access_rules."
echo; echo; echo; sleep 2;
echo "The following file path was provided for CYRAL_TF_FILE_PATH: ${CYRAL_TF_FILE_PATH}"
read -p "Would you like this script to append these lines to the .tf file ${CYRAL_TF_FILE_PATH}? [N/y] " -n 1 -r 
if [[  $REPLY =~ ^[Yy]$ ]]
then
    printf '%s\n\n' "${user_accounts_to_add[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${access_rules_to_add[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo 
    echo "You have chosen to manually edit your .tf file."
    echo
    echo "Please add the following lines to your .tf file:"
    printf '%s\n' "${user_accounts_to_add[@]}"
    printf '%s\n' "${access_rules_to_add[@]}"
fi

echo; echo; echo;
read -p "Before we procede, please ensure that the new resource definitions were added to your .tf file. Are you ready to continue? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

echo; echo; echo;
echo "The following cyral_repository_identity_maps will be removed:"
printf '%s\n' "${identity_maps_to_delete[@]}"
echo
echo "The following cyral_repository_local_accounts will be removed:"
printf '%s\n' "${local_accounts_to_delete[@]}"
echo

read -p "Would you like to procede with removing old? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

for ((i = 0; i < ${#local_accounts_to_delete[@]}; i++));do
    terraform state rm ${local_accounts_to_delete[$i]}
done

for ((i = 0; i < ${#identity_maps_to_delete[@]}; i++));do
    terraform state rm ${identity_maps_to_delete[$i]}
done

read -p "Now that your .tf file is cleaned up, are you ready to import migrated resources? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

for ((i = 0; i < ${#user_account_import_ids[@]}; i++));do
    terraform import ${user_account_import_ids[$i]}
done
for ((i = 0; i < ${#access_rule_import_ids[@]}; i++));do
    terraform import ${access_rule_import_ids[$i]}
done

