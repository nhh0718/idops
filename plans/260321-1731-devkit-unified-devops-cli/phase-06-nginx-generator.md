---
phase: 6
title: "Nginx Generator - idops nginx"
status: pending
effort: 2.5d
priority: P2
depends_on: [1]
---

# Phase 06: Nginx Generator - `idops nginx`

## Context Links
- [Plan Overview](plan.md) | [Phase 01](phase-01-project-setup.md)
- Go text/template: https://pkg.go.dev/text/template

## Overview
Interactive Nginx config generator. Templates cho common patterns (reverse proxy, static site, PHP-FPM, load balancer, WebSocket). Validate voi `nginx -t`, apply va reload.

## Key Insights
- Go `text/template` du manh cho Nginx config generation
- Templates embed vao binary voi `embed.FS` (Go 1.16+)
- `nginx -t` validate config khong can reload
- `nginx -s reload` reload gracefully (khong drop connections)
- SSL config can cert paths -> check file existence truoc khi generate

## Requirements

### Functional
- Template selection: interactive menu chon template type
- Interactive form: fill domain, upstream, port, SSL options
- Preview: hien generated config truoc khi save
- Save: ghi vao /etc/nginx/sites-available/ (hoac custom path)
- Validate: `nginx -t` check syntax
- Apply: symlink sites-enabled + `nginx -s reload`
- List: `idops nginx list` -> list generated configs
- Templates: reverse-proxy, static-site, php-fpm, load-balancer, websocket

### Non-functional
- Templates phai readable va maintainable
- Generated config co comments giai thich moi section
- Khong require root cho generate/preview, chi require root cho apply/reload

## Architecture

```
internal/cli/nginx.go             # Cobra commands
internal/nginx/generator.go       # Template rendering engine
internal/nginx/templates.go       # Template definitions (embed.FS)
internal/nginx/validator.go       # nginx -t wrapper
internal/nginx/types.go           # Config structs per template type
templates/nginx/                  # .tmpl files
```

## Related Code Files

### Files to Create
- `internal/cli/nginx.go`
- `internal/nginx/generator.go`
- `internal/nginx/templates.go`
- `internal/nginx/validator.go`
- `internal/nginx/types.go`
- `templates/nginx/reverse-proxy.tmpl`
- `templates/nginx/static-site.tmpl`
- `templates/nginx/php-fpm.tmpl`
- `templates/nginx/load-balancer.tmpl`
- `templates/nginx/websocket.tmpl`

## Implementation Steps

### 1. Define Config Types (`internal/nginx/types.go`)
```go
type BaseConfig struct {
    ServerName   string
    ListenPort   int
    SSLEnabled   bool
    SSLCertPath  string
    SSLKeyPath   string
    AccessLog    string
    ErrorLog     string
}

type ReverseProxyConfig struct {
    BaseConfig
    UpstreamHost string
    UpstreamPort int
    WebSocket    bool
}

type StaticSiteConfig struct {
    BaseConfig
    RootPath     string
    IndexFiles   []string
    EnableGzip   bool
    CacheMaxAge  int
}

type PHPFPMConfig struct {
    BaseConfig
    RootPath     string
    PHPSocket    string // /run/php/php8.2-fpm.sock
}

type LoadBalancerConfig struct {
    BaseConfig
    UpstreamName string
    Backends     []Backend
    Method       string // round-robin, least_conn, ip_hash
}

type Backend struct {
    Host   string
    Port   int
    Weight int
}
```

### 2. Embed Templates (`internal/nginx/templates.go`)
```go
//go:embed templates/nginx/*.tmpl
var templateFS embed.FS

func GetTemplate(name string) (*template.Template, error) {
    return template.ParseFS(templateFS, "templates/nginx/"+name+".tmpl")
}
```
**Note**: embed path relative to module root -> templates/ dir o root

### 3. Template Files

**`templates/nginx/reverse-proxy.tmpl`:**
```nginx
{{- if .SSLEnabled }}
server {
    listen 80;
    server_name {{ .ServerName }};
    return 301 https://$host$request_uri;
}
{{ end -}}
server {
    {{- if .SSLEnabled }}
    listen 443 ssl http2;
    ssl_certificate {{ .SSLCertPath }};
    ssl_certificate_key {{ .SSLKeyPath }};
    {{- else }}
    listen {{ .ListenPort }};
    {{- end }}
    server_name {{ .ServerName }};

    location / {
        proxy_pass http://{{ .UpstreamHost }}:{{ .UpstreamPort }};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        {{- if .WebSocket }}
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        {{- end }}
    }
}
```

### 4. Generator (`internal/nginx/generator.go`)
```go
func Generate(templateName string, config interface{}) (string, error) {
    tmpl, err := GetTemplate(templateName)
    if err != nil { return "", err }
    var buf bytes.Buffer
    err = tmpl.Execute(&buf, config)
    return buf.String(), err
}

func SaveConfig(content, outputPath string) error {
    return os.WriteFile(outputPath, []byte(content), 0644)
}
```

### 5. Interactive Form
```go
func InteractiveGenerate() (string, string, error) {
    // Step 1: Select template type
    var templateType string
    huh.NewSelect[string]().
        Title("Select Nginx config template").
        Options(
            huh.NewOption("Reverse Proxy", "reverse-proxy"),
            huh.NewOption("Static Site", "static-site"),
            huh.NewOption("PHP-FPM", "php-fpm"),
            huh.NewOption("Load Balancer", "load-balancer"),
            huh.NewOption("WebSocket Proxy", "websocket"),
        ).Value(&templateType).Run()

    // Step 2: Per-template form (domain, ports, SSL, etc.)
    // Step 3: Generate va return (content, filename)
}
```

### 6. Validator (`internal/nginx/validator.go`)
```go
func Validate() error {
    cmd := exec.Command("nginx", "-t")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("nginx config invalid:\n%s", string(output))
    }
    return nil
}

func Apply(configPath, sitesEnabled string) error {
    // 1. Symlink: sites-available -> sites-enabled
    linkPath := filepath.Join(sitesEnabled, filepath.Base(configPath))
    os.Symlink(configPath, linkPath)
    // 2. Validate
    if err := Validate(); err != nil {
        os.Remove(linkPath) // rollback
        return err
    }
    // 3. Reload
    return exec.Command("nginx", "-s", "reload").Run()
}
```

### 7. Cobra Commands (`internal/cli/nginx.go`)
```go
var nginxCmd = &cobra.Command{Use: "nginx", Short: "Nginx config generator", RunE: runNginx}
var nginxGenerateCmd = &cobra.Command{Use: "generate", Short: "Interactive config generation"}
var nginxValidateCmd = &cobra.Command{Use: "validate", Short: "Validate nginx config"}
var nginxApplyCmd = &cobra.Command{
    Use: "apply <config-file>",
    Short: "Enable config and reload nginx",
}
var nginxListCmd = &cobra.Command{Use: "list", Short: "List generated configs"}
```
- `idops nginx` hoac `idops nginx generate` -> interactive mode
- `idops nginx validate` -> run nginx -t
- `idops nginx apply <file>` -> symlink + validate + reload
- `idops nginx list` -> list configs trong sites-available
- `--output <path>` flag cho custom output path

## Todo List
- [ ] Define config types cho moi template
- [ ] Create 5 .tmpl template files
- [ ] Setup embed.FS cho templates
- [ ] Implement generator (template rendering)
- [ ] Implement interactive form cho moi template type
- [ ] Implement validator (nginx -t wrapper)
- [ ] Implement apply (symlink + validate + reload)
- [ ] Implement list command
- [ ] Implement Cobra commands + register
- [ ] Add SSL cert file existence check
- [ ] Implement preview mode (print config truoc khi save)
- [ ] Test voi real nginx installation
- [ ] Test rollback khi validation fail

## Success Criteria
- `idops nginx` hien interactive form, generate valid config
- Generated configs pass `nginx -t`
- 5 templates deu hoat dong chinh xac
- Apply command: symlink + validate + reload thanh cong
- Rollback khi validation fail (remove symlink)
- Preview mode hoat dong truoc khi save
- `--output` flag cho custom save path

## Risk Assessment
- **sudo required**: Apply/reload can root -> detect va suggest `sudo idops nginx apply`
- **Nginx not installed**: Check `nginx -v` truoc, hien clear error
- **Template bugs**: Generated config co the invalid -> always validate truoc apply
- **Path differences**: Debian sites-available vs CentOS conf.d -> configurable trong idops config

## Security Considerations
- Generated configs include security headers (X-Frame-Options, etc.) by default
- SSL config dung modern cipher suites
- Validate SSL cert/key paths exist truoc khi generate
- Apply command rollback neu validation fail (khong leave broken config)
- Warn neu running without SSL

## Next Steps
- Sau khi 6 phases xong: integration testing, README, GoReleaser release
- Future: nginx config analyzer, performance tuning templates (YAGNI)
