package main

import (
	"fmt"
	"os"

	"github.com/Shivam583-hue/TrueAPITester/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// mux := http.NewServeMux()
	// mux.Handle("/", &homeHandler{})

	m := model.New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	// fmt.Println("Server is running")
	// if err := http.ListenAndServe(":8080", mux); err != nil {
	// 	fmt.Println("An error occured")
	// }
}

// type homeHandler struct{}
//
// func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello World"))
// }
