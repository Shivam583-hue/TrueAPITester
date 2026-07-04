package model

import (
	"os"
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/store"
	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 80
	minHeight = 24
)

var httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func New() *Model {
	m := &Model{store: store.New(), help: newHelpModel()}

	path, err := store.DefaultPath()
	if err == nil {
		m.collectionPath = path
		if loaded, loadErr := store.Load(path); loadErr == nil {
			m.store = loaded
		} else if !os.IsNotExist(loadErr) {
			// Don't let a corrupt collection file silently vanish once we
			// save a fresh empty store over it.
			_ = os.Rename(path, path+".bak")
		}
	}

	// Auto-expand the full help screen exactly once: the first time this
	// user's collection file didn't already mark them as onboarded.
	m.help.ShowAll = !m.store.Onboarded()
	m.store.SetOnboarded(true)
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// -2 for the help bar's own left/right padding, so its wrapping
		// matches the width it's actually rendered at.
		if w := msg.Width - 2; w > 0 {
			m.help.Width = w
		} else {
			m.help.Width = 0
		}

	case responseMsg:
		exec := store.Execution{
			Timestamp: msg.timestamp,
			Status:    msg.resp.Status,
			Body:      msg.resp.Body,
			Duration:  msg.resp.Duration,
			Size:      msg.resp.Size,
		}
		for _, h := range msg.resp.Headers {
			exec.Headers = append(exec.Headers, store.Header{Key: h.Key, Value: h.Value})
		}
		for _, c := range msg.resp.Cookies {
			exec.Cookies = append(exec.Cookies, store.Header{Key: c.Key, Value: c.Value})
		}
		m.store.AppendExecution(msg.id, exec)
		if cur := m.activeRequest(); cur != nil && cur.ID == msg.id {
			m.resultScroll = 0
		}
		return m, nil

	case responseErrMsg:
		m.store.AppendExecution(msg.id, store.Execution{Timestamp: msg.timestamp, Error: msg.err.Error()})
		if cur := m.activeRequest(); cur != nil && cur.ID == msg.id {
			m.resultScroll = 0
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
					m.store.CreateRequest(name, "GET")
					m.requestCursor = m.store.Len() - 1
					m.editorScroll, m.resultScroll = 0, 0
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
				m.activeRequest().URI += string(msg.Runes)
			case tea.KeyBackspace:
				if len(m.activeRequest().URI) > 0 {
					runes := []rune(m.activeRequest().URI)
					m.activeRequest().URI = string(runes[:len(runes)-1])
				}
			}
			switch msg.String() {
			case "enter":
				m.focused = m.focused.Next()
			}
		}

		// Editor pane
		if m.focused == FocusEditor {
			switch m.activeRequest().EditorTab {
			case 0: // Body
				switch msg.String() {
				case "up":
					if m.editorScroll > 0 {
						m.editorScroll--
					}
				case "down":
					if m.editorScroll < m.editorMaxScroll() {
						m.editorScroll++
					}
				case "pgup":
					m.editorScroll -= m.paneBodyHeight()
					if m.editorScroll < 0 {
						m.editorScroll = 0
					}
				case "pgdown":
					m.editorScroll += m.paneBodyHeight()
					if max := m.editorMaxScroll(); m.editorScroll > max {
						m.editorScroll = max
					}
				default:
					switch msg.Type {
					case tea.KeyRunes, tea.KeySpace, tea.KeyBackspace, tea.KeyEnter:
						if msg.Type == tea.KeyEnter {
							m.activeRequest().Editor.Body += "\n"
						} else {
							editString(&m.activeRequest().Editor.Body, msg)
						}
						// keep the cursor (end of body) visible while typing
						m.editorScroll = m.editorMaxScroll()
					}
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
					*list = append(*list, store.Header{})
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
					a := &m.activeRequest().Editor.Auth
					a.Type = (a.Type + 1) % (store.AuthAPIKey + 1)
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
					if method == m.activeRequest().Method {
						m.activeRequest().Method = httpMethods[(i+1)%len(httpMethods)]
						break
					}
				}
			case "enter":
				m.focused = m.focused.Next()
			}
		}
		if m.focused == FocusResult {
			switch msg.String() {
			case "up":
				if m.resultScroll > 0 {
					m.resultScroll--
				}
			case "down":
				if m.resultScroll < m.resultMaxScroll() {
					m.resultScroll++
				}
			case "pgup":
				m.resultScroll -= m.paneBodyHeight()
				if m.resultScroll < 0 {
					m.resultScroll = 0
				}
			case "pgdown":
				m.resultScroll += m.paneBodyHeight()
				if max := m.resultMaxScroll(); m.resultScroll > max {
					m.resultScroll = max
				}
			case "[":
				if r := m.activeRequest(); r != nil && r.HistoryIndex > 0 {
					r.HistoryIndex--
					m.resultScroll = 0
				}
			case "]":
				if r := m.activeRequest(); r != nil && r.HistoryIndex < len(r.History)-1 {
					r.HistoryIndex++
					m.resultScroll = 0
				}
			}
		}

		if m.focused == FocusSidebar {
			switch msg.String() {
			case "n":
				m.namingRequest = true
				m.nameInput = ""
			case "d":
				if r := m.activeRequest(); r != nil {
					m.store.Delete(r.ID)
					if n := m.store.Len(); m.requestCursor >= n && m.requestCursor > 0 {
						m.requestCursor--
					}
					m.editorScroll, m.resultScroll = 0, 0
				}
			case "up":
				if m.requestCursor > 0 {
					m.requestCursor--
					m.editorScroll, m.resultScroll = 0, 0
				}
			case "down":
				if m.requestCursor < m.store.Len()-1 {
					m.requestCursor++
					m.editorScroll, m.resultScroll = 0, 0
				}
			}
		}
		switch msg.String() {
		case "tab":
			switch m.focused {
			case FocusEditor:
				m.activeRequest().EditorTab = (m.activeRequest().EditorTab + 1) % 4
				m.editorScroll = 0
			case FocusResult:
				m.activeRequest().ResultTab = (m.activeRequest().ResultTab + 1) % 4
				m.resultScroll = 0
			}
		case "ctrl+s":
			if m.activeRequest() != nil {
				return m, m.sendRequestCmd()
			}
		case "ctrl+w":
			m.SaveCollection()
		case "?":
			if !m.isTypingText() {
				m.help.ShowAll = !m.help.ShowAll
			}
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "right":
			if m.store.Len() > 0 {
				m.focused = m.focused.Next()
			}
		case "left":
			if m.store.Len() > 0 {
				m.focused = m.focused.Prev()
			}
		}
	}

	if !m.loaded {
		m.loaded = true
	}
	return m, nil
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

// SaveCollection persists the current collection to disk, if a path is
// available. Errors are non-fatal: this is a best-effort autosave/manual
// save, not a transaction the rest of the app depends on.
func (m *Model) SaveCollection() {
	if m.collectionPath == "" || m.store == nil {
		return
	}
	_ = m.store.Save(m.collectionPath)
}

// activeRequest returns the request under the sidebar cursor
func (m *Model) activeRequest() *store.Request {
	list := m.store.List()
	if len(list) == 0 {
		return nil
	}
	if m.requestCursor >= len(list) {
		m.requestCursor = len(list) - 1
	}
	return list[m.requestCursor]
}

// activeKVList returns the key-value list being edited on the current
// editor tab: request headers (tab 1) or query parameters (tab 2).
func (m *Model) activeKVList() *[]store.Header {
	if m.activeRequest().EditorTab == 2 {
		return &m.activeRequest().Editor.QueryParameters
	}
	return &m.activeRequest().Editor.ReqHeaders
}

// authFields returns the editable fields for the active request's auth type.
func (m *Model) authFields() []authField {
	a := &m.activeRequest().Editor.Auth
	switch a.Type {
	case store.AuthBearer:
		return []authField{{"Token", &a.Token}}
	case store.AuthBasic:
		return []authField{{"Username", &a.Username}, {"Password", &a.Password}}
	case store.AuthAPIKey:
		return []authField{{"Key Name", &a.KeyName}, {"Key Value", &a.KeyValue}}
	}
	return nil
}

// paneBodyHeight returns the lines available for scrollable body content in
// the editor/result panes: total minus uri row, help bar, borders, tab row,
// blank line.
func (m *Model) paneBodyHeight() int {
	h := m.height - uriHeight - m.helpHeight() - 4
	if h < 1 {
		h = 1
	}
	return h
}

func (m *Model) editorMaxScroll() int {
	body := m.activeRequest().Editor.Body + "█"
	return styles.MaxScroll(body, (m.width-sidebarWidth)/2-2, m.paneBodyHeight())
}

func (m *Model) resultMaxScroll() int {
	r := m.activeRequest()
	exec := r.CurrentExecution()
	body := resultTabContent(exec, r.ResultTab)
	h := m.paneBodyHeight()
	if len(r.History) > 0 {
		h -= 2 // status bar + blank line
	}
	mainWidth := m.width - sidebarWidth
	return styles.MaxScroll(body, mainWidth-mainWidth/2-2, h)
}
