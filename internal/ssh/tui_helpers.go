package ssh

import "github.com/charmbracelet/bubbles/list"

// hostsToItems converts SSHHost slice to list.Item slice with optional test results.
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

// currentHosts extracts SSHHost slice from current list items.
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

// reloadItemsWithResults reloads config and merges with existing test results.
func (m TUIModel) reloadItemsWithResults() []list.Item {
	hosts, err := LoadConfig(m.configPath)
	if err != nil {
		return m.list.Items()
	}
	return hostsToItems(hosts, m.testResults)
}

// orDash returns "-" if s is empty.
func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// orDefault returns def if s is empty.
func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
