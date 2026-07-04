package model

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
)

type Focus int

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBearer
	AuthBasic
	AuthAPIKey
)

type Auth struct {
	authtype AuthType

	// Bearere
	token string

	// Basic
	username string
	password string

	// API
	keyName  string
	keyValue string
}

type Response struct {
	Status   int
	Headers  []Header
	Cookies  []Cookie
	Body     string
	Duration time.Duration
	Size     int64

	editor Editor
}

type Header struct {
	Key   string
	Value string
}

type Editor struct {
	body            string
	reqHeaders      []Header
	queryParameters []Header
	auth            Auth
}

type Cookie struct {
	Name  string
	Value string
}

type Requests struct {
	title     string
	uri       string
	method    string
	response  Response
	editorTab int
	resultTab int

	editor Editor
}

type Model struct {
	quitting bool
	help     help.Model
	loaded   bool

	width  int
	height int

	focused Focus

	RequestsWidth  int
	RequestsHeight int

	UriWidth  int
	UriHeight int

	EditorWidth  int
	EditorHeight int

	ResultWidth  int
	ResultHeight int

	requests      []Requests
	requestCursor int
	namingRequest bool
	nameInput     string
}
