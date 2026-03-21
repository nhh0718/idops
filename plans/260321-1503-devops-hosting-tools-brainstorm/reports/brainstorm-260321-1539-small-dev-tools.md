# Brainstorm: 35 Small Dev Tools (1-3 ngày build)

**Target:** Solo dev, Go + React/Next.js, có Proxmox/VPS, thị trường VN + global.

---

## 1. DevOps / Sysadmin Tools

| # | Tên | Mô tả | Pain point | Ngày | Stack | Stars | Income | Tương tự |
|---|------|--------|------------|------|-------|-------|--------|----------|
| 1 | **portdog** | CLI scan port đang listen, hiện process/PID/user dạng bảng đẹp | `netstat`/`ss` output xấu, khó đọc | 1 | Go CLI | High | None | `lsof` nhưng UX tốt hơn |
| 2 | **dockstat** | Dashboard realtime cho Docker containers (CPU/RAM/Net) chạy trên terminal | `docker stats` thiếu history, không sort được | 2 | Go CLI (TUI) | High | None | lazydocker, ctop |
| 3 | **sshm** | SSH config manager — add/edit/list/connect hosts bằng TUI interactive | Sửa `~/.ssh/config` bằng tay dễ sai | 1 | Go CLI | Medium | None | storm, ssh-config-editor |
| 4 | **envsync** | So sánh `.env.example` vs `.env` thực tế, báo thiếu/thừa biến | Deploy thiếu env var → crash runtime | 1 | Go CLI | High | None | dotenv-linter |
| 5 | **logpretty** | Pipe JSON logs (structured logging) → output có màu, filter theo level/field | Đọc JSON log một dòng dài ngoằng | 1 | Go CLI | Medium | None | jq, humanlog |
| 6 | **certwatch** | Monitor SSL cert expiry cho list domains, alert qua Telegram/Discord | Cert hết hạn → site chết, quên gia hạn | 2 | Go + Web UI | Medium | Low | certbot, ssl-checker |
| 7 | **proxdash** | Lightweight dashboard cho Proxmox VE — xem VM/CT status, resource usage | Proxmox UI nặng, chậm trên mobile | 3 | Go API + Next.js | High | Medium | Proxmox built-in UI |
| 8 | **nginx-gen** | Interactive CLI tạo Nginx config (reverse proxy, SSL, rate limit) | Copy-paste config Nginx dễ sai syntax | 1 | Go CLI | High | None | nginxconfig.io |
| 9 | **syssnap** | Chụp snapshot system info (CPU/RAM/disk/network) → JSON/Markdown report | Debug server cần collect info nhanh gửi team | 1 | Go CLI | Medium | None | neofetch, inxi |

## 2. Web Dev Tools

| # | Tên | Mô tả | Pain point | Ngày | Stack | Stars | Income | Tương tự |
|---|------|--------|------------|------|-------|-------|--------|----------|
| 10 | **apimock** | Tạo mock REST API từ file JSON/YAML, hỗ trợ delay/random error | FE chờ BE xong API mới dev được | 1 | Go server | High | None | json-server, mockoon |
| 11 | **ogimage** | API generate Open Graph images từ template (title, desc, logo) | Mỗi blog post cần OG image, design tay mệt | 2 | Go API | High | Medium | vercel/og, satori |
| 12 | **colorcraft** | Web tool generate color palette, contrast checker, CSS export | Chọn màu accessible mất thời gian | 2 | Next.js | Medium | Low | coolors.co |
| 13 | **apidiff** | So sánh 2 API response (JSON diff), highlight thay đổi | Test API sau khi refactor, sợ break response | 1 | Next.js | Medium | None | jsondiff |
| 14 | **regexlab** | Regex builder interactive với real-time match + cheat sheet | Viết regex luôn phải mở regex101 | 2 | Next.js | Medium | Low | regex101 (nhưng simpler) |
| 15 | **svgmin** | Web tool optimize SVG — remove metadata, minify, preview before/after | SVG từ Figma export nặng gấp 3-5x | 1 | Next.js | Medium | None | svgo, svgomg |
| 16 | **faviconkit** | Upload 1 ảnh → generate full favicon set (ico, png các size, manifest) | Mỗi project phải resize favicon thủ công | 1 | Next.js + Go API | High | Low | realfavicongenerator |
| 17 | **tsgen** | Paste JSON → generate TypeScript interfaces/types tự động | Viết type cho API response bằng tay chán | 1 | Next.js | High | None | quicktype, json2ts |
| 18 | **readme-craft** | Interactive README.md generator cho GitHub repos | README template lặp đi lặp lại | 1 | Next.js | Medium | None | readme.so |

## 3. Developer Productivity

| # | Tên | Mô tả | Pain point | Ngày | Stack | Stars | Income | Tương tự |
|---|------|--------|------------|------|-------|-------|--------|----------|
| 19 | **gitstats** | CLI phân tích git repo: commits/ngày, top contributors, file thay đổi nhiều nhất | Muốn biết health của repo nhanh | 1 | Go CLI | Medium | None | git-quick-stats |
| 20 | **projinit** | Scaffolder tạo project Go/Next.js với boilerplate chuẩn (lint, CI, Docker) | Mỗi project mới setup lại từ đầu | 2 | Go CLI | Medium | None | cookiecutter, create-t3-app |
| 21 | **dotman** | Quản lý dotfiles bằng symlink + git, backup/restore 1 lệnh | Máy mới phải setup lại config | 2 | Go CLI | Medium | None | stow, chezmoi |
| 22 | **killport** | `killport 3000` — tìm và kill process đang chiếm port | "Port 3000 already in use" mỗi ngày | 0.5 | Go CLI | High | None | kill-port (npm) |
| 23 | **todo-cli** | CLI todo list lưu local per-project (trong `.todo` file) | Quản lý task nhỏ không cần mở Jira | 1 | Go CLI | Medium | None | taskwarrior, t |
| 24 | **gitclean** | Interactive CLI dọn branches đã merge, stale remotes, dangling objects | Repo local đầy branch cũ | 1 | Go CLI | Medium | None | git-trim |
| 25 | **snipbox** | CLI snippet manager — save/search/copy code snippets theo tag | Google lại cùng 1 snippet mỗi tuần | 2 | Go CLI | Medium | None | nap, pet |
| 26 | **cheatsh** | Offline cheatsheet viewer cho CLI tools (curl, docker, git, etc.) | Không nhớ flags, phải google | 2 | Go CLI (TUI) | High | None | tldr, cheat.sh |
| 27 | **timebox** | CLI Pomodoro timer + log thời gian vào file markdown | Track thời gian code không cần app nặng | 1 | Go CLI | Low | None | pomo |

## 4. SaaS / Monetizable Micro-tools

| # | Tên | Mô tả | Pain point | Ngày | Stack | Stars | Income | Tương tự |
|---|------|--------|------------|------|-------|-------|--------|----------|
| 28 | **qrcraft** | Web generate QR code đẹp (custom logo, màu, style) + API | QR code mặc định xấu, khách hàng cần branded | 2 | Next.js + Go API | Medium | Medium | qr-code-generator |
| 29 | **shortli** | URL shortener self-hosted với analytics (click, geo, device) | Bitly giới hạn free, muốn own data | 2 | Go API + Next.js | High | Medium | shlink, yourls |
| 30 | **jsonformat.dev** | Web format/validate/minify JSON/YAML/TOML — SEO domain đẹp | Dev search "json formatter" hàng ngày | 1 | Next.js | Medium | Low (ads) | jsonformatter.org |
| 31 | **crontab.guru-vn** | Cron expression builder UI tiếng Việt + giải thích | Dev VN không quen cron syntax | 1 | Next.js | Low | Low (ads) | crontab.guru |
| 32 | **screenshotapi** | API chụp screenshot website → PNG/PDF, dùng cho preview/thumbnail | Build link preview cần screenshot | 3 | Go API + Chromedp | Medium | High | screenshotapi.net |
| 33 | **upmon** | Uptime monitor đơn giản — check HTTP/TCP, alert Telegram | UptimeRobot free chỉ 50 monitors | 3 | Go + Next.js | High | Medium | uptimerobot, upptime |
| 34 | **pdf-tools.vn** | Merge/split/compress PDF online — target SEO tiếng Việt | Mỗi lần cần merge PDF phải dùng tool ngoại | 2 | Next.js + Go API | Low | Medium (ads) | ilovepdf |
| 35 | **mailcheck** | API validate email (syntax, MX, disposable check) — freemium | Spam signup bằng email fake | 2 | Go API | Medium | High | hunter.io, zerobounce |

---

## Đề xuất ưu tiên (Top 5 nên build trước)

1. **killport** (0.5 ngày) — Viral potential cực cao, ai cũng cần, build nhanh lấy momentum
2. **envsync** (1 ngày) — Giải quyết pain thật, dễ integrate vào CI/CD, nhiều stars
3. **faviconkit** (1 ngày) — SEO traffic tốt, mọi dev cần, income từ ads
4. **shortli** (2 ngày) — Self-hosted trend đang hot, monetizable, portfolio piece tốt
5. **proxdash** (3 ngày) — Niche Proxmox community đang thiếu tool nhẹ, differentiated

## Unresolved Questions

- Có muốn focus vào CLI tools (GitHub stars nhanh) hay web tools (income potential) trước?
- Proxmox API access level nào? Root hay limited user?
- Có domain sẵn cho các web tools không, hay cần mua mới?
- Target audience chính: dev VN hay global? (ảnh hưởng i18n strategy)
- Có muốn publish lên Homebrew/AUR cho CLI tools không?
