package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
)

func (m Model) renderEditor(width int, height int, focused bool, resp Response, activeTab int) string {
	var body string
	switch activeTab {
	case 0: // Body
		body = m.activeRequest().editor.body
		if focused {
			body += "█"
		}
		if body == "" {
			body = "Enter body text..."
		}

	case 1, 2: // Headers / Query
		list := m.activeRequest().editor.reqHeaders
		itemName := "header"
		if activeTab == 2 {
			list = m.activeRequest().editor.queryParameters
			itemName = "query param"
		}
		var rows []string
		for i, h := range list {
			line := h.Key + ": " + h.Value
			switch {
			case i == m.kvCursor && m.kvEditing:
				if m.kvOnValue {
					line = h.Key + ": " + h.Value + "█"
				} else {
					line = h.Key + "█: " + h.Value
				}
				rows = append(rows, styles.URLInputStyle.Render(line))
			case i == m.kvCursor:
				rows = append(rows, styles.ListItemSelectedStyle.Render(line))
			default:
				rows = append(rows, styles.ListItemStyle.Render(line))
			}
		}
		if len(rows) == 0 {
			rows = append(rows, styles.PlaceholderStyle.Render("Press n to add a "+itemName))
		}
		body = strings.Join(rows, "\n")

	case 3: // Auth
		a := m.activeRequest().editor.auth
		rows := []string{styles.ListItemStyle.Render("Type: "+a.authtype.String()) +
			styles.PlaceholderStyle.Render("  (t to change)")}
		for i, f := range m.authFields() {
			line := f.label + ": " + *f.value
			switch {
			case i == m.authCursor && m.authEditing:
				rows = append(rows, styles.URLInputStyle.Render(line+"█"))
			case i == m.authCursor:
				rows = append(rows, styles.ListItemSelectedStyle.Render(line))
			default:
				rows = append(rows, styles.ListItemStyle.Render(line))
			}
		}
		body = strings.Join(rows, "\n")
	}

	tabs := styles.RenderTabs([]string{"Body", "Headers", "Query", "Auth"}, activeTab)

	content := tabs + "\n\n" + body

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

	content := tabs + "\n\n" + body + "\n\n" + status

	return styles.TitledPane("Result", content, width, height, focused)
}

func (m Model) renderMethod(method string, width int, focused bool) string {
	text := styles.MethodStyle(method).Render(method)
	return styles.TitledPane("Method", text, width, 3, focused)
}

func (m Model) renderUri(uri string, width int, focused bool) string {
	text := uri
	style := styles.URLInputStyle

	if focused {
		text += "█"
	} else if text == "" {
		text = "Enter a URL..."
		style = styles.PlaceholderStyle
	}

	return styles.TitledPane("Uri", style.Render(text), width, 3, focused)
}

func (m Model) renderSidebar(width int, focused bool, height int) string {
	var rows []string
	for i, req := range m.requests {
		if i == m.requestCursor {
			rows = append(rows, styles.ListItemSelectedStyle.Render(req.title))
		} else {
			rows = append(rows, styles.ListItemStyle.Render(req.title))
		}
	}
	if m.namingRequest {
		rows = append(rows, styles.URLInputStyle.Render("> "+m.nameInput+"█"))
	} else if len(rows) == 0 {
		rows = append(rows, styles.PlaceholderStyle.Render("Press n to add a request"))
	}
	return styles.TitledPane("Requests", strings.Join(rows, "\n"), width, height, focused)
}
