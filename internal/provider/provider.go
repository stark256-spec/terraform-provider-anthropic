package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.ProviderWithFunctions = &AnthropicProvider{}

type AnthropicProvider struct{ version string }

type AnthropicProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	OrgID   types.String `tfsdk:"organization_id"`
	BaseURL types.String `tfsdk:"base_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &AnthropicProvider{version: version} }
}

func (p *AnthropicProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "anthropic"
	resp.Version = p.version
}

func (p *AnthropicProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Anthropic organization resources — workspaces, API keys, and members — via the Anthropic Admin API.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Anthropic admin API key. Can be set via `ANTHROPIC_API_KEY` env var.",
				Optional:            true,
				Sensitive:           true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Anthropic organization ID. Can be set via `ANTHROPIC_ORG_ID` env var.",
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Override the Anthropic API base URL. Defaults to `https://api.anthropic.com`.",
				Optional:            true,
			},
		},
	}
}

func (p *AnthropicProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg AnthropicProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if !cfg.APIKey.IsNull() {
		apiKey = cfg.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API key", "Set api_key in the provider block or ANTHROPIC_API_KEY env var.")
		return
	}

	orgID := os.Getenv("ANTHROPIC_ORG_ID")
	if !cfg.OrgID.IsNull() {
		orgID = cfg.OrgID.ValueString()
	}
	if orgID == "" {
		resp.Diagnostics.AddError("Missing organization_id", "Set organization_id in the provider block or ANTHROPIC_ORG_ID env var.")
		return
	}

	baseURL := ""
	if !cfg.BaseURL.IsNull() {
		baseURL = cfg.BaseURL.ValueString()
	}

	client := newClient(apiKey, orgID, baseURL)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AnthropicProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkspaceResource,
		NewAPIKeyResource,
	}
}

func (p *AnthropicProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AnthropicProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}
