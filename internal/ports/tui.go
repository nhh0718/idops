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

// refreshMsg is sent on watch-mode tick or manual refresh.
type refreshMsg struct{ infos []PortInfo }

// errMsg carries a scan error back to the model.
type errMsg struct{ err error }

// TUIModel is the Bubble Tea model for the ports TUI.
type TUIModel struct {
	table      table.Model
	filter     textinput.Model
	infos      []PortInfo   // full unfiltered dataset
	opts       ScanOptions  // scan options (protocol, port range)
	sortField  SortField
	watchMode  bool
	interval   time.Duration
	filtering  bool  // whether filter input is active
	lastErr    string
	width      int
	height     int
}

var sortLabels = []string{"Port", "PID", "Process", "Protocol"}

// NewTUIModel initialises the model. Call tea.NewProgram(model).Run() to start.
func NewTUIModel(opts ScanOptions, watchMode bool, interval time.Duration) TUIModel {
	fi := textinput.New()
	fi.Placeholder = "filter process or port..."
	fi.CharLimit = 64

	m := TUIModel{
		filter:    fi,
		opts:      opts,
		sortField: SortByPort,
		watchMode: watchMode,
		interval:  interval,
		width:     100,
		height:    24,
	}
	m.table = newPortsTable(nil, m.width)
	return m
}

// Init performs the first scan and, in watch mode, starts the tick.
func (m TUIModel) Init() tea.Cmd {
	return tea.Batch(doScan(m.opts), tea.EnterAltScreen)
}

// Update handles messages and key events.
func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.table = newPortsTable(m.visibleRows(), m.width)
		return m, nil

	case refreshMsg:
		m.infos = msg.infos
		m.lastErr = ""
		m.rebuildTable()
		if m.watchMode {
			return m, tea.Tick(m.interval, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
		return m, nil

	case errMsg:
		m.lastErr = msg.err.Error()
		return m, nil

	case tickMsg:
		return m, doScan(m.opts)

	case tea.KeyMsg:
		// Delegate keys to filter input when active.
		if m.filtering {
			return m.handleFilterKey(msg)
		}
		return m.handleKey(msg)
	}

	// Pass through to table when not filtering.
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

type tickMsg time.Time

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

	// Title bar.
	title := lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#7C3AED")).Render("idops ports")
	sortInfo := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render(fmt.Sprintf("  sort: %s", sortLabels[m.sortField]))
	watchInfo := ""
	if m.watchMode {
		watchInfo = lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981")).
			Render(fmt.Sprintf("  watch: %s", m.interval))
	}
	sb.WriteString(title + sortInfo + watchInfo + "\n")

	// Filter bar.
	if m.filtering {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#3B82F6")).Render("/") + " " + m.filter.View() + "\n")
	} else if m.filter.Value() != "" {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render(fmt.Sprintf("filter: %q  (press / to change)", m.filter.Value())) + "\n")
	}

	// Error line.
	if m.lastErr != "" {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#EF4444")).
			Render("error: "+m.lastErr) + "\n")
	}

	sb.WriteString(m.table.View() + "\n")

	// Help footer.
	help := lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render("q quit  s sort  / filter  r refresh")
	sb.WriteString(help)

	return sb.String()
}

// rebuildTable re-applies filter+sort and updates the embedded table.
func (m *TUIModel) rebuildTable() {
	rows := m.visibleRows()
	t := newPortsTable(rows, m.width)
	t.SetHeight(m.height - 5)
	m.table = t
}

// visibleRows returns filtered and sorted rows.
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

	rows := make([]table.Row, len(infos))
	for i, p := range infos {
		rows[i] = table.Row{
			p.Protocol,
			p.LocalAddr,
			fmt.Sprintf("%d", p.LocalPort),
			fmt.Sprintf("%d", p.PID),
			p.ProcessName,
			p.Status,
		}
	}
	return rows
}

// newPortsTable builds a fresh table model with given rows.
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

// doScan runs Scan in a goroutine and returns a Cmd.
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
