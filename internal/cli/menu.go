package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	menuTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED"))
	menuIdx   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")).Bold(true)
	menuName  = lipgloss.NewStyle().Bold(true)
	menuDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	menuVer   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
)

type menuEntry struct {
	icon string
	name string
	desc string
	cmd  string
}

var menuEntries = []menuEntry{
	{"🔍", "Port Scanner", "Quét port, kill process, theo dõi realtime", "ports"},
	{"🐳", "Docker Dashboard", "Quản lý container, start/stop/restart, xem logs", "docker"},
	{"🔑", "SSH Manager", "Quản lý SSH host, kết nối, test connection", "ssh"},
	{"📋", "Env Sync", "So sánh, đồng bộ, validate file .env", "env"},
	{"⚙️ ", "Nginx Generator", "Tạo config nginx từ template", "nginx"},
	{"🌐", "Web Dashboard", "Mở dashboard trên trình duyệt", "dashboard"},
}

// showMenuAndExecute displays a numbered menu and re-executes idops with selected command.
func showMenuAndExecute() error {
	fmt.Println()
	fmt.Println(menuTitle.Render("  idops - DevOps Toolkit") + "  " + menuVer.Render(version))
	fmt.Println()

	for i, entry := range menuEntries {
		num := menuIdx.Render(fmt.Sprintf("  [%d]", i+1))
		name := menuName.Render(fmt.Sprintf("%s %s", entry.icon, entry.name))
		desc := menuDesc.Render(entry.desc)
		fmt.Printf("%s %s  %s\n", num, name, desc)
	}

	fmt.Println()
	fmt.Print(menuDesc.Render(fmt.Sprintf("  Select [1-%d] or q to quit: ", len(menuEntries))))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "" {
		return nil
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(menuEntries) {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Render("  Invalid choice"))
		return nil
	}

	selected := menuEntries[idx-1].cmd
	fmt.Println()

	// Re-exec the binary with the selected subcommand.
	// This gives the subcommand a fresh terminal/stdin state,
	// which is required for Bubble Tea TUI programs to work.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable: %w", err)
	}

	cmd := exec.Command(execPath, selected)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
