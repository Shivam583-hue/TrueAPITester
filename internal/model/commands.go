package model

import (
	httpclient "github.com/Shivam583-hue/TrueAPITester/internal/httpClient"
	tea "github.com/charmbracelet/bubbletea"
)

type responseMsg struct {
	index int
	resp  httpclient.Response
}

type responseErrMsg struct {
	index int
	err   error
}

func toKVs(headers []Header) []httpclient.KV {
	kvs := make([]httpclient.KV, 0, len(headers))
	for _, h := range headers {
		kvs = append(kvs, httpclient.KV{Key: h.Key, Value: h.Value})
	}
	return kvs
}

// sendRequestCmd snapshots the active request and returns a command that
// sends it in the background, tagged with the request index so the
// response lands on the right request even if the cursor moves meanwhile.
func (m *Model) sendRequestCmd() tea.Cmd {
	index := m.requestCursor
	r := m.activeRequest()
	a := r.editor.auth
	req := httpclient.Request{
		Method:  r.method,
		URL:     r.uri,
		Body:    r.editor.body,
		Headers: toKVs(r.editor.reqHeaders),
		Query:   toKVs(r.editor.queryParameters),
		Auth: httpclient.Auth{
			// model.AuthType and httpclient.AuthType share the same iota order
			Type:     httpclient.AuthType(a.authtype),
			Token:    a.token,
			Username: a.username,
			Password: a.password,
			KeyName:  a.keyName,
			KeyValue: a.keyValue,
		},
	}
	return func() tea.Msg {
		resp, err := httpclient.Send(req)
		if err != nil {
			return responseErrMsg{index: index, err: err}
		}
		return responseMsg{index: index, resp: resp}
	}
}
