package paths

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/handlers"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

func Sign(m *stellar.Manager) *framework.Path {
	return &framework.Path{
		Pattern:      "accounts/" + framework.GenericNameRegex("publicKey") + "/sign",
		HelpSynopsis: "Sign a provided Stellar transaction envelope.",
		HelpDescription: `

    Sign a Stellar transaction envelope with the secret key of the specified account.

    `,
		Fields: map[string]*framework.FieldSchema{
			"publicKey": {
				Type:        framework.TypeString,
				Description: "The public key of the account to use for signing.",
			},
			"transaction": {
				Type:        framework.TypeString,
				Description: "The base64 encoded Stellar transaction envelope to sign.",
			},
			"network": {
				Type:        framework.TypeString,
				Description: "The network for the transaction ('Public' or 'Testnet').",
			},
		},
		ExistenceCheck: m.AccountExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: handlers.NewSignTxHandler(m),
		},
	}
}
