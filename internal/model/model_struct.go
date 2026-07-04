package model

import (
	"github.com/Shivam583-hue/TrueAPITester/internal/store"
	"github.com/charmbracelet/bubbles/help"
)

type Focus int

type authField struct {
	label string
	value *string
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

	store          *store.Store
	collectionPath string
	requestCursor  int
	namingRequest  bool
	nameInput      string

	kvCursor  int
	kvEditing bool
	kvOnValue bool

	authCursor  int
	authEditing bool

	editorScroll int
	resultScroll int
}
