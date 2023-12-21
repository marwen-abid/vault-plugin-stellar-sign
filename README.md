# Vault Plugin: Stellar Sign

## Overview
`vault-plugin-stellar-sign` is a Vault plugin designed to manage and sign Stellar blockchain transactions. It provides functionality to create and manage Stellar accounts and to sign transactions using stored account keys. 

This plugin is inspired by the work done in https://github.com/kaleido-io/vault-plugin-secrets-ethsign. 

## Building the Plugin
To build the plugin and generate the `stellar-sign` executable:

```bash
make all
```

## Running the Plugin with Docker

To start a local Vault instance with the `stellar-sign` plugin loaded, available at http://127.0.0.1:8200/:

```bash
make docker
```

## Endpoints

### Creating New Signing Account
Creates a new Stellar account and stores its key.

**Request:**
```bash
curl --location --request POST 'http://127.0.0.1:8200/v1/stellar/accounts' \
--header 'Authorization: Bearer root'
```

**Response**
```json
{
    "request_id": "0ff83110-354e-6f20-a5c2-7cb4e65084e4",
    "data": {
        "public_key": "GBG2QXP6SLJFEVWDUXE23JW2OKQIFTO2Q6ACDTNHJRKVPXAWQHKN7QKW"
    }
}
```

### List Existing Accounts
Lists all stored Stellar accounts.

**Request:**
```bash
curl --location 'http://localhost:8200/v1/stellar/accounts?list=true' \
--header 'Authorization: Bearer root'
```

**Response**
```json
{
  "request_id": "1e89ca73-c20f-efe0-6e46-1bf291e9bd4d",
  "data": {
    "keys": [
      "GAA3JQIL4BDWEF2RBHWVCMQM36DGR2FWLKBSSPSZG7ONS6KTLYP3BTKD",
      "GBG2QXP6SLJFEVWDUXE23JW2OKQIFTO2Q6ACDTNHJRKVPXAWQHKN7QKW"
    ]
  }
}
```

### Reading Individual Account
Retrieves details of a specific Stellar account.

**Request:**
```bash
curl --location 'http://localhost:8200/v1/stellar/accounts/GAA3JQIL4BDWEF2RBHWVCMQM36DGR2FWLKBSSPSZG7ONS6KTLYP3BTKD' \
--header 'Authorization: Bearer root'
```

**Response**
```json
{
  "request_id": "7690417c-2d4e-99dd-210a-9290e857f9da",
  "data": {
    "public_key": "GAA3JQIL4BDWEF2RBHWVCMQM36DGR2FWLKBSSPSZG7ONS6KTLYP3BTKD"
  }
}
```

### Delete an Account
Deletes a specified Stellar account.

**Request:**
```bash
curl --location --request DELETE 'http://localhost:8200/v1/stellar/accounts/GDLKO7LP4TCQBNA5UYL7CEP3H4GVATZHWV2FLTDDTEX6HH4S4RHLEBEY' \
--header 'Authorization: Bearer root'
```

**Response**
```
204 No Content 
```

### Sign a Transaction
Signs a Stellar transaction with the specified account.

**Request:**
```bash
curl --location 'http://127.0.0.1:8200/v1/stellar/accounts/GDSKR6UYBIYIU7GGVPIUZCX6C7EWG5VCRC2VCCH5NVFLBWMOSLBDBLHW/sign' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer root' \
--data '{
        "transaction": "AAAAAgAAAAATozPrNDRTqLO2WUflkFsbKLSQN79/VlhRpv7MMzePdgAAAGQAAMGGAAAAAQAAAAEAAAAAAAAAAAAAAABlhHryAAAAAAAAAAEAAAABAAAAABOjM+s0NFOos7ZZR+WQWxsotJA3v39WWFGm/swzN492AAAAAQAAAAB69J8A290AJGAqNy4f0QIXBG4NoPQm7B+vDdeR0AvXRQAAAAAAAAACVAvkAAAAAAAAAAAA",
        "network": "Testnet"
}'
```

**Response**
```json
{
  "request_id": "ff99bde8-1589-f6ea-b963-188730706bd4",
  "data": {
    "signed_transaction": "AAAAAgAAAAATozPrNDRTqLO2WUflkFsbKLSQN79/VlhRpv7MMzePdgAAAGQAAMGGAAAAAQAAAAEAAAAAAAAAAAAAAABlhHryAAAAAAAAAAEAAAABAAAAABOjM+s0NFOos7ZZR+WQWxsotJA3v39WWFGm/swzN492AAAAAQAAAAB69J8A290AJGAqNy4f0QIXBG4NoPQm7B+vDdeR0AvXRQAAAAAAAAACVAvkAAAAAAAAAAABU14fsAAAAEASISX5s51KVscnLwbjl/0kU8I47SmhFo+Ldn2+rugtHCBQmunwH994JhUDCT2ra0WrMyRDNoGYQL4xaKrh65YP"
  }
}

```
