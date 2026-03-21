package ports

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipglossv1 "github.com/charmbracelet/lipgloss"
)

type refreshMsg struct{ infos []PortInfo }
type errMsg struct{ err error }
type tickMsg time.Time
type clearStatusMsg struct{}

// TUIModel is the Bubble Tea model for the ports TUI.
type TUIModel struct {
	table     table.Model
	filter    textinput.Model
	infos     []PortInfo
	filtered  []PortInfo
	opts      ScanOptions
	sortField SortField
	watchMode bool
	interval  time.Duration
	filtering bool
	confirm   bool
	lastErr   string
	statusMsg string
	statusErr bool // true = error style, false = success style
	lastScan  time.Time
	width     int
	height    int
}

var sortLabels = []string{"Port", "PID", "Process", "Protocol"}

// NewTUIModel initialises the model.
func NewTUIModel(opts ScanOptions, watchMode bool, interval time.Duration) TUIModel {
	fi := textinput.New()
	fi.Placeholder = "filter process or port..."
	fi.CharLimit = 64

	return TUIModel{
		filter:    fi,
		opts:      opts,
		sortField: SortByPort,
		watchMode: watchMode,
		interval:  interval,
		width:     100,
		height:    24,
	}
}

func (m TUIModel) Init() tea.Cmd {
	return tea.Batch(doScan(m.opts), tea.EnterAltScreen)
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.rebuildTable()
		return m, nil

	case refreshMsg:
		m.infos = msg.infos
		m.lastErr = ""
		m.lastScan = time.Now()
		m.rebuildTable()
		if m.watchMode {
			return m, tea.Tick(m.interval, func(t time.Time) tea.Msg { return tickMsg(t) })
		}
		return m, nil

	case errMsg:
		m.lastErr = msg.err.Error()
		// Keep ticking even on error
		if m.watchMode {
			return m, tea.Tick(m.interval, func(t time.Time) tea.Msg { return tickMsg(t) })
		}
		return m, nil

	case tickMsg:
		return m, doScan(m.opts)

	case clearStatusMsg:
		if !m.confirm {
			m.statusMsg = ""
			m.statusErr = false
		}
		return m, nil

	case tea.KeyMsg:
		if m.confirm {
			return m.handleConfirm(msg)
		}
		if m.filtering {
			return m.handleFilterKey(msg)
		}
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// clearStatusAfter returns a command that clears status after delay.
func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return clearStatusMsg{} })
}

func (m *TUIModel) setStatus(msg string, isErr bool) tea.Cmd {
	m.statusMsg = msg
	m.statusErr = isErr
	if isErr {
		return clearStatusAfter(5 * time.Second)
	}
	return clearStatusAfter(3 * time.Second)
}

func (m TUIModel) handleConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		m.confirm = false
		selected := m.selectedPort()
		if selected != nil {
			if err := KillProcess(selected.PID); err != nil {
				cmd := m.setStatus("Kill failed: "+err.Error(), true)
				return m, cmd
			}
			cmd := m.setStatus(
				fmt.Sprintf("Killed %s (PID %d) on port %d", selected.ProcessName, selected.PID, selected.LocalPort),
				false,
			)
			return m, tea.Batch(cmd, doScan(m.opts))
		}
	default:
		m.confirm = false
		m.statusMsg = ""
	}
	return m, nil
}

func (m TUIModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "s":
		m.sortField = (m.sortField + 1) % 4
		m.rebuildTable()
	case "/":
		m.filtering = true
		m.filter.Focus()
	case "r":
		return m, doScan(m.opts)
	case "k":
		selected := m.selectedPort()
		if selected != nil && selected.PID > 0 {
			m.confirm = true
			m.statusMsg = fmt.Sprintf("Kill %s (PID %d) on port %d? [y/N]", selected.ProcessName, selected.PID, selected.LocalPort)
		} else {
			cmd := m.setStatus("No process selected or PID unavailable", true)
			return m, cmd
		}
		return m, nil
	case "o":
		selected := m.selectedPort()
		if selected != nil {
			if err := OpenInBrowser(selected.LocalPort); err != nil {
				cmd := m.setStatus("Browser failed: "+err.Error(), true)
				return m, cmd
			}
			cmd := m.setStatus(fmt.Sprintf("Opened http://localhost:%d", selected.LocalPort), false)
			return m, cmd
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m TUIModel) handleFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.filtering = false
		m.filter.Blur()
		m.rebuildTable()
		return m, nil
	}
	var cmd tea.Cmd
	m.filter, cmd = m.filter.Update(msg)
	m.rebuildTable()
	return m, cmd
}

// View renders the full TUI screen.
func (m TUIModel) View() string {
	var sb strings.Builder

	// Title bar with stats
	title := lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#7C3AED")).Render("idops ports")
	sortInfo := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render(fmt.Sprintf("  sort: %s", sortLabels[m.sortField]))
	stats := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render(fmt.Sprintf("  [%d ports", len(m.filtered)))
	if len(m.filtered) != len(m.infos) {
		stats += lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render(fmt.Sprintf("/%d total", len(m.infos)))
	}
	stats += lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).Render("]")

	watchInfo := ""
	if m.watchMode {
		watchInfo = lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981")).
			Render(fmt.Sprintf("  WATCH %s", m.interval))
	}
	sb.WriteString(title + sortInfo + stats + watchInfo + "\n")

	// Filter bar
	if m.filtering {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#3B82F6")).Render("/") + " " + m.filter.View() + "\n")
	} else if m.filter.Value() != "" {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render(fmt.Sprintf("filter: %q  (press / to change)", m.filter.Value())) + "\n")
	}

	// Error line
	if m.lastErr != "" {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#EF4444")).
			Render("error: "+m.lastErr) + "\n")
	}

	// Empty state
	if len(m.infos) == 0 && m.lastErr == "" && !m.lastScan.IsZero() {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#F59E0B")).
			Render("\n  No listening ports found.\n  Try running with elevated privileges for full results.\n") + "\n")
	} else {
		sb.WriteString(m.table.View() + "\n")
	}

	// Status message with appropriate color
	if m.statusMsg != "" {
		style := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981"))
		if m.confirm {
			style = lipglossv1.NewStyle().Foreground(lipglossv1.Color("#F59E0B")).Bold(true)
		} else if m.statusErr {
			style = lipglossv1.NewStyle().Foreground(lipglossv1.Color("#EF4444"))
		}
		sb.WriteString(style.Render(m.statusMsg) + "\n")
	}

	// Help footer
	help := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render("q quit  s sort  / filter  r refresh  k kill  o browser")
	sb.WriteString(help)

	return sb.String()
}

func (m *TUIModel) selectedPort() *PortInfo {
	row := m.table.SelectedRow()
	if row == nil || len(m.filtered) == 0 {
		return nil
	}
	idx := m.table.Cursor()
	if idx >= 0 && idx < len(m.filtered) {
		return &m.filtered[idx]
	}
	return nil
}

func (m *TUIModel) rebuildTable() {
	rows := m.visibleRows()
	t := newPortsTable(rows, m.width)
	h := m.height - 7
	if m.statusMsg != "" {
		h--
	}
	if m.filtering || m.filter.Value() != "" {
		h--
	}
	if h < 5 {
		h = 5
	}
	t.SetHeight(h)
	m.table = t
}

func (m *TUIModel) visibleRows() []table.Row {
	query := strings.ToLower(m.filter.Value())
	infos := make([]PortInfo, 0, len(m.infos))
	for _, p := range m.infos {
		if query == "" ||
			strings.Contains(strings.ToLower(p.ProcessName), query) ||
			strings.Contains(fmt.Sprintf("%d", p.LocalPort), query) {
			infos = append(infos, p)
		}
	}
	SortPortInfos(infos, m.sortField)
	m.filtered = infos

	rows := make([]table.Row, len(infos))
	for i, p := range infos {
		proto := p.Protocol
		rows[i] = table.Row{
			proto,
			p.LocalAddr,
			fmt.Sprintf("%d", p.LocalPort),
			fmt.Sprintf("%d", p.PID),
			p.ProcessName,
			p.Status,
		}
	}
	return rows
}

func newPortsTable(rows []table.Row, width int) table.Model {
	cols := []table.Column{
		{Title: "Proto", Width: 6},
		{Title: "Local Addr", Width: 16},
		{Title: "Port", Width: 7},
		{Title: "PID", Width: 8},
		{Title: "Process", Width: max(width-55, 16)},
		{Title: "Status", Width: 8},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(18),
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

func doScan(opts ScanOptions) tea.Cmd {
	return func() tea.Msg {
		infos, err := Scan(context.Background(), opts)
		if err != nil {
			return errMsg{err}
		}
		return refreshMsg{infos}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
