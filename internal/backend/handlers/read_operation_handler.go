package handlers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

type ReadAccountHandler struct {
	manager *stellar.Manager
}

func NewReadAccountHandler(m *stellar.Manager) *ReadAccountHandler {
	return &ReadAccountHandler{manager: m}
}

func (h *ReadAccountHandler) Handler() framework.OperationFunc {
	return h.manager.ReadAccount
}

func (h *ReadAccountHandler) Properties() framework.OperationProperties {
	return framework.OperationProperties{
		Summary:     "Reads a Stellar account",
		Description: "Retrieves information about a specified Stellar account using its public key.",
		Examples: []framework.RequestExample{
			{
				Description: "Read a Stellar account",
				Data: map[string]interface{}{
					"publicKey": "publicKeyToRead",
				},
				Response: &framework.Response{
					Description: "Successful retrieval of the Stellar account",
					MediaType:   "application/json",
					Fields: map[string]*framework.FieldSchema{
						"public_key": {
							Type:        framework.TypeString,
							Description: "The public key of the Stellar account",
						},
					},
					Example: &logical.Response{
						Data: map[string]interface{}{
							"public_key": "ExamplePublicKey",
						},
					},
				},
			},
		},
	}
}
