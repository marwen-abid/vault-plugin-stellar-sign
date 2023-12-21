package paths

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/handlers"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

func CreateAndList(m *stellar.Manager) *framework.Path {
	return &framework.Path{
		Pattern:      "accounts/?",
		HelpSynopsis: "List all the Stellar accounts maintained by the plugin backend and create new accounts.",
		HelpDescription: `

    LIST - list all accounts
    POST - create a new account

    `,
		Fields: map[string]*framework.FieldSchema{
			"secret_key": {
				Type:        framework.TypeString,
				Description: "Base64 encoded string representing the Stellar secret key. If provided, the request will import this key instead of generating a new one. The secret key is used to sign transactions and should be kept private.",
				Default:     "",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation:   handlers.NewListAccountsHandler(m),
			logical.UpdateOperation: handlers.NewCreateAccountHandler(m),
		},
	}
}
