package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/store"
	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderEditor(width int, height int, focused bool, activeTab int) string {
	var body string
	switch activeTab {
	case 0: // Body
		body = m.activeRequest().Editor.Body
		if focused {
			body += "█"
		}
		if body == "" {
			body = "Enter body text..."
		}

	case 1, 2: // Headers / Query
		list := m.activeRequest().Editor.ReqHeaders
		itemName := "header"
		if activeTab == 2 {
			list = m.activeRequest().Editor.QueryParameters
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
		a := m.activeRequest().Editor.Auth
		rows := []string{styles.ListItemStyle.Render("Type: "+a.Type.String()) +
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

	bodyH := height - 4 // borders + tab row + blank line
	offset := 0
	switch activeTab {
	case 0:
		offset = m.editorScroll
	case 1, 2:
		// keep the key-value cursor visible
		if m.kvCursor >= bodyH {
			offset = m.kvCursor - bodyH + 1
		}
	}
	body = styles.ScrollView(body, width-2, bodyH, offset)

	tabs := styles.RenderTabs([]string{"Body", "Headers", "Query", "Auth"}, activeTab)

	content := tabs + "\n\n" + body

	return styles.TitledPane("Editor", content, width, height, focused)
}

func resultTabContent(resp store.Execution, tab int) string {
	if resp.Error != "" {
		return lipgloss.NewStyle().Foreground(styles.Red).Render("Error: " + resp.Error)
	}
	switch tab {
	case 0: // Pretty
		return highlightJSON(resp.Body)
	case 1: // Raw
		return resp.Body
	case 2: // Headers
		var rows []string
		for _, h := range resp.Headers {
			rows = append(rows, styles.RenderHeaderRow(h.Key, h.Value))
		}
		return strings.Join(rows, "\n\n")
	case 3: // Cookies
		var rows []string
		for _, c := range resp.Cookies {
			rows = append(rows, styles.RenderHeaderRow(c.Key, c.Value))
		}
		return strings.Join(rows, "\n\n")
	}
	return ""
}

func (m Model) renderResult(r *store.Request, width, height int, focused bool) string {
	activeTab := r.ResultTab
	resp := r.CurrentExecution()

	tabs := styles.RenderTabs([]string{"Pretty", "Raw", "Headers", "Cookies"}, activeTab)

	body := resultTabContent(resp, activeTab)

	bodyH := height - 4 // borders + tab row + blank line
	var status string
	if len(r.History) > 0 {
		status = fmt.Sprintf("Run %d/%d  %s", r.HistoryIndex+1, len(r.History), resp.Timestamp.Format("15:04:05"))
		if resp.Error == "" {
			status += fmt.Sprintf("   Status: %s  Time: %dms  Size: %dB",
				styles.StatusCodeStyle(resp.Status).Render(strconv.Itoa(resp.Status)),
				resp.Duration.Milliseconds(), resp.Size)
		}
		bodyH -= 2 // blank line + status bar
	}
	body = styles.ScrollView(body, width-2, bodyH, m.resultScroll)

	content := tabs + "\n\n" + body
	if status != "" {
		content += "\n\n" + status
	}

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
	for i, req := range m.store.List() {
		if i == m.requestCursor {
			rows = append(rows, styles.ListItemSelectedStyle.Render(req.Title))
		} else {
			rows = append(rows, styles.ListItemStyle.Render(req.Title))
		}
	}
	if m.namingRequest {
		rows = append(rows, styles.URLInputStyle.Render("> "+m.nameInput+"█"))
	} else if len(rows) == 0 {
		rows = append(rows, styles.PlaceholderStyle.Render("Press n to add a request"))
	}
	return styles.TitledPane("Requests", strings.Join(rows, "\n"), width, height, focused)
}
