package cli

import (
	"bufio"
	"fmt"
	"os"
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
	{"🔍", "Port Scanner", "Scan ports, kill processes, watch mode", "ports"},
	{"🐳", "Docker Dashboard", "Container stats, start/stop/restart, logs", "docker"},
	{"🔑", "SSH Manager", "Manage SSH hosts, connect, test connections", "ssh"},
	{"📋", "Env Sync", "Compare, sync, validate .env files", "env"},
	{"⚙️ ", "Nginx Generator", "Generate nginx configs from templates", "nginx"},
}

// showMenuAndExecute displays a numbered menu and runs the selected subcommand.
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
	fmt.Print(menuDesc.Render("  Select [1-5] or q to quit: "))

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

	// Find and execute the subcommand
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == selected {
			sub.SetArgs([]string{})
			return sub.Execute()
		}
	}
	return fmt.Errorf("command %q not found", selected)
}
