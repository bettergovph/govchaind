#!/bin/sh
set -e

# Exit if the MONIKER environment variable is not set.
if [ -z "$MONIKER" ]; then
  echo "Error: The MONIKER environment variable is not set."
  echo "Please set it using -e MONIKER='your-moniker' with docker run, or in your docker-compose.yml."
  exit 1
fi

# Ensure the config directory exists
mkdir -p "/home/nonroot/.govchain/config"
chown nonroot:nonroot "/home/nonroot/.govchain/config"

# Check if genesis.json exists. If not, initialize and then replace genesis.json
if [ ! -f "/home/nonroot/.govchain/config/genesis.json" ]; then
  # Define the path for the actual genesis.json
  ACTUAL_GENESIS_PATH="/home/nonroot/.govchain/config/genesis.json"
  TEMP_GENESIS_PATH="/tmp/genesis.json" # Use a temporary path for download

  # Download the actual genesis.json to a temporary location
  echo "📥 Downloading actual genesis file to extract chain-id..."
  env 'HOME=/home/nonroot' curl -sL "https://raw.githubusercontent.com/bettergovph/govchain/refs/heads/main/genesis.json" -o "$TEMP_GENESIS_PATH"

  # Extract chain-id from the downloaded genesis.json
  CHAIN_ID=$(jq -r '.chain_id' "$TEMP_GENESIS_PATH")
  if [ -z "$CHAIN_ID" ]; then
    echo "Error: Could not extract chain_id from downloaded genesis file."
    exit 1
  fi
  echo "🌐 Extracted chain-id: $CHAIN_ID"

  # Initialize the node using the extracted chain-id
  echo "🔧 Initializing node with chain-id: $CHAIN_ID..."
  env 'HOME=/home/nonroot' govchaind init "$MONIKER" --chain-id "$CHAIN_ID" --home "/home/nonroot/.govchain"

  # Replace the dummy genesis.json created by init with the actual downloaded one
  echo "Replacing dummy genesis.json with the actual genesis file..."
  mv "$TEMP_GENESIS_PATH" "$ACTUAL_GENESIS_PATH"

  # Set the minimum gas price in app.toml
  env 'HOME=/home/nonroot' sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0stake"/' "/home/nonroot/.govchain/config/app.toml"
fi

# IP address configuration: Manual > Tailscale > Public IP
if [ -n "$EXTERNAL_IP" ]; then
  # 1. Use manually provided IP if available
  echo "Using manually provided external IP: $EXTERNAL_IP. Updating config.toml..."
  sed -i "s/^external_address = \"\"/external_address = \"$EXTERNAL_IP:26656\"/" "/home/nonroot/.govchain/config/config.toml"
  echo "config.toml updated with manual IP."
else
  # 2. If no manual IP, start automatic detection: Try Tailscale first
  echo "Waiting for Tailscale IP file from sidecar..."
  TS_IP_FILE="/var/run/tailscale-ip/ts_ip"
  TS_IP=""
  for i in $(seq 1 30); do
    if [ -f "$TS_IP_FILE" ]; then
      TS_IP=$(cat "$TS_IP_FILE")
      if [ -n "$TS_IP" ]; then
        echo "Tailscale IP found: $TS_IP"
        break
      fi
    fi
    echo "Waiting for Tailscale IP file... ($i/30)"
    sleep 2
  done

  if [ -n "$TS_IP" ]; then
    echo "Tailscale IP: $TS_IP. Updating config.toml..."
    sed -i "s/^external_address = \"\"/external_address = \"$TS_IP:26656\"/" "/home/nonroot/.govchain/config/config.toml"
    echo "config.toml updated with Tailscale IP."
  else
    # 3. If Tailscale fails, try to discover public IP
    echo "Could not get Tailscale IP. Attempting to discover public IP for VPS setup..."
    PUBLIC_IP=$(curl -s api.ipify.org)

    if [ -n "$PUBLIC_IP" ]; then
      echo "Discovered public IP: $PUBLIC_IP. Updating config.toml..."
      sed -i "s/^external_address = \"\"/external_address = \"$PUBLIC_IP:26656\"/" "/home/nonroot/.govchain/config/config.toml"
      echo "config.toml updated with public IP."
    else
      echo "Could not discover public IP. Proceeding without updating external_address."
    fi
  fi
fi

# Add persistent peer to config.toml
PERSISTENT_PEER="4d153c889d9f0f4b670d2f548994fcdde208240e@157.90.134.175:26656"
sed -i "s/^persistent_peers = \"\"/persistent_peers = \"$PERSISTENT_PEER\"/" "/home/nonroot/.govchain/config/config.toml"
echo "Persistent peer added to config.toml."

# Start govchaind
govchaind "$@"