provider "tfe" {
  token    = var.tf_team_token
  hostname = var.tf_hostname
}

terraform {
  required_version = ">=1.0"
}