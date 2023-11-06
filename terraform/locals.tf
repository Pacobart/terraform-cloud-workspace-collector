locals {

  workspace_vars = flatten([
    for k, v in var.workspaces : [
      for var_key, attributes in v.variables : {
        id             = "${k}-${var_key}"
        workspace_name = k
        category       = attributes.category
        value          = attributes.value
        variable_key   = var_key
        hcl            = attributes.hcl
        description    = attributes.description
        sensitive      = attributes.sensitive
      }
    ]
  ])

  team_mapping = flatten([
    for workspace_name, workspace_vars in var.workspaces : [
      for team_name, attributes in coalesce(workspace_vars.teams, {}) : {
        id             = "${workspace_name}-${team_name}"
        workspace_name = workspace_name
        team_name      = team_name
        access         = attributes.access
      }
    ]
  ])
}