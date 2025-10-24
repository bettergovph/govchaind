# Security Policy

## Overview

Thank you for helping secure Open GovChain, OpenGovChain is a decentralized blockchain network designed to store and manage government datasets with complete transparency. Built as a foundational platform, agencies can extend it with custom modules for tokenomics, governance, financial transactions, and other blockchain utilities based on their specific needs.

Security is a shared responsibility - every validator, developer, and contributor plays a vital role in maintaining network integrity.

## Network Security Guidelines

### Firewall and Network Access

Validators and full nodes must implement strict access controls for network ports:

| Port | Purpose | Access | Notes |
|------|---------|--------|-------|
| 26656 | P2P | Public | Required for validator/peer connections (TCP) |
| 26657 | RPC | Restricted | Admin access only, never public, TLS recommended |
| 1317 | REST API | Local | Localhost or authenticated reverse proxy |
| 9090 | gRPC | Local | Internal services only, mutual TLS recommended |
| 26660 | Metrics | Local | Prometheus metrics (`127.0.0.1` only) |

#### Recommended Firewall Configuration

For cloud providers (AWS, GCP, OVH etc.):
```bash
# Allow P2P port only
- Ingress: TCP 26656 (from 0.0.0.0/0)
- Egress: All

# Restrict admin ports to trusted IPs
- Ingress: TCP 26657 (from trusted_ip_1, trusted_ip_2)
- Ingress: TCP 1317 (from trusted_ip_1, trusted_ip_2)
```

For host firewall (UFW):
```bash
# P2P port
ufw allow 26656/tcp

# Admin ports (replace with your IPs)
ufw allow from trusted_ip_1 to any port 26657
ufw allow from trusted_ip_1 to any port 1317

ufw enable
```

## Validator Security 

### Key Management

1. **Critical Key Protection**
   - Never expose `priv_validator_key.json` to public networks
   - Store mnemonics offline in secure, encrypted storage (e.g., encrypted USB drive)
   - Consider using Tendermint KMS for remote signing
   - Regularly audit key access logs

2. **Key Recovery Protocol**
   - Test key recovery procedures before production use
   - Document emergency key rotation steps
   - Store encrypted backups with proper access controls
   - Maintain separate backup recovery keys

### Infrastructure Security

1. **System Hardening**
   - Use dedicated hardware/VPS (no shared hosting)
   - Regular security updates and patches
   - Implement secure SSH configuration:
     ```bash
     # /etc/ssh/sshd_config
     PermitRootLogin no
     PasswordAuthentication no
     PubkeyAuthentication yes
     
     # Restrict SSH access to trusted IPs
     AllowUsers user@trusted_ip_1 user@trusted_ip_2
     
     # Alternative: using TCP wrappers
     # In /etc/hosts.allow:
     sshd: trusted_ip_1, trusted_ip_2
     # In /etc/hosts.deny:
     sshd: ALL
     ```
   - SSH can be restricted (VPN or Trusted IP) 
   - Enable system auditing and logging
   - Monitor resource usage and set alerts

2. **Container Security**
   - Scan container images for vulnerabilities:
     ```bash
     trivy image govchaind:latest
     docker scout quickview govchaind:latest
     ```
   - Implement least privilege principles:
     - Run containers as non-root user (UID 1001)
     - Use read-only root filesystem
     - Drop all capabilities and add only required ones
     - Set SELinux/AppArmor profiles
     - Configure seccomp policies
   - Secure container networking:
     - Use Tailscale for encrypted overlay networking
     - Restrict exposed ports to necessary services only
     - Implement network policies for container-to-container communication

## Security Checklist

### Validator Setup
- [ ] Hardware security module configured
- [ ] Network ports properly restricted
- [ ] SSH hardening implemented
- [ ] Monitoring systems active
- [ ] Alert systems configured

### Node Operation
- [ ] Regular software updates
- [ ] Log monitoring enabled
- [ ] Resource usage tracked
- [ ] Peer connections verified
- [ ] Key backup secure
- [ ] Recovery procedures documented

## Additional Resources

- [Cosmos SDK Security Guidelines](https://docs.cosmos.network/main/architecture/adr-006-secret-store-replacement)
- [Validator Security Best Practices](https://hub.cosmos.network/main/validators/security)

---

*This security policy is regularly updated. Last update: 2025-10-12*
