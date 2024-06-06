resource "cyral_policy" "this" {
  name = "My first policy"
  description = "This is my first policy"
  enabled = true
  data = ["EMAIL"]
  metadata_tags = ["Risk Level 1"]
}
