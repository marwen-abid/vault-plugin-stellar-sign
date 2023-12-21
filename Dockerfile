FROM alpine:latest as builder

# Install Go and build dependencies
RUN apk add --no-cache go gcc musl-dev

WORKDIR /go/src/app
COPY . .
RUN go build -o stellar-sign

# Use Vault's official image
FROM hashicorp/vault:latest

# Copy the compiled plugin from the builder stage
COPY --from=builder /go/src/app/stellar-sign /vault/plugins/stellar-sign

# Expose port 8200 for Vault
EXPOSE 8200

## Start Vault with the plugin directory specified
CMD ["vault", "server", "-dev", "-dev-root-token-id=root", "-dev-plugin-dir=/vault/plugins"]
