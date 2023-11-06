# terraform-cloud-workspace-collector
collect terraform cloud workspace info


## Usage

`tfc-collect <terraform_organization_name>

This will create `workspaces.tfvars` file with all the workspaces added as a map of maps.

## Notes

This expects the terraform credentials token to exist in the following locations based on Operating System and currently only supports "app.terraform.io"

`windows` = $HOME\AppData\Roaming\terraform.d\credentials.tfrc.json
`linux|mac` = $HOME/.terraform.d/credentials.tfrc.json

Example credentials file:
```
{
  "credentials": {
    "app.terraform.io": {
      "token": "TOKEN"
    }
  }
}
```