package model

import (
	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if v := m.SpecialView(); v != "" {
		return v
	}
	const (
		sidebarWidth = 30
		methodWidth  = 12
		uriHeight    = 3
	)

	mainWidth := m.width - sidebarWidth
	editorWidth := mainWidth / 2
	resultWidth := mainWidth - editorWidth
	mainHeight := m.height - uriHeight

	Sidebar := m.renderSidebar(sidebarWidth, m.focused == FocusSidebar, m.height)

	if len(m.requests) == 0 {
		placeholder := styles.PlaceholderStyle.
			Width(mainWidth).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Press n to create a new request")
		return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, placeholder)
	}

	Method := m.renderMethod(m.method, methodWidth, m.focused == FocusMethod)
	Uri := m.renderUri(m.uri, mainWidth-methodWidth, m.focused == FocusUri)
	UriRow := lipgloss.JoinHorizontal(lipgloss.Left, Method, Uri)
	Editor := m.renderEditor(editorWidth, mainHeight, m.focused == FocusEditor, m.response, m.editorTab)
	Result := m.renderResult(m.response, m.resultTab, resultWidth, mainHeight, m.focused == FocusResult)
	EditorandResultContent := lipgloss.JoinHorizontal(lipgloss.Left, Editor, Result)

	UriAndContent := lipgloss.JoinVertical(lipgloss.Top, UriRow, EditorandResultContent)

	return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, UriAndContent)
}
