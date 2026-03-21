---
phase: 02
title: "SaaS Platform Tools"
priority: Core Business
status: pending
---

# Phase 02: SaaS Platform Tools - Kinh Doanh Hosting/VPS/Domain

## A. Billing & Client Management

### Lựa Chọn Nền Tảng

| Platform | Giá | Ưu | Nhược | Target |
|----------|-----|-----|-------|--------|
| **WHMCS** | $29.95/mo (≤250 clients) | Ecosystem lớn, plugin nhiều | Đắt, legacy code | Large reseller |
| **FOSSBilling** | FREE (open-source) | Modern, active dev, Apache 2.0 | Beta, ít plugin | Indie/SME |
| **Blesta** | Flat fee | Security-first, no client cap | Ít theme | Security-conscious |
| **HostBill** | Custom | 500+ integrations | Proprietary, đắt | Enterprise |

**Khuyến nghị:** FOSSBilling cho khởi đầu (free, modern), migrate lên WHMCS khi scale.

### Ý Tưởng: Custom Billing Platform

- **Mô tả:** Lightweight billing cho hosting VN (tích hợp VietQR, MoMo, ZaloPay)
- **Target:** Hosting provider nhỏ tại VN
- **Difficulty:** Hard
- **USP:** Payment gateway VN-native, tiếng Việt, giá rẻ hơn WHMCS
- **Stack:** Laravel/Go + PostgreSQL + VietQR API + SePay
- **Competitor gap:** WHMCS không hỗ trợ tốt payment VN

---

## B. Control Panel (Server Management)

### So Sánh Panel Hiện Có

| Panel | Giá | Đặc Điểm | Phù Hợp |
|-------|-----|---------|---------|
| **HestiaCP** | FREE | Lightweight, CLI-driven, Apache+Nginx | VPS nhỏ, simple hosting |
| **CyberPanel** | FREE | LiteSpeed, REST API, Docker support | DevOps, performance |
| **CloudPanel** | FREE tier | ARM support (-30% cost), fast | Cost-optimized |
| **Coolify** | FREE | Docker Compose, modern UI, Heroku-like | Modern dev teams |
| **CapRover** | FREE | Simple PaaS, Docker-native | SMB |
| **Dokku** | FREE | Heroku-on-VPS, minimal resources | Budget-conscious |

**Khuyến nghị:** CyberPanel (power users) > HestiaCP (simple) > Coolify (modern teams)

### Ý Tưởng: VN-Focused Hosting Panel

- **Mô tả:** Control panel tối ưu cho thị trường VN (tiếng Việt, VPS local providers)
- **Features:** Website management, SSL auto, DNS, Email, Database, File Manager
- **Difficulty:** Hard (6+ tháng)
- **USP:** Tích hợp VN cloud providers (Viettel IDC, VNPT, FPT), UI tiếng Việt
- **Stack:** Go backend + React + Docker

---

## C. Monitoring & Alerting

### Công Cụ Hiện Có

| Tool | Giá | Đặc Điểm |
|------|-----|---------|
| **Uptime Kuma** | FREE (self-hosted) | 95+ notification channels, Docker |
| **Grafana + Prometheus** | FREE (self-hosted) | Visualization + metrics standard |
| **Netdata** | FREE tier | Per-second metrics, AI anomaly |
| **Zabbix** | FREE (self-hosted) | Enterprise-grade, learning curve cao |

### Ý Tưởng: Hosting Client Dashboard

- **Mô tả:** Dashboard cho khách hàng hosting xem resource usage, uptime, invoices
- **Target:** Hosting provider cung cấp cho end-user
- **Difficulty:** Medium
- **USP:** White-label, embed vào website hosting provider
- **Stack:** React + Prometheus exporter + REST API
- **Revenue:** Bán cho hosting providers ($10-50/mo)

---

## D. Reseller Management

### Ý Tưởng: Reseller Portal

- **Mô tả:** Multi-tenant portal cho đại lý bán lại hosting/VPS/domain
- **Features:**
  - Quản lý sub-accounts
  - Custom pricing/branding
  - Auto-provisioning VPS
  - Invoice & payment tracking
  - Domain reselling integration
- **Target:** Hosting resellers, web agencies
- **Difficulty:** Hard
- **USP:** All-in-one (billing + provisioning + domain + support ticket)
- **Stack:** Laravel/Next.js + PostgreSQL + Proxmox API + domain registrar APIs

---

## E. Support Ticket System

### Lựa Chọn

| Tool | Giá | Tích Hợp |
|------|-----|---------|
| **osTicket** | FREE | PHP, MySQL, email piping |
| **FreeScout** | FREE | Laravel, HelpScout-like |
| **Chatwoot** | FREE | Multi-channel, modern |
| **WHMCS built-in** | Included | Direct billing integration |

### Ý Tưởng: AI-Powered Support Bot

- **Mô tả:** Chatbot tự động trả lời câu hỏi hosting phổ biến (DNS, SSL, email setup)
- **Difficulty:** Medium
- **USP:** Giảm 50-70% ticket volume, tiếng Việt
- **Stack:** LLM API + RAG (knowledge base) + Chatwoot integration

## Ưu Tiên Triển Khai

1. FOSSBilling setup + VN payment integration (Quick Win)
2. Client resource dashboard (white-label)
3. Reseller portal MVP
4. AI support bot
5. Custom control panel (long-term)
