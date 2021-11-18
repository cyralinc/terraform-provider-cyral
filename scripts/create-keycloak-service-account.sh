[ ! -z ${CYRAL_CONTROL_PLANE+x} ] || read -p "Enter control plane DNS: " CYRAL_CONTROL_PLANE

[ ! -z ${CYRAL_TOKEN+x} ] || read -p "Enter the JWT token obtained from the CP (use the UI or API): " CYRAL_TOKEN
HEADER="Authorization:Bearer $CYRAL_TOKEN"

# Get role ids necessary to run the provider
echo "Getting role IDs..."
ROLE_IDS=$(curl -X GET https://$CYRAL_CONTROL_PLANE:8000/v1/users/roles -H "$HEADER" -H "Content-type:Application/JSON" --fail --show-error)
if [[ $? -ne 0 ]]
then
	echo "Error getting the role IDs."
	exit 1
fi

ROLE_IDS=$(echo $ROLE_IDS | jq '[.roles | map(select(.name | contains("Modify Integrations", "Modify Policies", "Modify Roles","Modify Sidecars and Repositories", "View Sidecars and Repositories", "Modify Users"))) | .[].id]' -c)

if [[ $? -ne 0 ]]
then
	echo "Error unmarshalling the role IDs"
	echo "$ROLE_IDS"
	exit 1
fi

# Create/update a service account for Terraform and return the necessary parameters
# client_id and client_secret that will be used in the provider
curl -X POST https://$CYRAL_CONTROL_PLANE:8000/v1/users/serviceAccounts \
  -d '{"displayName":"terraform","roleIds":'"$ROLE_IDS"'}' \
  -H "$HEADER" \
  -H "Content-type:Application/JSON" | jq
if [[ $? -ne 0 ]]
then
	echo "Error creating the service account"
	exit 1
fi

echo "Use both \`clientId\` and \`clientSecret\` returned above to set up your provider."
echo "See the provider documentation in https://github.com/cyralinc/terraform-provider-cyral/blob/main/docs/index.md how to set those two parameters."
echo
