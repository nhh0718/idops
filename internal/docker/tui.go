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

// dashModel is the Bubble Tea model for the Docker dashboard.
type dashModel struct {
	cli        *dockerclient.Client
	table      table.Model
	filter     textinput.Model
	containers []ContainerInfo
	filtered   []ContainerInfo
	filtering  bool
	err        error
	width      int
}

var (
	stateGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981"))
	stateRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	stateYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
)

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

	return dashModel{cli: cli, table: t, filter: fi}
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
		m.table.SetHeight(msg.Height - 6)

	case tickMsg:
		return m, tea.Batch(fetchContainers(m.cli), tickEvery())

	case containersMsg:
		m.containers = msg
		m.applyFilter()

	case errMsg:
		m.err = msg

	case tea.KeyMsg:
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
		out := m.filtered[:0]
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
		rows = append(rows, table.Row{c.ID, c.Name, c.Image, c.State, c.Status, cpu, mem})
	}
	m.table.SetRows(rows)
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
			go func() { //nolint
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_ = StartContainer(ctx, m.cli, selected.ID)
			}()
		}
	case "x":
		if selected != nil {
			go func() { //nolint
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_ = StopContainer(ctx, m.cli, selected.ID)
			}()
		}
	case "r":
		if selected != nil {
			go func() { //nolint
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_ = RestartContainer(ctx, m.cli, selected.ID)
			}()
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

func (m dashModel) View() string {
	if m.err != nil {
		return stateRed.Render("Error: "+m.err.Error()) + "\n"
	}

	stateColored := func(state string) string {
		switch strings.ToLower(state) {
		case "running":
			return stateGreen.Render(state)
		case "exited":
			return stateRed.Render(state)
		case "paused":
			return stateYellow.Render(state)
		default:
			return state
		}
	}

	// Re-render rows with colored state column for display
	_ = stateColored

	header := titleStyle.Render("Docker Dashboard")
	var filterLine string
	if m.filtering {
		filterLine = "Filter: " + m.filter.View()
	} else if m.filter.Value() != "" {
		filterLine = helpStyle.Render("Filter: " + m.filter.Value())
	}

	help := helpStyle.Render("s:start  x:stop  r:restart  l:logs  /:filter  q:quit")

	parts := []string{header}
	if filterLine != "" {
		parts = append(parts, filterLine)
	}
	parts = append(parts, m.table.View(), help)
	return strings.Join(parts, "\n")
}
