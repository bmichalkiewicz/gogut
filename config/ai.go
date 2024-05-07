package config

// OpenAI Compabilities
const (
	commonKey         = "settings.key"
	commonModel       = "settings.model"
	commonURL         = "settings.url"
	commonTemperature = "settings.temperature"
	commonMaxTokens   = "settings.max_tokens"
)

type AIConfig struct {
	key         string
	model       string
	url         string
	temperature float64
	maxTokens   int
}

func (c AIConfig) GetKey() string {
	return c.key
}

func (c AIConfig) GetModel() string {
	return c.model
}

func (c AIConfig) GetURL() string {
	return c.url
}

func (c AIConfig) GetTemperature() float64 {
	return c.temperature
}

func (c AIConfig) GetMaxTokens() int {
	return c.maxTokens
}
