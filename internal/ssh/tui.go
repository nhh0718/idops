package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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

// clearStatusMsg signals the status bar should be cleared.
type clearStatusMsg struct{}

// connectDoneMsg carries the result of a connect re-exec.
type connectDoneMsg struct{ err error }

// clearStatusAfter returns a Cmd that fires clearStatusMsg after d.
func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return clearStatusMsg{} })
}

// hostItem implements list.Item for SSHHost.
type hostItem struct {
	host       SSHHost
	testResult *TestResult // nil = untested
}

func (h hostItem) Title() string {
	indicator := ""
	if h.testResult != nil {
		if h.testResult.Success {
			indicator = "  ✓"
		} else {
			indicator = "  ✗"
		}
	}
	return h.host.Name + indicator
}

func (h hostItem) Description() string {
	return fmt.Sprintf("%s@%s:%s", orDash(h.host.User), orDash(h.host.Hostname), orDefault(h.host.Port, "22"))
}

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

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipglossv1.Color("#FFFFFF")).
		Background(lipglossv1.Color("#7C3AED"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipglossv1.Color("#D1D5DB")).
		Background(lipglossv1.Color("#7C3AED"))

	l := list.New(hostsToItems(hosts, nil), delegate, 80, 20)
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
	// Global messages handled regardless of mode.
	switch msg.(type) {
	case clearStatusMsg:
		m.status = ""
		return m, nil
	}

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
	case connectDoneMsg:
		if msg.err != nil {
			m.status = "Connect error: " + msg.err.Error()
		} else {
			m.status = "Session ended"
		}
		return m, clearStatusAfter(3 * time.Second)

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

		case "c":
			if sel, ok := m.list.SelectedItem().(hostItem); ok {
				execPath, _ := os.Executable()
				cmd := exec.Command(execPath, "ssh", "connect", sel.host.Name)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return connectDoneMsg{err}
				})
			}
			return m, nil

		case "t":
			hosts := m.currentHosts()
			results := TestAllConnections(hosts, defaultTimeout)
			for _, r := range results {
				m.testResults[r.Host.Name] = r
			}
			m.status = fmt.Sprintf("Tested %d host(s)", len(results))
			m.list.SetItems(m.reloadItemsWithResults())
			return m, clearStatusAfter(3 * time.Second)
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
					m.list.SetItems(m.reloadItemsWithResults())
				}
			}
			m.mode = modeList
			return m, clearStatusAfter(3 * time.Second)
		case "n", "N", "esc":
			m.mode = modeList
		}
	}
	return m, nil
}

func (m TUIModel) View() string {
	switch m.mode {
	case modeAdd, modeEdit:
		return m.renderForm()
	case modeConfirmDelete:
		sel, _ := m.list.SelectedItem().(hostItem)
		h := sel.host
		prompt := fmt.Sprintf("\nDelete %s (%s@%s:%s)? [y/N] ",
			h.Name, orDash(h.User), orDash(h.Hostname), orDefault(h.Port, "22"))
		return lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#EF4444")).Render(prompt)
	default:
		return m.renderList()
	}
}

func (m TUIModel) renderList() string {
	var sb strings.Builder

	// Empty state
	if len(m.list.Items()) == 0 {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render("\n  No SSH hosts found in config\n"))
	} else {
		sb.WriteString(m.list.View())
	}

	if m.status != "" {
		sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981")).Render(m.status))
	}
	sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render("a add  e edit  d delete  c connect  t test  q quit"))
	return sb.String()
}

// --- helpers ---

func hostsToItems(hosts []SSHHost, results map[string]TestResult) []list.Item {
	items := make([]list.Item, len(hosts))
	for i, h := range hosts {
		item := hostItem{host: h}
		if results != nil {
			if r, ok := results[h.Name]; ok {
				r := r // capture
				item.testResult = &r
			}
		}
		items[i] = item
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

func (m TUIModel) reloadItemsWithResults() []list.Item {
	hosts, err := LoadConfig(m.configPath)
	if err != nil {
		return m.list.Items()
	}
	return hostsToItems(hosts, m.testResults)
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
