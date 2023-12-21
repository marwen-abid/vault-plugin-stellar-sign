package handlers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

type CreateAccountHandler struct {
	manager *stellar.Manager
}

func NewCreateAccountHandler(m *stellar.Manager) *CreateAccountHandler {
	return &CreateAccountHandler{manager: m}
}

func (h *CreateAccountHandler) Handler() framework.OperationFunc {
	return h.manager.CreateAccount
}

func (h *CreateAccountHandler) Properties() framework.OperationProperties {
	return framework.OperationProperties{
		Summary:     "Creates a new Stellar account",
		Description: "Generates a new Stellar account and stores it in the backend.",
		Examples: []framework.RequestExample{
			{
				Description: "Create a new Stellar account",
				Data: map[string]interface{}{
					"publicKey": "NewPublicKey",
				},
				Response: &framework.Response{
					Description: "Successful creation of a new Stellar account",
					MediaType:   "application/json",
					Fields: map[string]*framework.FieldSchema{
						"public_key": {
							Type:        framework.TypeString,
							Description: "The public key of the newly created Stellar account",
						},
						"secret_key": {
							Type:        framework.TypeString,
							Description: "The secret key of the newly created Stellar account",
						},
					},
					Example: &logical.Response{
						Data: map[string]interface{}{
							"public_key": "NewPublicKey",
							"secret_key": "NewSecretKey",
						},
					},
				},
			},
		},
	}
}
