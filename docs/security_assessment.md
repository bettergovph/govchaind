# Security Assessment Summary

This document captures the primary security weaknesses identified during the review and the remediation steps implemented in this change-set. A hands-on proof-of-concept that recreates the original vulnerability and demonstrates the hardened mitigation is available in [`docs/poc/genesis_poisoning.md`](poc/genesis_poisoning.md).

## 1. Insecure Genesis File Bootstrap

### Issue
- Both the Docker entrypoint (`docker-entrypoint.sh`) and the volunteer helper script (`scripts/join-as-volunteer.sh`) downloaded the network genesis file directly from the `main` branch without integrity checks.
- The downloads used silent `curl` invocations that ignored HTTP errors and offered no tamper detection. Any man-in-the-middle attack, compromised CDN, or unexpected upstream change could silently deliver a malicious genesis file. A poisoned genesis file allows an attacker to modify the validator set or bootstrap nodes onto a hostile fork.

### Remediation
- Added strict download options (`curl -fSL --retry 3 --retry-delay 2`) so HTTP failures surface immediately.
- Introduced optional `GENESIS_URL` and `GENESIS_SHA256` overrides and checksum verification in both scripts.
- Added dependency checks for `curl`, `jq`, and `sha256sum` to fail fast when prerequisites are missing.
- Hardened scripts with `set -euo pipefail`, temporary-file handling, and cleanup traps to avoid partial writes.

### Operator Action Items
- Publish an official `GENESIS_SHA256` alongside every release and configure the environment variables (or script arguments) accordingly.
- When hosting alternative genesis mirrors, ensure the checksum matches the canonical release artifact.
- Consider pinning `GENESIS_URL` to an immutable tag or release asset rather than a moving branch reference.

## 2. General Operational Guidance

- Monitor CI/CD pipelines to ensure scripts remain executable and dependencies such as `jq` and `sha256sum` stay available in container images.
- Document the expected checksum management process in validator onboarding materials so volunteers consistently verify genesis integrity.

The combination of these mitigations closes the immediate loophole around silent genesis tampering and establishes a clearer operational posture for secure node provisioning.
