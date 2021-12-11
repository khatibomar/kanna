package browser

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
