package main

import (
	"fmt"
	"os"

	"github.com/Shivam583-hue/TrueAPITester/model"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := model.New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
