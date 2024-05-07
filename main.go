package main

import (
	"log"

	"github.com/bmichalkiewicz/gogut/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	input, err := ui.NewUIInput()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tea.NewProgram(ui.NewUI(input)).Run(); err != nil {
		log.Fatal(err)
	}
}
