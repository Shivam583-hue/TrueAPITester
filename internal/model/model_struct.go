package model

import (
	"github.com/Shivam583-hue/TrueAPITester/internal/store"
	"github.com/charmbracelet/bubbles/help"
)

type Focus int

// authField pairs a display label with a pointer to the string it edits.
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

	// key-value editing state (Headers / Query tabs)
	kvCursor  int
	kvEditing bool
	kvOnValue bool

	// auth editing state
	authCursor  int
	authEditing bool

	// scroll offsets for the editor Body tab and the result pane
	editorScroll int
	resultScroll int
}
