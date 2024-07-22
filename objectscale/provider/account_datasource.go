package provider

import (
	"context"
	"fmt"
	"terraform-provider-objectscale/client"
	"terraform-provider-objectscale/objectscale/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AccountDataSource{}

func NewAccountDataSource() datasource.DataSource {
	return &AccountDataSource{}
}

type AccountDataSource struct {
	client *client.Client
}

func (d *AccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

// Schema describes the data source arguments.
func (d *AccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Account.",
		Description:         "Account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Identifier",
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
			"accounts": schema.ListNestedAttribute{
				Description:         "List of Accounts",
				MarkdownDescription: "List of Accounts",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Description:         "AccountId.",
							MarkdownDescription: "AccountId.",
							Computed:            true,
						},
						"alias": schema.StringAttribute{
							Description:         "Alias.",
							MarkdownDescription: "Alias.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							Description:         "Description.",
							MarkdownDescription: "Description.",
							Computed:            true,
						},
						"encryption_enabled": schema.BoolAttribute{
							Description:         "EncryptionEnabled.",
							MarkdownDescription: "EncryptionEnabled.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	obsClient, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = obsClient
}

func (d *AccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.AccountDatasource
	var plan models.AccountDatasource

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accounts, err := d.client.ManagementClient.ListAccounts()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting the list of accounts",
			err.Error(),
		)
		return
	}

	var accountList []models.AccountDatasourceEntity
	// Convert from json to terraform model
	for _, account := range accounts {
		entity := models.AccountDatasourceEntity{
			AccountId:         types.StringValue(account.AccountId),
			Alias:             types.StringValue(account.Alias),
			Description:       types.StringValue(account.Description),
			EncryptionEnabled: types.BoolValue(account.EncryptionEnabled),
		}
		accountList = append(accountList, entity)
	}

	state.ID = types.StringValue("account_datasource")
	state.Accounts = accountList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
