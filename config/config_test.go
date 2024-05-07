package config

import (
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("NewConfig", testNewConfig)
	t.Run("WriteConfig", testWriteConfig)
}

func setupConfig(t *testing.T) {
	t.Helper()

	err := config.Set(commonKey, "test_key")
	require.NoError(t, err)

	err = config.Set(commonModel, openai.GPT3Dot5Turbo)
	require.NoError(t, err)

	err = config.Set(commonURL, "test_url")
	require.NoError(t, err)

	err = config.Set(commonTemperature, 0.2)
	require.NoError(t, err)

	err = config.Set(commonMaxTokens, 2000)
	require.NoError(t, err)

	err = config.Set(userDefaultPromptMode, "exec")
	require.NoError(t, err)

	err = config.Set(userPreferences, "test_preferences")
	require.NoError(t, err)

	bytes, err := config.Marshal(parser)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile("/tmp/config.yaml", bytes, 0644))
}

func cleanup(t *testing.T) {
	t.Helper()
	require.NoError(t, os.Remove("/tmp/config.yaml"))
}

func testNewConfig(t *testing.T) {
	setupConfig(t)

	cfg, err := NewConfig("/tmp/config.yaml")
	require.NoError(t, err)

	assert.Equal(t, "test_key", cfg.GetAIConfig().GetKey())
	assert.Equal(t, openai.GPT3Dot5Turbo, cfg.GetAIConfig().GetModel())
	assert.Equal(t, "test_url", cfg.GetAIConfig().GetURL())
	assert.Equal(t, 0.2, cfg.GetAIConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAIConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())

}

func testWriteConfig(t *testing.T) {
	setupConfig(t)
	defer cleanup(t)

	cfg, err := WriteConfig("new_test_key", "/tmp/config.yaml", true)
	require.NoError(t, err)

	assert.Equal(t, "new_test_key", cfg.GetAIConfig().GetKey())
	assert.Equal(t, openai.GPT3Dot5Turbo, cfg.GetAIConfig().GetModel())
	assert.Equal(t, "test_url", cfg.GetAIConfig().GetURL())
	assert.Equal(t, 0.2, cfg.GetAIConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAIConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())

	assert.Equal(t, "new_test_key", config.Get(commonKey))
	assert.Equal(t, openai.GPT3Dot5Turbo, config.Get(commonModel))
	assert.Equal(t, "test_url", config.Get(commonURL))
	assert.Equal(t, 0.2, config.Get(commonTemperature))
	assert.Equal(t, 2000, config.Get(commonMaxTokens))
	assert.Equal(t, "exec", config.Get(userDefaultPromptMode))
	assert.Equal(t, "test_preferences", config.Get(userPreferences))

}
