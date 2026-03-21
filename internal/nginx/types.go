package nginx

// BaseConfig holds common nginx server block settings.
type BaseConfig struct {
	ServerName  string
	ListenPort  int
	SSLEnabled  bool
	SSLCertPath string
	SSLKeyPath  string
	AccessLog   string
	ErrorLog    string
}

// ReverseProxyConfig for upstream proxy sites.
type ReverseProxyConfig struct {
	BaseConfig
	UpstreamHost string
	UpstreamPort int
	WebSocket    bool
}

// StaticSiteConfig for serving static files.
type StaticSiteConfig struct {
	BaseConfig
	RootPath   string
	IndexFiles []string
	EnableGzip bool
	CacheMaxAge int
}

// PHPFPMConfig for PHP applications.
type PHPFPMConfig struct {
	BaseConfig
	RootPath  string
	PHPSocket string
}

// LoadBalancerConfig for upstream load balancing.
type LoadBalancerConfig struct {
	BaseConfig
	UpstreamName string
	Backends     []Backend
	Method       string // round-robin, least_conn, ip_hash
}

// Backend represents a single upstream server.
type Backend struct {
	Host   string
	Port   int
	Weight int
}

// TemplateType enumerates available nginx templates.
type TemplateType string

const (
	TemplateReverseProxy TemplateType = "reverse-proxy"
	TemplateStaticSite   TemplateType = "static-site"
	TemplatePHPFPM       TemplateType = "php-fpm"
	TemplateLoadBalancer TemplateType = "load-balancer"
	TemplateWebSocket    TemplateType = "websocket"
)
