# --- Build Stage ---
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /src

# Copy go module and sum files
COPY go.mod go.sum ./

# Copy only necessary source code
COPY app/ app/
COPY cmd/ cmd/
COPY x/ x/
COPY proto/ proto/
COPY docs/ docs/
COPY testutil/ testutil/

# Download dependencies
RUN go version
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 go build -o /bin/govchaind ./cmd/govchaind

# --- Runtime Stage ---
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates dos2unix curl jq

# Create a non-root user
RUN adduser -D -u 1001 nonroot

# Copy the binary from the builder stage
COPY --from=builder /bin/govchaind /usr/local/bin/govchaind

# Copy and make executable the entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN dos2unix /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Set the user to nonroot
USER nonroot



# Expose the P2P and RPC ports
EXPOSE 26656 26657

HEALTHCHECK --interval=30s --timeout=10s --retries=5 \
    CMD curl --fail http://localhost:26657/status || exit 1

# Define a volume for the govchaind data
VOLUME /home/nonroot/.govchain
RUN mkdir -p /home/nonroot/.govchain && chown -R nonroot:nonroot /home/nonroot/.govchain




ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["start"]
