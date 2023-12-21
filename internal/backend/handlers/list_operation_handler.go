package handlers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

type ListAccountsHandler struct {
	manager *stellar.Manager
}

func NewListAccountsHandler(m *stellar.Manager) *ListAccountsHandler {
	return &ListAccountsHandler{manager: m}
}

func (h *ListAccountsHandler) Handler() framework.OperationFunc {
	return h.manager.ListAccounts
}

func (h *ListAccountsHandler) Properties() framework.OperationProperties {
	return framework.OperationProperties{
		Summary:     "Lists Stellar accounts",
		Description: "Retrieves a list of all Stellar accounts stored in the backend.",
		Examples: []framework.RequestExample{
			{
				Description: "List all Stellar accounts",
				Response: &framework.Response{
					Description: "Successful retrieval of Stellar account list",
					MediaType:   "application/json",
					Fields: map[string]*framework.FieldSchema{
						"accounts": {
							Type:        framework.TypeSlice,
							Description: "A list of stored Stellar account public keys",
						},
					},
					Example: &logical.Response{
						Data: map[string]interface{}{
							"accounts": []string{"PublicKey1", "PublicKey2"},
						},
					},
				},
			},
		},
	}
}
