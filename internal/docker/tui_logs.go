package docker

import (
	"context"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
)

// logModel is the Bubble Tea model for viewing container logs.
type logModel struct {
	viewport    viewport.Model
	containerID string
	title       string
	ready       bool
	err         error
}

// NewLogViewer creates a log viewer model for the given container.
func NewLogViewer(ctx context.Context, cli *dockerclient.Client, containerID, name string) (*logModel, error) {
	logs, err := fetchLogs(ctx, cli, containerID)
	if err != nil {
		return nil, err
	}

	vp := viewport.New(120, 30)
	vp.SetContent(logs)

	return &logModel{
		viewport:    vp,
		containerID: containerID,
		title:       name,
		ready:       true,
	}, nil
}

func fetchLogs(ctx context.Context, cli *dockerclient.Client, id string) (string, error) {
	out, err := cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	})
	if err != nil {
		return "", err
	}
	defer out.Close()

	raw, err := io.ReadAll(out)
	if err != nil {
		return "", err
	}

	// Strip Docker multiplexed stream header bytes (8-byte prefix per frame).
	return stripDockerStreamHeader(raw), nil
}

// stripDockerStreamHeader removes the 8-byte multiplexed stream headers from Docker log output.
func stripDockerStreamHeader(data []byte) string {
	var sb strings.Builder
	i := 0
	for i < len(data) {
		if i+8 > len(data) {
			break
		}
		// Byte 0: stream type (1=stdout, 2=stderr). Bytes 4-7: frame size.
		size := int(data[i+4])<<24 | int(data[i+5])<<16 | int(data[i+6])<<8 | int(data[i+7])
		i += 8
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		sb.Write(data[i:end])
		i = end
	}
	return sb.String()
}

func (m logModel) Init() tea.Cmd { return nil }

func (m logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m logModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7C3AED"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280"))
	title := titleStyle.Render("Logs: " + m.title)
	help := helpStyle.Render("q/Esc: back")
	return title + "  " + help + "\n" + m.viewport.View()
}
