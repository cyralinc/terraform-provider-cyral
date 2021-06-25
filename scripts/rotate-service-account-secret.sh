[ ! -z ${CYRAL_CONTROL_PLANE+x} ] || read -p "Enter control plane DNS: " CYRAL_CONTROL_PLANE

[ ! -z ${CYRAL_TOKEN+x} ] || read -p "Enter the JWT token obtained from the CP (use the UI or API): " CYRAL_TOKEN
HEADER="Authorization:Bearer $CYRAL_TOKEN"

[ ! -z ${CYRAL_CLIENT_ID+x} ] || read -p "Enter the client ID for the target service account: " CYRAL_CLIENT_ID

# Rotate the secret
curl -X POST https://$CYRAL_CONTROL_PLANE:8000/v1/users/serviceAccounts/$CYRAL_CLIENT_ID/rotateSecret /
  -H "$HEADER" /
  -H "Content-type:Application/JSON" | jq
if [[ $? -ne 0 ]]
then
	echo "Error creating the service account"
	exit 1
fi

echo "Use both \`clientId\` and \`clientSecret\` returned above to set up your provider."
echo "See the provider documentation in https://github.com/cyralinc/terraform-provider-cyral/blob/main/doc/provider.md how to set those two parameters."
echo
