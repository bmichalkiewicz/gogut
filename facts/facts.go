package facts

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bmichalkiewicz/gogut/run"
	"github.com/mitchellh/go-homedir"
)

const applicationName = "GoGut"

type Analysis struct {
	operatingSystem OperatingSystem
	distribution    string
	shell           string
	homeDirectory   string
	username        string
	editor          string
	configFile      string
	configPath      string
}

func (a *Analysis) GetApplicationName() string {
	return applicationName
}

func (a *Analysis) GetOperatingSystem() OperatingSystem {
	return a.operatingSystem
}

func (a *Analysis) GetDistribution() string {
	return a.distribution
}

func (a *Analysis) GetShell() string {
	return a.shell
}

func (a *Analysis) GetHomeDirectory() string {
	return a.homeDirectory
}

func (a *Analysis) GetUsername() string {
	return a.username
}

func (a *Analysis) GetEditor() string {
	return a.editor
}

func (a *Analysis) GetConfigFile() string {
	return a.configFile
}

func (a *Analysis) GetConfigPath() string {
	return a.configPath
}

func Analyse() *Analysis {
	return &Analysis{
		operatingSystem: GetOperatingSystem(),
		distribution:    GetDistribution(),
		shell:           GetShell(),
		homeDirectory:   GetHomeDirectory(),
		username:        GetUsername(),
		editor:          GetEditor(),
		configFile:      GetConfigFile(),
		configPath:      GetConfigPath(),
	}
}

func GetOperatingSystem() OperatingSystem {
	switch runtime.GOOS {
	case "linux":
		return LinuxOperatingSystem
	case "darwin":
		return MacOperatingSystem
	case "windows":
		return WindowsOperatingSystem
	default:
		return UnknownOperatingSystem
	}
}

func GetDistribution() string {
	dist, err := run.RunCommand("lsb_release", "-sd")
	if err != nil {
		return ""
	}
	return strings.Trim(dist, "\"\n")
}

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ""
	}
	parts := strings.Split(shell, "/")
	return parts[len(parts)-1]
}

func GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

func GetUsername() string {
	name := os.Getenv("USER")
	if name == "" {
		return ""
	}

	return name
}

func GetEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "vi"
	}

	return editor
}

func GetConfigFile() string {
	return fmt.Sprintf(
		"%s/config.yaml",
		GetConfigPath(),
	)
}

func GetConfigPath() string {
	return fmt.Sprintf(
		"%s/.%s",
		GetHomeDirectory(),
		strings.ToLower(applicationName),
	)
}
