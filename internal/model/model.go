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

// activeKVList returns the key-value list being edited on the current
// editor tab: request headers (tab 1) or query parameters (tab 2).
func (m *Model) activeKVList() *[]Header {
	if m.activeRequest().editorTab == 2 {
		return &m.activeRequest().editor.queryParameters
	}
	return &m.activeRequest().editor.reqHeaders
}

// authFields returns the editable fields for the active request's auth type.
func (m *Model) authFields() []authField {
	a := &m.activeRequest().editor.auth
	switch a.authtype {
	case AuthBearer:
		return []authField{{"Token", &a.token}}
	case AuthBasic:
		return []authField{{"Username", &a.username}, {"Password", &a.password}}
	case AuthAPIKey:
		return []authField{{"Key Name", &a.keyName}, {"Key Value", &a.keyValue}}
	}
	return nil
}

// editString applies rune/space/backspace input to a string field.
func editString(s *string, msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		*s += string(msg.Runes)
	case tea.KeySpace:
		*s += " "
	case tea.KeyBackspace:
		if len(*s) > 0 {
			runes := []rune(*s)
			*s = string(runes[:len(runes)-1])
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case responseMsg:
		if msg.index < len(m.requests) {
			r := &m.requests[msg.index]
			r.response = Response{
				Status:   msg.resp.Status,
				Body:     msg.resp.Body,
				Duration: msg.resp.Duration,
				Size:     msg.resp.Size,
			}
			for _, h := range msg.resp.Headers {
				r.response.Headers = append(r.response.Headers, Header{Key: h.Key, Value: h.Value})
			}
			for _, c := range msg.resp.Cookies {
				r.response.Cookies = append(r.response.Cookies, Cookie{Name: c.Key, Value: c.Value})
			}
		}
		return m, nil

	case responseErrMsg:
		if msg.index < len(m.requests) {
			m.requests[msg.index].response = Response{Error: msg.err.Error()}
		}
		return m, nil

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

		// Editor pane
		if m.focused == FocusEditor {
			switch m.activeRequest().editorTab {
			case 0: // Body
				if msg.Type == tea.KeyEnter {
					m.activeRequest().editor.body += "\n"
				} else {
					editString(&m.activeRequest().editor.body, msg)
				}

			case 1, 2: // Headers / Query
				list := m.activeKVList()
				if m.kvCursor >= len(*list) {
					m.kvCursor = 0
					if len(*list) > 0 {
						m.kvCursor = len(*list) - 1
					}
				}
				if m.kvEditing {
					row := &(*list)[m.kvCursor]
					target := &row.Key
					if m.kvOnValue {
						target = &row.Value
					}
					switch msg.Type {
					case tea.KeyEnter:
						if m.kvOnValue {
							m.kvEditing = false
							m.kvOnValue = false
						} else {
							m.kvOnValue = true
						}
					case tea.KeyEsc:
						m.kvEditing = false
						m.kvOnValue = false
					case tea.KeyCtrlC:
						m.quitting = true
						return m, tea.Quit
					default:
						editString(target, msg)
					}
					return m, nil
				}
				switch msg.String() {
				case "n":
					*list = append(*list, Header{})
					m.kvCursor = len(*list) - 1
					m.kvEditing = true
					m.kvOnValue = false
				case "d":
					if len(*list) > 0 {
						*list = append((*list)[:m.kvCursor], (*list)[m.kvCursor+1:]...)
						if m.kvCursor >= len(*list) && m.kvCursor > 0 {
							m.kvCursor--
						}
					}
				case "enter":
					if len(*list) > 0 {
						m.kvEditing = true
						m.kvOnValue = false
					}
				case "up":
					if m.kvCursor > 0 {
						m.kvCursor--
					}
				case "down":
					if m.kvCursor < len(*list)-1 {
						m.kvCursor++
					}
				}

			case 3: // Auth
				fields := m.authFields()
				if m.authCursor >= len(fields) {
					m.authCursor = 0
				}
				if m.authEditing {
					switch msg.Type {
					case tea.KeyEnter, tea.KeyEsc:
						m.authEditing = false
					case tea.KeyCtrlC:
						m.quitting = true
						return m, tea.Quit
					default:
						editString(fields[m.authCursor].value, msg)
					}
					return m, nil
				}
				switch msg.String() {
				case "t":
					a := &m.activeRequest().editor.auth
					a.authtype = (a.authtype + 1) % (AuthAPIKey + 1)
					m.authCursor = 0
				case "enter":
					if len(fields) > 0 {
						m.authEditing = true
					}
				case "up":
					if m.authCursor > 0 {
						m.authCursor--
					}
				case "down":
					if m.authCursor < len(fields)-1 {
						m.authCursor++
					}
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
		case "ctrl+s":
			if m.activeRequest() != nil {
				return m, m.sendRequestCmd()
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
