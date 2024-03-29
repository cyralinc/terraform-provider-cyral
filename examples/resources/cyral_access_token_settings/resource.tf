resource "cyral_access_token_settings" "current" {
  max_validity = "72000s"
  default_validity = "36000s"
  max_number_of_tokens_per_user = 3
  offline_token_validation = true
  token_length = 16
}
