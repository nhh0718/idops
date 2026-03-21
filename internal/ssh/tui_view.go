package ssh

import (
	"fmt"
	"strings"

	lipglossv1 "github.com/charmbracelet/lipgloss"
)

// View dispatches rendering based on current mode.
func (m TUIModel) View() string {
	switch m.mode {
	case modeAdd, modeEdit:
		return m.renderForm()
	case modeConfirmDelete:
		sel, _ := m.list.SelectedItem().(hostItem)
		h := sel.host
		prompt := fmt.Sprintf("\n  ⚠️  Xóa host '%s' (%s@%s:%s)? [y/N] ",
			h.Name, orDash(h.User), orDash(h.Hostname), orDefault(h.Port, "22"))
		return lipglossv1.NewStyle().Bold(true).Foreground(lipglossv1.Color("#EF4444")).Render(prompt)
	default:
		return m.renderList()
	}
}

// renderList renders the main host list view with status bar and help footer.
func (m TUIModel) renderList() string {
	var sb strings.Builder

	// Empty state
	if len(m.list.Items()) == 0 {
		sb.WriteString(lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render("\n  Không tìm thấy SSH host nào. Nhấn 'a' để thêm mới.\n"))
	} else {
		sb.WriteString(m.list.View())
	}

	// Status bar
	if m.status != "" {
		sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#10B981")).Render("  "+m.status))
	}

	// Host count
	total := len(m.list.Items())
	if total > 0 {
		sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
			Render(fmt.Sprintf("  %d host(s)", total)))
	}

	// Help footer
	sb.WriteString("\n" + lipglossv1.NewStyle().Foreground(lipglossv1.Color("#6B7280")).
		Render("  a:Thêm  e:Sửa  d:Xóa  c:Kết nối  t:Test  q:Thoát"))
	return sb.String()
}
