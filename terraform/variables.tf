variable "tf_team_token" {
  description = "The API token used should be able to create and configure workspaces variables"
  default     = ""
  sensitive   = true
}

variable "workspaces" {
  type = map(object({
    reponame         = string
    description      = string
    branchname       = string
    agent            = string
    project_id       = string
    variableset_name = optional(string)
    variables = map(object({
      category    = string
      description = string
      sensitive   = optional(bool, false)
      hcl         = optional(bool, false)
      value       = any
    }))
    teams = optional(map(object({
      access = string
    })))
  }))
}

variable "tf_hostname" {
  description = "The Terraform Cloud or Enterprise hostname."
  default     = "app.terraform.io"
}

variable "tf_organization" {
  type = string
}

variable "github_organization" {
  type = string
}
