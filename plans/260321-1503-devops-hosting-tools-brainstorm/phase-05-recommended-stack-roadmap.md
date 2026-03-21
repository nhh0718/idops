---
phase: 05
title: "Recommended Stack & Roadmap"
priority: Strategy
status: pending
---

# Phase 05: Stack Gợi Ý & Lộ Trình Triển Khai

## Recommended Tech Stack

### Core Infrastructure

| Layer | Công Cụ | Lý Do |
|-------|---------|------|
| Hypervisor | **Proxmox VE** | All-in-one (KVM+LXC), REST API, HA, free |
| Container Orchestration | **K3s** | Lightweight K8s (~40MB), edge-ready |
| Reverse Proxy | **Traefik** hoặc **Caddy** | Auto-SSL, Docker-native |
| CI/CD | **Woodpecker** | Self-hosted, Docker-native, low RAM |
| IaC | **Terraform + Ansible** | Declarative + config management |

### Application Stack

| Layer | Công Cụ | Lý Do |
|-------|---------|------|
| Backend API | **Go** | Performance, single binary, concurrency |
| Frontend | **Next.js** hoặc **Astro** | SSR/SSG cho SEO tools, React ecosystem |
| Database | **PostgreSQL** | Reliable, JSON support, extensions |
| Cache | **Redis** | Rate limiting, session, cache |
| Time-series | **TimescaleDB** | Metrics storage (extension PostgreSQL) |
| Search | **Meilisearch** | Full-text search cho domain marketplace |

### Monitoring & Observability

| Layer | Công Cụ | Lý Do |
|-------|---------|------|
| Metrics | **Prometheus** | Standard, 75% adoption |
| Logs | **Loki** | Lightweight, Prometheus-like |
| Visualization | **Grafana** | Best dashboards |
| Uptime | **Uptime Kuma** | Simple, 95+ notifications |
| Alerting | **Grafana Alerting** | Unified với dashboards |

### Billing & Payment

| Layer | Công Cụ | Lý Do |
|-------|---------|------|
| Billing Platform | **FOSSBilling** (start) → **Custom** (scale) | Free, modern, migration path |
| International | **Paddle** | MoR, tax handled |
| Vietnam | **SePay + VietQR** | Local payment, bank transfer |
| Backup payment | **Stripe** | Fallback, card payments |

### Backup & Security

| Layer | Công Cụ | Lý Do |
|-------|---------|------|
| Backup | **Restic** | Dedup, S3, encrypted |
| Container scan | **Trivy** | Lightweight, CI-friendly |
| IDS/IPS | **CrowdSec** | Community-driven |
| Firewall | **nftables** | Modern iptables replacement |

---

## Lộ Trình Triển Khai

### Quarter 1 (Tháng 1-3): Foundation — Quick Wins

**Mục tiêu:** Launch 10+ free tools → bắt đầu thu traffic SEO

| Tuần | Task | Output |
|------|------|--------|
| 1-2 | Setup infra: Proxmox + K3s + CI/CD | Production environment |
| 3-4 | Build 5 Easy tools (DNS, WHOIS, IP, SSL, Domain Check) | 5 tools live |
| 5-6 | Build 3 Medium tools (Speed Test, Port Scan, Email Test) | 8 tools live |
| 7-8 | Build API layer + rate limiting | Public API |
| 9-10 | SEO optimization + content marketing | Organic traffic start |
| 11-12 | FOSSBilling setup + VN payment | Billing ready |

**KPIs:** 10+ tools live, 1K+ organic visits/month, billing system operational

### Quarter 2 (Tháng 4-6): Growth — Monetization

| Tuần | Task | Output |
|------|------|--------|
| 1-3 | Unified Dashboard (DNS+SSL+Uptime) | Premium feature |
| 4-6 | VPS Provisioning (Proxmox integration) | Sell VPS |
| 7-9 | Reseller portal MVP | B2B channel |
| 10-12 | CLI tools (Server Harden, Health Check, Deploy) | Developer adoption |

**KPIs:** 10K+ visits/month, first paying customers, VPS sales

### Quarter 3 (Tháng 7-9): Scale — Platform

| Tuần | Task | Output |
|------|------|--------|
| 1-4 | Custom billing system (thay FOSSBilling) | Own billing |
| 5-8 | Domain reselling integration | Full hosting stack |
| 9-12 | AI support bot + client dashboard | Reduce support cost |

**KPIs:** 50K+ visits/month, 100+ paying customers, positive unit economics

### Quarter 4 (Tháng 10-12): Optimize — Ecosystem

| Tuần | Task | Output |
|------|------|--------|
| 1-4 | Marketplace (themes, plugins, apps) | Ecosystem |
| 5-8 | White-label solution cho resellers | B2B revenue |
| 9-12 | Mobile app + advanced monitoring | Complete platform |

**KPIs:** 100K+ visits/month, 500+ customers, profitable

---

## Ước Tính Chi Phí Infra (Khởi Đầu)

| Item | Chi phí/tháng | Notes |
|------|--------------|-------|
| Proxmox server (dedicated) | $50-100 | OVH/Hetzner |
| Domain + SSL | $15/năm | Porkbun |
| S3 storage (backup) | $5-10 | Backblaze B2 |
| Email (transactional) | $0-10 | Resend free tier |
| DNS (Cloudflare) | $0 | Free plan |
| **Total** | **~$70-120/mo** | |

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Competitor (MXToolbox, etc.) | High | Focus VN market + unified UX |
| SEO takes too long | Medium | Paid ads supplement, community building |
| Proxmox scaling limits | Medium | Multi-node cluster, cloud hybrid |
| Payment gateway issues VN | Medium | Multiple providers, manual fallback |
| Legal (WHOIS data scraping) | Low-Medium | Use official APIs, comply with ICANN |

## Câu Hỏi Cần Trả Lời Trước Khi Bắt Đầu

1. Domain chính cho platform? (vd: devtools.vn, hostkit.vn, ...)
2. Team size? Solo dev hay có team?
3. Budget ban đầu cho infra?
4. Target market chính: VN hay international?
5. Ưu tiên free tools trước (traffic) hay billing trước (revenue)?
