package model

const (
	FocusSidebar Focus = iota
	FocusMethod
	FocusUri
	FocusEditor
	FocusResult
)

func (f Focus) String() string {
	switch f {
	case FocusSidebar:
		return "Sidebar"
	case FocusMethod:
		return "Method"
	case FocusUri:
		return "Uri"
	case FocusEditor:
		return "Editor"
	case FocusResult:
		return "Result"
	default:
		return "Unknown"
	}
}

func (f Focus) Next() Focus {
	return (f + 1) % (FocusResult + 1)
}

func (f Focus) Prev() Focus {
	return (f - 1 + FocusResult + 1) % (FocusResult + 1)
}
