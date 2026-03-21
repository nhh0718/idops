---
phase: 03
title: "CLI/DevOps Tools"
priority: Infrastructure
status: pending
---

# Phase 03: CLI/DevOps Tools

## A. Server Provisioning & IaC

### Công Cụ Hiện Có

| Công Cụ | Mục Đích | Difficulty | Notes |
|---------|---------|------------|-------|
| **Terraform/OpenTofu** | IaC declarative, state-based | Medium | Multi-cloud, HCL |
| **Ansible** | Config management, agentless | Medium | SSH-based, YAML |
| **Pulumi** | IaC bằng programming language | Medium | TypeScript/Python/Go |

### Ý Tưởng CLI Tools

#### 1. VPS Provisioner CLI
- **Mô tả:** CLI tool tự động tạo VPS trên Proxmox/cloud providers, cài đặt stack
- **Command:** `vps-cli create --os ubuntu22 --ram 4g --cpu 2 --stack lamp`
- **Target:** Hosting provider, DevOps engineer
- **Difficulty:** Medium
- **Stack:** Go CLI + Proxmox API + Ansible playbooks
- **USP:** One-command VPS provisioning với preset stacks (LAMP, LEMP, Docker, K3s)

#### 2. Server Hardening CLI
- **Mô tả:** Auto-harden server mới (SSH config, firewall, fail2ban, auto-update)
- **Command:** `server-harden --profile web-server --ssh-port 2222`
- **Target:** Tất cả
- **Difficulty:** Easy-Medium
- **Stack:** Bash/Go + iptables/nftables + sshd config
- **USP:** Opinionated security profiles (web, database, mail server)

#### 3. DNS Zone Manager CLI
- **Mô tả:** CLI quản lý DNS records (add/remove/update) cho nhiều providers
- **Command:** `dns-cli set A example.com 1.2.3.4 --provider cloudflare`
- **Target:** DevOps, hosting admin
- **Difficulty:** Medium
- **Stack:** Go + multi-provider SDK (Cloudflare, Route53, DigitalOcean DNS)
- **USP:** Unified interface cho mọi DNS provider

---

## B. CI/CD & Deployment

### Công Cụ Hiện Có

| Công Cụ | Self-hosted | Ưu Điểm |
|---------|------------|---------|
| **Drone** | Yes | Container-native, RAM thấp |
| **Woodpecker** | Yes | Fork Drone, community |
| **Argo CD** | Yes | GitOps cho K8s |

### Ý Tưởng

#### 4. Deploy CLI (Zero-Config Deployment)
- **Mô tả:** Deploy app lên VPS bằng 1 command (detect framework, build, deploy)
- **Command:** `deploy push --target vps1.example.com`
- **Target:** Indie dev, startup
- **Difficulty:** Medium-Hard
- **Stack:** Go CLI + Docker + SSH + Caddy/Traefik reverse proxy
- **USP:** Heroku-like DX nhưng trên own VPS, auto-SSL, zero-downtime deploy
- **Competitor:** Coolify, Kamal, CapRover (nhưng CLI-first)

#### 5. Multi-Server Sync Tool
- **Mô tả:** Sync config/files giữa nhiều servers (rsync on steroids)
- **Command:** `sync push ./config --servers production`
- **Target:** Hosting admin
- **Difficulty:** Medium
- **Stack:** Go + rsync + SSH tunneling
- **USP:** Server groups, dry-run preview, rollback support

---

## C. Backup & Recovery

### Công Cụ Hiện Có

| Công Cụ | Đặc Điểm | Phù Hợp |
|---------|---------|---------|
| **Restic** | Dedup, S3-compatible, encrypted | VPS fleet |
| **BorgBackup** | Performance cao, local | High-throughput |
| **Velero** | Kubernetes backup | K8s clusters |

### Ý Tưởng

#### 6. Backup Orchestrator
- **Mô tả:** Quản lý backup cho fleet VPS (schedule, verify, restore, report)
- **Command:** `backup-ctl schedule --server all --daily --retain 30d --target s3://bucket`
- **Target:** Hosting provider
- **Difficulty:** Medium
- **Stack:** Go + Restic + S3 API + cron
- **USP:** Multi-server fleet management, backup verification, email reports

---

## D. Monitoring & Security

### Ý Tưởng

#### 7. Server Health CLI
- **Mô tả:** Quick health check cho server (CPU, RAM, disk, top processes, open ports)
- **Command:** `health-check --server vps1 --alert telegram`
- **Target:** Tất cả
- **Difficulty:** Easy
- **Stack:** Go/Bash + SSH + Prometheus node_exporter
- **USP:** Instant diagnostics, no agent install needed

#### 8. Container Security Scanner
- **Mô tả:** Scan Docker images trước khi deploy (CVE, secrets, best practices)
- **Command:** `scan-image myapp:latest --fail-on high`
- **Target:** DevOps, CI/CD pipeline
- **Difficulty:** Medium
- **Stack:** Go + Trivy integration + custom rules
- **USP:** CI/CD friendly, custom policy support, Vietnamese report

#### 9. Firewall Manager
- **Mô tả:** Unified firewall management (UFW/iptables/nftables) cho fleet
- **Command:** `fw allow 80,443 --server-group web && fw deny 3306 --except 10.0.0.0/8`
- **Target:** Hosting admin
- **Difficulty:** Medium
- **Stack:** Go + SSH + nftables
- **USP:** Fleet-wide firewall rules, audit log, rollback

---

## E. Xu Hướng 2025-2026 (Tools Cần Chú Ý)

1. **Platform Engineering CLI** — Internal Developer Platform tools
2. **GitOps Dashboard** — Argo CD web UI improvements
3. **AI-powered Incident Response** — Auto-detect & suggest fixes
4. **Edge Deployment** — Deploy to edge locations (Fly.io model)

## Ưu Tiên Triển Khai

1. Server Hardening CLI (Easy, instant value)
2. Server Health CLI (Easy, daily use)
3. VPS Provisioner CLI (Medium, core business)
4. Deploy CLI (Medium-Hard, high developer adoption)
5. Backup Orchestrator (Medium, essential for hosting)
6. DNS Zone Manager CLI (Medium, DevOps utility)
