#!/usr/bin/env bash

# Set the color variable
red='\033[0;31m'
green='\033[0;32m'
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

if [[ -z $CYRAL_TF_CONTROL_PLANE || -z $CYRAL_TF_CLIENT_ID || -z $CYRAL_TF_CLIENT_SECRET ]];
then
    echo "All of the following Control Plane configuration environment"
    echo "variables must be set:"
    echo "- CYRAL_TF_CONTROL_PLANE"
    echo "- CYRAL_TF_CLIENT_ID"
    echo "- CYRAL_TF_CLIENT_SECRET"
    exit
fi

echo -e "${green}Welcome to Cyral's Terraform Provider Version 4.0 Migration script!${clear}"
echo
echo \
"This script will create new resource definitions for the
cyral_repository and cyral_repository_binding resource that
will be migrated. Additionally, cyral_sidecar_listener resources
will be created, which are now required to bind sidecars to
repositories."
echo
echo -e "${green}Please set CYRAL_TF_FILE_PATH equal to the file path of your .tf file.${clear}"
if [ -z "$CYRAL_TF_FILE_PATH" ]
then
    echo -e "${red}CYRAL_TF_FILE_PATH has not been set. Please set it and run the script again.${clear}"
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
else
    echo -e "CYRAL_TF_FILE_PATH is set to ${green}'$CYRAL_TF_FILE_PATH'${clear}"
fi
echo
echo "Searching for cyral_repository and cyral_repository_binding resources to migrate..."
echo

terraform state pull > terraform.tfstate.cyral.migration.backup
cp ${CYRAL_TF_FILE_PATH} cyral_terraform_migration_backup_configuration.txt

repo_import_args=()
binding_import_args=()
access_gateway_import_args=()
listener_import_args=()

repo_resource_defs=()
binding_resource_defs=()
access_gateway_resource_defs=()

declare -A repo_resource_foreach_defs_map
declare -A binding_resource_foreach_defs_map
declare -A access_gateway_resource_foreach_defs_map
declare -A listener_resource_foreach_defs_map

declare -A repo_resource_id_to_address_map
declare -A binding_resource_id_to_address_map
declare -A listener_resource_id_to_address_map
declare -A sidecar_resource_id_to_address_map

access_gateway_resource_addresses=()

repos_to_delete=()
bindings_to_delete=()

# Regex to determine if a resource was defined using a
# for_each in the terraform configuration.
# (E.g.: cyral_repository.all_repos["mysql"])
foreach_regex=*"[\""*

# Create array of repo and binding resource to migrate
resources_to_migrate=($(terraform state list | grep "cyral_repository\.\|cyral_repository_binding\.\|cyral_sidecar\."))
# Store terraform state JSON representation
tf_state_json=$(terraform show -json | jq ".values.root_module.resources[]")

# Find all cyral_repository and cyral_repository_binding
for resource_address in ${resources_to_migrate[@]}; do
  if [[ $resource_address == cyral_repository.* ]]
  then
    # We will need to delete this repo from the .tf file, so we
    # store the full resource address.
    repos_to_delete+=($resource_address)
    # Escape the double quotes so we find it using jq
    escaped_resource_address=$(sed -e 's/\"/\\"/g'<<<$resource_address)
    # Get repo ID. This will be used to import the updated repo.
    repo_id=($(jq -r "select(.address == \"$escaped_resource_address\") | .values.id"<<<$tf_state_json))
    # Save an empty resource definition for a new repo, so that
    # it can be added to the .tf file. We also check if this
    # resource was defined using a for_each (E.g.: cyral_repository.all_repos["mysql"]).
    if [[ "$resource_address" == $foreach_regex ]]; then # If the resource was defined using a for_each
        resource_name=$(sed -e 's/\[[^][]*\]//'<<<${resource_address##"cyral_repository."})
        repo_resource_foreach_defs_map[resource_name]="resource \"cyral_repository\" \"${resource_name}\" {}"
    else # If its a normal resource definition
        repo_resource_defs+=("resource \"cyral_repository\" \"${resource_address##"cyral_repository."}\" {}")
    fi
    # Store import resource address and repo ID, which will be the argument to terraform import.
    repo_import_args+=("$resource_address $repo_id")
    repo_resource_id_to_address_map[${repo_id}]=${resource_address}
  elif [[ $resource_address == cyral_repository_binding.* ]]
  then
    # We will need to delete this sidecar binding from the .tf file, store its name
    bindings_to_delete+=($resource_address)
    # Escape the double quotes so we find it using jq
    escaped_resource_address=$(sed -e 's/\"/\\"/g'<<<$resource_address)
    # Get ids required to import binding resource.
    id_values_array=($(jq -r "select(.address == \"$escaped_resource_address\") | .values.sidecar_id, .values.repository_id, .values.sidecar_as_idp_access_gateway"<<<$tf_state_json))
    # Construct import ID for the repository binding that was migrated in CP.
    import_id="${id_values_array[0]}/${id_values_array[1]}"
    if [[ "$resource_address" == $foreach_regex ]]; then  # If the resource was defined using a for_each
        binding_resource_name=$(sed -e 's/\[[^][]*\]//'<<<${resource_address##"cyral_repository_binding."})
        # Save empty resource definition for the binding, so that it can be added to the .tf file
        binding_resource_foreach_defs_map[binding_resource_name]="resource \"cyral_repository_binding\" \"${binding_resource_name}\" {}"
        # If the binding is an access gateway, we will need to import it.
        if [[ ${id_values_array[2]} == true ]]
        then
            # Construct a name for the access gateway resource based on the binding resource name.
            access_gateway_resource_name="${binding_resource_name}_all_access_gateways"
            # Extract the binding key. For 'cyral_repository_binding.all_bindings["binding_0"]'
            # the binding key would be "binding_0", for example.
            binding_key=$(echo $resource_address | awk -F '[\\[\\]]' '{print $2}')
            # Contruct full resource address for the access gateway.
            access_gateway_resource_address="cyral_repository_access_gateway.${access_gateway_resource_name}[$binding_key]"
            # Save empty resource definition for the access gateway, so that it can be added to the .tf file
            access_gateway_resource_foreach_defs_map[access_gateway_resource_name]="resource \"cyral_repository_access_gateway\" \"${access_gateway_resource_name}\" {}"
            # Store import argument and address for access gateway
            access_gateway_import_args+=("$access_gateway_resource_address ${id_values_array[1]}")
            access_gateway_resource_addresses+=(${access_gateway_resource_address})
        fi
    else # If its a normal resource definition
        binding_resource_name=${resource_address##"cyral_repository_binding."}
        # Save empty resource definition for the binding, so that it can be added to the .tf file
        binding_resource_defs+=("resource \"cyral_repository_binding\" \"${binding_resource_name}\" {}")
        # If the binding is an access gateway, we will need to import it.
        if [[ ${id_values_array[2]} == true ]]
        then
            # Construct a name for the access gateway resource based on the binding resource name.
            access_gateway_resource_name="${binding_resource_name}_access_gateway"
            # Contruct full resource name for the access gateway.
            access_gateway_resource_address="cyral_repository_access_gateway.${access_gateway_resource_name}"
            # Save empty resource definition for the access gateway, so that it can be added to the .tf file
            access_gateway_resource_defs+="resource \"cyral_repository_access_gateway\" \"${access_gateway_resource_name}\" {}"
            # Store import argument and name for access gateway
            access_gateway_import_args+=("$access_gateway_resource_address ${id_values_array[1]}")
            access_gateway_resource_addresses+=(${access_gateway_resource_address})
        fi
    fi
    # Store import argument and name for binding
    binding_import_args+=("$resource_address $import_id")
    binding_resource_id_to_address_map[${id_values_array[1]}]=${resource_address}
  elif [[ $resource_address == cyral_sidecar.* ]]
  then
    # Get sidecar ID.
    sidecar_id=($(jq -r "select(.address == \"$resource_address\") | .values.id"<<<$tf_state_json))
    # Store full resource name with ID. We need this to replace references to the sidecar id
    # with the sidecar resource's full resource name.
    sidecar_resource_id_to_address_map[${sidecar_id}]=${resource_address}
  fi
done

echo "Found ${#repos_to_delete[@]} cyral_repository resources to migrate."
echo "Found ${#bindings_to_delete[@]} cyral_repository_binding resources to migrate."
echo
read -p "Would you like this script to append these lines to the .tf file ${CYRAL_TF_FILE_PATH}? [N/y] " -n 1 -r
if [[  $REPLY =~ ^[Yy]$ ]]
then
    printf '\n\n' >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${repo_resource_foreach_defs_map[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${binding_resource_foreach_defs_map[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${access_gateway_resource_foreach_defs_map[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${repo_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${binding_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
    printf '%s\n\n' "${access_gateway_resource_defs[@]}" >> ${CYRAL_TF_FILE_PATH}
else
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1
fi

echo; echo; echo;
echo "Now its time to upgrade your Cyral Terraform Provider to version 4!"
echo
echo -e "${green}Before we proceed, you will need to do the following${clear}:
    1.  Open your Terraform .tf configuration file.
    2.  Change the version number of the cyral provider in the required_providers
        section of your .tf configuration file to '~>4.0'. It should look like this:
${green}    cyral = {
                source  = \"cyralinc/cyral\"
                version = \"~>4.0\"
            }${clear}
    3. Ensure that new empty resource definitions were added to the end of
        your .tf file. The definitions will look like this:
            Repositories
            resource \"cyral_repository\" \"<resource_name>\" {}

            Respository Bindings
            resource \"cyral_respository_binding\" \"<resource_name>\" {}

            Respository Access Gateways
            resource \"cyral_respository_access_gateway\" \"<resource_name>\" {}
    ************************************************
    *                                              *
    *        ${red}IMPORTANT STEP PLEASE DONT SKIP${clear}       *
    *                                              *
    ************************************************
    4.  Find all non-empty references to cyral_repository and cyral_repository_binding
        resources in your .tf file and remove the entire resource definition
        for each one. Please leave the empty resource definitions that were added by
        this script."
echo
read -p "Are you ready to upgrade Terraform? [N/y] " -n 1 -r
echo    # move to a new line
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Exiting..."
    [[ "$0" = "$BASH_SOURCE" ]] && exit 1 || return 1 # handle exits from shell or function but don't exit interactive shell
fi

terraform init -upgrade

repos="[$(IFS=" "; echo "${repos_to_delete[*]}")]"
terraform state rm ${repos:1:${#repos}-2}

bindings="[$(IFS=" "; echo "${bindings_to_delete[*]}")]"
terraform state rm ${bindings:1:${#bindings}-2}

echo
echo "Importing the following cyral_repository resources into your Terraform state:"
printf '%s\n' "${repo_resource_id_to_address_map[@]}"
echo
echo "Importing the following cyral_repository_binding resource into your Terraform state:"
printf '%s\n' "${binding_resource_id_to_address_map[@]}"
echo
echo "Importing the following cyral_repository_access_gateway resources into your Terraform state:"
printf '%s\n' "${access_gateway_resource_addresses[@]}"
echo

for repo in "${repo_import_args[@]}";do
    terraform import $repo
done
for binding in "${binding_import_args[@]}";do
    terraform import $binding
done
for access_gateway in "${access_gateway_import_args[@]}";do
    terraform import $access_gateway
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

declare -A sidecar_resource_id_to_listener_ids_map

NEW_LINE=$'\n'

for binding in "${migrated_bindings[@]}"; do
    # Get ids required to import listener resource.
    binding_resource_address=${binding##"cyral_repository_binding."}
    # Escape the double quotes so we find it using jq
    escaped_binding_resource_address=$(sed -e 's/\"/\\"/g'<<<$binding_resource_address)
    binding_json=$(jq -r "select(.address == \"cyral_repository_binding.${escaped_binding_resource_address}\")"<<<"$tf_state_json")
    sidecar_id=$(jq -r ".values.sidecar_id"<<<"$binding_json")
    listener_ids=$(jq -r ".values.listener_binding[] | .listener_id"<<<"$binding_json")
    current_listener_ids=${sidecar_resource_id_to_listener_ids_map[$sidecar_id]}
    # Concat listener ids using a line break separator
    sidecar_resource_id_to_listener_ids_map[$sidecar_id]="$current_listener_ids${NEW_LINE}$listener_ids"
done

for sidecar_id in ${!sidecar_resource_id_to_listener_ids_map[@]}; do
    listener_ids=${sidecar_resource_id_to_listener_ids_map[$sidecar_id]}
    SAVEIFS=$IFS   # Save current IFS (Internal Field Separator)
    IFS=$NEW_LINE  # Change IFS to newline char
    listener_ids=($listener_ids) # split the `listener_ids` string into an array
    IFS=$SAVEIFS   # Restore original IFS

    # Get sidecar listeners so that we can create the listener resource name.
    # Since we have a listener for each repo type and port combination, the name
    # should follow the format ${sidecar_resource_name}_listener_${type}_${port}
    access_token=$(
        curl --location --request POST \
        "https://${CYRAL_TF_CONTROL_PLANE}/v1/users/oidc/token" \
        --header "Content-Type: application/x-www-form-urlencoded" \
        --data-urlencode "grant_type=client_credentials" \
        --data-urlencode "clientId=${CYRAL_TF_CLIENT_ID}" \
        --data-urlencode "clientSecret=${CYRAL_TF_CLIENT_SECRET}" | jq -r '.access_token'
    )
    sidecar_listeners_json=$(
        curl --location --request GET \
        "https://${CYRAL_TF_CONTROL_PLANE}/v1/sidecars/${sidecar_id}/listeners" \
        --header "Authorization: Bearer ${access_token}"
    )

    sidecar_resource_address=${sidecar_resource_id_to_address_map[${sidecar_id}]}
    sidecar_resource_name=$(sed -e 's/\[[^][]*\]//'<<<${sidecar_resource_address##"cyral_sidecar."})
    # Save empty resource definition for all the sidecar listeners, so that it can be added to the .tf file
    all_listeners_resource_name="${sidecar_resource_name}_all_listeners"
    listener_resource_foreach_defs_map[all_listeners_resource_name]="resource \"cyral_sidecar_listener\" \"${all_listeners_resource_name}\" {}"

    for i in "${!listener_ids[@]}"; do
        # Check to ensure that listener name & resource has not already been stored
        if ! [[ -v listener_resource_id_to_address_map[${listener_ids[$i]}] ]]
        then
            listener_json=$(
                jq -r ".listenerConfigs[] | select(.id == \"${listener_ids[$i]}\")" \
                <<<"$sidecar_listeners_json"
            )
            read listener_type listener_port  < <(
                echo $(jq -r ".repoTypes[0], .address.port"<<<"$listener_json")
            )
            # Construct import ID for the listener that was created during CP migration.
            import_id="${sidecar_id}/${listener_ids[$i]}"
            # Store import argument and name for listener
            listener_resource_address="cyral_sidecar_listener.${all_listeners_resource_name}[\"${listener_type}_${listener_port}\"]"
            listener_import_args+=("${listener_resource_address} ${import_id}")
            listener_resource_id_to_address_map[${listener_ids[$i]}]=${listener_resource_address}
        fi
    done
done

echo "Found ${#listener_resource_id_to_address_map[@]} cyral_sidecar_listener resources to import."
echo
echo -e "Appending empty resource definitions to the file: ${green}${CYRAL_TF_FILE_PATH}${clear}"

printf '\n\n' >> ${CYRAL_TF_FILE_PATH}
printf '%s\n\n' "${listener_resource_foreach_defs_map[@]}" >> ${CYRAL_TF_FILE_PATH}

echo "Importing the following cyral_sidecar_listeners into your Terraform state:"
printf '%s\n' "${listener_resource_id_to_address_map[@]}"
echo

for listener in "${listener_import_args[@]}";do
    terraform import $listener
done

# Once resources have been imported, we need to replace the empty resource definitions
# (which are required to exist prior to import), with resource definitions containing
# the correct values. We do this by printing the resource definition from the terraform
# state, then remove any computed values. Finally, we append the new resource definitions
# to a seperate .tf file, containing all of the resources that were migrated.
# OBS: For resources defined using for_each, we just append an empty resource definition
# since its expected that users will have to configure the resources manually based
# on their specific terraform configuration.

MIGRATED_RESOURCES_CONFIG_FILE="migrated_repositories_bindings_access_gateways_listeners"

# Write for_each resources definitions
for repo_key in ${!repo_resource_foreach_defs_map[@]};do
    echo ${repo_resource_foreach_defs_map[$repo_key]%%"}"}$'\n\tfor_each={}\n}' >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done
for binding_key in ${!binding_resource_foreach_defs_map[@]};do
    echo ${binding_resource_foreach_defs_map[$binding_key]%%"}"}$'\n\tfor_each={}\n}' >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done
for access_gateway_key in ${!access_gateway_resource_foreach_defs_map[@]};do
    echo ${access_gateway_resource_foreach_defs_map[$access_gateway_key]%%"}"}$'\n\tfor_each={}\n}' >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done
for listener_key in ${!listener_resource_foreach_defs_map[@]};do
    echo ${listener_resource_foreach_defs_map[$listener_key]%%"}"}$'\n\tfor_each={}\n}' >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done

for repo in ${repo_resource_id_to_address_map[@]};do
    if [[ "$repo" != $foreach_regex ]]; then # Skip if its a resource defined using a for_each
        terraform state show -no-color $repo | grep -v "   id " >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
    fi
done
for binding in ${binding_resource_id_to_address_map[@]};do
    if [[ "$binding" != $foreach_regex ]]; then # # Skip if its a resource defined using a for_each
        terraform state show -no-color $binding | grep -v "   binding_id" | grep -v "   id " >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
    fi
done
for access_gateway in ${access_gateway_resource_addresses[@]};do
    if [[ "$access_gateway" != $foreach_regex ]]; then # Skip if its a resource defined using a for_each
        terraform state show -no-color $access_gateway | grep -v "   id " >> ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
    fi
done

# Replace resource IDs with full resource names.
for repo_id in ${!repo_resource_id_to_address_map[@]};do
    sed -i.bak "s/[[:space:]]*repository_id[[:space:]]*=[[:space:]]*\"${repo_id}\"/   repository_id = ${repo_resource_id_to_address_map[${repo_id}]}.id/g" ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done

for binding_id in ${!binding_resource_id_to_address_map[@]};do
    sed -i.bak "s/[[:space:]]*binding_id[[:space:]]*=[[:space:]]*\"${binding_id}\"/   binding_id = ${binding_resource_id_to_address_map[${binding_id}]}.binding_id/g" ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done

for listener_id in ${!listener_resource_id_to_address_map[@]};do
    sed -i.bak "s/[[:space:]]*listener_id[[:space:]]*=[[:space:]]*\"${listener_id}\"/   listener_id = ${listener_resource_id_to_address_map[${listener_id}]}.listener_id/g" ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done

for sidecar_id in ${!sidecar_resource_id_to_address_map[@]};do
    sed -i.bak "s/[[:space:]]*sidecar_id[[:space:]]*=[[:space:]]*\"${sidecar_id}\"/   sidecar_id = ${sidecar_resource_id_to_address_map[${sidecar_id}]}.id/g" ${MIGRATED_RESOURCES_CONFIG_FILE}.txt
done

mv ${MIGRATED_RESOURCES_CONFIG_FILE}.txt ${MIGRATED_RESOURCES_CONFIG_FILE}.tf
rm ${MIGRATED_RESOURCES_CONFIG_FILE}.txt.bak
terraform fmt

echo; echo; echo;
echo "Now that the Terraform state is up-to-date, let's clean up your .tf file."
echo "This script created a new .tf file containing the resources definitions"
echo "for the new resources that were migrated into your Terraform state."
echo "The new .tf files is called:"
echo
echo "${MIGRATED_RESOURCES_CONFIG_FILE}.tf"
echo
echo
echo "It is finally time to remove the empty resources from your .tf files"
echo "and finish configuring the remaining resources."
echo "Please perform the following actions: "
echo
echo "  1.  Remove the empty resource definitions for the cyral_repository"
echo "      cyral_repository_binding, cyral_repository_access_gateway and"
echo "      cyral_sidecar_listener resources that were added to the"
echo "      the end of your .tf file, which is named:"
echo "      ${CYRAL_TF_FILE_PATH}"
echo "  2.  Find all the resources defined with a for_each in the file"
echo "      '${MIGRATED_RESOURCES_CONFIG_FILE}.tf'"
echo "      and finish configuring them according to their new schema."
echo "      Please refer to the migration guide to find examples on "
echo "      how to configure these resources."
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
echo "      - ${MIGRATED_RESOURCES_CONFIG_FILE}.tf"
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
