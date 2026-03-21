package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	menuTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).MarginBottom(1)
	menuItem  = lipgloss.NewStyle().PaddingLeft(2)
	menuSel   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7C3AED")).Bold(true)
	menuDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	menuHelp  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).MarginTop(1)
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

type menuModel struct {
	cursor int
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuEntries)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	s := menuTitle.Render("idops - DevOps Toolkit") + "\n\n"

	for i, entry := range menuEntries {
		line := fmt.Sprintf("%s %s  %s", entry.icon, entry.name, menuDesc.Render(entry.desc))
		if i == m.cursor {
			s += menuSel.Render("▸ " + line) + "\n"
		} else {
			s += menuItem.Render("  " + line) + "\n"
		}
	}

	s += "\n" + menuHelp.Render("  ↑↓/jk navigate • enter select • q quit")
	return s
}

// runMenu shows interactive menu and returns selected command name.
func runMenu() (string, error) {
	m := menuModel{}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	final := result.(menuModel)
	return menuEntries[final.cursor].cmd, nil
}

// showMenuAndExecute runs the interactive menu then dispatches to selected subcommand.
func showMenuAndExecute() error {
	cmd, err := runMenu()
	if err != nil {
		return err
	}
	if cmd == "" {
		return nil
	}

	// Find and execute the subcommand
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == cmd {
			sub.SetArgs(os.Args[1:]) // pass remaining args
			return sub.Execute()
		}
	}
	return fmt.Errorf("command %q not found", cmd)
}
