package ui

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var loadingMessages = []string{
	"Reticulating splines...",
	"Counting backwards from Infinity",
	"Spinning the hamster…",
	"Generating witty dialog…",
	"Swapping time and space…",
	"At least you're not on hold…",
	"Awakening the minions…",
	"Summoning internet fairies…",
	"Looking for the 'Any' key…",
	"Building a fort…",
	"Grabbing extra pixels…",
	"Aligning covfefe levels…",
	"Entertaining electrons…",
	"Distilling the essence of pure flavor…",
	"Trying to sort by duck…",
	"Convincing AI not to turn evil…",
	"Hold on, I saw this in a cartoon once…",
	"Proving P=NP…",
	"Entangling superstrings…",
	"Untangling the Internet…",
	"Downloading more RAM…",
	"Deciding what message to display next…",
	"Looking for sense of humor…",
}

type Spinner struct {
	message string
	spinner spinner.Model
}

func NewSpinner() *Spinner {
	spin := spinner.New()
	spin.Spinner = spinner.MiniDot

	return &Spinner{
		message: loadingMessages[rand.Intn(len(loadingMessages))],
		spinner: spin,
	}
}

func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	var updateCmd tea.Cmd
	s.spinner, updateCmd = s.spinner.Update(msg)

	return s, updateCmd
}

func (s *Spinner) View() string {
	return fmt.Sprintf(
		"\n  %s %s...",
		s.spinner.View(),
		s.spinner.Style.Render(s.message),
	)
}

func (s *Spinner) Tick() tea.Msg {
	return s.spinner.Tick()
}
