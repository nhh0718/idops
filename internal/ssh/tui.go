package ssh

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipglossv1 "github.com/charmbracelet/lipgloss"
)

// viewMode tracks which screen is active.
type viewMode int

const (
	modeList viewMode = iota
	modeAdd
	modeEdit
	modeConfirmDelete
)

// hostItem implements list.Item for SSHHost.
type hostItem struct{ host SSHHost }

func (h hostItem) Title() string       { return h.host.Name }
func (h hostItem) Description() string { return fmt.Sprintf("%s@%s:%s", orDash(h.host.User), orDash(h.host.Hostname), orDefault(h.host.Port, "22")) }
func (h hostItem) FilterValue() string { return h.host.Name }

// formField labels for add/edit.
var formFields = []string{"Name", "Hostname", "Port", "User", "IdentityFile", "ProxyJump"}

// TUIModel is the root Bubble Tea model for the SSH manager.
type TUIModel struct {
	configPath  string
	list        list.Model
	inputs      []textinput.Model
	focusIdx    int
	mode        viewMode
	editingName string // original name when editing
	status      string
	testResults map[string]TestResult
}

// NewTUIModel creates a TUIModel loaded from the given SSH config path.
func NewTUIModel(configPath string) (TUIModel, error) {
	hosts, err := LoadConfig(configPath)
	if err != nil {
		return TUIModel{}, err
	}

	items := hostsToItems(hosts)
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipglossv1.Color("#FFFFFF")).
		Background(lipglossv1.Color("#7C3AED"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipglossv1.Color("#D1D5DB")).
		Background(lipglossv1.Color("#7C3AED"))

	l := list.New(items, delegate, 80, 20)
	l.Title = "SSH Manager"
	l.Styles.Title = lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#7C3AED"))

	return TUIModel{
		configPath:  configPath,
		list:        l,
		testResults: make(map[string]TestResult),
	}, nil
}

func (m TUIModel) Init() tea.Cmd { return nil }

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeAdd, modeEdit:
		return m.updateForm(msg)
	case modeConfirmDelete:
		return m.updateConfirm(msg)
	default:
		return m.updateList(msg)
	}
}

func (m TUIModel) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "a":
			m.mode = modeAdd
			m.inputs = makeInputs(SSHHost{})
			m.focusIdx = 0
			m.inputs[0].Focus()
			return m, nil
		case "e":
			if sel, ok := m.list.SelectedItem().(hostItem); ok {
				m.mode = modeEdit
				m.editingName = sel.host.Name
				m.inputs = makeInputs(sel.host)
				m.focusIdx = 0
				m.inputs[0].Focus()
			}
			return m, nil
		case "d":
			if _, ok := m.list.SelectedItem().(hostItem); ok {
				m.mode = modeConfirmDelete
			}
			return m, nil
		case "t":
			hosts := m.currentHosts()
			results := TestAllConnections(hosts, defaultTimeout)
			for _, r := range results {
				m.testResults[r.Host.Name] = r
			}
			m.status = fmt.Sprintf("Tested %d host(s)", len(results))
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TUIModel) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "y", "Y":
			if sel, ok := m.list.SelectedItem().(hostItem); ok {
				if err := DeleteHost(m.configPath, sel.host.Name); err != nil {
					m.status = "Delete failed: " + err.Error()
				} else {
					m.status = "Deleted " + sel.host.Name
					m.list.SetItems(m.reloadItems())
				}
			}
			m.mode = modeList
		case "n", "N", "esc":
			m.mode = modeList
		}
	}
	return m, nil
}

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
			// Save
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
				m.status = "Saved " + host.Name
				m.list.SetItems(m.reloadItems())
			}
			m.mode = modeList
			return m, nil
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

func (m TUIModel) View() string {
	switch m.mode {
	case modeAdd, modeEdit:
		return m.renderForm()
	case modeConfirmDelete:
		sel, _ := m.list.SelectedItem().(hostItem)
		return fmt.Sprintf("\nDelete host %q? [y/N] ", sel.host.Name)
	default:
		return m.renderList()
	}
}

func (m TUIModel) renderList() string {
	var sb strings.Builder
	sb.WriteString(m.list.View())
	if m.status != "" {
		sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981")).Render(m.status))
	}
	sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).Render("a add  e edit  d del  t test  q quit"))
	return sb.String()
}

func (m TUIModel) renderForm() string {
	title := "Add Host"
	if m.mode == modeEdit {
		title = "Edit Host"
	}
	var sb strings.Builder
	sb.WriteString(lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#7C3AED")).Render(title) + "\n\n")
	for i, inp := range m.inputs {
		sb.WriteString(lipglossv1.NewStyle().Bold(i == m.focusIdx).Render(formFields[i]+": ") + inp.View() + "\n")
	}
	sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).Render("tab/enter next  shift+tab prev  esc cancel"))
	return sb.String()
}

// --- helpers ---

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

func inputsToHost(inputs []textinput.Model) SSHHost {
	return SSHHost{
		Name:         inputs[0].Value(),
		Hostname:     inputs[1].Value(),
		Port:         inputs[2].Value(),
		User:         inputs[3].Value(),
		IdentityFile: inputs[4].Value(),
		ProxyJump:    inputs[5].Value(),
		Options:      make(map[string]string),
	}
}

func hostsToItems(hosts []SSHHost) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		items[i] = hostItem{h}
	}
	return items
}

func (m TUIModel) currentHosts() []SSHHost {
	items := m.list.Items()
	hosts := make([]SSHHost, 0, len(items))
	for _, it := range items {
		if h, ok := it.(hostItem); ok {
			hosts = append(hosts, h.host)
		}
	}
	return hosts
}

func (m TUIModel) reloadItems() []list.Item {
	hosts, err := LoadConfig(m.configPath)
	if err != nil {
		return m.list.Items()
	}
	return hostsToItems(hosts)
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
