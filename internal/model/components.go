package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Shivam583-hue/TrueAPITester/internal/styles"
)

func highlightJSON(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}

	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber()

	var sb strings.Builder
	if err := writeJSONValue(dec, &sb, 0); err != nil {
		return raw
	}
	return sb.String()
}

func writeJSONContainer(dec *json.Decoder, sb *strings.Builder, indent int, openChar, closeChar rune, isObject bool) error {
	sb.WriteString(styles.JSONPunctStyle.Render(string(openChar)))

	first := true
	for dec.More() {
		if !first {
			sb.WriteString(styles.JSONPunctStyle.Render(","))
		}
		first = false
		sb.WriteString("\n" + strings.Repeat("  ", indent+1))

		if isObject {
			keyTok, err := dec.Token()
			if err != nil {
				return err
			}
			key, _ := keyTok.(string)
			sb.WriteString(styles.JSONKeyStyle.Render(strconv.Quote(key)))
			sb.WriteString(styles.JSONPunctStyle.Render(": "))
		}

		if err := writeJSONValue(dec, sb, indent+1); err != nil {
			return err
		}
	}

	if _, err := dec.Token(); err != nil {
		return err
	}
	if !first {
		sb.WriteString("\n" + strings.Repeat("  ", indent))
	}
	sb.WriteString(styles.JSONPunctStyle.Render(string(closeChar)))
	return nil
}

func writeJSONValue(dec *json.Decoder, sb *strings.Builder, indent int) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}

	switch t := tok.(type) {
	case json.Delim:
		switch t {
		case '{':
			return writeJSONContainer(dec, sb, indent, '{', '}', true)
		case '[':
			return writeJSONContainer(dec, sb, indent, '[', ']', false)
		}
	case string:
		sb.WriteString(styles.JSONStringStyle.Render(strconv.Quote(t)))
	case json.Number:
		sb.WriteString(styles.JSONNumberStyle.Render(t.String()))
	case bool:
		sb.WriteString(styles.JSONBoolStyle.Render(strconv.FormatBool(t)))
	case nil:
		sb.WriteString(styles.JSONNullStyle.Render("null"))
	}
	return nil
}

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

func (m Model) renderPreview(resp Response, activeTab int, width, height int, focused bool) string {
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

	return styles.TitledPane("Preview", content, width, height, focused)
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
