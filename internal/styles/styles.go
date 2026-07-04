package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	Base   = lipgloss.Color("#1e1e2e")
	Mantle = lipgloss.Color("#181825")
	Crust  = lipgloss.Color("#11111b")

	Text     = lipgloss.Color("#cdd6f4")
	Subtext1 = lipgloss.Color("#bac2de")
	Subtext0 = lipgloss.Color("#a6adc8")
	Overlay2 = lipgloss.Color("#9399b2")
	Overlay1 = lipgloss.Color("#7f849c")
	Overlay0 = lipgloss.Color("#6c7086")
	Surface2 = lipgloss.Color("#585b70")
	Surface1 = lipgloss.Color("#45475a")
	Surface0 = lipgloss.Color("#313244")

	Rosewater = lipgloss.Color("#f5e0dc")
	Flamingo  = lipgloss.Color("#f2cdcd")
	Pink      = lipgloss.Color("#f5c2e7")
	Mauve     = lipgloss.Color("#cba6f7")
	Red       = lipgloss.Color("#f38ba8")
	Maroon    = lipgloss.Color("#eba0ac")
	Peach     = lipgloss.Color("#fab387")
	Yellow    = lipgloss.Color("#f9e2af")
	Green     = lipgloss.Color("#a6e3a1")
	Teal      = lipgloss.Color("#94e2d5")
	Sky       = lipgloss.Color("#89dceb")
	Sapphire  = lipgloss.Color("#74c7ec")
	Blue      = lipgloss.Color("#89b4fa")
	Lavender  = lipgloss.Color("#b4befe")

	SelectionBg = lipgloss.Color("#4a3040")
)

var AppStyle = lipgloss.NewStyle().Background(Base).Foreground(Text)

var MethodColors = map[string]lipgloss.Color{
	"GET":     Blue,
	"POST":    Peach,
	"PUT":     Yellow,
	"PATCH":   Teal,
	"DELETE":  Red,
	"HEAD":    Overlay1,
	"OPTIONS": Overlay1,
}

func MethodStyle(method string) lipgloss.Style {
	c, ok := MethodColors[strings.ToUpper(method)]
	if !ok {
		c = Overlay1
	}
	return lipgloss.NewStyle().Bold(true).Foreground(c)
}

var (
	ListItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	ListItemSelectedStyle = lipgloss.NewStyle().
				Foreground(Rosewater).
				Background(SelectionBg).
				Bold(true).
				Padding(0, 1)
)

var (
	TabActiveStyle = lipgloss.NewStyle().
			Foreground(Text).
			Bold(true).
			Underline(true)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(Overlay0)

	tabSepStyle = lipgloss.NewStyle().Foreground(Surface2)
)

func RenderTabs(labels []string, activeIndex int) string {
	parts := make([]string, len(labels))
	for i, l := range labels {
		if i == activeIndex {
			parts[i] = TabActiveStyle.Render(l)
		} else {
			parts[i] = TabInactiveStyle.Render(l)
		}
	}
	return strings.Join(parts, tabSepStyle.Render(" │ "))
}

var (
	borderInactive = Surface1
	titleInactive  = Subtext0
	borderActive   = Mauve
	titleActive    = Peach
)

func PaneStyle(focused bool) lipgloss.Style {
	c := borderInactive
	if focused {
		c = borderActive
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(c).
		Background(Base).
		Padding(0, 1)
}

func TitledPane(title, body string, width, height int, focused bool) string {
	bc, tc := borderInactive, titleInactive
	if focused {
		bc, tc = borderActive, titleActive
	}
	border := lipgloss.NewStyle().Foreground(bc)
	titleStyle := lipgloss.NewStyle().Foreground(tc).Bold(true)

	b := lipgloss.RoundedBorder()
	innerW := width - 2
	if innerW < 0 {
		innerW = 0
	}

	titleTxt := " " + title
	fill := innerW - lipgloss.Width(titleTxt)
	if fill < 0 {
		fill = 0
	}

	top := border.Render(b.TopLeft) +
		titleStyle.Render(titleTxt) +
		border.Render(strings.Repeat(b.Top, fill)+b.TopRight)

	bottom := border.Render(b.BottomLeft + strings.Repeat(b.Bottom, innerW) + b.BottomRight)

	content := lipgloss.NewStyle().
		Width(innerW).
		Height(height - 2).
		Render(body)

	rows := strings.Split(content, "\n")
	// hard clip: lipgloss Height pads short content but never cuts long
	// content, so overflowing rows would break the pane layout
	if maxRows := height - 2; maxRows >= 0 && len(rows) > maxRows {
		rows = rows[:maxRows]
	}
	for i, r := range rows {
		rows[i] = border.Render(b.Left) + r + border.Render(b.Right)
	}

	return top + "\n" + strings.Join(rows, "\n") + "\n" + bottom
}

// ScrollView wraps body to width, then returns the window of height lines
// starting at offset (clamped to the content).
func ScrollView(body string, width, height, offset int) string {
	if width < 1 || height < 1 {
		return ""
	}
	lines := strings.Split(lipgloss.NewStyle().Width(width).Render(body), "\n")
	max := len(lines) - height
	if max < 0 {
		max = 0
	}
	if offset > max {
		offset = max
	}
	if offset < 0 {
		offset = 0
	}
	end := offset + height
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[offset:end], "\n")
}

// MaxScroll returns the largest useful scroll offset for body wrapped to
// width in a window of height lines.
func MaxScroll(body string, width, height int) int {
	if width < 1 {
		return 0
	}
	n := strings.Count(lipgloss.NewStyle().Width(width).Render(body), "\n") + 1
	if max := n - height; max > 0 {
		return max
	}
	return 0
}

var ModeColors = map[string]lipgloss.Color{
	"NORMAL":  Blue,
	"INSERT":  Green,
	"VISUAL":  Mauve,
	"COMMAND": Peach,
}

func StatusModeBadge(mode string) string {
	c, ok := ModeColors[strings.ToUpper(mode)]
	if !ok {
		c = Overlay1
	}
	return lipgloss.NewStyle().
		Background(c).
		Foreground(Crust).
		Bold(true).
		Padding(0, 1).
		Render(strings.ToUpper(mode))
}

var StatusSegmentStyle = lipgloss.NewStyle().
	Background(Surface0).
	Foreground(Peach).
	Padding(0, 1)

func StatusCodeStyle(code int) lipgloss.Style {
	switch {
	case code >= 200 && code < 300:
		return lipgloss.NewStyle().Foreground(Green).Bold(true)
	case code >= 300 && code < 400:
		return lipgloss.NewStyle().Foreground(Yellow).Bold(true)
	case code >= 400 && code < 500:
		return lipgloss.NewStyle().Foreground(Peach).Bold(true)
	case code >= 500:
		return lipgloss.NewStyle().Foreground(Red).Bold(true)
	default:
		return lipgloss.NewStyle().Foreground(Overlay1)
	}
}

var (
	HeaderKeyStyle   = lipgloss.NewStyle().Foreground(Peach).Bold(true)
	HeaderValueStyle = lipgloss.NewStyle().Foreground(Text)
)

func RenderHeaderRow(key, value string) string {
	return HeaderKeyStyle.Render(key) + "\n" + HeaderValueStyle.Render(value)
}

var URLInputStyle = lipgloss.NewStyle().Foreground(Sky)

var PlaceholderStyle = lipgloss.NewStyle().Foreground(Overlay1).Italic(true)

var (
	JSONKeyStyle    = lipgloss.NewStyle().Foreground(Blue)
	JSONStringStyle = lipgloss.NewStyle().Foreground(Green)
	JSONNumberStyle = lipgloss.NewStyle().Foreground(Peach)
	JSONBoolStyle   = lipgloss.NewStyle().Foreground(Mauve)
	JSONNullStyle   = lipgloss.NewStyle().Foreground(Overlay1).Italic(true)
	JSONPunctStyle  = lipgloss.NewStyle().Foreground(Overlay0)
)

var TooSmallStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("196")).
	Bold(true).
	Align(lipgloss.Center, lipgloss.Center)
