package docker

import (
	"fmt"
	"strings"
)

// colorState returns a styled state string.
func colorState(state string) string {
	switch strings.ToLower(state) {
	case "running":
		return stateGreen.Render(state)
	case "exited", "dead":
		return stateRed.Render(state)
	case "paused":
		return stateYellow.Render(state)
	default:
		return state
	}
}

// View renders the dashboard TUI.
func (m dashModel) View() string {
	var sb strings.Builder

	// Title + stats
	running := 0
	for _, c := range m.containers {
		if strings.ToLower(c.State) == "running" {
			running++
		}
	}
	title := titleStyle.Render("idops docker")
	stats := helpStyle.Render(fmt.Sprintf("  [%d containers, %d running]", len(m.containers), running))
	sb.WriteString(title + stats + "\n")

	// Filter bar
	if m.filtering {
		sb.WriteString(helpStyle.Render("/") + " " + m.filter.View() + "\n")
	} else if m.filter.Value() != "" {
		sb.WriteString(helpStyle.Render(fmt.Sprintf("filter: %q  (/ to change)", m.filter.Value())) + "\n")
	}

	// Error
	if m.err != nil {
		sb.WriteString(statusErr.Render("Docker error: "+m.err.Error()) + "\n")
	}

	// Empty state
	if len(m.containers) == 0 && m.err == nil {
		sb.WriteString(stateYellow.Render("\n  No containers found. Is Docker running?\n") + "\n")
	} else {
		sb.WriteString(m.table.View() + "\n")
	}

	// Confirm prompt
	if m.confirm {
		action := map[string]string{"x": "Stop", "r": "Restart", "d": "REMOVE"}[m.confirmKey]
		sb.WriteString(stateYellow.Bold(true).Render(
			fmt.Sprintf("  %s container %q? [y/N]", action, m.confirmName)) + "\n")
	}

	// Status message
	if m.statusMsg != "" && !m.confirm {
		style := statusOK
		if m.statusIsErr {
			style = statusErr
		}
		sb.WriteString(style.Render("  "+m.statusMsg) + "\n")
	}

	// Help footer
	sb.WriteString(helpStyle.Render("q quit  s start  x stop  r restart  d remove  l logs  / filter"))

	return sb.String()
}
