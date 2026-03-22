package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// runEnvMenu shows interactive sub-menu when `idops env` is called without subcommand.
func runEnvMenu(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).Render("  📋 Env Sync - Quản lý file môi trường"))
	fmt.Println()

	entries := []struct {
		key  string
		name string
		desc string
	}{
		{"1", "Compare", "So sánh .env.example với .env"},
		{"2", "Sync", "Đồng bộ biến thiếu từ .env.example"},
		{"3", "Validate", "Kiểm tra format .env"},
		{"4", "Init", "Tạo .env từ .env.example"},
		{"5", "Show", "Hiển thị .env (ẩn secrets)"},
	}

	for _, e := range entries {
		num := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).Render(fmt.Sprintf("  [%s]", e.key))
		name := lipgloss.NewStyle().Bold(true).Render(e.name)
		desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render(e.desc)
		fmt.Printf("%s %s  %s\n", num, name, desc)
	}

	fmt.Println()
	fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("  Chọn [1-5], b: quay lại, q: thoát: "))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "Q" {
		os.Exit(0)
	}
	if input == "b" || input == "B" || input == "" {
		return nil // return to main menu
	}

	subCmds := []string{"compare", "sync", "validate", "init", "show"}
	idx := 0
	if _, err := fmt.Sscanf(input, "%d", &idx); err != nil || idx < 1 || idx > len(subCmds) {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("  Lựa chọn không hợp lệ"))
		return nil
	}

	// Re-exec with subcommand for fresh stdin
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable: %w", err)
	}
	c := exec.Command(execPath, "env", subCmds[idx-1])
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
