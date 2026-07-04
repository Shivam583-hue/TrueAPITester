package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 80
	minHeight = 24
)

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func New() *Model {
	return &Model{method: "GET"}
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
			switch msg.String() {
			case "enter":
				m.focused = m.focused.Next()
			}
		}
		if m.focused == FocusMethod {
			switch msg.String() {
			case "m":
				for i, method := range httpMethods {
					if method == m.method {
						m.method = httpMethods[(i+1)%len(httpMethods)]
						break
					}
				}
			case "enter":
				m.focused = m.focused.Next()
			}
		}
		switch msg.String() {
		case "tab":
			switch m.focused {
			case FocusEditor:
				m.editorTab = (m.editorTab + 1) % 4
			case FocusResult:
				m.resultTab = (m.resultTab + 1) % 4
			}
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
