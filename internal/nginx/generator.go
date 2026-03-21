package nginx

import (
	"bytes"
	"fmt"
	"os"
)

// Generate renders the named template with the given config data.
func Generate(templateName string, config interface{}) (string, error) {
	tmpl, err := GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return buf.String(), nil
}

// SaveConfig writes generated nginx config to a file.
func SaveConfig(content, outputPath string) error {
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	return nil
}

// Preview prints the generated config to stdout.
func Preview(content string) {
	fmt.Println("# --- Generated Nginx Config ---")
	fmt.Println(content)
	fmt.Println("# --- End Config ---")
}
