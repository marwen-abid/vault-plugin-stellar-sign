package handlers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

type DeleteAccountHandler struct {
	manager *stellar.Manager
}

func NewDeleteAccountHandler(m *stellar.Manager) *DeleteAccountHandler {
	return &DeleteAccountHandler{manager: m}
}

func (h *DeleteAccountHandler) Handler() framework.OperationFunc {
	return h.manager.DeleteAccount
}

func (h *DeleteAccountHandler) Properties() framework.OperationProperties {
	return framework.OperationProperties{
		Summary:     "Deletes a Stellar account",
		Description: "Removes a specified Stellar account from storage using its public key.",
		Examples: []framework.RequestExample{
			{
				Description: "Delete a Stellar account",
				Data: map[string]interface{}{
					"publicKey": "GASYNOBIVOZZJGBH6C5K2FPJQ2RPN2QBJ7OO6OXENIG5S2IKRALFTKIA",
				},
				Response: &framework.Response{
					Description: "Successful deletion of the Stellar account",
					MediaType:   "application/json",
				},
			},
		},
	}
}
