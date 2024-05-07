package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAiConfig(t *testing.T) {
	t.Run("GetKey", testGetKey)
	t.Run("GetProxy", testGetProxy)
	t.Run("GetTemperature", testGetTemperature)
	t.Run("GetMaxTokens", testGetMaxTokens)
}

func testGetKey(t *testing.T) {
	expectedKey := "test_key"
	aiConfig := AIConfig{key: expectedKey}

	actualKey := aiConfig.GetKey()

	assert.Equal(t, expectedKey, actualKey, "The two keys should be the same.")
}

func testGetProxy(t *testing.T) {
	expectedURL := "test_url"
	aiConfig := AIConfig{url: expectedURL}

	actualURL := aiConfig.GetURL()

	assert.Equal(t, expectedURL, actualURL, "The two proxies should be the same.")
}

func testGetTemperature(t *testing.T) {
	expectedTemperature := 0.7
	aiConfig := AIConfig{temperature: expectedTemperature}

	actualTemperature := aiConfig.GetTemperature()

	assert.Equal(t, expectedTemperature, actualTemperature, "The two temperatures should be the same.")
}

func testGetMaxTokens(t *testing.T) {
	expectedMaxTokens := 2000
	aiConfig := AIConfig{maxTokens: expectedMaxTokens}

	actualMaxTokens := aiConfig.GetMaxTokens()

	assert.Equal(t, expectedMaxTokens, actualMaxTokens, "The two maxTokens should be the same.")
}
