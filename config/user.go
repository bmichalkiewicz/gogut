package config

const (
	userDefaultPromptMode = "user.default_prompt_mode"
	userPreferences       = "user.preferences"
)

type UserConfig struct {
	defaultPromptMode string
	preferences       string
}

func (c UserConfig) GetDefaultPromptMode() string {
	return c.defaultPromptMode
}

func (c UserConfig) GetPreferences() string {
	return c.preferences
}
