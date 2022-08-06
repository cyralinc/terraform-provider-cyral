data "cyral_role" "admin_roles" {
  # Optional. Filter roles with name that matches regular expression.
  name = "^.*Admin$"
}
