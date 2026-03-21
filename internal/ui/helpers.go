package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	lipglossv1 "github.com/charmbracelet/lipgloss"
)

// NewSpinner creates a spinner with the primary color theme.
func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipglossv1.NewStyle().Foreground(lipglossv1.Color("#7C3AED"))
	return s
}

// RenderError formats an error message with the error style.
func RenderError(msg string) string {
	return ErrorStyle.Render(fmt.Sprintf("Error: %s", msg))
}

// RenderSuccess formats a success message.
func RenderSuccess(msg string) string {
	return SuccessStyle.Render(fmt.Sprintf("✓ %s", msg))
}

// RenderWarning formats a warning message.
func RenderWarning(msg string) string {
	return WarningStyle.Render(fmt.Sprintf("! %s", msg))
}

// RenderInfo formats an info message.
func RenderInfo(msg string) string {
	return InfoStyle.Render(fmt.Sprintf("ℹ %s", msg))
}
