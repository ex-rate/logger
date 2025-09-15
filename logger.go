package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Level представляет уровень логирования
type Level = logrus.Level

// Уровни логирования
const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
)

// OutputType определяет тип вывода логов
type OutputType string

const (
	ConsoleOutput OutputType = "console"
	FileOutput    OutputType = "file"
	BothOutput    OutputType = "both"
)

// Config конфигурация логгера
type Config struct {
	Level    Level      `yaml:"level"`
	Output   OutputType `yaml:"output"`
	FilePath string     `yaml:"file_path"`
	Format   string     `yaml:"format"` // json или text
}

// Logger основной логгер приложения
type Logger struct {
	logger      *logrus.Logger
	serviceName string
}

// New создает новый родительский логгер
func New(config Config) (*Logger, error) {
	logger := logrus.New()

	// Устанавливаем уровень логирования
	logger.SetLevel(config.Level)

	// Настраиваем формат вывода
	if err := setupFormatter(logger, config); err != nil {
		return nil, fmt.Errorf("failed to setup formatter: %w", err)
	}

	// Настраиваем вывод
	if err := setupOutput(logger, config); err != nil {
		return nil, fmt.Errorf("failed to setup output: %w", err)
	}

	return &Logger{
		logger:      logger,
		serviceName: "", // Родительский логгер без имени сервиса
	}, nil
}

// setupFormatter настраивает формат вывода логов
func setupFormatter(logger *logrus.Logger, config Config) error {
	// Для консоли всегда используем текстовый формат
	// Для файла - JSON формат
	switch config.Output {
	case ConsoleOutput, BothOutput:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	case FileOutput:
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return nil
}

// setupOutput настраивает вывод логов
func setupOutput(logger *logrus.Logger, config Config) error {
	var writers []io.Writer

	switch config.Output {
	case ConsoleOutput:
		writers = append(writers, os.Stdout)

	case FileOutput:
		if config.FilePath == "" {
			return fmt.Errorf("file path is required for file output")
		}

		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, file)

	case BothOutput:
		writers = append(writers, os.Stdout)

		if config.FilePath != "" {
			file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}
			writers = append(writers, file)
		}

	default:
		return fmt.Errorf("unsupported output type: %s", config.Output)
	}

	// Устанавливаем множественный вывод
	if len(writers) > 1 {
		logger.SetOutput(io.MultiWriter(writers...))
	} else {
		logger.SetOutput(writers[0])
	}

	return nil
}

// withFields добавляет стандартные поля к логу
func (l *Logger) withFields() *logrus.Entry {
	fields := make(map[string]interface{})
	fields["service"] = l.serviceName

	// Добавляем информацию о вызывающей функции
	if pc, file, line, ok := runtime.Caller(2); ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			fields["func"] = fn.Name()
		}
		fields["file"] = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	return l.logger.WithFields(fields)
}

// WithService создает новый логгер с указанным именем сервиса
func (l *Logger) WithService(serviceName string) *Logger {
	return &Logger{
		logger:      l.logger,
		serviceName: serviceName,
	}
}

// WithGroup создает новый логгер с дополнительной группой
func (l *Logger) WithGroup(group string) *Logger {
	serviceName := l.serviceName
	if serviceName != "" {
		serviceName = fmt.Sprintf("%s.%s", serviceName, group)
	} else {
		serviceName = group
	}

	return &Logger{
		logger:      l.logger,
		serviceName: serviceName,
	}
}

// Debug логирует сообщение на уровне Debug
func (l *Logger) Debug(args ...interface{}) {
	l.withFields().Debug(args...)
}

// Debugf логирует форматированное сообщение на уровне Debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.withFields().Debugf(format, args...)
}

// Info логирует сообщение на уровне Info
func (l *Logger) Info(args ...interface{}) {
	l.withFields().Info(args...)
}

// Infof логирует форматированное сообщение на уровне Info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.withFields().Infof(format, args...)
}

// Warn логирует сообщение на уровне Warn
func (l *Logger) Warn(args ...interface{}) {
	l.withFields().Warn(args...)
}

// Warnf логирует форматированное сообщение на уровне Warn
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.withFields().Warnf(format, args...)
}

// Error логирует сообщение на уровне Error
func (l *Logger) Error(args ...interface{}) {
	l.withFields().Error(args...)
}

// Errorf логирует форматированное сообщение на уровне Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.withFields().Errorf(format, args...)
}

// Fatal логирует сообщение на уровне Fatal и завершает программу
func (l *Logger) Fatal(args ...interface{}) {
	l.withFields().Fatal(args...)
}

// Fatalf логирует форматированное сообщение на уровне Fatal и завершает программу
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.withFields().Fatalf(format, args...)
}

// Panic логирует сообщение на уровне Panic и вызывает панику
func (l *Logger) Panic(args ...interface{}) {
	l.withFields().Panic(args...)
}

// Panicf логирует форматированное сообщение на уровне Panic и вызывает панику
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.withFields().Panicf(format, args...)
}

// WithField добавляет поле к логу
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.withFields().WithField(key, value)
}

// WithFields добавляет несколько полей к логу
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.withFields().WithFields(fields)
}

// WithError добавляет ошибку к логу
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.withFields().WithError(err)
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level Level) {
	l.logger.SetLevel(level)
}

// GetLevel возвращает текущий уровень логирования
func (l *Logger) GetLevel() Level {
	return l.logger.GetLevel()
}
