package ui

import (
	"github.com/charmbracelet/bubbles/table"
	lipglossv1 "github.com/charmbracelet/lipgloss"
)

// NewTable creates a styled Bubble Tea table with default theme colors.
func NewTable(columns []table.Column, rows []table.Row, height int) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipglossv1.NormalBorder()).
		BorderForeground(lipglossv1.Color("#6B7280")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipglossv1.Color("#FFFFFF")).
		Background(lipglossv1.Color("#7C3AED")).
		Bold(false)

	t.SetStyles(s)
	return t
}
