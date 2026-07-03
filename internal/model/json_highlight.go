package model

import (
	"encoding/json"
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
