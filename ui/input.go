package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
)

type UIInput struct {
	runMode    RunMode
	promptMode PromptMode
	args       string
	pipe       string
}

func getPipeData() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting stat: %s", err)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var buf []byte
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			buf = append(buf, scanner.Bytes()...)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		pipe := strings.TrimSpace(string(buf))

		return pipe, nil
	}

	return "", nil
}

func NewUIInput() (*UIInput, error) {
	flags := flag.NewFlagSet("GoGut", flag.ExitOnError)

	exec := flags.Bool("exec", false, "Run with exec mode")
	chat := flags.Bool("prompt", false, "Run with chat mode")
	debug := flags.Bool("debug", false, "Debug mode")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("error with flags parsing: %s", err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	args := flags.Args()

	runMode := ReplMode
	if len(args) > 0 {
		runMode = CliMode
	}

	var promptMode PromptMode

	switch {
	case !*exec && *chat:
		promptMode = ChatPromptMode
	case *exec && !*chat:
		promptMode = ExecPromptMode
	default:
		promptMode = DefaultPromptMode
	}

	pipe, err := getPipeData()
	if err != nil {
		return nil, fmt.Errorf("error getting data from pipe: %s", err)
	}

	return &UIInput{
		runMode:    runMode,
		promptMode: promptMode,
		args:       strings.Join(args, " "),
		pipe:       pipe,
	}, nil
}

func (i *UIInput) GetRunMode() RunMode {
	return i.runMode
}

func (i *UIInput) GetPromptMode() PromptMode {
	return i.promptMode
}

func (i *UIInput) GetArgs() string {
	return i.args
}

func (i *UIInput) GetPipe() string {
	return i.pipe
}
