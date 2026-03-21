package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nhh0718/idops/internal/nginx"
	"github.com/nhh0718/idops/internal/ui"
	"github.com/spf13/cobra"
)

var nginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "Nginx config generator",
	RunE:  runNginxGenerate,
}

var nginxValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate nginx configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := nginx.ValidateNginxConfig(); err != nil {
			return err
		}
		fmt.Println(ui.RenderSuccess("Nginx configuration is valid"))
		return nil
	},
}

var nginxApplyCmd = &cobra.Command{
	Use:   "apply <config-file>",
	Short: "Enable config and reload nginx",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sitesEnabled, _ := cmd.Flags().GetString("sites-enabled")
		return nginx.Apply(args[0], sitesEnabled)
	},
}

var nginxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List nginx configs",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		configs, err := nginx.ListConfigs(dir)
		if err != nil {
			return err
		}
		if len(configs) == 0 {
			fmt.Println("No configs found")
			return nil
		}
		for _, c := range configs {
			fmt.Printf("  %s\n", c)
		}
		return nil
	},
}

func init() {
	nginxApplyCmd.Flags().String("sites-enabled", "/etc/nginx/sites-enabled", "sites-enabled directory")
	nginxListCmd.Flags().String("dir", "/etc/nginx/sites-available", "config directory")
	nginxCmd.Flags().String("output", "", "output file path (default: stdout/preview)")

	nginxCmd.AddCommand(nginxValidateCmd, nginxApplyCmd, nginxListCmd)
	rootCmd.AddCommand(nginxCmd)
}

func runNginxGenerate(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Select template
	fmt.Println("Select Nginx config template:")
	templates := []struct{ label, value string }{
		{"Reverse Proxy", "reverse-proxy"},
		{"Static Site", "static-site"},
		{"PHP-FPM", "php-fpm"},
		{"Load Balancer", "load-balancer"},
		{"WebSocket Proxy", "websocket"},
	}
	for i, t := range templates {
		fmt.Printf("  %d) %s\n", i+1, t.label)
	}
	fmt.Print("\nChoice [1-5]: ")
	choice, _ := reader.ReadString('\n')
	idx, _ := strconv.Atoi(strings.TrimSpace(choice))
	if idx < 1 || idx > 5 {
		return fmt.Errorf("invalid choice")
	}
	tmplName := templates[idx-1].value

	// Step 2: Collect common fields
	domain := prompt(reader, "Server name (domain)", "example.com")
	portStr := prompt(reader, "Listen port", "80")
	port, _ := strconv.Atoi(portStr)
	ssl := promptBool(reader, "Enable SSL?", false)

	base := nginx.BaseConfig{
		ServerName: domain,
		ListenPort: port,
		SSLEnabled: ssl,
	}
	if ssl {
		base.SSLCertPath = prompt(reader, "SSL cert path", "/etc/ssl/certs/cert.pem")
		base.SSLKeyPath = prompt(reader, "SSL key path", "/etc/ssl/private/key.pem")
	}

	// Step 3: Template-specific fields
	var config interface{}
	switch tmplName {
	case "reverse-proxy":
		config = buildReverseProxy(reader, base)
	case "static-site":
		config = buildStaticSite(reader, base)
	case "php-fpm":
		config = buildPHPFPM(reader, base)
	case "load-balancer":
		config = buildLoadBalancer(reader, base)
	case "websocket":
		config = buildWebSocket(reader, base)
	}

	// Step 4: Generate and preview
	content, err := nginx.Generate(tmplName, config)
	if err != nil {
		return err
	}
	nginx.Preview(content)

	// Step 5: Save if output specified
	output, _ := cmd.Flags().GetString("output")
	if output != "" {
		if err := nginx.SaveConfig(content, output); err != nil {
			return err
		}
		fmt.Println(ui.RenderSuccess("Saved to " + output))
	} else {
		if promptBool(reader, "Save to file?", false) {
			path := prompt(reader, "Output path", domain+".conf")
			if err := nginx.SaveConfig(content, path); err != nil {
				return err
			}
			fmt.Println(ui.RenderSuccess("Saved to " + path))
		}
	}
	return nil
}

func buildReverseProxy(r *bufio.Reader, base nginx.BaseConfig) nginx.ReverseProxyConfig {
	host := prompt(r, "Upstream host", "127.0.0.1")
	portStr := prompt(r, "Upstream port", "3000")
	port, _ := strconv.Atoi(portStr)
	ws := promptBool(r, "WebSocket support?", false)
	return nginx.ReverseProxyConfig{BaseConfig: base, UpstreamHost: host, UpstreamPort: port, WebSocket: ws}
}

func buildStaticSite(r *bufio.Reader, base nginx.BaseConfig) nginx.StaticSiteConfig {
	root := prompt(r, "Document root", "/var/www/html")
	gzip := promptBool(r, "Enable Gzip?", true)
	cacheStr := prompt(r, "Cache max-age (days, 0=disabled)", "30")
	cache, _ := strconv.Atoi(cacheStr)
	return nginx.StaticSiteConfig{BaseConfig: base, RootPath: root, IndexFiles: []string{"index.html", "index.htm"}, EnableGzip: gzip, CacheMaxAge: cache}
}

func buildPHPFPM(r *bufio.Reader, base nginx.BaseConfig) nginx.PHPFPMConfig {
	root := prompt(r, "Document root", "/var/www/html")
	sock := prompt(r, "PHP-FPM socket", "/run/php/php8.2-fpm.sock")
	return nginx.PHPFPMConfig{BaseConfig: base, RootPath: root, PHPSocket: sock}
}

func buildLoadBalancer(r *bufio.Reader, base nginx.BaseConfig) nginx.LoadBalancerConfig {
	name := prompt(r, "Upstream name", "backend")
	method := prompt(r, "Method (round-robin/least_conn/ip_hash)", "round-robin")
	countStr := prompt(r, "Number of backends", "2")
	count, _ := strconv.Atoi(countStr)
	var backends []nginx.Backend
	for i := 0; i < count; i++ {
		host := prompt(r, fmt.Sprintf("Backend %d host", i+1), "127.0.0.1")
		portStr := prompt(r, fmt.Sprintf("Backend %d port", i+1), fmt.Sprintf("%d", 3000+i))
		port, _ := strconv.Atoi(portStr)
		backends = append(backends, nginx.Backend{Host: host, Port: port, Weight: 1})
	}
	return nginx.LoadBalancerConfig{BaseConfig: base, UpstreamName: name, Backends: backends, Method: method}
}

func buildWebSocket(r *bufio.Reader, base nginx.BaseConfig) nginx.ReverseProxyConfig {
	host := prompt(r, "WebSocket upstream host", "127.0.0.1")
	portStr := prompt(r, "WebSocket upstream port", "8080")
	port, _ := strconv.Atoi(portStr)
	return nginx.ReverseProxyConfig{BaseConfig: base, UpstreamHost: host, UpstreamPort: port, WebSocket: true}
}

func prompt(r *bufio.Reader, label, defaultVal string) string {
	fmt.Printf("  %s [%s]: ", label, defaultVal)
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func promptBool(r *bufio.Reader, label string, defaultVal bool) bool {
	def := "n"
	if defaultVal {
		def = "y"
	}
	fmt.Printf("  %s [%s]: ", label, def)
	input, _ := r.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultVal
	}
	return input == "y" || input == "yes"
}
