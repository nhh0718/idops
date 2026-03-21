package docker

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	dockerclient "github.com/docker/docker/client"
)

func (m dashModel) handleConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.confirm = false
		action := m.confirmKey
		id, name := m.confirmID, m.confirmName
		m.confirmID, m.confirmName, m.confirmKey = "", "", ""
		return m, runAction(m.cli, action, id, name)
	default:
		// Any other key cancels
		m.confirm = false
		m.statusMsg = "Cancelled"
		m.statusIsErr = false
		m.confirmID, m.confirmName, m.confirmKey = "", "", ""
		return m, clearStatusAfter(2 * time.Second)
	}
}

func (m dashModel) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "esc":
		m.filtering = false
		m.filter.Blur()
		m.applyFilter()
		return m, nil
	}
	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	m.applyFilter()
	return m, cmd
}

func (m dashModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	selected := m.selectedContainer()
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "/":
		m.filtering = true
		m.filter.Focus()
		return m, textinput.Blink
	case "s":
		if selected != nil {
			return m, runAction(m.cli, "s", selected.ID, selected.Name)
		}
	case "x":
		if selected != nil {
			m.confirm, m.confirmKey = true, "x"
			m.confirmID, m.confirmName = selected.ID, selected.Name
			return m, nil
		}
	case "r":
		if selected != nil {
			m.confirm, m.confirmKey = true, "r"
			m.confirmID, m.confirmName = selected.ID, selected.Name
			return m, nil
		}
	case "d":
		if selected != nil {
			m.confirm, m.confirmKey = true, "d"
			m.confirmID, m.confirmName = selected.ID, selected.Name
			return m, nil
		}
	case "l":
		if selected != nil {
			return m, launchLogViewer(m.cli, selected.ID, selected.Name)
		}
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m dashModel) selectedContainer() *ContainerInfo {
	idx := m.table.Cursor()
	if idx >= 0 && idx < len(m.filtered) {
		c := m.filtered[idx]
		return &c
	}
	return nil
}

func launchLogViewer(cli *dockerclient.Client, id, name string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		lm, err := NewLogViewer(ctx, cli, id, name)
		if err != nil {
			return errMsg(err)
		}
		p := tea.NewProgram(lm, tea.WithAltScreen())
		_, _ = p.Run()
		return nil
	}
}
