package backend

import (
	"context"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"vault-plugin-stellar-sign/internal/backend/paths"
	"vault-plugin-stellar-sign/internal/backend/stellar"
)

// Factory returns the Backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := newBackend()
	if err != nil {
		return nil, err
	}
	if err = b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns the Backend
func newBackend() (*Backend, error) {
	var b Backend
	b.Backend = &framework.Backend{
		Help: "",
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"accounts/",
			},
		},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}

	stellarManager := stellar.NewManager(b.Logger())

	b.Paths = framework.PathAppend(
		stellarPaths(stellarManager),
	)

	return &b, nil
}

// Backend implements the Backend for this plugin
type Backend struct {
	*framework.Backend
}

func stellarPaths(sm *stellar.Manager) []*framework.Path {
	return []*framework.Path{
		paths.CreateAndList(sm),
		paths.ReadAndDelete(sm),
		paths.Sign(sm),
	}
}
