terraform {
  required_providers {
    anthropic = {
      source  = "stark256-spec/anthropic"
      version = "~> 1.0"
    }
  }
}

provider "anthropic" {
  api_key         = var.anthropic_api_key   # or ANTHROPIC_API_KEY env var
  organization_id = var.anthropic_org_id    # or ANTHROPIC_ORG_ID env var
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