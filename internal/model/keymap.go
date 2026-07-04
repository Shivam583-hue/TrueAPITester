package model

import (
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var (
	keyNew       = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new request"))
	keyDelete    = key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete"))
	keyUpDown    = key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑/↓", "navigate"))
	keyFocusNext = key.NewBinding(key.WithKeys("right"), key.WithHelp("→", "next pane"))
	keyFocusPrev = key.NewBinding(key.WithKeys("left"), key.WithHelp("←", "prev pane"))
	keyTab       = key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch tab"))
	keyEnter     = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm/edit"))
	keyEsc       = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	keySend      = key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "send request"))
	keySave      = key.NewBinding(key.WithKeys("ctrl+w"), key.WithHelp("ctrl+w", "save collection"))
	keyHistory   = key.NewBinding(key.WithKeys("[", "]"), key.WithHelp("[/]", "prev/next run"))
	keyMethod    = key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "change method"))
	keyAuthType  = key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "change auth type"))
	keyScroll    = key.NewBinding(key.WithKeys("up", "down", "pgup", "pgdown"), key.WithHelp("↑/↓/pgup/pgdn", "scroll"))
	keyHelp      = key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
	keyQuit      = key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit"))
)

type keyMap struct {
	focused     Focus
	editorTab   int
	hasRequests bool
}

func (k keyMap) ShortHelp() []key.Binding {
	switch {
	case !k.hasRequests:
		return []key.Binding{keyNew, keyHelp, keyQuit}
	case k.focused == FocusSidebar:
		return []key.Binding{keyNew, keyDelete, keyUpDown, keyFocusNext, keyHelp, keyQuit}
	case k.focused == FocusMethod:
		return []key.Binding{keyMethod, keyEnter, keyFocusNext, keyFocusPrev, keyHelp}
	case k.focused == FocusUri:
		return []key.Binding{keyEnter, keyFocusNext, keyFocusPrev, keySend, keyHelp}
	case k.focused == FocusEditor && k.editorTab == 3: // Auth
		return []key.Binding{keyAuthType, keyEnter, keyTab, keySend, keyHelp}
	case k.focused == FocusEditor && (k.editorTab == 1 || k.editorTab == 2): // Headers/Query
		return []key.Binding{keyNew, keyDelete, keyEnter, keyTab, keySend, keyHelp}
	case k.focused == FocusEditor: // Body
		return []key.Binding{keyTab, keySend, keyFocusNext, keyHelp}
	case k.focused == FocusResult:
		return []key.Binding{keyHistory, keyScroll, keyTab, keyFocusPrev, keyHelp}
	default:
		return []key.Binding{keySend, keyHelp, keyQuit}
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{keyFocusNext, keyFocusPrev, keyTab, keyUpDown},
		{keyNew, keyDelete, keyEnter, keyEsc},
		{keyMethod, keyAuthType, keyHistory, keyScroll},
		{keySend, keySave, keyHelp, keyQuit},
	}
}

func newHelpModel() help.Model {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(styles.Peach).Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(styles.Subtext0)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(styles.Overlay0)
	h.Styles.Ellipsis = lipgloss.NewStyle().Foreground(styles.Overlay0)
	h.Styles.FullKey = h.Styles.ShortKey
	h.Styles.FullDesc = h.Styles.ShortDesc
	h.Styles.FullSeparator = h.Styles.ShortSeparator
	return h
}

func (m *Model) helpKeyMap() keyMap {
	var editorTab int
	if r := m.activeRequest(); r != nil {
		editorTab = r.EditorTab
	}
	return keyMap{focused: m.focused, editorTab: editorTab, hasRequests: m.store.Len() > 0}
}

// helpBar renders the fully styled help bar (border, background, width) for
// the current terminal width and UI state.
func (m *Model) helpBar() string {
	content := m.help.View(m.helpKeyMap())
	return styles.HelpBarStyle.Width(m.width).Render(content)
}

func (m *Model) helpHeight() int {
	return strings.Count(m.helpBar(), "\n") + 1
}

func (m *Model) isTypingText() bool {
	if m.focused == FocusUri {
		return true
	}
	if m.focused == FocusEditor {
		if r := m.activeRequest(); r != nil && r.EditorTab == 0 {
			return true
		}
	}
	return false
}
