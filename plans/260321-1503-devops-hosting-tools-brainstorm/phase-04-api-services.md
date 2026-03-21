---
phase: 04
title: "API Services"
priority: Integration
status: pending
---

# Phase 04: API Services

## A. Domain & DNS APIs

### 1. Domain Availability Check API
- **Endpoint:** `GET /api/domain/check?q=example.com`
- **Mô tả:** Check domain available + giá + multi-TLD (.com, .net, .vn, .io, ...)
- **Target:** Domain reseller, app dev
- **Difficulty:** Medium
- **Stack:** Go + WHOIS protocol + registrar APIs (Namecheap, Porkbun, VNNIC)
- **Competitor:** Namecheap API, Domainr
- **USP:** Hỗ trợ .vn domains (ít API nào support), batch check, pricing comparison
- **Revenue:** Freemium API ($0 for 100 calls/day, $19/mo for 10K)

### 2. DNS Management API
- **Endpoint:** `POST /api/dns/records` (CRUD operations)
- **Mô tả:** Unified API quản lý DNS records trên nhiều providers
- **Target:** SaaS dev, hosting provider
- **Difficulty:** Medium-Hard
- **Stack:** Go + Cloudflare/Route53/DigitalOcean SDKs
- **Competitor:** Cloudflare API (single provider)
- **USP:** Multi-provider abstraction layer, migration tool giữa providers

### 3. WHOIS API
- **Endpoint:** `GET /api/whois?domain=example.com`
- **Mô tả:** Structured WHOIS data (JSON format, parsed)
- **Target:** Security researcher, domain investor
- **Difficulty:** Medium
- **Stack:** Go + WHOIS protocol parser
- **Competitor:** WhoisXML API ($30/mo)
- **USP:** Free tier generous hơn, Vietnamese domain support

---

## B. SSL/TLS APIs

### 4. SSL Certificate API
- **Endpoint:** `GET /api/ssl/check?host=example.com`
- **Mô tả:** Check SSL cert info, chain validation, expiry, grade
- **Difficulty:** Easy-Medium
- **Stack:** Go + crypto/tls library
- **USP:** Bulk check, expiry webhook alerts

### 5. Let's Encrypt Automation API
- **Endpoint:** `POST /api/ssl/issue` (auto DNS-01 challenge)
- **Mô tả:** Issue/renew SSL certs automatically via ACME
- **Difficulty:** Medium
- **Stack:** Go + lego ACME library + DNS provider APIs
- **Competitor:** Certbot CLI
- **USP:** API-first (không cần SSH vào server), multi-domain, wildcard support

---

## C. Server & Infrastructure APIs

### 6. VPS Provisioning API
- **Endpoint:** `POST /api/vps/create`
- **Mô tả:** Tạo/quản lý VPS trên Proxmox hoặc cloud providers
- **Difficulty:** Hard
- **Stack:** Go + Proxmox API + libvirt
- **Competitor:** Cloud provider APIs (AWS, DO, Vultr)
- **USP:** Multi-hypervisor support, Proxmox-native, on-premise friendly

### 7. Server Monitoring API
- **Endpoint:** `GET /api/monitor/server/{id}/metrics`
- **Mô tả:** Real-time server metrics (CPU, RAM, disk, network, processes)
- **Difficulty:** Medium
- **Stack:** Go + Prometheus remote_write + node_exporter
- **Competitor:** Datadog API, New Relic
- **USP:** Self-hosted option, per-customer metering cho hosting providers

---

## D. Payment & Billing APIs

### 8. VN Payment Gateway Wrapper
- **Endpoint:** `POST /api/payment/create-order`
- **Mô tả:** Unified payment API cho VN (VietQR, MoMo, ZaloPay, VNPAY, bank transfer)
- **Difficulty:** Medium-Hard
- **Stack:** Go/Node.js + SePay + MoMo/ZaloPay SDKs
- **Competitor:** SePay (single), payment gateway riêng lẻ
- **USP:** Unified interface cho tất cả payment method VN, auto-reconciliation
- **Revenue:** Transaction fee 0.5-1% hoặc flat $29/mo

### 9. Subscription & Metering API
- **Endpoint:** `POST /api/billing/usage`
- **Mô tả:** Track resource usage (bandwidth, storage, compute) → auto invoice
- **Difficulty:** Hard
- **Stack:** Go + PostgreSQL + time-series DB (TimescaleDB)
- **Competitor:** Stripe Billing, Paddle (không hỗ trợ metering tốt cho hosting)
- **USP:** Hosting-specific metering (bandwidth, disk, CPU hours)

---

## E. Utility APIs

### 10. IP Geolocation API
- **Endpoint:** `GET /api/ip/lookup?ip=1.2.3.4`
- **Mô tả:** IP → country, city, ISP, ASN, timezone
- **Difficulty:** Easy
- **Stack:** Go + MaxMind GeoLite2 DB
- **Competitor:** ipinfo.io, ip-api.com
- **USP:** Free tier 10K/day, VN IP data chính xác hơn

### 11. Website Screenshot API
- **Endpoint:** `GET /api/screenshot?url=example.com`
- **Mô tả:** Capture website screenshot (full page, thumbnail)
- **Difficulty:** Medium
- **Stack:** Go + Playwright/Chromium headless
- **Competitor:** screenshotapi.net
- **USP:** Bundled với các tool khác, preview cho domain marketplace

### 12. Email Validation API
- **Endpoint:** `GET /api/email/validate?email=user@example.com`
- **Mô tả:** Validate email (syntax, MX record, disposable check, SMTP verify)
- **Difficulty:** Medium
- **Stack:** Go + DNS lookup + SMTP handshake
- **Competitor:** ZeroBounce, Hunter.io
- **USP:** Bulk validation, VN email provider support

## Kiến Trúc API Chung

```
API Gateway (Traefik/Kong)
├── Auth Service (JWT + API keys)
├── Rate Limiter (Redis)
├── Domain Services
├── SSL Services
├── Server Services
├── Payment Services
├── Utility Services
└── Usage Metering → Billing
```

**Tech Stack chung:** Go (performance) + PostgreSQL + Redis + Docker + Traefik

## Ưu Tiên Triển Khai

1. IP Geolocation + WHOIS + Domain Check (Easy, high traffic)
2. SSL Certificate Check + Email Validation (Easy-Medium)
3. VN Payment Gateway Wrapper (Medium, core revenue)
4. DNS Management API (Medium, DevOps value)
5. VPS Provisioning API (Hard, core business)
6. Subscription Metering (Hard, long-term)
