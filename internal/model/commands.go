package model

import (
	"time"

	httpclient "github.com/Shivam583-hue/TrueAPITester/internal/httpClient"
	"github.com/Shivam583-hue/TrueAPITester/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

type responseMsg struct {
	id        store.RequestID
	timestamp time.Time
	resp      httpclient.Response
}

type responseErrMsg struct {
	id        store.RequestID
	timestamp time.Time
	err       error
}

func toKVs(headers []store.Header) []httpclient.KV {
	kvs := make([]httpclient.KV, 0, len(headers))
	for _, h := range headers {
		kvs = append(kvs, httpclient.KV{Key: h.Key, Value: h.Value})
	}
	return kvs
}

func (m *Model) sendRequestCmd() tea.Cmd {
	r := m.activeRequest()
	id := r.ID
	a := r.Editor.Auth
	req := httpclient.Request{
		Method:  r.Method,
		URL:     r.URI,
		Body:    r.Editor.Body,
		Headers: toKVs(r.Editor.ReqHeaders),
		Query:   toKVs(r.Editor.QueryParameters),
		Auth: httpclient.Auth{
			Type:     httpclient.AuthType(a.Type),
			Token:    a.Token,
			Username: a.Username,
			Password: a.Password,
			KeyName:  a.KeyName,
			KeyValue: a.KeyValue,
		},
	}
	return func() tea.Msg {
		resp, err := httpclient.Send(req)
		if err != nil {
			return responseErrMsg{id: id, timestamp: time.Now(), err: err}
		}
		return responseMsg{id: id, timestamp: time.Now(), resp: resp}
	}
}
