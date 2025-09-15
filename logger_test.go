package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid console config",
			config: Config{
				Level:  InfoLevel,
				Output: ConsoleOutput,
				Format: "text",
			},
			wantErr: false,
		},
		{
			name: "valid json config",
			config: Config{
				Level:  DebugLevel,
				Output: ConsoleOutput,
				Format: "json",
			},
			wantErr: false,
		},
		{
			name: "file output without path",
			config: Config{
				Level:  InfoLevel,
				Output: FileOutput,
				Format: "text",
			},
			wantErr: true,
		},
		{
			name: "invalid output type",
			config: Config{
				Level:  InfoLevel,
				Output: "invalid",
				Format: "text",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.Equal(t, tt.config.Level, logger.GetLevel())
			}
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	config := Config{
		Level:  DebugLevel,
		Output: ConsoleOutput,
		Format: "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	entry := logger.WithField("key", "value")
	assert.NotNil(t, entry)

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	entry = logger.WithFields(fields)
	assert.NotNil(t, entry)
}

func TestLogger_WithError(t *testing.T) {
	config := Config{
		Level:  DebugLevel,
		Output: ConsoleOutput,
		Format: "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	testErr := assert.AnError
	entry := logger.WithError(testErr)
	assert.NotNil(t, entry)
}

func TestLogger_WithGroup(t *testing.T) {
	config := Config{
		Level:  DebugLevel,
		Output: ConsoleOutput,
		Format: "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	groupLogger := logger.WithGroup("subservice")
	assert.NotNil(t, groupLogger)
	assert.Equal(t, "subservice", groupLogger.serviceName)
}

func TestLogger_WithService(t *testing.T) {
	config := Config{
		Level:  DebugLevel,
		Output: ConsoleOutput,
		Format: "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	serviceName := "service"

	serviceLogger := logger.WithService(serviceName)
	require.NotNil(t, serviceLogger)
	assert.Equal(t, serviceName, serviceLogger.serviceName)
}

func TestLogger_Levels(t *testing.T) {
	config := Config{
		Level:  DebugLevel,
		Output: ConsoleOutput,
		Format: "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	// Тестируем все уровни логирования
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Проверяем, что уровень можно изменить
	logger.SetLevel(InfoLevel)
	assert.Equal(t, InfoLevel, logger.GetLevel())
}

func TestLogger_FileOutput(t *testing.T) {
	tempFile := t.TempDir() + "/test.log"

	config := Config{
		Level:    InfoLevel,
		Output:   FileOutput,
		FilePath: tempFile,
		Format:   "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	logger.Info("test message")

	// Проверяем, что файл создался и содержит сообщение
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

func TestLogger_BothOutput(t *testing.T) {
	tempFile := t.TempDir() + "/test.log"

	config := Config{
		Level:    InfoLevel,
		Output:   BothOutput,
		FilePath: tempFile,
		Format:   "text",
	}

	logger, err := New(config)
	require.NoError(t, err)

	logger.Info("test message")

	// Проверяем, что файл создался и содержит сообщение
	content, err := os.ReadFile(tempFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}
