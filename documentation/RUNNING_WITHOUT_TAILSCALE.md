# Running a GovChain Node without Tailscale

This guide explains how to run a `govchaind` node using Docker without the Tailscale sidecar. This is useful for local development or for environments where you are managing network exposure manually.

## Prerequisites

- **Docker**: Ensure Docker is installed and running on your system.
- **GovChain Docker Image**: You must have the `govchaind` Docker image available locally. You can build it from the source using `make docker-build` or pull it from the GitHub Container Registry.

---

## Option 1: Using a Base Docker Compose File (without Tailscale)

This is the recommended method for running your node without Tailscale, as it is easy to manage. You can choose between a production setup (pulling a pre-built image) or a local development setup (building the image locally).

### 1.1. For Production/Deployment (pulling image from GHCR)

Use the `docker-compose.prod.yaml` file.

```yaml
# docker-compose.prod.yaml
version: '3.8'

services:
  govchaind:
    image: ghcr.io/bettergovph/govchaind:latest
    ports:
      - "26656:26656"
      - "26657:26657"
    volumes:
      - govchaind_data:/home/nonroot/.govchain
    environment:
      # Add any production specific environment variables here
      # For example:
      # - NODE_ENV=production
volumes:
  govchaind_data:
```

**How to Use:**
To start the node:
```bash
docker compose -f docker-compose.prod.yaml up -d
```

### 1.2. For Local Development (building image locally)

Use the `docker-compose.local.yaml` file.

```yaml
# docker-compose.local.yaml
version: '3.8'

services:
  govchaind:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "26656:26656"
      - "26657:26657"
    volumes:
      - govchaind_data:/home/nonroot/.govchain
    environment:
      # Add any local development specific environment variables here
      # For example:
      # - DEBUG=true
    # command: start --log_level debug # Example of overriding default command
volumes:
  govchaind_data:
```

**How to Use:**
To build and start the node:
```bash
docker compose -f docker-compose.local.yaml build
docker compose -f docker-compose.local.yaml up -d
```

3.  **Monitor Logs**:
    ```bash
    docker logs -f govchaind-node
    ```

4.  **Stop the Node**:
    ```bash
    docker-compose -f docker-compose.standalone.yml down
    ```

---

## Option 2: Using `docker run`

If you prefer not to use Docker Compose, you can run the container directly with the `docker run` command.

### Command

```bash
docker run -d --name govchaind-node \
  -p 26656:26656 \
  -p 26657:26657 \
  -p 1317:1317 \
  -p 9090:9090 \
  -v govchaind-data:/home/nonroot/.govchain \
  -e MONIKER="My Standalone Node" \ # Customize your node's name
  -e EXTERNAL_IP="YOUR_PUBLIC_IP_OR_DOMAIN" \ # Optional: Manually set external IP
  ghcr.io/bettergovph/govchaind:latest
```

### Explanation

- `-d`: Runs the container in detached mode.
- `--name govchaind-node`: Assigns a name to the container for easy reference.
- `-p ...`: Maps the necessary ports from your host machine to the container.
- `-v govchaind-data:/home/nonroot/.govchain`: Creates a named volume `govchaind-data` to persist blockchain data across container restarts.
- `-e MONIKER="..."`: Sets the moniker for your node. This overrides the default set in the Docker image.
- `-e EXTERNAL_IP="..."`: **(New)** Manually sets the external IP address for the node. If this is provided, the entrypoint script will use it directly, bypassing automatic detection.
- `ghcr.io/bettergovph/govchaind:latest`: The Docker image to run.

---

## How It Works: IP Address and Chain ID Configuration

The `docker-entrypoint.sh` script inside the container is designed to be resilient and versatile in configuring your node's network identity and chain ID:

### Moniker Configuration

The node's moniker (name) is set via the `MONIKER` environment variable. If not explicitly provided at runtime (e.g., via `-e MONIKER=...`), it defaults to "GovChain Validator" as defined in the Docker image.

### External IP Address Configuration

1.  **Manual Override**: If the `EXTERNAL_IP` environment variable is set (e.g., `-e EXTERNAL_IP="YOUR_PUBLIC_IP_OR_DOMAIN"`), the script will use this value directly to configure `external_address` in `config.toml`.
2.  **Tailscale Detection**: If `EXTERNAL_IP` is not set, the script first checks for a Tailscale sidecar. If a Tailscale IP is found, it will be used as the `external_address`.
3.  **Public IP Discovery**: If neither `EXTERNAL_IP` nor a Tailscale IP is detected, the script automatically attempts to discover the host machine's public IP address using an external service. If found, it configures `external_address` in `config.toml` with that IP (e.g., `external_address = "YOUR_PUBLIC_IP:26656"`). This is essential for VPS or public cloud deployments.
4.  **No Configuration**: If no IP can be determined through any of the above methods, the node will start without a pre-configured `external_address`.

### Chain ID Configuration

The `chain_id` for your node is now dynamically extracted from the `genesis.json` file that is downloaded during the node's initialization. This ensures that your node always uses the correct and authoritative `chain_id` for the network it is joining, removing any need for manual configuration of this value.
