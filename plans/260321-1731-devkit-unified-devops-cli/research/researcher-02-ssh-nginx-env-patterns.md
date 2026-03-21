# Nghiên Cứu: SSH Config, Nginx Config Generation, và .env File Handling trong Go

**Ngày:** 2026-03-21 | **Tác giả:** Researcher Agent | **Trạng thái:** Hoàn thành

---

## 1. SSH Config Parser/Manager

### Thư viện chính
- **kevinburke/ssh_config**: Parser chuyên biệt cho ~/.ssh/config, bảo toàn comment khi parse.
  - Hỗ trợ: Host, HostName, Port, User, IdentityFile, ProxyJump
  - Không hỗ trợ Match directive hiện tại
  - API: `Get(host, key)`, `GetAll(host, key)` cho directive lặp lại

### SSH Connection Testing
- Dùng `x/crypto/ssh.Dial("tcp", "addr:port", config)` để thiết lập kết nối
- **Cẩn trọng:** Timeout chỉ áp dụng ở TCP handshake, không áp dụng SSH auth phase
- **Giải pháp:** Dùng `net.DialTimeout` + `ssh.NewClientConn` để kiểm soát chi tiết timeouts
- Cần `ClientConfig` với user, auth method (PublicKey/Password), HostKeyCallback

### UX Pattern
- Load config từ ~/.ssh/config, parse host patterns
- List available hosts với filter fuzzy search
- Validate kết nối với progress indicator (timeout 5-10s)
- Show ngôn ngữ lỗi: "Host unreachable", "Auth failed", "Timeout"

---

## 2. Nginx Config Generation

### Cấu trúc cơ bản
- **Upstream blocks**: Định nghĩa backend servers cho load balancing
- **Server blocks**: Virtual hosts, directives chính
- **Location blocks**: Route-specific configuration (reverse proxy, static files)
- Phổ biến: reverse proxy, static site, PHP-FPM, WebSocket proxy, load balancer

### Approach dùng Go templates
- Sử dụng Go `text/template` để generate nginx configs
- NGINX Instance Manager dùng JSON schema + Go templates cho augment templates
- Loadcat project (Go-based) demo config generation + reload mechanism
- Template variables: upstream addresses, ports, SSL certs, domain names

### Config Validation
- Chạy `nginx -t -c /path/to/nginx.conf` (dùng os/exec)
- Kiểm tra syntax + mở toàn bộ files referenced via include
- Exit code 0 = OK, non-zero = error
- Parse stderr để detect lỗi cụ thể

### Let's Encrypt Integration
- Cần hook cho certbot renewal (post-renewal, pre-renewal)
- Update nginx config với new cert paths, reload nginx
- Consider automation: cron job hoặc systemd timer

---

## 3. .env File Handling

### Parsing Libraries
- **joho/godotenv**: Port Ruby dotenv, load .env vào env vars, đơn giản & phổ biến
- **hashicorp/go-envparse**: Minimal allocations, hỗ trợ JSON strings, better tested
- Cả hai support multiple .env files, override behavior

### Multi-environment Support
- Convention: `FOO_ENV` env var xác định environment (development/staging/production)
- Load order: `.env.{environment}` sau `.env` để override
- Fallback to "development" nếu không set FOO_ENV

### Validation & Comparison
- **Fail fast:** Validate critical vars at startup (DB URL, secret keys)
- **Diff algorithm:** Compare .env.example vs .env để detect missing keys
  - Dùng map[string]string để track required vs present
  - Suggest missing keys with default values
- **Validation rules:** Kiểm tra empty values, duplicate keys, special chars, spaces in keys

### Interactive Prompts
- **charmbracelet/huh**: Build terminal forms với validation
  - Form.Run() orchestrates multi-group experience
  - ValidateMinLength, ValidateMaxLength, custom validators
  - Accessibility mode cho screen readers
  - Dynamic forms dựa trên input trước đó

---

## 4. Go CLI Architecture (Multi-tool)

### Directory Layout
```
cmd/
├── devkit/            # Main CLI entry point
├── devkit-ssh/        # SSH subcommand
└── devkit-nginx/      # Nginx subcommand

internal/
├── app/               # Core application logic (shared)
├── config/            # Config parsing & validation
├── ssh/               # SSH manager, connection testing
├── nginx/             # Nginx config generation, validation
├── env/               # .env file handling
└── ui/                # TUI components (shared)

pkg/
├── template/          # Go template helpers
└── exec/              # Wrapper cho os/exec (nginx -t, restart)
```

### Patterns từ lazydocker
- Separate `pkg/commands` cho business logic (Docker operations)
- `pkg/config` cho user configuration
- Shared config loader, reusable components
- Minimal main.go, delegate to internal/app

### Build & Release
- **GoReleaser:** Cross-compile (Linux, Darwin, Windows)
- **Install methods:** Homebrew tap, AUR, curl installer, GitHub releases
- Version flags: `-ldflags "-X main.Version=..."`

---

## 5. Tóm tắt khuyến nghị

| Thành phần | Thư viện/Approach | Ghi chú |
|-----------|------------------|--------|
| SSH Parser | kevinburke/ssh_config | Production-ready, bảo toàn comments |
| SSH Testing | x/crypto/ssh + net.DialTimeout | Handle timeout carefully, context cancellation |
| Nginx Config | Go text/template | JSON schema cho flexible templates |
| Nginx Validation | os/exec + `nginx -t` | Parse stderr cho error messages |
| .env Parsing | joho/godotenv | Simple, widely used |
| .env Validation | Custom map + hashicorp/go-envparse | Better error messages |
| Interactive Forms | charmbracelet/huh | V2 API stable, good accessibility |
| CLI Structure | cmd/internal/pkg | Scalable, follow Go conventions |

---

## Unresolved Questions

- Cần authentication strategy nào cho SSH key management? (encrypted storage?)
- Nginx config hot-reload vs full restart strategy?
- Multi-environment .env merging algorithm chi tiết?
- TUI component sharing mechanism giữa SSH/Nginx subtools?

---

**Sources:**
- [kevinburke/ssh_config](https://github.com/kevinburke/ssh_config)
- [golang.org/x/crypto/ssh](https://pkg.go.dev/golang.org/x/crypto/ssh)
- [joho/godotenv](https://github.com/joho/godotenv)
- [hashicorp/go-envparse](https://github.com/hashicorp/go-envparse)
- [charmbracelet/huh](https://github.com/charmbracelet/huh)
- [NGINX Load Balancing](https://nginx.org/en/docs/http/load_balancing.html)
- [Nginx Config Testing](https://nginx.org/en/docs/switches.html)
- [Go CLI Structure](https://www.bytesizego.com/blog/structure-go-cli-app)
