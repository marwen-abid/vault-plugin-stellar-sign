package backend

import (
	"context"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCreateStellarAccountWithProvidedSecretKey tests account creation with a provided secret key.
func TestCreateStellarAccountWithProvidedSecretKey(t *testing.T) {
	b, storage := getTestBackendAndStorage(t)

	// Generate a keypair for testing
	pair, _ := keypair.Random()
	testSecretKey := pair.Seed()

	// Create a request with a provided secret key
	data := map[string]interface{}{
		"secret_key": testSecretKey,
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "accounts",
		Data:      data,
		Storage:   storage,
	}

	resp, err := b.HandleRequest(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, pair.Address(), resp.Data["public_key"])
}

// TestCreateStellarAccountWithoutProvidedSecretKey tests account creation without a provided secret key.
func TestCreateStellarAccountWithoutProvidedSecretKey(t *testing.T) {
	b, storage := getTestBackendAndStorage(t)

	// Create a request without a provided secret key
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "accounts",
		Data:      map[string]interface{}{},
		Storage:   storage,
	}

	resp, err := b.HandleRequest(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Data["public_key"])
	_, err = keypair.Parse(resp.Data["public_key"].(string))
	assert.NoError(t, err)
	assert.Empty(t, resp.Data["secret_key"])

}

// TestListStellarAccountsEmpty tests listing Stellar accounts when none have been created.
func TestListStellarAccountsEmpty(t *testing.T) {
	b, storage := getTestBackendAndStorage(t)

	// Create a request to list accounts
	req := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "accounts",
		Storage:   storage,
	}

	resp, err := b.HandleRequest(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Data["keys"])
}

// TestListStellarAccounts tests listing Stellar accounts after creating some.
func TestListStellarAccounts(t *testing.T) {
	b, storage := getTestBackendAndStorage(t)

	// Create a couple of Stellar accounts
	for i := 0; i < 2; i++ {
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "accounts",
			Data:      map[string]interface{}{},
			Storage:   storage,
		}

		_, err := b.HandleRequest(context.Background(), req)
		assert.NoError(t, err)
	}

	// Now list the accounts
	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "accounts",
		Storage:   storage,
	}

	listResp, err := b.HandleRequest(context.Background(), listReq)
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Len(t, listResp.Data["keys"], 2)
}

func TestSignStellarTx(t *testing.T) {
	b, storage := getTestBackendAndStorage(t)

	// Create a couple of Stellar accounts
	for i := 0; i < 2; i++ {
		req := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "accounts",
			Data:      map[string]interface{}{},
			Storage:   storage,
		}

		_, err := b.HandleRequest(context.Background(), req)
		assert.NoError(t, err)
	}

	// Now list the accounts
	listReq := &logical.Request{
		Operation: logical.ListOperation,
		Path:      "accounts",
		Storage:   storage,
	}

	listResp, err := b.HandleRequest(context.Background(), listReq)
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Len(t, listResp.Data["keys"], 2)

	// Create a request to sign a transaction
	for i := 0; i < 2; i++ {
		data := map[string]interface{}{
			"transaction": "AAAAAgAAAAATozPrNDRTqLO2WUflkFsbKLSQN79/VlhRpv7MMzePdgAAAGQAAMGGAAAAAQAAAAEAAAAAAAAAAAAAAABlhHryAAAAAAAAAAEAAAABAAAAABOjM+s0NFOos7ZZR+WQWxsotJA3v39WWFGm/swzN492AAAAAQAAAAB69J8A290AJGAqNy4f0QIXBG4NoPQm7B+vDdeR0AvXRQAAAAAAAAACVAvkAAAAAAAAAAAA",
			"network":     "Testnet",
		}

		req := &logical.Request{
			Operation: logical.CreateOperation,
			Path:      "accounts/" + listResp.Data["keys"].([]string)[i] + "/sign",
			Data:      data,
			Storage:   storage,
		}

		resp, innerErr := b.HandleRequest(context.Background(), req)
		assert.NoError(t, innerErr)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Data["signed_transaction"])
	}

}

// getTestBackendAndStorage is a helper function to create a Backend and in-memory storage for testing.
func getTestBackendAndStorage(t *testing.T) (logical.Backend, logical.Storage) {
	config := logical.TestBackendConfig()
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	storage := &logical.InmemStorage{}
	return b, storage
}
