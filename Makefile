
BINARY_NAME=stellar-sign
SRC_GOFILES := $(shell find . -name '*.go' -print)
.DELETE_ON_ERROR:

all: build test
test: deps
		go test  ./... -cover -coverprofile=coverage.txt -covermode=atomic
stellarsign: ${SRC_GOFILES}
		go build -o ${BINARY_NAME} -ldflags "-X main.GitCommit=$(LABEL)" .
build: stellarsign
clean: 
		go clean
		rm -f ${BINARY_NAME}
deps:
		go get
docker:
	# Start up Docker Compose services
	docker-compose up --build -d

	# Wait for Vault to become ready
	@echo "Waiting for Vault to start..."
	@while ! docker exec vault_stellar_sign vault status > /dev/null 2>&1; do \
		sleep 1; \
	done

	# Enable the plugin
	docker exec vault_stellar_sign vault secrets enable -path=stellar -description="Stellar Wallet" -plugin-name=stellar-sign plugin
