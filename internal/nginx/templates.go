package nginx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

// GetTemplate loads a named nginx template from the templates directory.
func GetTemplate(name string) (*template.Template, error) {
	// Find templates dir relative to executable or working directory
	paths := []string{
		filepath.Join("templates", "nginx", name+".tmpl"),
		filepath.Join(execDir(), "templates", "nginx", name+".tmpl"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return template.ParseFiles(path)
		}
	}

	return nil, fmt.Errorf("template %q not found", name)
}

// AvailableTemplates returns the list of template names.
func AvailableTemplates() []string {
	return []string{
		"reverse-proxy",
		"static-site",
		"php-fpm",
		"load-balancer",
		"websocket",
	}
}

func execDir() string {
	if runtime.GOOS == "windows" {
		exe, err := os.Executable()
		if err != nil {
			return "."
		}
		return filepath.Dir(exe)
	}
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	// Resolve symlinks
	real, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return filepath.Dir(exe)
	}
	return filepath.Dir(real)
}
