package config

import "fmt"

// ConfigFileNotfoundError error when config file hasn't been found
type ConfigFileNotfoundError struct{}

func (e ConfigFileNotfoundError) Error() string {
	return fmt.Sprintln("config file hasn't been found")
}
