package provider

import (
	"context"
	"fmt"
	"terraform-provider-objectscale/internal/client"
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
						"inactive": schema.BoolAttribute{
							Description:         "Inactive.",
							MarkdownDescription: "Inactive.",
							Computed:            true,
						},
						"internal": schema.BoolAttribute{
							Description:         "Internal.",
							MarkdownDescription: "Internal.",
							Computed:            true,
						},
						"default_data_services_vpool": schema.StringAttribute{
							Description:         "DefaultDataServicesVpool.",
							MarkdownDescription: "DefaultDataServicesVpool.",
							Required:            true,
						},
						"allowed_vpools_list": schema.ListAttribute{
							Description:         "AllowedVpoolsList.",
							MarkdownDescription: "AllowedVpoolsList.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"disallowed_vpools_list": schema.ListAttribute{
							Description:         "DisallowedVpoolsList.",
							MarkdownDescription: "DisallowedVpoolsList.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"namespace_admins": schema.StringAttribute{
							Description:         "NamespaceAdmins.",
							MarkdownDescription: "NamespaceAdmins.",
							Computed:            true,
						},
						"user_mapping": schema.ListNestedAttribute{
							Description:         "UserMapping.",
							MarkdownDescription: "UserMapping.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"domain": schema.StringAttribute{
										Description:         "Domain",
										MarkdownDescription: "Domain",
										Computed:            true,
									},
									"groups": schema.ListAttribute{
										Description:         "Groups.",
										MarkdownDescription: "Groups.",
										Computed:            true,
										ElementType:         types.StringType,
									},
									"attributes": schema.ListNestedAttribute{
										Description:         "Attributes.",
										MarkdownDescription: "Attributes.",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Description:         "Key",
													MarkdownDescription: "Key",
													Computed:            true,
												},
												"value": schema.ListAttribute{
													Description:         "Value.",
													MarkdownDescription: "Value.",
													Computed:            true,
													ElementType:         types.StringType,
												},
											},
										},
									},
								},
							},
						},
						"is_encryption_enabled": schema.BoolAttribute{
							Description:         "IsEncryptionEnabled.",
							MarkdownDescription: "IsEncryptionEnabled.",
							Computed:            true,
						},
						"default_bucket_block_size": schema.Int64Attribute{
							Description:         "DefaultBucketBlockSize.",
							MarkdownDescription: "DefaultBucketBlockSize.",
							Computed:            true,
						},
						"external_group_admins": schema.StringAttribute{
							Description:         "ExternalGroupAdmins.",
							MarkdownDescription: "ExternalGroupAdmins.",
							Computed:            true,
						},
						"is_stale_allowed": schema.BoolAttribute{
							Description:         "IsStaleAllowed.",
							MarkdownDescription: "IsStaleAllowed.",
							Computed:            true,
						},
						"is_object_lock_with_ado_allowed": schema.BoolAttribute{
							Description:         "IsObjectLockWithAdoAllowed.",
							MarkdownDescription: "IsObjectLockWithAdoAllowed.",
							Computed:            true,
						},
						"is_compliance_enabled": schema.BoolAttribute{
							Description:         "IsComplianceEnabled.",
							MarkdownDescription: "IsComplianceEnabled.",
							Computed:            true,
						},
						"notification_size": schema.Int64Attribute{
							Description:         "NotificationSize.",
							MarkdownDescription: "NotificationSize.",
							Computed:            true,
						},
						"block_size": schema.Int64Attribute{
							Description:         "BlockSize.",
							MarkdownDescription: "BlockSize.",
							Computed:            true,
						},
						"notification_size_in_count": schema.Int64Attribute{
							Description:         "NotificationSizeInCount.",
							MarkdownDescription: "NotificationSizeInCount.",
							Computed:            true,
						},
						"block_size_in_count": schema.Int64Attribute{
							Description:         "BlockSizeInCount.",
							MarkdownDescription: "BlockSizeInCount.",
							Computed:            true,
						},
						"default_audit_delete_expiration": schema.Int64Attribute{
							Description:         "DefaultAuditDeleteExpiration.",
							MarkdownDescription: "DefaultAuditDeleteExpiration.",
							Computed:            true,
						},
						"retention_classes": schema.SingleNestedAttribute{
							Description:         "RetentionClasses.",
							MarkdownDescription: "RetentionClasses.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"retention_class": schema.ListNestedAttribute{
									Description:         "RetentionClass.",
									MarkdownDescription: "RetentionClass.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description:         "Name",
												MarkdownDescription: "Name",
												Computed:            true,
											},
											"period": schema.Int64Attribute{
												Description:         "Period.",
												MarkdownDescription: "Period.",
												Computed:            true,
											},
										},
									},
								},
							},
						},
						"root_user_name": schema.StringAttribute{
							Description:         "RootUserName.",
							MarkdownDescription: "RootUserName.",
							Computed:            true,
						},
						"root_user_password": schema.StringAttribute{
							Description:         "RootUserPassword.",
							MarkdownDescription: "RootUserPassword.",
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
