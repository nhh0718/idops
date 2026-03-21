package ssh

import (
	"fmt"
	"os"
	"os/exec"
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
	name := h.host.Name
	if h.testResult != nil {
		if h.testResult.Success {
			name += "  ✓ OK"
		} else {
			name += "  ✗ FAIL"
		}
	}
	return name
}

func (h hostItem) Description() string {
	user := orDash(h.host.User)
	host := orDash(h.host.Hostname)
	port := orDefault(h.host.Port, "22")
	desc := fmt.Sprintf("  Host: %s  |  User: %s  |  Port: %s", host, user, port)
	if h.host.IdentityFile != "" {
		desc += fmt.Sprintf("  |  Key: %s", h.host.IdentityFile)
	}
	return desc
}

func (h hostItem) FilterValue() string { return h.host.Name }

// formField labels for add/edit.
var formFields = []string{"Tên host", "Địa chỉ (IP/domain)", "Port", "User", "Key file", "ProxyJump"}

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
	l.Title = "🔑 SSH Manager - Quản lý kết nối SSH"
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


