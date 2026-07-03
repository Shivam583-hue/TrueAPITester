package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
)

func (m Model) renderEditor(width int, height int, focused bool, resp Response, activeTab int) string {
	tabs := styles.RenderTabs([]string{"Body", "Headers", "Query", "Auth"}, activeTab)

	var body string

	status := fmt.Sprintf("Status: %s",
		styles.StatusCodeStyle(resp.Status).Render(strconv.Itoa(resp.Status)))
	// status := fmt.Sprintf("Status: %s  Time: %dms  Size: %s",
	// 	styles.StatusCodeStyle(resp.Status).Render(strconv.Itoa(resp.Status)))
	// resp.TimeMs, resp.SizeStr)

	content := tabs + "\n\n" + body + "\n\n" + status

	return styles.TitledPane("Editor", content, width, height, focused)
}

func (m Model) renderResult(resp Response, activeTab int, width, height int, focused bool) string {
	tabs := styles.RenderTabs([]string{"Pretty", "Raw", "Headers", "Cookies"}, activeTab)

	var body string
	switch activeTab {
	case 0: // Pretty
		body = highlightJSON(resp.Body) // tokenizer + JSON*Style
	case 1: // Raw
		body = resp.Body
	case 2: // Headers
		var rows []string
		for _, h := range resp.Headers {
			rows = append(rows, styles.RenderHeaderRow(h.Key, h.Value))
		}
		body = strings.Join(rows, "\n\n")
	case 3: // Cookies
		// same shape as Headers for now
	}

	status := fmt.Sprintf("Status: %s",
		styles.StatusCodeStyle(resp.Status).Render(strconv.Itoa(resp.Status)))
	// status := fmt.Sprintf("Status: %s  Time: %dms  Size: %s",
	// 	styles.StatusCodeStyle(resp.Status).Render(strconv.Itoa(resp.Status)))
	// resp.TimeMs, resp.SizeStr)

	content := tabs + "\n\n" + body + "\n\n" + status

	return styles.TitledPane("Result", content, width, height, focused)
}

func (m Model) renderUri(uri string, width int, focused bool) string {
	text := uri
	style := styles.URLInputStyle

	if text == "" {
		text = "Enter a URL..."
		style = styles.PlaceholderStyle
	}

	return styles.TitledPane("Uri", style.Render(text), width, 3, focused)
}

func (m Model) renderSidebar(width int, focused bool, height int) string {
	text := "haha"
	style := styles.URLInputStyle

	return styles.TitledPane("Requests", style.Render(text), width, height, focused)
}
