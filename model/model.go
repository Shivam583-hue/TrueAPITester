package model

import (
	"fmt"

	"github.com/Shivam583-hue/TrueAPITester/styles"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	minWidth  = 80
	minHeight = 24
)

type Model struct {
	quitting bool
	help     help.Model
	loaded   bool

	width  int
	height int

	choices  []string
	cursor   int
	selected map[int]struct{}
}

func New() *Model {
	return &Model{
		choices: []string{
			"Buy carrots",
			"Buy celery",
			"Buy kohlrabi",
		},
		selected: make(map[int]struct{}),
	}
}

func (m *Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if !m.loaded {
		return "Loading..."
	}
	if m.quitting {
		return ""
	}

	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf(
			"Terminal too small: %dx%d\nMinimum required: %dx%d\nPlease resize your terminal.",
			m.width, m.height, minWidth, minHeight,
		)
		return styles.TooSmallStyle.
			Width(m.width).
			Height(m.height).
			Render(msg)
	}
	return ""
}
