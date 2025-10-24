#!/bin/bash

# govchain Volunteer Node Setup
# Allows volunteers to join the network as validators without tokens

set -euo pipefail

if ! command -v curl >/dev/null 2>&1; then
    echo "❌ Error: curl is required but was not found in PATH." >&2
    exit 1
fi

echo "================================"
echo "govchain Volunteer Node Setup"
echo "================================"
echo ""

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <node-name> <genesis-file-url> [genesis-sha256]"
    echo ""
    echo "Example:"
    echo "  $0 volunteer-node-1 https://raw.githubusercontent.com/org/govchain/main/genesis.json <sha256>"
    echo ""
    exit 1
fi

NODE_NAME="$1"
GENESIS_URL="$2"
GENESIS_SHA256="${3:-}"

echo "📝 Node Name: $NODE_NAME"
echo "🌐 Genesis URL: $GENESIS_URL"
echo ""

# Locate govchaind binary
if [ -f "./build/govchaind" ]; then
    GOVCHAIND="./build/govchaind"
elif command -v govchaind >/dev/null 2>&1; then
    GOVCHAIND="$(command -v govchaind)"
else
    echo "❌ Error: govchaind binary not found."
    echo "Please build or install it first, e.g.:"
    echo "  ignite chain build"
    echo "or ensure 'govchaind' is in your PATH."
    exit 1
fi

echo "⚙️ Using govchaind binary at: $GOVCHAIND"
echo ""

# Initialize node
echo "🔧 Initializing volunteer node..."
"$GOVCHAIND" init "$NODE_NAME" --chain-id govchain

# Download genesis file
echo "📥 Downloading genesis file..."
TEMP_GENESIS=$(mktemp /tmp/genesis.XXXXXX)

cleanup() {
    [ -f "$TEMP_GENESIS" ] && rm -f "$TEMP_GENESIS"
}
trap cleanup EXIT

if ! curl -fSL --retry 3 --retry-delay 2 "$GENESIS_URL" -o "$TEMP_GENESIS"; then
    echo "❌ Error: failed to download genesis from $GENESIS_URL" >&2
    exit 1
fi

if [ -n "$GENESIS_SHA256" ]; then
    if ! command -v sha256sum >/dev/null 2>&1; then
        echo "❌ Error: sha256sum is required when providing a genesis checksum." >&2
        exit 1
    fi
    echo "🔐 Verifying genesis checksum..."
    echo "$GENESIS_SHA256  $TEMP_GENESIS" | sha256sum -c -
else
    echo "⚠️  Warning: No genesis checksum provided. Skipping integrity verification."
fi

mv "$TEMP_GENESIS" "$HOME/.govchain/config/genesis.json"
trap - EXIT

# Create validator key
echo "🔑 Creating validator key..."
"$GOVCHAIND" keys add validator --keyring-backend test

# Get validator address
VALIDATOR_ADDR=$("$GOVCHAIND" keys show validator -a --keyring-backend test)

echo ""
echo "✅ Volunteer node setup complete!"
echo "================================"
echo ""
echo "📝 Node Details:"
echo "  Name: $NODE_NAME"
echo "  Validator Address: $VALIDATOR_ADDR"
echo ""
echo "🚀 To start your volunteer node:"
echo "  $GOVCHAIND start"
echo ""
echo "🌐 To become a validator:"
echo "  $GOVCHAIND tx staking create-validator \\"
echo "    --amount=1000000stake \\"
echo "    --pubkey=\$($GOVCHAIND tendermint show-validator) \\"
echo "    --moniker=\"$NODE_NAME\" \\"
echo "    --chain-id=govchain \\"
echo "    --from=validator \\"
echo "    --keyring-backend=test"
echo ""
