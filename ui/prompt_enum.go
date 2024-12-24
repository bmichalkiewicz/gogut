package ui

type PromptMode int

const (
	ExecPromptMode PromptMode = iota
	ChatPromptMode
	ConfigPromptMode
	DefaultPromptMode
)

func (pm PromptMode) String() string {
	switch pm {
	case ExecPromptMode:
		return "exec"
	case ChatPromptMode:
		return "chat"
	case ConfigPromptMode:
		return "config"
	default:
		return "default"
	}
}

func GetPromptModeFromString(s string) PromptMode {
	switch s {
	case "exec":
		return ExecPromptMode
	case "chat":
		return ChatPromptMode
	case "config":
		return ConfigPromptMode
	default:
		return DefaultPromptMode
	}
}

type RunMode int

const (
	CliMode RunMode = iota
	ReplMode
)

func (m RunMode) String() string {
	switch m {
	case CliMode:
		return "cli"
	case ReplMode:
		return "repl"
	}
	return ""
}
