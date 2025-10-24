# Joining the GovChain Node with Docker Compose and Tailscale

This document provides a comprehensive guide on how to set up and run a `govchaind` node using Docker Compose, with integrated Tailscale for secure and dynamic networking. This setup ensures your node can join the GovChain network even behind NAT, and maintains persistent Tailscale state across container restarts.

## 1. Prerequisites

Before you begin, ensure you have the following installed and configured:

*   **Docker**: Make sure Docker Desktop (for macOS/Windows) or Docker Engine (for Linux) is installed and running on your system.
*   **Docker Compose**: This is usually bundled with Docker Desktop or can be installed separately.
*   **Git**: To clone the `govchaind` project repository.
*   **A Tailscale Account**: You'll need an active Tailscale account to generate authentication keys.

## 2. Tailscale Setup

Tailscale provides a secure mesh VPN that simplifies network configuration.

### 2.1. Sign Up for Tailscale

If you don't already have one, sign up for a free Tailscale account at [https://tailscale.com/](https://tailscale.com/). You can use your existing Google, Microsoft, or GitHub account.

### 2.2. Generate a Reusable Authentication Key

For non-interactive Docker environments, it's best practice to use a reusable authentication key.

1.  Log in to your Tailscale admin console: [https://login.tailscale.com/admin/](https://login.tailscale.com/admin/).
2.  Navigate to the **Auth keys** section (usually found in the left-hand menu under "Settings" or directly accessible via [https://login.tailscale.com/admin/settings/authkeys](https://login.tailscale.com/admin/settings/authkeys)).
3.  Click on **"Generate auth key"** or a similar button.
4.  Ensure you select options for a **reusable** key and, if available, set an appropriate **expiration** (e.g., 90 days) and **tags** for better organization.
5.  Copy the generated key. It will start with `tskey-auth-`. **Treat this key like a password; do not share it publicly or commit it directly to your repository.**

## 3. Project Setup

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

3.  **Populate `.env` with your Tailscale Auth Key**:
    Open the newly created `.env` file and replace the placeholder with your actual Tailscale authentication key:
    ```dotenv
    # .env
    TS_AUTHKEY=tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY
    TS_STATE_DIR=/var/lib/tailscale
    ```
    **Make sure to replace `tskey-auth-YOUR_ACTUAL_TAILSCALE_AUTH_KEY` with the key you generated in step 2.2.**

    *   `TS_STATE_DIR`: This variable is used to specify the directory where Tailscale stores its state files. By mapping this to a Docker volume, Tailscale's configuration and identity will persist across container restarts, preventing the need to re-authenticate every time.

## 4. Running the GovChain Node with Docker Compose and Tailscale

To run your `govchaind` node with Tailscale, you will combine a base Docker Compose configuration (either for local development or production) with the `docker-compose.tailscale.yaml` file. This modular approach allows you to opt-in to Tailscale networking when needed.

### 4.1. Base Docker Compose Files

Choose one of the following base configurations:

*   **For Production/Deployment (pulling image from GHCR):** Use `docker-compose.prod.yaml`.
*   **For Local Development (building image locally):** Use `docker-compose.local.yaml`.

### 4.2. Tailscale Docker Compose File

The `docker-compose.tailscale.yaml` file defines the `tailscale` sidecar service and extends the `govchaind` service to integrate with Tailscale networking.

### 4.3. Starting the Services

To start your `govchaind` node with Tailscale, use the `docker compose` command with multiple `-f` flags:

**For Production/Deployment with Tailscale:**
```bash
docker compose -f docker-compose.prod.yaml -f docker-compose.tailscale.yaml up -d
```

**For Local Development with Tailscale:**
```bash
docker compose -f docker-compose.local.yaml -f docker-compose.tailscale.yaml up -d
```

This command orchestrates two services: `govchaind-node` (your blockchain node) and `tailscale` (a sidecar container for Tailscale connectivity).

### 4.4. Build the Docker Image (for local development with Tailscale)

If you are using `docker-compose.local.yaml` with Tailscale, you will need to build the `govchaind` Docker image first:
```bash
docker compose -f docker-compose.local.yaml -f docker-compose.tailscale.yaml build
```

### 4.5. Start the Docker Compose Services

Now, start the `govchaind` and `tailscale` services using Docker Compose:
```bash
# For Production/Deployment with Tailscale
docker compose -f docker-compose.prod.yaml -f docker-compose.tailscale.yaml up -d

# For Local Development with Tailscale
docker compose -f docker-compose.local.yaml -f docker-compose.tailscale.yaml up -d
```
This command executes `docker-compose up -d`, which starts the containers in detached mode (in the background).

### 4.3. Initial Tailscale Authentication (if required)

If the `TS_AUTHKEY` you provided in your `.env` file is invalid, expired, or if you didn't provide one, the `tailscale` container will require manual authentication.

1.  **Monitor Tailscale Logs for the Authentication URL**:
    Immediately after running `make docker-up`, check the logs of the `tailscale` container. You'll be looking for a URL that starts with `https://login.tailscale.com/a/...`:
    ```bash
    docker-compose logs tailscale
    ```
    Scroll through the output. You should see lines similar to this:
    ```
    tailscale-sidecar  |
    tailscale-sidecar  | To authenticate, visit:
    tailscale-sidecar  |
    tailscale-sidecar  |         https://login.tailscale.com/a/YOUR_AUTH_URL_HERE
    tailscale-sidecar  |
    ```

2.  **Authenticate in your Web Browser**:
    Copy the full URL (`https://login.tailscale.com/a/YOUR_AUTH_URL_HERE`) from the logs and paste it into your web browser. Follow the prompts to log in to your Tailscale account and authorize the new device (your Docker container).

3.  **Verify Tailscale Connection**:
    Once authorized, Tailscale will connect. You can continue monitoring the `tailscale` container logs (`docker-compose logs tailscale`) until you see messages indicating it's in a `Running` state and has received an IP address.

    You can also check the Tailscale status inside the container:
    ```bash
    docker exec -it tailscale tailscale status
    ```
    This should show your container's Tailscale IP and connection status.

### 4.4. Verify GovChain Node Configuration

The `docker-entrypoint.sh` script for the `govchaind-node` is designed to dynamically configure the node's network identity and chain ID:

### Moniker Configuration

The node's moniker (name) is set via the `MONIKER` environment variable. If not explicitly provided at runtime (e.g., via `-e MONIKER=...`), it defaults to "GovChain Validator" as defined in the Docker image.

### External IP Address Configuration

1.  **Manual Override**: If the `EXTERNAL_IP` environment variable is set (e.g., `-e EXTERNAL_IP="YOUR_PUBLIC_IP_OR_DOMAIN"`), the script will use this value directly to configure `external_address` in `config.toml`. This can be useful if you need to advertise a specific IP that differs from the Tailscale IP.
2.  **Tailscale Detection**: If `EXTERNAL_IP` is not set, the script checks for a Tailscale sidecar. If a Tailscale IP is found, it will be used as the `external_address`.
3.  **Public IP Discovery**: If neither `EXTERNAL_IP` nor a Tailscale IP is detected, the script automatically attempts to discover the host machine's public IP address using an external service. If found, it configures `external_address` in `config.toml` with that IP (e.g., `external_address = "YOUR_PUBLIC_IP:26656"`).
4.  **No Configuration**: If no IP can be determined through any of the above methods, the node will start without a pre-configured `external_address`.

### Chain ID Configuration

The `chain_id` for your node is now dynamically extracted from the `genesis.json` file that is downloaded during the node's initialization. This ensures that your node always uses the correct and authoritative `chain_id` for the network it is joining, removing any need for manual configuration of this value.
## 5. Persistent Tailscale State

The `docker-compose.yml` configuration includes a named volume `tailscale-state` that is mounted to `/var/lib/tailscale` inside the `tailscale` container. This ensures that Tailscale's machine key and other state information are preserved. If you stop and restart your `tailscale` container (or even remove and recreate it, as long as the `tailscale-state` volume is not removed), it will retain its identity on your Tailscale network and should not require re-authentication.

## 6. Important Notes and Troubleshooting

*   **Security of `TS_AUTHKEY`**: Always keep your `TS_AUTHKEY` secure. If you suspect it has been compromised, revoke it immediately from your Tailscale admin console. For production environments, consider using ephemeral keys or more advanced secrets management solutions.
*   **Checking Logs**: The `docker-compose logs <service_name>` command is your primary tool for debugging. Use `docker-compose logs -f <service_name>` to follow logs in real-time.
*   **Stopping and Cleaning Up**:
    To stop the services:
    ```bash
    make docker-down
    ```
    To stop and remove containers, networks, and volumes (including persistent data):
    ```bash
    make docker-clean
    ```
    **Be cautious with `docker-clean` as it removes all blockchain data stored in volumes.**

By following these steps, you will have a `govchaind` node running securely with Tailscale, ready to join the GovChain network.
