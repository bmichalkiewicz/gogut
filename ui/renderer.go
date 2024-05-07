package ui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	execColor    = "#ffa657"
	configColor  = "#ffffff"
	chatColor    = "#66b3ff"
	scriptColor  = "#bca0dd"
	helpColor    = "#aaaaaa"
	errorColor   = "#cc3333"
	warningColor = "#ffcc00"
	successColor = "#46b946"
)

type Renderer struct {
	contentRenderer *glamour.TermRenderer
	successRenderer lipgloss.Style
	warningRenderer lipgloss.Style
	errorRenderer   lipgloss.Style
	helpRenderer    lipgloss.Style
}

func NewRenderer(options ...glamour.TermRendererOption) *Renderer {
	contentRenderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return nil
	}

	successRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(successColor))
	warningRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(warningColor))
	errorRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(errorColor))
	helpRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(helpColor)).Italic(true)

	return &Renderer{
		contentRenderer: contentRenderer,
		successRenderer: successRenderer,
		warningRenderer: warningRenderer,
		errorRenderer:   errorRenderer,
		helpRenderer:    helpRenderer,
	}
}

func (r *Renderer) RenderContent(in string) string {
	out, _ := r.contentRenderer.Render(in)

	return out
}

func (r *Renderer) RenderSuccess(in string) string {
	return r.successRenderer.Render(in)
}

func (r *Renderer) RenderWarning(in string) string {
	return r.warningRenderer.Render(in)
}

func (r *Renderer) RenderError(in string) string {
	return r.errorRenderer.Render(in)
}

func (r *Renderer) RenderHelp(in string) string {
	return r.helpRenderer.Render(in)
}

func (r *Renderer) RenderConfigMessage() string {
	var sb strings.Builder

	sb.WriteString("Welcome! 👋  \n\n")
	sb.WriteString("I cannot find a configuration file, please enter an `API key` \n\n")
	sb.WriteString("If you want to use OpenAI compatibility service/local models, write here `API key` or just press enter to pass it...")

	return sb.String()
}

func (r *Renderer) RenderHelpMessage() string {
	var sb strings.Builder

	sb.WriteString("**Help**\n")
	sb.WriteString("- `↑`/`↓` : navigate in history\n")
	sb.WriteString("- `tab`   : switch between `🚀 exec` and `💬 chat` prompt modes\n")
	sb.WriteString("- `ctrl+h`: show help\n")
	sb.WriteString("- `ctrl+s`: edit settings\n")
	sb.WriteString("- `ctrl+r`: clear terminal and reset discussion history\n")
	sb.WriteString("- `ctrl+l`: clear terminal but keep discussion history\n")
	sb.WriteString("- `ctrl+c`: exit or interrupt command execution\n")

	return sb.String()
}
