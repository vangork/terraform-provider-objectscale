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
	"fmt"
	"terraform-provider-objectscale/internal/client"
	"terraform-provider-objectscale/internal/helper"
	"terraform-provider-objectscale/internal/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NamespaceResource{}
var _ resource.ResourceWithImportState = &NamespaceResource{}

func NewNamespaceResource() resource.Resource {
	return &NamespaceResource{}
}

// NamespaceResource defines the resource implementation.
type NamespaceResource struct {
	client *client.Client
}

func (r *NamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (r *NamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Namespace.",
		Description:         "Namespace.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description:         "Name.",
				MarkdownDescription: "Name.",
				Required:            true,
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
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"user_mapping": schema.ListNestedAttribute{
				Description:         "UserMapping.",
				MarkdownDescription: "UserMapping.",
				Optional:            true,
				Computed:            true,
				Default: listdefault.StaticValue(types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
					"domain": types.StringType,
					"groups": types.ListType{ElemType: types.StringType},
					"attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
						"key":   types.StringType,
						"value": types.ListType{ElemType: types.StringType},
					}}},
				}}, []attr.Value{})),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain": schema.StringAttribute{
							Description:         "Domain",
							MarkdownDescription: "Domain",
							Optional:            true,
							Computed:            true,
						},
						"groups": schema.ListAttribute{
							Description:         "Groups.",
							MarkdownDescription: "Groups.",
							Optional:            true,
							Computed:            true,
							ElementType:         types.StringType,
						},
						"attributes": schema.ListNestedAttribute{
							Description:         "Attributes.",
							MarkdownDescription: "Attributes.",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description:         "Key",
										MarkdownDescription: "Key",
										Optional:            true,
										Computed:            true,
									},
									"value": schema.ListAttribute{
										Description:         "Value.",
										MarkdownDescription: "Value.",
										Optional:            true,
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
				Optional:            true,
				Default:             booldefault.StaticBool(false),
			},
			"default_bucket_block_size": schema.Int64Attribute{
				Description:         "DefaultBucketBlockSize.",
				MarkdownDescription: "DefaultBucketBlockSize.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"external_group_admins": schema.StringAttribute{
				Description:         "ExternalGroupAdmins.",
				MarkdownDescription: "ExternalGroupAdmins.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"is_stale_allowed": schema.BoolAttribute{
				Description:         "IsStaleAllowed.",
				MarkdownDescription: "IsStaleAllowed.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_object_lock_with_ado_allowed": schema.BoolAttribute{
				Description:         "IsObjectLockWithAdoAllowed.",
				MarkdownDescription: "IsObjectLockWithAdoAllowed.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_compliance_enabled": schema.BoolAttribute{
				Description:         "IsComplianceEnabled.",
				MarkdownDescription: "IsComplianceEnabled.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"notification_size": schema.Int64Attribute{
				Description:         "NotificationSize.",
				MarkdownDescription: "NotificationSize.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"block_size": schema.Int64Attribute{
				Description:         "BlockSize.",
				MarkdownDescription: "BlockSize.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"notification_size_in_count": schema.Int64Attribute{
				Description:         "NotificationSizeInCount.",
				MarkdownDescription: "NotificationSizeInCount.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"block_size_in_count": schema.Int64Attribute{
				Description:         "BlockSizeInCount.",
				MarkdownDescription: "BlockSizeInCount.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(-1),
			},
			"default_audit_delete_expiration": schema.Int64Attribute{
				Description:         "DefaultAuditDeleteExpiration.",
				MarkdownDescription: "DefaultAuditDeleteExpiration.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"retention_classes": schema.SingleNestedAttribute{
				Description:         "RetentionClasses.",
				MarkdownDescription: "RetentionClasses.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"retention_class": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
								"name":   types.StringType,
								"period": types.Int64Type,
							}},
							},
						},
						map[string]attr.Value{
							"retention_class": types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"name":   types.StringType,
								"period": types.Int64Type,
							}}, []attr.Value{}),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"retention_class": schema.ListNestedAttribute{
						Description:         "RetentionClass.",
						MarkdownDescription: "RetentionClass.",
						Required:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description:         "Name",
									MarkdownDescription: "Name",
									Required:            true,
								},
								"period": schema.Int64Attribute{
									Description:         "Period.",
									MarkdownDescription: "Period.",
									Required:            true,
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
	}
}

func (r *NamespaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "creating namespace")
	var plan models.NamespaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	namespace, err := helper.BuildNamespaceFromPlan(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building namespace from plan", err.Error())
		return
	}

	namespace, err = r.client.ManagementClient.CreateNamespace(namespace)

	if err != nil {
		resp.Diagnostics.AddError("Error creating namespace", err.Error())
		return
	}

	data := models.NamespaceEntity{}
	err = helper.CopyFields(ctx, namespace, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting created namespace",
			err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "reading namespace")
	var data models.NamespaceEntity

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	namespace, err := r.client.ManagementClient.GetNamespace(data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error reading namespace", err.Error())
		return
	}

	err = helper.CopyFields(ctx, namespace, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting read namespace",
			err.Error(),
		)
		return
	}
	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "updating namespace")
	var plan models.NamespaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	namespace, err := helper.BuildNamespaceFromPlan(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building namespace from plan", err.Error())
		return
	}
	// To update the Id whose value does not exist in the plan from the state.
	// For the rest of the non-existing fields' value in the plan won't impact the update result,
	// as the update API would check the difference of the local value and remote value internally,
	// it would use the Id to retrieve the remote value,
	// and the non change value won't trigger the update
	// For the update API, it should use the same get API to get the remote value,
	// so just to refer get API definition to make sure all the required fields have the value assigned
	var data models.NamespaceEntity
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	namespace.Id = data.Id.ValueString()

	// TODO: To prevent the non-updatable fields from being changed

	_, err = r.client.ManagementClient.UpdateNamespace(namespace)
	if err != nil {
		resp.Diagnostics.AddError("Error updating namespace", err.Error())
		return
	}

	namespace, err = r.client.ManagementClient.GetNamespace(namespace.Id)

	if err != nil {
		resp.Diagnostics.AddError("Error reading namespace", err.Error())
		return
	}

	err = helper.CopyFields(ctx, namespace, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting read namespace",
			err.Error(),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "deleting namespace")
	var data models.NamespaceEntity

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ManagementClient.DeleteNamespace(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting namespace",
			err.Error(),
		)
	}
}

func (r *NamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "importing namespace")
	id := req.ID

	namespace, err := r.client.ManagementClient.GetNamespace(id)

	if err != nil {
		resp.Diagnostics.AddError("Error reading namespace", err.Error())
		return
	}

	data := models.NamespaceEntity{}
	err = helper.CopyFields(ctx, namespace, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting imported namespace",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
