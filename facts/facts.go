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

	return strings.Trim(strings.Trim(dist, "\n"), "\"")
}

func GetShell() string {
	shell, err := run.RunCommand("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	split := strings.Split(strings.Trim(strings.Trim(shell, "\n"), "\""), "/")

	return split[len(split)-1]
}

func GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

func GetUsername() string {
	name, err := run.RunCommand("echo", os.Getenv("USER"))
	if err != nil {
		return ""
	}

	return strings.Trim(name, "\n")
}

func GetEditor() string {
	name, err := run.RunCommand("echo", os.Getenv("EDITOR"))
	if err != nil {
		return "nano"
	}

	return strings.Trim(name, "\n")
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
