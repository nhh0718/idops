package ssh

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipglossv1 "github.com/charmbracelet/lipgloss"
)

func (m TUIModel) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "esc":
			m.mode = modeList
			return m, nil

		case "enter", "tab":
			if m.focusIdx < len(m.inputs)-1 {
				m.inputs[m.focusIdx].Blur()
				m.focusIdx++
				m.inputs[m.focusIdx].Focus()
				return m, nil
			}
			// Last field — validate and save.
			if validationErr := validateFormInputs(m.inputs); validationErr != "" {
				m.status = validationErr
				return m, clearStatusAfter(3 * time.Second)
			}
			host := inputsToHost(m.inputs)
			var err error
			if m.mode == modeEdit {
				err = UpdateHost(m.configPath, m.editingName, host)
			} else {
				err = AddHost(m.configPath, host)
			}
			if err != nil {
				m.status = "Save failed: " + err.Error()
			} else {
				m.status = "Host saved: " + host.Name
				m.list.SetItems(m.reloadItemsWithResults())
			}
			m.mode = modeList
			return m, clearStatusAfter(3 * time.Second)

		case "shift+tab":
			if m.focusIdx > 0 {
				m.inputs[m.focusIdx].Blur()
				m.focusIdx--
				m.inputs[m.focusIdx].Focus()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
	return m, cmd
}

func (m TUIModel) renderForm() string {
	title := "Add Host"
	if m.mode == modeEdit {
		title = "Edit Host"
	}

	var sb strings.Builder
	sb.WriteString(lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#7C3AED")).Render(title) + "\n\n")

	for i, inp := range m.inputs {
		label := lipglossv1.NewStyle().Bold(i == m.focusIdx).Render(formFields[i] + ": ")
		sb.WriteString(label + inp.View() + "\n")
	}

	if m.status != "" {
		sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#EF4444")).Render(m.status))
	}
	sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render("tab/enter next  shift+tab prev  esc cancel"))
	return sb.String()
}

// validateFormInputs returns a non-empty error string if inputs are invalid.
func validateFormInputs(inputs []textinput.Model) string {
	name := strings.TrimSpace(inputs[0].Value())
	if name == "" {
		return "Name is required"
	}

	portStr := strings.TrimSpace(inputs[2].Value())
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil || p < 1 || p > 65535 {
			return fmt.Sprintf("Port must be a number between 1 and 65535 (got %q)", portStr)
		}
	}

	return ""
}

// makeInputs builds textinput models pre-filled from host.
func makeInputs(h SSHHost) []textinput.Model {
	vals := []string{h.Name, h.Hostname, h.Port, h.User, h.IdentityFile, h.ProxyJump}
	inputs := make([]textinput.Model, len(formFields))
	for i, label := range formFields {
		t := textinput.New()
		t.Placeholder = label
		t.SetValue(vals[i])
		inputs[i] = t
	}
	return inputs
}

// inputsToHost converts form input values into an SSHHost.
func inputsToHost(inputs []textinput.Model) SSHHost {
	return SSHHost{
		Name:         strings.TrimSpace(inputs[0].Value()),
		Hostname:     strings.TrimSpace(inputs[1].Value()),
		Port:         strings.TrimSpace(inputs[2].Value()),
		User:         strings.TrimSpace(inputs[3].Value()),
		IdentityFile: strings.TrimSpace(inputs[4].Value()),
		ProxyJump:    strings.TrimSpace(inputs[5].Value()),
		Options:      make(map[string]string),
	}
}
