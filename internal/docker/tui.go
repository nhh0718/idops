package docker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	dockerclient "github.com/docker/docker/client"
)

type tickMsg time.Time
type containersMsg []ContainerInfo
type errMsg error
type actionResultMsg struct{ text string; isErr bool }
type clearStatusMsg struct{}

var (
	stateGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	stateRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	stateYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	statusOK    = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	statusErr   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
)

// dashModel is the Bubble Tea model for the Docker dashboard.
type dashModel struct {
	cli         *dockerclient.Client
	table       table.Model
	filter      textinput.Model
	containers  []ContainerInfo
	filtered    []ContainerInfo
	filtering   bool
	confirm     bool
	confirmKey  string // "x", "r", "d"
	confirmName string
	confirmID   string
	statusMsg   string
	statusIsErr bool
	err         error
	width       int
}

func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return clearStatusMsg{} })
}

// NewDashboard creates the initial dashboard model.
func NewDashboard(cli *dockerclient.Client) dashModel {
	cols := []table.Column{
		{Title: "ID", Width: 13},
		{Title: "Name", Width: 22},
		{Title: "Image", Width: 28},
		{Title: "State", Width: 10},
		{Title: "Status", Width: 20},
		{Title: "CPU%", Width: 7},
		{Title: "Mem%", Width: 7},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(15),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#6B7280")).BorderBottom(true).Bold(true)
	s.Selected = s.Selected.Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7C3AED")).Bold(false)
	t.SetStyles(s)

	fi := textinput.New()
	fi.Placeholder = "filter..."
	fi.CharLimit = 40

	return dashModel{cli: cli, table: t, filter: fi, filtered: []ContainerInfo{}}
}

func (m dashModel) Init() tea.Cmd {
	return tea.Batch(fetchContainers(m.cli), tickEvery())
}

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func fetchContainers(cli *dockerclient.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		list, err := ListContainers(ctx, cli)
		if err != nil {
			return errMsg(err)
		}
		return containersMsg(list)
	}
}

func (m dashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.table.SetHeight(msg.Height - 7)

	case tickMsg:
		// Always reschedule tick even if fetch fails
		return m, tea.Batch(fetchContainers(m.cli), tickEvery())

	case containersMsg:
		m.containers = msg
		m.err = nil
		m.applyFilter()

	case errMsg:
		m.err = msg
		// Keep ticking — do not stop refresh on error

	case actionResultMsg:
		m.statusMsg = msg.text
		m.statusIsErr = msg.isErr
		return m, tea.Batch(clearStatusAfter(3*time.Second), fetchContainers(m.cli))

	case clearStatusMsg:
		m.statusMsg = ""
		m.statusIsErr = false

	case tea.KeyMsg:
		if m.confirm {
			return m.handleConfirmKey(msg)
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

func (m *dashModel) applyFilter() {
	query := strings.ToLower(m.filter.Value())
	if query == "" {
		m.filtered = m.containers
	} else {
		out := make([]ContainerInfo, 0)
		for _, c := range m.containers {
			if strings.Contains(strings.ToLower(c.Name), query) ||
				strings.Contains(strings.ToLower(c.Image), query) {
				out = append(out, c)
			}
		}
		m.filtered = out
	}
	rows := make([]table.Row, 0, len(m.filtered))
	for _, c := range m.filtered {
		cpu, mem := "-", "-"
		if c.Stats != nil {
			cpu = fmt.Sprintf("%.1f", c.Stats.CPUPercent)
			mem = fmt.Sprintf("%.1f", c.Stats.MemPercent)
		}
		state := colorState(c.State)
		rows = append(rows, table.Row{c.ID, c.Name, c.Image, state, c.Status, cpu, mem})
	}
	m.table.SetRows(rows)
}

