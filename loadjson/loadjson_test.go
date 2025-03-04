package loadjson_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"sm_parser_go/loadjson"
)

func createTestConfigFile(t *testing.T, config map[string]interface{}) string {
	// Создаем временный файл
	tempFile, err := os.CreateTemp("", "config-*.json")
	require.NoError(t, err)
	defer tempFile.Close()
	
	// Записываем JSON в файл
	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(config)
	require.NoError(t, err)
	
	return tempFile.Name()
}

func TestLoadJSON_ValidConfig(t *testing.T) {
	// Подготовка тестовых данных
	testConfig := map[string]interface{}{
		"cookies_yt": map[string]string{
			"cookie1": "value1",
			"cookie2": "value2",
		},
		"headers_yt": map[string]string{
			"User-Agent": "test-agent",
			"Accept":     "application/json",
		},
	}
	
	configPath := createTestConfigFile(t, testConfig)
	defer os.Remove(configPath)
	
	// Выполнение тестируемой функции
	fetchConfig, err := loadjson.LoadJSON(configPath)
	
	// Проверка результатов
	require.NoError(t, err)
	require.NotNil(t, fetchConfig)
	
	// Проверка cookies
	assert.Equal(t, "value1", fetchConfig.Cookies["cookie1"])
	assert.Equal(t, "value2", fetchConfig.Cookies["cookie2"])
	assert.Len(t, fetchConfig.Cookies, 2)
	
	// Проверка headers
	assert.Equal(t, "test-agent", fetchConfig.Headers["User-Agent"])
	assert.Equal(t, "application/json", fetchConfig.Headers["Accept"])
	assert.Len(t, fetchConfig.Headers, 2)
}

func TestLoadJSON_FileNotExists(t *testing.T) {
	// Использование несуществующего пути
	nonExistentPath := "/path/to/nonexistent/config.json"
	
	// Выполнение тестируемой функции
	fetchConfig, err := loadjson.LoadJSON(nonExistentPath)
	
	// Проверка результатов
	assert.Error(t, err)
	assert.Nil(t, fetchConfig)
	assert.Contains(t, err.Error(), "failed to open config file")
}

func TestLoadJSON_InvalidJSON(t *testing.T) {
	// Создаем временный файл с некорректным JSON
	tempFile, err := os.CreateTemp("", "invalid-config-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Записываем некорректный JSON
	_, err = tempFile.WriteString("{invalid json}")
	require.NoError(t, err)
	tempFile.Close()
	
	// Выполнение тестируемой функции
	fetchConfig, err := loadjson.LoadJSON(tempFile.Name())
	
	// Проверка результатов
	assert.Error(t, err)
	assert.Nil(t, fetchConfig)
	assert.Contains(t, err.Error(), "failed to decode config")
}

func TestLoadJSON_EmptyConfig(t *testing.T) {
	// Подготовка пустого конфига
	emptyConfig := map[string]interface{}{}
	
	configPath := createTestConfigFile(t, emptyConfig)
	defer os.Remove(configPath)
	
	// Выполнение тестируемой функции
	fetchConfig, err := loadjson.LoadJSON(configPath)
	
	// Проверка результатов
	require.NoError(t, err)
	require.NotNil(t, fetchConfig)
	
	assert.Empty(t, fetchConfig.Cookies)
	assert.Empty(t, fetchConfig.Headers)
}

func TestLoadJSON_PartialConfig(t *testing.T) {
	// Подготовка частичного конфига (только cookies)
	partialConfig := map[string]interface{}{
		"cookies_yt": map[string]string{
			"cookie1": "value1",
		},
	}
	
	configPath := createTestConfigFile(t, partialConfig)
	defer os.Remove(configPath)
	
	// Выполнение тестируемой функции
	fetchConfig, err := loadjson.LoadJSON(configPath)
	
	// Проверка результатов
	require.NoError(t, err)
	require.NotNil(t, fetchConfig)
	
	assert.Equal(t, "value1", fetchConfig.Cookies["cookie1"])
	assert.Len(t, fetchConfig.Cookies, 1)
	assert.Empty(t, fetchConfig.Headers)
}

func TestLoadJSON_AbsolutePath(t *testing.T) {
	// Подготовка тестовых данных
	testConfig := map[string]interface{}{
		"cookies_yt": map[string]string{"cookie1": "value1"},
		"headers_yt": map[string]string{"User-Agent": "test-agent"},
	}
	
	configPath := createTestConfigFile(t, testConfig)
	defer os.Remove(configPath)
	
	// Получение абсолютного пути
	absPath, err := filepath.Abs(configPath)
	require.NoError(t, err)
	
	// Выполнение тестируемой функции с абсолютным путем
	fetchConfig, err := loadjson.LoadJSON(absPath)
	
	// Проверка результатов
	require.NoError(t, err)
	require.NotNil(t, fetchConfig)
	assert.Equal(t, "value1", fetchConfig.Cookies["cookie1"])
	assert.Equal(t, "test-agent", fetchConfig.Headers["User-Agent"])
}