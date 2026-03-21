package ui

import (
	lipgloss "charm.land/lipgloss/v2"
)

// Color palette for consistent TUI styling.
var (
	Primary = lipgloss.Color("#7C3AED")
	Success = lipgloss.Color("#10B981")
	Warning = lipgloss.Color("#F59E0B")
	Error   = lipgloss.Color("#EF4444")
	Muted   = lipgloss.Color("#6B7280")
	Info    = lipgloss.Color("#3B82F6")
)

// Reusable styles.
var (
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(Primary)
	ErrorStyle   = lipgloss.NewStyle().Foreground(Error)
	SuccessStyle = lipgloss.NewStyle().Foreground(Success)
	WarningStyle = lipgloss.NewStyle().Foreground(Warning)
	MutedStyle   = lipgloss.NewStyle().Foreground(Muted)
	InfoStyle    = lipgloss.NewStyle().Foreground(Info)
	BoldStyle    = lipgloss.NewStyle().Bold(true)
)
