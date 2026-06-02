# terraform-provider-anthropic

Terraform provider for Anthropic admin API — workspaces and API keys.

## Usage

```hcl
terraform {
  required_providers {
    anthropic = {
      source  = "stark256-spec/anthropic"
      version = "~> 1.0"
    }
  }
}

provider "anthropic" {
  api_key         = var.anthropic_api_key
  organization_id = var.anthropic_org_id
}

resource "anthropic_workspace" "eng" {
  name         = "engineering"
  display_name = "Engineering"
}

resource "anthropic_api_key" "ci" {
  name         = "ci-pipeline"
  workspace_id = anthropic_workspace.eng.id
}

output "ci_key" {
  value     = anthropic_api_key.ci.secret_key
  sensitive = true
}
```

## Authentication

Set your API key via the `api_key` argument or the environment variable shown in the provider schema.

## Resources

| Resource | Description |
|----------|-------------|
| `anthropic_workspace` / `anthropic_project` / `anthropic_team` | Isolated environment |
| `anthropic_api_key` | API key scoped to a workspace/project |

## License

Apache 2.0
