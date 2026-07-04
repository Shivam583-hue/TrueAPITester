package model

import (
	"strings"

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

	helpBar := m.helpBar()
	helpLines := strings.Count(helpBar, "\n") + 1

	mainWidth := m.width - sidebarWidth
	editorWidth := mainWidth / 2
	resultWidth := mainWidth - editorWidth
	bodyHeight := m.height - helpLines
	mainHeight := bodyHeight - uriHeight

	Sidebar := m.renderSidebar(sidebarWidth, m.focused == FocusSidebar, bodyHeight)

	var main string
	if m.store.Len() == 0 {
		placeholder := styles.PlaceholderStyle.
			Width(mainWidth).
			Height(bodyHeight).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Press n to create a new request")
		main = lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, placeholder)
	} else {
		r := m.activeRequest()
		Method := m.renderMethod(r.Method, methodWidth, m.focused == FocusMethod)
		Uri := m.renderUri(r.URI, mainWidth-methodWidth, m.focused == FocusUri)
		UriRow := lipgloss.JoinHorizontal(lipgloss.Left, Method, Uri)
		Editor := m.renderEditor(editorWidth, mainHeight, m.focused == FocusEditor, r.EditorTab)
		Result := m.renderResult(r, resultWidth, mainHeight, m.focused == FocusResult)
		EditorandResultContent := lipgloss.JoinHorizontal(lipgloss.Left, Editor, Result)

		UriAndContent := lipgloss.JoinVertical(lipgloss.Top, UriRow, EditorandResultContent)
		main = lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, UriAndContent)
	}

	return lipgloss.JoinVertical(lipgloss.Top, main, helpBar)
}
