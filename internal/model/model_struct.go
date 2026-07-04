package model

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
)

type Focus int

type Response struct {
	Status   int
	Headers  []Header
	Cookies  []Cookie
	Body     string
	Duration time.Duration
	Size     int64
}

type Header struct {
	Key   string
	Value string
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
