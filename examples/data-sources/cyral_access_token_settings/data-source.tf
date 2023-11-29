data "cyral_access_token_settings" "token_settings" {}

output "max_validity" {
  value = data.cyral_access_token_settings.token_settings.max_validity
}

output "default_validity" {
  value = data.cyral_access_token_settings.token_settings.default_validity
}

output "max_number_of_tokens_per_user" {
  value = data.cyral_access_token_settings.token_settings.max_number_of_tokens_per_user
}

output "offline_token_validation" {
  value = data.cyral_access_token_settings.token_settings.offline_token_validation
}

output "token_length" {
  value = data.cyral_access_token_settings.token_settings.token_length
}
