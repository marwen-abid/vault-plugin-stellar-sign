package stellar

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

// Account is the structure of a Stellar account
type Account struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key,omitempty"`
}

type Manager struct {
	logger hclog.Logger
}

func NewManager(logger hclog.Logger) *Manager {
	return &Manager{logger: logger}
}

func (m *Manager) ListAccounts(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// List all the stored accounts under the "stellar/accounts/" path
	accountList, err := req.Storage.List(ctx, "stellar/accounts/")
	if err != nil {
		m.logger.Error("Failed to list stellar accounts", "error", err)
		return nil, fmt.Errorf("failed to list stellar accounts: %s", err)
	}

	// Return the list of accounts
	return logical.ListResponse(accountList), nil
}

func (m *Manager) CreateAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	secretKeyInput := data.Get("secret_key").(string)
	var pair *keypair.Full
	var err error

	if secretKeyInput != "" {
		pair, err = keypair.ParseFull(secretKeyInput)
		if err != nil {
			m.logger.Error("Error parsing input secret key", "error", err)
			return nil, fmt.Errorf("error parsing input secret key")
		}
	} else {
		pair, err = keypair.Random()
		if err != nil {
			m.logger.Error("Error generating new keypair", "error", err)
			return nil, fmt.Errorf("error generating new keypair")
		}
	}

	publicKey := pair.Address()
	secretKey := pair.Seed()

	accountPath := fmt.Sprintf("stellar/accounts/%s", publicKey)

	accountJSON := &Account{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}

	entry, _ := logical.StorageEntryJSON(accountPath, accountJSON)
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		m.logger.Error("Failed to save the new stellar account to storage", "error", err)
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"public_key": accountJSON.PublicKey,
		},
	}, nil
}

func (m *Manager) ReadAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if req.Operation != logical.ReadOperation {
		return nil, fmt.Errorf("unsupported operation: %s", req.Operation)
	}

	publicKey := data.Get("publicKey").(string)
	if publicKey == "" {
		return nil, fmt.Errorf("missing public key")
	}

	m.logger.Info("Retrieving Stellar account for public key", "publicKey", publicKey)
	account, err := m.retrieveAccount(ctx, req.Storage, publicKey)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, fmt.Errorf("stellar account does not exist")
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"public_key": account.PublicKey,
		},
	}, nil
}

func (m *Manager) DeleteAccount(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	publicKey := data.Get("publicKey").(string)

	account, err := m.retrieveAccount(ctx, req.Storage, publicKey)
	if err != nil {
		m.logger.Error("Failed to retrieve the Stellar account by public key", "publicKey", publicKey, "error", err)
		return nil, err
	}
	if account == nil {
		return nil, nil
	}
	if err = req.Storage.Delete(ctx, fmt.Sprintf("stellar/accounts/%s", account.PublicKey)); err != nil {
		m.logger.Error("Failed to delete the Stellar account from storage", "publicKey", publicKey, "error", err)
		return nil, err
	}
	return nil, nil
}

type signRequest struct {
	publicKey         string
	txEnvelopeBase64  string
	networkPassphrase string
}

func validateSignRequest(data *framework.FieldData) (*signRequest, error) {
	publicKey := data.Get("publicKey").(string)
	if publicKey == "" {
		return nil, fmt.Errorf("publicKey must be provided")
	}

	txEnvelopeBase64 := data.Get("transaction").(string)
	if txEnvelopeBase64 == "" {
		return nil, fmt.Errorf("transaction must be provided")
	}

	networkParam := data.Get("network").(string)
	var networkPassphrase string
	switch networkParam {
	case "Public":
		networkPassphrase = network.PublicNetworkPassphrase
	case "Testnet":
		networkPassphrase = network.TestNetworkPassphrase
	default:
		return nil, fmt.Errorf("invalid network: %s", networkParam)
	}

	return &signRequest{
		publicKey:         publicKey,
		txEnvelopeBase64:  txEnvelopeBase64,
		networkPassphrase: networkPassphrase,
	}, nil
}

func (m *Manager) SignTx(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	sr, err := validateSignRequest(data)
	if err != nil {
		return nil, err
	}

	// Retrieve the account from storage
	account, err := m.retrieveAccount(ctx, req.Storage, sr.publicKey)
	if err != nil {
		m.logger.Error("Error retrieving account", "error", err)
		return nil, fmt.Errorf("error retrieving account: %s", err)
	}
	if account == nil {
		return nil, fmt.Errorf("account not found")
	}

	signedTxBase64, errSign := m.sign(account, sr.txEnvelopeBase64, sr.networkPassphrase)
	if errSign != nil {
		m.logger.Error("Error signing transaction", "error", errSign)
		return nil, fmt.Errorf("error signing transaction: %s", errSign)
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signed_transaction": signedTxBase64,
		},
	}, nil
}

func (m *Manager) sign(account *Account, txEnvelopeBase64 string, networkPassphrase string) (string, error) {
	// Decode the transaction envelope
	txEnvelope, err := txnbuild.TransactionFromXDR(txEnvelopeBase64)
	if err != nil {
		m.logger.Error("Error decoding transaction envelope", "error", err)
		return "", fmt.Errorf("error decoding transaction envelope: %s", err)
	}

	// Convert to a Transaction object
	tx, ok := txEnvelope.Transaction()
	if !ok {
		return "", fmt.Errorf("failed to convert to Transaction object")
	}

	// Sign the transaction
	kp, err := keypair.ParseFull(account.SecretKey)
	if err != nil {
		m.logger.Error("Error parsing keypair", "error", err)
		return "", fmt.Errorf("error parsing keypair: %s", err)
	}

	signedTx, err := tx.Sign(networkPassphrase, kp)
	if err != nil {
		m.logger.Error("Error signing transaction", "error", err)
		return "", fmt.Errorf("error signing transaction: %s", err)
	}

	// Convert the signed transaction to base64
	signedTxBase64, err := signedTx.Base64()
	if err != nil {
		m.logger.Error("Error encoding signed transaction", "error", err)
		return "", fmt.Errorf("error encoding signed transaction: %s", err)
	}

	return signedTxBase64, nil
}

func (m *Manager) retrieveAccount(ctx context.Context, storage logical.Storage, publicKey string) (*Account, error) {
	_, err := keypair.ParseAddress(publicKey)
	if err != nil {
		m.logger.Error("Failed to retrieve the account, invalid Stellar public key", "publicKey", publicKey, "error", err)
		return nil, fmt.Errorf("failed to retrieve the account, invalid Stellar public key: %s", err)
	}

	path := fmt.Sprintf("stellar/accounts/%s", publicKey)
	entry, err := storage.Get(ctx, path)
	if err != nil {
		m.logger.Error("Failed to retrieve the account by public key", "path", path, "error", err)
		return nil, err
	}
	if entry == nil {
		// Could not find the corresponding account for the public key
		return nil, nil
	}
	var account Account
	_ = entry.DecodeJSON(&account)
	return &account, nil
}

func (m *Manager) AccountExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		m.logger.Error("Path existence check failed", err)
		return false, fmt.Errorf("existence check failed: %v", err)
	}

	return out != nil, nil
}
