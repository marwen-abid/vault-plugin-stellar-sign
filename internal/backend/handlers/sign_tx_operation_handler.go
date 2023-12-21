package handlers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

type SignTxHandler struct {
	manager *stellar.Manager
}

func NewSignTxHandler(m *stellar.Manager) *SignTxHandler {
	return &SignTxHandler{manager: m}
}

func (h *SignTxHandler) Handler() framework.OperationFunc {
	return h.manager.SignTx
}

func (h *SignTxHandler) Properties() framework.OperationProperties {
	return framework.OperationProperties{
		Summary: "Signs a Stellar transaction envelope",
		Description: "This operation signs a provided Stellar transaction envelope using the secret key " +
			"of the specified account. The transaction is specified in a base64-encoded format, " +
			"and the account is identified by its public key.",
		Examples: []framework.RequestExample{
			{
				Description: "Sign a transaction with a specific account",
				Data: map[string]interface{}{
					"publicKey":   "GATBMIXGZKJGSEVJQH7D2ZP3A4UQ4WKB3X5H3C6KHPGJRH4B3U5UJ6CH",
					"transaction": "base64EncodedTransactionEnvelope",
					"network":     "Public",
				},
				Response: &framework.Response{
					Description: "Successful signing of the Stellar transaction",
					MediaType:   "application/json",
					Fields: map[string]*framework.FieldSchema{
						"signed_transaction": {
							Type:        framework.TypeString,
							Description: "The base64 encoded signed Stellar transaction envelope",
						},
					},
					Example: &logical.Response{
						Data: map[string]interface{}{
							"signed_transaction": "base64EncodedSignedTransactionEnvelope",
						},
					},
				},
			},
		},
	}
}
