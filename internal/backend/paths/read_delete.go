package paths

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/handlers"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

func ReadAndDelete(m *stellar.Manager) *framework.Path {
	return &framework.Path{
		Pattern:      "accounts/" + framework.GenericNameRegex("publicKey"),
		HelpSynopsis: "Create, get or delete a Stellar account by publicKey",
		HelpDescription: `
			GET - return the account by the publicKey
			DELETE - deletes the account by the publicKey`,
		Fields: map[string]*framework.FieldSchema{
			"publicKey": {Type: framework.TypeString},
		},
		ExistenceCheck: m.AccountExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation:   handlers.NewReadAccountHandler(m),
			logical.DeleteOperation: handlers.NewDeleteAccountHandler(m),
		},
	}
}
