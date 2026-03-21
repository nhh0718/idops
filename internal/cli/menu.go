package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var (
	menuTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED")).MarginBottom(1)
	menuItem  = lipgloss.NewStyle().PaddingLeft(2)
	menuSel   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7C3AED")).Bold(true)
	menuDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	menuHelp  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).MarginTop(1)
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

type menuModel struct {
	cursor   int
	quitting bool // true = user pressed q (exit), false = user pressed enter (select)
}

func (m menuModel) Init() tea.Cmd { return nil }

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
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
			m.quitting = false
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m menuModel) View() string {
	if m.quitting {
		return ""
	}

	s := menuTitle.Render("idops - DevOps Toolkit") + "  " + menuVer.Render(version) + "\n\n"

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

// showMenuAndExecute runs the interactive menu then dispatches to selected subcommand.
func showMenuAndExecute() error {
	m := menuModel{}
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return err
	}

	final := result.(menuModel)

	// User pressed q/esc/ctrl+c — exit cleanly
	if final.quitting {
		return nil
	}

	selected := menuEntries[final.cursor].cmd

	// Find and execute the subcommand
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == selected {
			sub.SetArgs([]string{}) // no extra args from menu
			return sub.Execute()
		}
	}
	return fmt.Errorf("command %q not found", selected)
}
