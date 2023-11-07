# terraform-cloud-workspace-collector
collect terraform cloud workspace info


## Usage

`tfc-collect <terraform_organization_name>

This will create the following files:
- `workspaces.tfvars` file with all the workspaces added as a map of maps.
- `imports.tf` file with all the import blocks for use with terraform 1.5 and newer
- `import.sh` file with all the import CLI commands for use with terraform 1.4 and older

## Flags
- `--debug` | `-d` = Debug mode. Prints api calls to std to determine if there are issues in responses

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