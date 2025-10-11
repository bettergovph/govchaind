# Local Development with Docker Compose

This guide provides comprehensive instructions for setting up and running your `govchaind` node for local development using Docker Compose. You can choose to run your node either directly or with Tailscale integration for secure networking.

## 1. Prerequisites

Before you begin, ensure you have the following installed and configured:

*   **Docker**: Make sure Docker Desktop (for macOS/Windows) or Docker Engine (for Linux) is installed and running on your system.
*   **Docker Compose**: This is usually bundled with Docker Desktop or can be installed separately.
*   **Git**: To clone the `govchaind` project repository.
*   **A Tailscale Account (Optional)**: Only if you plan to use Tailscale for local development.

## 2. Project Setup

1.  **Clone the `govchaind` repository**:
    ```bash
    git clone <repository_url>
    cd govchaind
    ```
    (Replace `<repository_url>` with the actual URL of the `govchaind` repository.)

2.  **Create your `.env` file**:
    The project includes a `.env.sample` file to guide you. Create a new file named `.env` in the root directory of the `govchaind` project by copying the sample:
    ```bash
    cp .env.sample .env
    ```

3.  **Populate `.env`**:
    Open the newly created `.env` file and populate it with necessary environment variables.
    *   `MONIKER`: Set your node's moniker (name).
    *   `TS_AUTHKEY` (Optional): If using Tailscale, replace the placeholder with your actual Tailscale authentication key (e.g., `tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY`).
    ```dotenv
    # .env
    MONIKER="My Local Dev Node"
    TS_AUTHKEY=tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY # Only if using Tailscale
    ```

## 3. Running for Local Development

You have two main options for running your `govchaind` node locally: without Tailscale, or with Tailscale integration.

### 3.1. Local Development without Tailscale

This setup uses `docker-compose.local.yaml` to build the `govchaind` image from your local source code and run it.

**`docker-compose.local.yaml`:**
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

1.  **Build the Docker Image:**
    ```bash
    docker compose -f docker-compose.local.yaml build
    ```
2.  **Start the Node:**
    ```bash
    docker compose -f docker-compose.local.yaml up -d
    ```

### 3.2. Local Development with Tailscale

This setup combines `docker-compose.local.yaml` with `docker-compose.tailscale.yaml` to build your `govchaind` image locally and integrate it with a Tailscale sidecar for secure networking.

**`docker-compose.tailscale.yaml`:**
```yaml
# docker-compose.tailscale.yaml
version: '3.8'

services:
  govchaind:
    volumes:
      - tailscale-ip-share:/var/run/tailscale-ip
    environment:
      - MONIKER=${MONIKER}
    networks:
      - govchain_network
    depends_on:
      - tailscale

  tailscale:
    image: tailscale/tailscale
    container_name: tailscale-sidecar
    cap_add:
      - NET_ADMIN
    devices:
      - /dev/net/tun
    volumes:
      - tailscale-state:/var/lib/tailscale
      - tailscale-ip-share:/var/run/tailscale-ip
    environment:
      - TS_AUTHKEY=${TS_AUTHKEY}
      - TS_STATE_DIR=/var/lib/tailscale
      - TS_HOSTNAME=my-node-ts
    networks:
      - govchain_network
    command: sh -c "tailscaled --state=/var/lib/tailscale/tailscaled.state --socket=/var/run/tailscale/tailscaled.sock & tailscale up --authkey=$TS_AUTHKEY --hostname=$TS_HOSTNAME --accept-routes && tailscale ip -4 > /var/run/tailscale-ip/ts_ip && sleep infinity"

networks:
  govchain_network:
    driver: bridge

volumes:
  tailscale-state:
  tailscale-ip-share:
```

**How to Use:**

1.  **Build the Docker Image:**
    ```bash
    docker compose -f docker-compose.local.yaml -f docker-compose.tailscale.yaml build
    ```
2.  **Start the Node:**
    ```bash
    docker compose -f docker-compose.local.yaml -f docker-compose.tailscale.yaml up -d
    ```

## 4. Common Operations

*   **Monitor Logs:**
    ```bash
    docker compose -f docker-compose.local.yaml [-f docker-compose.tailscale.yaml] logs -f govchaind
    ```
*   **Stop Services:**
    ```bash
    docker compose -f docker-compose.local.yaml [-f docker-compose.tailscale.yaml] down
    ```
*   **Clean Up (Stop and remove containers, networks, and volumes):**
    ```bash
    docker compose -f docker-compose.local.yaml [-f docker-compose.tailscale.yaml] down -v
    ```
    **Be cautious with `down -v` as it removes all blockchain data stored in volumes.**

## 5. Initial Tailscale Authentication (if required)

If you are using Tailscale and the `TS_AUTHKEY` you provided in your `.env` file is invalid, expired, or if you didn't provide one, the `tailscale` container will require manual authentication. Refer to the [Joining the GovChain Node with Docker Compose and Tailscale](DOCKER_COMPOSE_TAILSCALE.md) documentation for detailed authentication steps.