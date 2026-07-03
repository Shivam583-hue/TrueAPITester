package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	purple = lipgloss.Color("135")

	FocusedCellStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("255")).
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("55")).
				Bold(true).
				Align(lipgloss.Center, lipgloss.Center)

	CellStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple).
			Foreground(purple).
			Align(lipgloss.Center, lipgloss.Center)
	MonthTitleStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 2)

	YearStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(purple)

	HelpBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			Align(lipgloss.Center)

	PanelStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple)

	PanelTitleStyle = lipgloss.NewStyle().
			Background(purple).
			Foreground(lipgloss.Color("255")).
			Bold(true).
			Padding(0, 1)

	TaskStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			PaddingLeft(1)

	TaskDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			PaddingLeft(1)

	EmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			PaddingLeft(1)

	TooSmallStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Align(lipgloss.Center, lipgloss.Center)
)
