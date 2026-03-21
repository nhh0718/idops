---
phase: 01
title: "Free Tools (SEO/Marketing)"
priority: Quick Win
status: pending
---

# Phase 01: Free Tools - Thu Hút Traffic qua SEO

> **Chiến lược:** Xây dựng bộ công cụ miễn phí → thu hút organic traffic từ dev/sysadmin → convert sang khách hàng hosting/VPS/domain.

## Danh Sách Ý Tưởng Công Cụ

### Tier 1 — Easy (1-3 ngày/tool) | Quick Win

| # | Tên Công Cụ | Mô Tả | SEO Value | Competitor |
|---|-------------|-------|-----------|------------|
| 1 | **DNS Lookup** | Tra cứu bản ghi DNS (A, AAAA, MX, TXT, CNAME, NS) | Rất cao | DNSChecker, MXToolbox |
| 2 | **WHOIS Lookup** | Tra cứu thông tin domain (registrar, expiry, nameserver) | Cao | whois.domaintools.com |
| 3 | **IP Lookup / Geolocation** | Tra IP → quốc gia, ISP, ASN, timezone | Rất cao | WhatIsMyIP, ipinfo.io |
| 4 | **SSL Certificate Checker** | Kiểm tra SSL validity, chain, expiry date | Cao | SSLShopper, SSLLabs |
| 5 | **HTTP Header Checker** | Kiểm tra response headers, security headers | Trung bình | SecurityHeaders.com |
| 6 | **Domain Availability Checker** | Check domain khả dụng (multi-TLD) | Rất cao | Namecheap, GoDaddy |
| 7 | **Subnet/CIDR Calculator** | Tính subnet mask, IP range, wildcard | Trung bình | subnet-calculator.com |
| 8 | **Base64/URL/HTML Encoder** | Encode/decode utility cho dev | Trung bình | base64encode.org |
| 9 | **Password Generator** | Tạo mật khẩu mạnh với tùy chọn | Cao | 1Password generator |
| 10 | **UUID/ULID Generator** | Tạo unique ID cho dev | Trung bình | uuidgenerator.net |

**Tech Stack:** Next.js/Astro + serverless API routes. Không cần DB cho hầu hết tools.

### Tier 2 — Medium (3-7 ngày/tool) | Giá Trị Cao

| # | Tên Công Cụ | Mô Tả | SEO Value | Competitor |
|---|-------------|-------|-----------|------------|
| 11 | **Website Speed Test** | Lighthouse-based performance audit | Rất cao | GTmetrix, PageSpeed |
| 12 | **DNS Propagation Checker** | Check DNS từ nhiều location toàn cầu | Rất cao | whatsmydns.net |
| 13 | **Port Scanner** | Scan open ports trên server (top 100 ports) | Cao | HackerTarget, nmap.online |
| 14 | **Email Deliverability Test** | Kiểm tra SPF, DKIM, DMARC config | Cao | MXToolbox, mail-tester |
| 15 | **Uptime Monitor (Free Tier)** | Monitor 5 URLs miễn phí, alert qua email/Telegram | Rất cao | UptimeRobot, BetterStack |
| 16 | **SSL Expiry Dashboard** | Track SSL certificates, alert trước khi hết hạn | Cao | TrackSSL (paid) |
| 17 | **Website Technology Detector** | Phát hiện tech stack của website (CMS, framework) | Cao | BuiltWith, Wappalyzer |
| 18 | **Redirect Checker** | Trace redirect chain (301/302/meta refresh) | Trung bình | redirect-checker.org |
| 19 | **Broken Link Checker** | Crawl site tìm link chết | Cao | brokenlinkcheck.com |
| 20 | **Cron Expression Generator** | UI tạo cron expression với preview | Trung bình | crontab.guru |

**Tech Stack:** Next.js + worker queues (cho speed test, port scan). Redis cache cho kết quả.

### Tier 3 — Hard (1-2 tuần) | Giá Trị Rất Cao

| # | Tên Công Cụ | Mô Tả | SEO Value | Competitor |
|---|-------------|-------|-----------|------------|
| 21 | **Unified Dashboard** | DNS+SSL+Uptime+Speed trong 1 dashboard | Rất cao | Không có (gap!) |
| 22 | **Domain Expiry Monitor** | Track domain expiry + WHOIS change detection | Cao | DomainIQ (limited) |
| 23 | **Email Blacklist Checker** | Check IP/domain trên 100+ blacklists | Cao | MXToolbox (limited free) |
| 24 | **Website Security Scanner** | Scan XSS, SQLi, outdated headers | Cao | Qualys, Mozilla Observatory |
| 25 | **API Performance Tester** | Load test endpoint + response time graph | Trung bình | Postman (khác segment) |

**Tech Stack:** Go/Rust backend (performance), PostgreSQL, React dashboard, background workers.

## USP / Điểm Khác Biệt

1. **All-in-one:** Tất cả tools trên 1 domain (không cần nhớ 10 website khác nhau)
2. **Tiếng Việt:** Hỗ trợ tiếng Việt → capture thị trường VN (ít competitor)
3. **No signup cho basic:** Dùng ngay không cần đăng ký
4. **API miễn phí:** Cung cấp API (rate-limited) cho mỗi tool → developer adoption
5. **Modern UX:** Clean, fast, dark mode, mobile-friendly (nhiều tool cũ UX kém)

## Monetization

- **Free:** Unlimited web UI, API 100 calls/day
- **Pro ($9/mo):** API 10K calls/day, webhook alerts, SSL dashboard
- **Business ($29/mo):** White-label, API 100K calls/day, team features

## Ưu Tiên Triển Khai (Gợi ý)

1. DNS Lookup + WHOIS + IP Lookup + Domain Checker (4 tools cùng lúc, share infra)
2. SSL Checker + HTTP Header Checker
3. Speed Test + Uptime Monitor (cần worker infrastructure)
4. Unified Dashboard (kết hợp tất cả tools đã xây)
