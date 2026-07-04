package model

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 80
	minHeight = 24
)

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func New() *Model {
	return &Model{}
}

// activeRequest returns the request under the sidebar cursor, or nil when
// there are no requests. All pane reads/writes should go through this.
func (m *Model) activeRequest() *Requests {
	if len(m.requests) == 0 {
		return nil
	}
	return &m.requests[m.requestCursor]
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
		if m.namingRequest {
			switch msg.Type {
			case tea.KeyRunes:
				m.nameInput += string(msg.Runes)
			case tea.KeySpace:
				m.nameInput += " "
			case tea.KeyBackspace:
				if len(m.nameInput) > 0 {
					runes := []rune(m.nameInput)
					m.nameInput = string(runes[:len(runes)-1])
				}
			case tea.KeyEnter:
				if name := strings.TrimSpace(m.nameInput); name != "" {
					m.requests = append(m.requests, Requests{title: name, method: "GET"})
					m.requestCursor = len(m.requests) - 1
				}
				m.nameInput = ""
				m.namingRequest = false
			case tea.KeyEsc:
				m.nameInput = ""
				m.namingRequest = false
			case tea.KeyCtrlC:
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

		// uri
		if m.focused == FocusUri {
			switch msg.Type {
			case tea.KeyRunes:
				m.activeRequest().uri += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.activeRequest().uri) > 0 {
					runes := []rune(m.activeRequest().uri)
					m.activeRequest().uri = string(runes[:len(runes)-1])
				}
			}
			switch msg.String() {
			case "enter":
				m.focused = m.focused.Next()
			}
		}

		// Editor tab, body
		if m.focused == FocusEditor {
			if m.activeRequest().editorTab == 0 {
				switch msg.Type {
				case tea.KeyRunes:
					m.activeRequest().editor.body += string(msg.Runes)
				case tea.KeySpace:
					m.activeRequest().editor.body += string(" ")
				case tea.KeyBackspace:
					if len(m.activeRequest().editor.body) > 0 {
						runes := []rune(m.activeRequest().editor.body)
						m.activeRequest().editor.body = string(runes[:len(runes)-1])
					}
				case tea.KeyEnter:
					m.activeRequest().editor.body += "\n"
				}
			}
		}

		if m.focused == FocusMethod {
			switch msg.String() {
			case "m":
				for i, method := range httpMethods {
					if method == m.activeRequest().method {
						m.activeRequest().method = httpMethods[(i+1)%len(httpMethods)]
						break
					}
				}
			case "enter":
				m.focused = m.focused.Next()
			}
		}
		if m.focused == FocusSidebar {
			switch msg.String() {
			case "n":
				m.namingRequest = true
				m.nameInput = ""
			case "d":
				if len(m.requests) > 0 {
					m.requests = append(m.requests[:m.requestCursor], m.requests[m.requestCursor+1:]...)
					if m.requestCursor >= len(m.requests) && m.requestCursor > 0 {
						m.requestCursor--
					}
				}
			case "up":
				if m.requestCursor > 0 {
					m.requestCursor--
				}
			case "down":
				if m.requestCursor < len(m.requests)-1 {
					m.requestCursor++
				}
			}
		}
		switch msg.String() {
		case "tab":
			switch m.focused {
			case FocusEditor:
				m.activeRequest().editorTab = (m.activeRequest().editorTab + 1) % 4
			case FocusResult:
				m.activeRequest().resultTab = (m.activeRequest().resultTab + 1) % 4
			}
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "right":
			if len(m.requests) > 0 {
				m.focused = m.focused.Next()
			}
		case "left":
			if len(m.requests) > 0 {
				m.focused = m.focused.Prev()
			}
		}
	}

	if !m.loaded {
		m.loaded = true
	}
	return m, nil
}
