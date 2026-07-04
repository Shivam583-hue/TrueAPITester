package model

import (
	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

const (
	sidebarWidth = 30
	methodWidth  = 12
	uriHeight    = 3
)

func (m *Model) View() string {
	if v := m.SpecialView(); v != "" {
		return v
	}

	mainWidth := m.width - sidebarWidth
	editorWidth := mainWidth / 2
	resultWidth := mainWidth - editorWidth
	mainHeight := m.height - uriHeight

	Sidebar := m.renderSidebar(sidebarWidth, m.focused == FocusSidebar, m.height)

	if m.store.Len() == 0 {
		placeholder := styles.PlaceholderStyle.
			Width(mainWidth).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Press n to create a new request")
		return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, placeholder)
	}

	r := m.activeRequest()
	Method := m.renderMethod(r.Method, methodWidth, m.focused == FocusMethod)
	Uri := m.renderUri(r.URI, mainWidth-methodWidth, m.focused == FocusUri)
	UriRow := lipgloss.JoinHorizontal(lipgloss.Left, Method, Uri)
	Editor := m.renderEditor(editorWidth, mainHeight, m.focused == FocusEditor, r.EditorTab)
	Result := m.renderResult(r, resultWidth, mainHeight, m.focused == FocusResult)
	EditorandResultContent := lipgloss.JoinHorizontal(lipgloss.Left, Editor, Result)

	UriAndContent := lipgloss.JoinVertical(lipgloss.Top, UriRow, EditorandResultContent)

	return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, UriAndContent)
}
