resource "tfe_workspace" "managed_ws" {
  for_each = var.workspaces

  name                          = each.key
  description                   = each.value.description
  organization                  = var.tf_organization
  auto_apply                    = false
  assessments_enabled           = false
  file_triggers_enabled         = false
  execution_mode                = "agent"
  agent_pool_id                 = data.tfe_agent_pool.agent-pool[each.key].id
  structured_run_output_enabled = false
  project_id                    = data.tfe_project.proj[each.key].id

  force_delete = true
  vcs_repo {
    identifier     = each.value.reponame
    oauth_token_id = data.tfe_oauth_client.default.oauth_token_id
    branch         = each.value.branchname
  }
}

resource "tfe_variable" "var" {
  for_each = { for x in local.workspace_vars : x.id => x }

  key          = each.value.variable_key
  value        = each.value.value
  category     = each.value.category
  hcl          = each.value.hcl
  description  = each.value.description
  sensitive    = each.value.sensitive
  workspace_id = tfe_workspace.managed_ws[each.value.workspace_name].id
}

resource "tfe_workspace_variable_set" "ws_var_set" {
  for_each = {
    for key, value in var.workspaces :
    key => value
    if coalesce(value.variableset_name, "skip") != "skip"
  }

  variable_set_id = data.tfe_variable_set.var_set[each.key].id
  workspace_id    = tfe_workspace.managed_ws[each.key].id
}

resource "tfe_team_access" "team_access" {
  for_each = { for x in local.team_mapping : x.id => x }

  workspace_id = tfe_workspace.managed_ws[each.value.workspace_name].id
  access       = each.value.access
  team_id      = data.tfe_team.team[each.key].id
}