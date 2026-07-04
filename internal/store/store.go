package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RequestID int64

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBearer
	AuthBasic
	AuthAPIKey
)

func (t AuthType) String() string {
	switch t {
	case AuthBearer:
		return "Bearer"
	case AuthBasic:
		return "Basic"
	case AuthAPIKey:
		return "API Key"
	default:
		return "None"
	}
}

type Auth struct {
	Type     AuthType `json:"type"`
	Token    string   `json:"token,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	KeyName  string   `json:"keyName,omitempty"`
	KeyValue string   `json:"keyValue,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Editor struct {
	Body            string   `json:"body,omitempty"`
	ReqHeaders      []Header `json:"reqHeaders,omitempty"`
	QueryParameters []Header `json:"queryParameters,omitempty"`
	Auth            Auth     `json:"auth"`
}

// Execution is a single past run of a request.
type Execution struct {
	Timestamp time.Time     `json:"timestamp"`
	Status    int           `json:"status,omitempty"`
	Headers   []Header      `json:"headers,omitempty"`
	Cookies   []Header      `json:"cookies,omitempty"`
	Body      string        `json:"body,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
	Size      int64         `json:"size,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// maxHistory bounds how many past executions are kept per request, so
// repeated test runs don't grow the collection file without limit.
const maxHistory = 50

// Request is one saved API request definition plus its run history.
//
// EditorTab/ResultTab/HistoryIndex are UI navigation state, not really
// "data" - they're kept here (tagged json:"-") rather than in a second
// map keyed by ID, so the UI only ever needs one lookup per request.
type Request struct {
	ID     RequestID `json:"id"`
	Title  string    `json:"title"`
	URI    string    `json:"uri,omitempty"`
	Method string    `json:"method"`
	Editor Editor    `json:"editor"`

	History []Execution `json:"history,omitempty"`

	EditorTab    int `json:"-"`
	ResultTab    int `json:"-"`
	HistoryIndex int `json:"-"`
}

// CurrentExecution returns the run currently selected via HistoryIndex,
// clamped into range, or the zero value if there's no history yet.
func (r *Request) CurrentExecution() Execution {
	if len(r.History) == 0 {
		return Execution{}
	}
	idx := r.HistoryIndex
	if idx < 0 || idx >= len(r.History) {
		idx = len(r.History) - 1
	}
	return r.History[idx]
}

// Store holds the request collection. Structural changes (create, delete,
// list, append-execution) are mutex-guarded; in-place field edits on a
// *Request returned by Get/List happen directly, since callers only ever
// mutate from Bubble Tea's single Update goroutine.
type Store struct {
	mu        sync.Mutex
	nextID    RequestID
	order     []RequestID
	byID      map[RequestID]*Request
	onboarded bool
}

// Onboarded reports whether the user has already been shown the app's
// first-run help screen.
func (s *Store) Onboarded() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.onboarded
}

// SetOnboarded marks the first-run help screen as seen, so it won't
// auto-expand again once this is persisted.
func (s *Store) SetOnboarded(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onboarded = v
}

func New() *Store {
	return &Store{byID: make(map[RequestID]*Request)}
}

func (s *Store) CreateRequest(title, method string) RequestID {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := s.nextID
	s.byID[id] = &Request{ID: id, Title: title, Method: method, HistoryIndex: -1}
	s.order = append(s.order, id)
	return id
}

func (s *Store) Delete(id RequestID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.byID, id)
	for i, oid := range s.order {
		if oid == id {
			s.order = append(s.order[:i], s.order[i+1:]...)
			break
		}
	}
}

func (s *Store) Get(id RequestID) (*Request, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.byID[id]
	return r, ok
}

// List returns requests in display order. The returned pointers alias the
// store's internal state, so callers can edit fields in place.
func (s *Store) List() []*Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Request, len(s.order))
	for i, id := range s.order {
		out[i] = s.byID[id]
	}
	return out
}

func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.order)
}

// AppendExecution records a finished run and points HistoryIndex at it.
func (s *Store) AppendExecution(id RequestID, exec Execution) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.byID[id]
	if !ok {
		return
	}
	r.History = append(r.History, exec)
	if len(r.History) > maxHistory {
		r.History = r.History[len(r.History)-maxHistory:]
	}
	r.HistoryIndex = len(r.History) - 1
}

type persistedFile struct {
	NextID    RequestID `json:"nextId"`
	Onboarded bool      `json:"onboarded"`
	Requests  []Request `json:"requests"`
}

// Save writes the collection to path as JSON, creating parent directories
// as needed.
func (s *Store) Save(path string) error {
	s.mu.Lock()
	pf := persistedFile{NextID: s.nextID, Onboarded: s.onboarded}
	for _, id := range s.order {
		pf.Requests = append(pf.Requests, *s.byID[id])
	}
	s.mu.Unlock()

	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a collection previously written by Save. Callers should treat
// a non-nil error as "start from an empty store" rather than fatal.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf persistedFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	s := New()
	s.nextID = pf.NextID
	s.onboarded = pf.Onboarded
	for _, r := range pf.Requests {
		req := r
		req.HistoryIndex = len(req.History) - 1
		s.byID[req.ID] = &req
		s.order = append(s.order, req.ID)
	}
	return s, nil
}

// DefaultPath returns the collection file location under the user's
// standard config directory.
func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "trueapitester", "collection.json"), nil
}
