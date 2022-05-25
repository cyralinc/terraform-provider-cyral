# Parses metadata Base64 document to SAML configuration.
data "cyral_saml_configuration" "some_data_source_name" {
  base_64_saml_metadata_document = "some_metadata_base_64_document"
}

# Parses metadata URL to SAML configuration.
data "cyral_saml_configuration" "some_data_source_name" {
  saml_metadata_url = "some_metadata_url"
}
