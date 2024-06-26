package models

import "github.com/hashicorp/terraform-plugin-framework/types"


type AccountDatasource struct {
	ID             types.String               `tfsdk:"id"`
	Accounts       []AccountDatasourceEntity  `tfsdk:"accounts"`
}

type AccountDatasourceEntity struct {
	AccountId         types.String `tfsdk:"account_id"`
	Alias             types.String `tfsdk:"alias"`
	Description       types.String `tfsdk:"description"`
	EncryptionEnabled types.Bool   `tfsdk:"encryption_enabled"`
}
