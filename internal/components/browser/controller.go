package browser

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, nil
		}
		switch msg.String() {
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "M", "m":
			if !m.isAcceptingRequests {
				break
			}
			m.isAcceptingRequests = false
			go func(m *model) {
				animes, err := getAnimeData(m.cfg, m.offset+m.limit, m.limit)
				m.err = err
				m.newContentChan <- animes
			}(&m)
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case errMsg:
		m.err = msg
		return m, nil
	}
	var cmd tea.Cmd
	select {
	case animes := <-m.newContentChan:
		for _, anime := range animes {
			m.list.InsertItem(len(m.list.Items()), anime)
		}
		m.offset += m.limit
		m.isAcceptingRequests = true
	default:
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
