data "tfe_oauth_client" "default" {
  organization     = var.github_organization
  service_provider = "github"
}

data "tfe_project" "proj" {
  for_each = var.workspaces

  name         = each.value.project_id
  organization = var.tf_organization
}

data "tfe_agent_pool" "agent-pool" {
  for_each = var.workspaces

  name         = each.value.agent
  organization = var.tf_organization
}

data "tfe_variable_set" "var_set" {
  for_each = {
    for key, value in var.workspaces :
    key => value
    if coalesce(value.variableset_name, "skip") != "skip"
  }

  name         = each.value.variableset_name
  organization = var.tf_organization
}

data "tfe_team" "team" {
  for_each = { for x in local.team_mapping : x.id => x }

  name         = each.value.team_name
  organization = var.tf_organization
}