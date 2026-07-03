package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minWidth  = 80
	minHeight = 24
)

const (
	FocusSidebar Focus = iota
	FocusUri
	FocusEditor
	FocusPreview
)

func New() *Model {
	return &Model{}
}

func (f Focus) String() string {
	switch f {
	case FocusSidebar:
		return "Sidebar"
	case FocusUri:
		return "Uri"
	case FocusEditor:
		return "Editor"
	case FocusPreview:
		return "Preview"
	default:
		return "Unknown"
	}
}

func (f Focus) Next() Focus {
	return (f + 1) % (FocusPreview + 1)
}

func (f Focus) Prev() Focus {
	return (f - 1 + FocusPreview + 1) % (FocusPreview + 1)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.focused == FocusUri {
			switch msg.Type {
			case tea.KeyRunes:
				m.uri += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.uri) > 0 {
					runes := []rune(m.uri)
					m.uri = string(runes[:len(runes)-1])
				}
			}
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "right":
			m.focused = m.focused.Next()
		case "left":
			m.focused = m.focused.Prev()

		}
	}

	if !m.loaded {
		m.loaded = true
	}
	return m, nil
}

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
	previewWidth := mainWidth - editorWidth
	mainHeight := m.height - uriHeight

	Sidebar := m.renderSidebar(sidebarWidth, m.focused == FocusSidebar, m.height)
	Uri := m.renderUri(m.uri, mainWidth, m.focused == FocusUri)
	Editor := m.renderEditor(editorWidth, mainHeight, m.focused == FocusEditor, m.response, m.editorTab)
	Preview := m.renderPreview(m.response, m.previewTab, previewWidth, mainHeight, m.focused == FocusPreview) // render
	EditorandPreviewContent := lipgloss.JoinHorizontal(lipgloss.Left, Editor, Preview)

	UriAndContent := lipgloss.JoinVertical(lipgloss.Top, Uri, EditorandPreviewContent)

	return lipgloss.JoinHorizontal(lipgloss.Left, Sidebar, UriAndContent)
}
