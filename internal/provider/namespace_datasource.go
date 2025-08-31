package provider

import (
	"context"
	"fmt"
	"terraform-provider-objectscale/client"
	"terraform-provider-objectscale/internal/helper"
	"terraform-provider-objectscale/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NamespaceDataSource{}

func NewNamespaceDataSource() datasource.DataSource {
	return &NamespaceDataSource{}
}

type NamespaceDataSource struct {
	client *client.Client
}

func (d *NamespaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

// Schema describes the data source arguments.
func (d *NamespaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Namespace.",
		Description:         "Namespace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"namespaces": schema.ListNestedAttribute{
				Description:         "List of Namespaces",
				MarkdownDescription: "List of Namespaces",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description:         "Name.",
							MarkdownDescription: "Name.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description:         "Id.",
							MarkdownDescription: "Id.",
							Computed:            true,
						},
						"global": schema.BoolAttribute{
							Description:         "Global.",
							MarkdownDescription: "Global.",
							Computed:            true,
						},
						"remote": schema.BoolAttribute{
							Description:         "Remote.",
							MarkdownDescription: "Remote.",
							Computed:            true,
						},
						"link": schema.SingleNestedAttribute{
							Description:         "Link.",
							MarkdownDescription: "Link.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"rel": schema.StringAttribute{
									Description:         "Rel.",
									MarkdownDescription: "Rel.",
									Computed:            true,
								},
								"href": schema.StringAttribute{
									Description:         "Href.",
									MarkdownDescription: "Href.",
									Computed:            true,
								},
							},
						},
						"creation_time": schema.Int64Attribute{
							Description:         "CreationTime.",
							MarkdownDescription: "CreationTime.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NamespaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *NamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.NamespaceDatasourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	namespaces, err := d.client.ManagementClient.ListNamespaces("")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting the list of namespaces",
			err.Error(),
		)
		return
	}

	var namespaceList []models.NamespaceEntity
	// Convert from json to terraform model
	for _, namespace := range namespaces {
		entity := models.NamespaceEntity{}
		err := helper.CopyFields(ctx, namespace, &entity)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting namespaces",
				err.Error(),
			)
			return
		}
		namespaceList = append(namespaceList, entity)
	}

	// hardcoding a response value to save into the Terraform state.
	data.ID = types.StringValue("namespace_datasource")
	data.Namespaces = namespaceList

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read namespace data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
