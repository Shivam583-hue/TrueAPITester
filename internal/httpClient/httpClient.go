package httpclient

import (
	"context"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

type KV struct {
	Key   string
	Value string
}

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBearer
	AuthBasic
	AuthAPIKey
)

type Auth struct {
	Type     AuthType
	Token    string
	Username string
	Password string
	KeyName  string
	KeyValue string
}

type Request struct {
	Method  string
	URL     string
	Body    string
	Headers []KV
	Query   []KV
	Auth    Auth
}

type Response struct {
	Status   int
	Headers  []KV
	Cookies  []KV
	Body     string
	Duration time.Duration
	Size     int64
}

func Send(r Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var bodyReader io.Reader
	if r.Body != "" {
		bodyReader = strings.NewReader(r.Body)
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, bodyReader)
	if err != nil {
		return Response{}, err
	}

	q := req.URL.Query()
	for _, p := range r.Query {
		if p.Key != "" {
			q.Add(p.Key, p.Value)
		}
	}
	req.URL.RawQuery = q.Encode()

	for _, h := range r.Headers {
		if h.Key != "" {
			req.Header.Set(h.Key, h.Value)
		}
	}

	switch r.Auth.Type {
	case AuthBearer:
		req.Header.Set("Authorization", "Bearer "+r.Auth.Token)
	case AuthBasic:
		req.SetBasicAuth(r.Auth.Username, r.Auth.Password)
	case AuthAPIKey:
		if r.Auth.KeyName != "" {
			req.Header.Set(r.Auth.KeyName, r.Auth.KeyValue)
		}
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	out := Response{
		Status:   resp.StatusCode,
		Body:     string(body),
		Duration: time.Since(start),
		Size:     int64(len(body)),
	}
	for key, vals := range resp.Header {
		for _, v := range vals {
			out.Headers = append(out.Headers, KV{Key: key, Value: v})
		}
	}
	// map iteration order is random; keep the display stable
	sort.Slice(out.Headers, func(i, j int) bool { return out.Headers[i].Key < out.Headers[j].Key })
	for _, c := range resp.Cookies() {
		out.Cookies = append(out.Cookies, KV{Key: c.Name, Value: c.Value})
	}
	return out, nil
}
