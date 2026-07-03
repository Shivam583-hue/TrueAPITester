package model

import (
	"fmt"

	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
)

func (m Model) SpecialView() string {
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
