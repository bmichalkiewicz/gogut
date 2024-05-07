package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/bmichalkiewicz/gogut/facts"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	config = koanf.New(".")
	parser = yaml.Parser()
)

type Config struct {
	common AIConfig
	user   UserConfig
	facts  *facts.Analysis
}

func (c *Config) GetAIConfig() AIConfig {
	return c.common
}

func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

func (c *Config) GetSystemConfig() *facts.Analysis {
	return c.facts
}

func NewConfig(configFile string) (*Config, error) {
	facts := facts.Analyse()

	err := config.Load(env.Provider("common_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "common_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from envs: %v", err)
	}

	if err := config.Load(file.Provider(configFile), parser); err != nil {
		return nil, ConfigFileNotfoundError{}
	}

	return &Config{
		common: AIConfig{
			key:         config.String(commonKey),
			model:       config.String(commonModel),
			url:         config.String(commonURL),
			temperature: config.Float64(commonTemperature),
			maxTokens:   config.Int(commonMaxTokens),
		},
		user: UserConfig{
			defaultPromptMode: config.String(userDefaultPromptMode),
			preferences:       config.String(userPreferences),
		},
		facts: facts,
	}, nil
}

func WriteConfig(APIKey, configFile string, save bool) (*Config, error) {
	facts := facts.Analyse()

	// openai defaults
	defaults := map[string]interface{}{
		commonURL:             "",
		commonTemperature:     0.2,
		commonMaxTokens:       1000,
		commonModel:           "",
		userDefaultPromptMode: "exec",
		userPreferences:       "",
	}

	err := config.Set(commonKey, APIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to set APIKey in config: %v", err)
	}

	for i, option := range defaults {
		if !config.Exists(i) {
			err := config.Set(i, option)
			if err != nil {
				return nil, fmt.Errorf("failed to set %s in config: %v", i, err)
			}
		}
	}

	if save {
		bytes, err := config.Marshal(parser)
		if err != nil {
			return nil, fmt.Errorf("marshaling parser failed: %v", err)
		}

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			pathExist := os.Mkdir(facts.GetConfigPath(), os.ModeDir)
			if pathExist == nil {
				return nil, fmt.Errorf("problem with creating folder: %v", err)
			}
		}

		err = os.WriteFile(configFile, bytes, 0644)
		if err != nil {
			return nil, fmt.Errorf("problem with writing to file: %v", err)
		}
	}

	return NewConfig(configFile)
}
