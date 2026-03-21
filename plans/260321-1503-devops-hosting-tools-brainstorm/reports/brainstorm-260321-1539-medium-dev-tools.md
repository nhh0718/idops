# Brainstorm: 21 Medium-Sized Dev Tools (1-2 tuan build)

**Target:** Solo dev, Go + React/Next.js, Proxmox/VPS, VN market first -> global.
**Scope:** Lon hon small tools, co dashboard/platform feel, nhung van realistic cho 1 nguoi.

---

## 1. DevOps Platforms (6 ideas)

### 1.1 DeployDock
- **Mo ta:** Self-hosted deployment platform kieu Vercel/Coolify lite. Push code -> auto build -> deploy len VPS.
- **Pain point:** Coolify qua nang, Vercel ton tien, dev muon own infra nhung khong muon viet CI/CD scripts.
- **Thoi gian:** 2 tuan | **Stack:** Go backend + Docker API + Next.js dashboard + Git webhooks
- **MVP:** (1) GitHub webhook listener (2) Dockerfile auto-detect and build (3) Container orchestration (4) Domain routing qua Traefik/Caddy (5) Deploy logs realtime (6) Rollback 1-click (7) Env vars management
- **Stars:** Very High | **Income:** High (hosted tier) | **Alt:** Coolify, CapRover, Dokku
- **USP:** Nhe hon Coolify 10x, Go single binary, focus UX

### 1.2 FleetPulse
- **Mo ta:** Multi-server monitoring dashboard. Agent nhe cai tren moi server, push metrics ve central. Alert Telegram/Discord.
- **Pain point:** Grafana+Prometheus setup phuc tap, Netdata chi single-node, thieu fleet overview.
- **Thoi gian:** 2 tuan | **Stack:** Go agent + Go API + WebSocket + Next.js + SQLite
- **MVP:** (1) Go agent <5MB collect CPU/RAM/disk/net (2) Central dashboard multi-server (3) Realtime charts (4) Alert rules threshold-based (5) Telegram/Discord notify (6) Server health score (7) SSH quick-connect
- **Stars:** Very High | **Income:** Medium | **Alt:** Netdata, Uptime Kuma, Grafana
- **USP:** Zero-config agent, 1 binary cai xong

### 1.3 CronHub
- **Mo ta:** Cron job management platform. Tao, schedule, monitor cron jobs qua Web UI. HTTP calls, scripts, webhooks.
- **Pain point:** crontab -e khong co logging/monitoring, job fail khong biet.
- **Thoi gian:** 1 tuan | **Stack:** Go scheduler + Next.js UI + SQLite
- **MVP:** (1) Web UI tao/edit cron jobs (2) HTTP trigger support (3) Execution history + logs (4) Failure alerts (5) Job dependency chains (6) Cron expression builder tieng Viet
- **Stars:** High | **Income:** Medium | **Alt:** Healthchecks.io, EasyCron
- **USP:** Self-hosted, UI dep, ho tro VN

### 1.4 ProxPanel
- **Mo ta:** Modern dashboard thay the Proxmox UI cho mobile/tablet. Quan ly VMs, containers, backup qua giao dien responsive.
- **Pain point:** Proxmox UI cu, khong responsive, cham tren mobile.
- **Thoi gian:** 2 tuan | **Stack:** Go API (proxy Proxmox REST API) + Next.js PWA
- **MVP:** (1) VM/CT list + status (2) Start/stop/restart/console (3) Resource charts (4) Backup management (5) Firewall rule editor (6) Mobile-first responsive (7) Multi-node support
- **Stars:** Very High | **Income:** Medium | **Alt:** Proxmox built-in
- **USP:** Mobile-first, modern UI, PWA offline

### 1.5 GitRelay
- **Mo ta:** Lightweight GitOps engine. Watch Git repo -> detect changes -> auto apply (Docker Compose, K3s manifests, scripts).
- **Pain point:** ArgoCD/Flux qua nang cho small teams. Muon GitOps nhung khong can full K8s.
- **Thoi gian:** 2 tuan | **Stack:** Go daemon + Git polling + Next.js dashboard
- **MVP:** (1) Watch multiple git repos (2) Auto-detect docker-compose/k3s changes (3) Apply voi rollback (4) Sync status dashboard (5) Manual sync (6) Webhook + polling (7) Diff preview
- **Stars:** High | **Income:** Low | **Alt:** ArgoCD, Flux, Portainer GitOps
- **USP:** No K8s required, Docker Compose native

### 1.6 TunnelGate
- **Mo ta:** Self-hosted ngrok alternative. Expose local dev server ra internet qua tunnel server tren VPS. Custom domains, auto TLS.
- **Pain point:** ngrok free gioi han, cloudflared phuc tap, muon own tunnel server.
- **Thoi gian:** 2 tuan | **Stack:** Go tunnel server + Go CLI client + LetsEncrypt
- **MVP:** (1) TCP/HTTP tunnel (2) Custom subdomains (3) Auto TLS (4) Web dashboard active tunnels (5) Auth tokens (6) Bandwidth monitoring (7) Request inspector
- **Stars:** Very High | **Income:** High (hosted service) | **Alt:** ngrok, bore, tunnelto
- **USP:** Self-hosted, unlimited tunnels, Go single binary

---

## 2. Developer Tools and Platforms (5 ideas)

### 2.1 CodePaste
- **Mo ta:** Self-hosted pastebin cho developers. Syntax highlighting, expiration, embed anywhere. API cho CI/CD output sharing.
- **Pain point:** GitHub Gist cham, Pastebin xau co ads, khong self-hosted.
- **Thoi gian:** 1 tuan | **Stack:** Go API + Next.js + Monaco Editor + SQLite
- **MVP:** (1) Syntax highlight 50+ languages (2) Expiration (1h/1d/1w/never) (3) Password protection (4) API endpoint for CI/CD (5) Embed snippet widget (6) Burn after read mode
- **Stars:** High | **Income:** Low | **Alt:** dpaste, PrivateBin
- **USP:** Dev-focused, API-first, embeddable

### 2.2 APIForge
- **Mo ta:** API development workspace. Design OpenAPI schema, auto mock server, test endpoints. Lightweight Postman alternative.
- **Pain point:** Postman nang 500MB+, Insomnia bi enshittified, can tool nhe self-hosted.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Next.js + Monaco + OpenAPI parser
- **MVP:** (1) Collection/request management (2) Env variables + secrets (3) OpenAPI import/export (4) Auto mock server tu schema (5) Code generation (Go, TS, Python) (6) Response diff testing (7) Team sharing via git
- **Stars:** Very High | **Income:** Medium | **Alt:** Postman, Insomnia, Hoppscotch
- **USP:** Self-hosted, fast, OpenAPI-native

### 2.3 DocuForge
- **Mo ta:** Documentation platform cho teams/open-source. Markdown-based, versioning, search. Self-hosted Gitbook alternative.
- **Pain point:** Gitbook paywall, Notion khong phai docs tool, Docusaurus can build.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Next.js + MDX + Bleve (full-text search)
- **MVP:** (1) Markdown/MDX editor WYSIWYG (2) Sidebar nav auto-generate (3) Full-text search (4) Version/branch support (5) Dark mode (6) Custom domain (7) API docs auto-gen tu OpenAPI
- **Stars:** Very High | **Income:** Medium (hosted) | **Alt:** Gitbook, Notion, Docusaurus
- **USP:** Self-hosted, zero build step, instant preview

### 2.4 WebhookHub
- **Mo ta:** Webhook debugging and routing platform. Nhan webhooks, inspect payload, forward/transform/replay.
- **Pain point:** Debug webhook tu Stripe/VNPay/GitHub cuc kho, khong thay payload, khong replay duoc.
- **Thoi gian:** 1 tuan | **Stack:** Go API + WebSocket + Next.js + SQLite
- **MVP:** (1) Unique endpoint URLs (2) Realtime payload inspector (3) Replay requests (4) Forward/transform rules (5) Filter and search history (6) Auto-respond custom status
- **Stars:** High | **Income:** Medium | **Alt:** webhook.site, RequestBin
- **USP:** Self-hosted, unlimited, transform rules

### 2.5 EnvVault
- **Mo ta:** Quan ly env vars cho teams. Encrypted storage, role-based access, sync .env, inject vao Docker/CI. Self-hosted Doppler.
- **Pain point:** Chia se .env qua Slack/email nguy hiem. Moi dev env khac nhau.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Go CLI + Next.js + AES-256 + SQLite
- **MVP:** (1) Project/environment management (2) Encrypted storage (3) CLI sync (envvault pull) (4) Role-based access (5) Audit log (6) Docker integration (7) Git-ignored .env generation
- **Stars:** Very High | **Income:** High (team tier) | **Alt:** Doppler, Vault, Infisical
- **USP:** Don gian hon Vault 100x, UI dep, Go binary

---

## 3. SaaS Products (5 ideas)

### 3.1 StatusPage.vn
- **Mo ta:** Status page builder cho doanh nghiep VN. Trang trang thai dich vu tieng Viet, custom domain.
- **Pain point:** Statuspage.io dat ($29+/mo), khong tieng Viet, startup VN can giai phap re.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Next.js + PostgreSQL + Uptime checker
- **MVP:** (1) Public status page + custom domain (2) Uptime monitoring tich hop (3) Incident management (4) Subscriber notify (email/Telegram) (5) Tieng Viet native (6) Embed widget (7) API
- **Stars:** Medium | **Income:** High ($5-15/mo/page) | **Alt:** Statuspage.io, Cachet
- **USP:** Gia re, tieng Viet, target SMB VN

### 3.2 FormPilot
- **Mo ta:** Form builder + data collection. Drag-drop tao form, responses, webhook integration. Typeform alternative.
- **Pain point:** Typeform dat, Google Form xau, can form dep + webhook cho devs.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Next.js + DnD Kit + PostgreSQL
- **MVP:** (1) Drag-drop form builder (2) 10+ field types (3) Conditional logic (4) Webhook/response (5) Analytics/export (6) Custom thank-you page (7) Embed and share
- **Stars:** Medium | **Income:** High (freemium 100 responses free) | **Alt:** Typeform, Tally
- **USP:** Self-hostable, developer-friendly webhooks

### 3.3 MailForge
- **Mo ta:** Transactional email API cho devs VN. Gui email qua API (OTP, invoice, notification), templates, tracking.
- **Pain point:** SendGrid/Mailgun kho setup cho dev VN (payment, DNS), SES can AWS account.
- **Thoi gian:** 2 tuan | **Stack:** Go API + SMTP relay + Next.js dashboard + PostgreSQL
- **MVP:** (1) REST API gui email (2) Template engine (Handlebars) (3) Delivery tracking (open/click) (4) Domain verification (5) Dashboard analytics (6) Webhook events (7) Rate limiting
- **Stars:** Medium | **Income:** High (pay-per-email) | **Alt:** SendGrid, Resend, Mailgun
- **USP:** De dung cho dev VN, thanh toan VND

### 3.4 LinkBio.vn
- **Mo ta:** Linktree alternative cho VN. Landing page ca nhan, links, bio, portfolio. Ho tro VN payment (donate).
- **Pain point:** Linktree khong ho tro VND, KOLs/freelancers VN can link-in-bio tool local.
- **Thoi gian:** 1 tuan | **Stack:** Go API + Next.js + PostgreSQL
- **MVP:** (1) Custom bio page (2) Drag-drop link ordering (3) Click analytics (4) Custom themes (5) VN payment (MoMo/ZaloPay donate) (6) SEO optimization (7) Custom domain
- **Stars:** Low | **Income:** Medium ($2-5/mo premium) | **Alt:** Linktree, Bento
- **USP:** VN payment, tieng Viet, gia re

### 3.5 FeedbackFish.vn
- **Mo ta:** Widget thu feedback/bug report embed vao website. Screenshot tu dong, mood rating, dashboard quan ly.
- **Pain point:** Thu feedback kho, email/form dai, user luoi. Can widget nhe 1-click.
- **Thoi gian:** 1 tuan | **Stack:** Go API + Next.js widget (embeddable) + PostgreSQL
- **MVP:** (1) Embeddable widget (<5KB) (2) Screenshot capture tu dong (3) Mood rating (emoji) (4) Category tagging (5) Dashboard quan ly (6) Slack/Telegram notify (7) Public roadmap vote
- **Stars:** Medium | **Income:** Medium (freemium) | **Alt:** Canny, FeedbackFish
- **USP:** Nhe, de embed, tieng Viet

---

## 4. Open-Source Community Projects (5 ideas)

### 4.1 PasteBoard
- **Mo ta:** Universal clipboard sharing giua devices. Copy text/file may A -> paste may B. E2E encrypted, no account.
- **Pain point:** AirDrop chi Apple, KDE Connect chi Linux, khong co cross-platform simple tool.
- **Thoi gian:** 2 tuan | **Stack:** Go server + WebRTC/WebSocket + Next.js PWA + E2E encryption
- **MVP:** (1) Text/file sharing qua room code (2) E2E encryption (3) No account required (4) QR code connect (5) PWA mobile (6) Auto-expire rooms (7) CLI client
- **Stars:** Very High | **Income:** Low | **Alt:** Snapdrop, LocalSend
- **USP:** Cloud relay fallback, CLI + web + mobile

### 4.2 ScreenCast
- **Mo ta:** Self-hosted screen recording -> GIF/WebM tool. Record, convert, share link. Loom alternative cho devs.
- **Pain point:** Loom paywall, Giphy Capture macOS only, can record and share nhanh cho code review.
- **Thoi gian:** 2 tuan | **Stack:** Go API + FFmpeg + Next.js + MediaRecorder API
- **MVP:** (1) Browser-based recording (screen/tab/camera) (2) Auto convert GIF/WebM/MP4 (3) Shareable links (4) Trim/cut editor (5) Embed markdown (6) Password protection (7) Auto-expire
- **Stars:** Very High | **Income:** Medium (hosted) | **Alt:** Loom, ShareX, Kap
- **USP:** Self-hosted, no account, dev-focused (embed in PR/issues)

### 4.3 FileDrop
- **Mo ta:** Self-hosted file sharing. Upload -> get link, expiration, password, download limit. WeTransfer alternative.
- **Pain point:** WeTransfer 2GB limit free, Google Drive can account, can self-hosted file share.
- **Thoi gian:** 1 tuan | **Stack:** Go API + Next.js + S3-compatible storage (MinIO)
- **MVP:** (1) Drag-drop upload (2) Shareable links (3) Password protection (4) Expiration + download limit (5) Chunked upload progress (6) Storage quota management (7) Admin dashboard
- **Stars:** Very High | **Income:** Low | **Alt:** transfer.sh, Firefox Send (dead)
- **USP:** Go single binary, S3/local storage, modern UI

### 4.4 BookmarkOS
- **Mo ta:** Self-hosted bookmark manager voi full-text search. Save URL -> auto-fetch title/screenshot/content -> search everything.
- **Pain point:** Browser bookmarks khong searchable, Raindrop paywall for full-text search.
- **Thoi gian:** 2 tuan | **Stack:** Go API + Headless Chrome + Bleve (search) + Next.js + SQLite
- **MVP:** (1) Browser extension save bookmark (2) Auto-fetch metadata + screenshot (3) Full-text search (4) Tag and collection (5) Import tu browser/Raindrop (6) Public sharing (7) API
- **Stars:** Very High | **Income:** Medium | **Alt:** Raindrop.io, Linkding, Shiori
- **USP:** Full-text search free, screenshot preview, Go single binary

### 4.5 Waitlist
- **Mo ta:** Open-source waitlist/launch page builder. Coming soon page, collect emails, referral system, analytics.
- **Pain point:** Moi side project can waitlist page, build lai tu dau moi lan.
- **Thoi gian:** 1 tuan | **Stack:** Go API + Next.js + SQLite
- **MVP:** (1) Launch page templates (2) Email collection + verification (3) Referral tracking (move up queue) (4) Analytics (signups/day, sources) (5) Custom domain (6) Export CSV (7) Webhook on signup
- **Stars:** High | **Income:** Low | **Alt:** LaunchRock, Waitlist.me
- **USP:** Self-hosted, co referral system, free

---

## Tong hop va De xuat uu tien

| Rank | Project | Thoi gian | Stars | Income | Ly do uu tien |
|------|---------|-----------|-------|--------|----------------|
| 1 | **TunnelGate** | 2w | Very High | High | Viral potential cuc cao, ai cung can, monetizable |
| 2 | **DeployDock** | 2w | Very High | High | Self-hosted PaaS trend nong, portfolio showpiece |
| 3 | **FileDrop** | 1w | Very High | Low | Build nhanh, stars de, dung ngay personal |
| 4 | **EnvVault** | 2w | Very High | High | Pain point that, team tool, B2B potential |
| 5 | **StatusPage.vn** | 2w | Medium | High | VN niche chua ai lam, recurring revenue |

---

## Unresolved Questions

1. Co muon focus category nao truoc? (DevOps vs SaaS vs Open-source)
2. VPS specs hien tai? (RAM/CPU anh huong toi Docker build, FFmpeg, headless Chrome)
3. Co san domain .vn cho VN-market products khong?
4. Muc do san sang handle payments? (Stripe, MoMo, VNPay experience?)
5. Co muon build 1 umbrella brand cho tat ca tools hay moi tool standalone?
