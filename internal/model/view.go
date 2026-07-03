package model

import "github.com/charmbracelet/lipgloss"

func (m *Model) View() string {
	if v := m.SpecialView(); v != "" {
		return v
	}
	const (
		sidebarWidth = 30
		uriHeight    = 3
	)

	mainWidth := m.width - sidebarWidth
	editorWidth := mainWidth / 2
	resultWidth := mainWidth - editorWidth
	mainHeight := m.height - uriHeight

	Sidebar := m.renderSidebar(sidebarWidth, m.focused == FocusSidebar, m.height)
	Uri := m.renderUri(m.uri, mainWidth, m.focused == FocusUri)
	Editor := m.renderEditor(editorWidth, mainHeight, m.focused == FocusEditor, m.response, m.editorTab)
	Result := m.renderResult(m.response, m.resultTab, resultWidth, mainHeight, m.focused == FocusResult)
	EditorandResultContent := lipgloss.JoinHorizontal(lipgloss.Left, Editor, Result)

	UriAndContent := lipgloss.JoinVertical(lipgloss.Top, Uri, EditorandResultContent)

	return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, UriAndContent)
}
