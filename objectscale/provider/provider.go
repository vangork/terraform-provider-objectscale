/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"terraform-provider-objectscale/client"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ObsProvider satisfies various provider interfaces.
var _ provider.Provider = &ObsProvider{}

// ObsProvider defines the provider implementation.
type ObsProvider struct {
	// client can contain the upstream provider SDK or HTTP client used to
	// communicate with the upstream service. Resource and DataSource
	// implementations can then make calls using this client.

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Data describes the provider data model.
type Data struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

// Metadata describes the provider arguments.
func (p *ObsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "objectscale"
	resp.Version = p.version
}

// Schema describes the provider arguments.
func (p *ObsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.",
		Description:         "The Terraform provider for Dell PowerScale can be used to interact with a Dell PowerScale array in order to manage the array resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The API endpoint, ex. https://172.17.177.230:8080",
				Description:         "The API endpoint, ex. https://172.17.177.230:8080",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username",
				Description:         "The username",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password",
				Description:         "The password",
				Required:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "whether to skip SSL validation",
				Description:         "whether to skip SSL validation",
				Required:            true,
			},
		},
	}
}

// Configure configures the provider.
func (p *ObsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Data

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	obsClient, err := client.NewClient(
		data.Endpoint.ValueString(),
		data.Insecure.ValueBool(),
		data.Username.ValueString(),
		data.Password.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create objectscale client",
			err.Error(),
		)
		return
	}

	// client configuration for data sources and resources
	resp.DataSourceData = obsClient
	resp.ResourceData = obsClient
}

// Resources describes the provider resources.
func (p *ObsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources describes the provider data sources.
func (p *ObsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
	}
}

// New returns a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ObsProvider{
			version: version,
		}
	}
}
