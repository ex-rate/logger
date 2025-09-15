# Logger Package

Пакет логгера для приложения, построенный на основе библиотеки [logrus](https://github.com/sirupsen/logrus).

## Возможности

- ✅ Название сервиса в каждом логе
- ✅ Настраиваемый вывод: консоль, файл или оба
- ✅ Поддержка всех уровней логирования (Trace, Debug, Info, Warn, Error, Fatal, Panic)
- ✅ Автоматическая запись функции, в которой происходит лог
- ✅ Поддержка структурированного логирования (JSON и текстовый формат)
- ✅ Добавление дополнительных полей к логам
- ✅ Логирование ошибок с контекстом

## Установка

Установка:

```bash
go get  github.com/exrate/logger
```

## Быстрый старт

```go
package main

import (
    "github.com/exrate/logger"
)

func main() {
    // Создание логгера
    config := logger.Config{
        ServiceName: "my-service",
        Level:       logger.InfoLevel,
        Output:      logger.ConsoleOutput,
        Format:      "text",
    }
    
    log, err := logger.New(config)
    if err != nil {
        panic(err)
    }
    
    // Использование
    log.Info("Service started")
    log.WithField("port", 8080).Info("Server listening")
}
```

## Конфигурация

### Config

```go
type Config struct {
    ServiceName string     // Название сервиса
    Level       Level      // Уровень логирования
    Output      OutputType // Тип вывода
    FilePath    string     // Путь к файлу (для FileOutput и BothOutput)
    Format      string     // Формат: "json" или "text"
}
```

### Уровни логирования

```go
const (
    TraceLevel = logrus.TraceLevel
    DebugLevel = logrus.DebugLevel
    InfoLevel  = logrus.InfoLevel
    WarnLevel  = logrus.WarnLevel
    ErrorLevel = logrus.ErrorLevel
    FatalLevel = logrus.FatalLevel
    PanicLevel = logrus.PanicLevel
)
```

### Типы вывода

```go
const (
    ConsoleOutput OutputType = "console" // Только консоль
    FileOutput    OutputType = "file"    // Только файл
    BothOutput    OutputType = "both"    // Консоль и файл
)
```

## Примеры использования

### Базовое логирование

```go
log.Info("Information message")
log.Debug("Debug information")
log.Warn("Warning message")
log.Error("Error occurred")
log.Fatal("Fatal error - application will exit")
log.Panic("Panic - application will panic")
```

### Форматированное логирование

```go
log.Infof("User %s logged in from %s", username, ipAddress)
log.Errorf("Failed to connect to database: %v", err)
```

### Логирование с полями

```go
// Одно поле
log.WithField("user_id", "12345").Info("User logged in")

// Несколько полей
fields := map[string]interface{}{
    "user_id":    "12345",
    "action":     "login",
    "ip_address": "192.168.1.1",
}
log.WithFields(fields).Info("User action performed")
```

### Логирование ошибок

```go
if err != nil {
    log.WithError(err).Error("Failed to process request")
}
```

### Изменение уровня логирования

```go
log.SetLevel(logger.DebugLevel)
currentLevel := log.GetLevel()
```

## Примеры конфигурации

### Консольный вывод

```go
config := logger.Config{
    ServiceName: "api-server",
    Level:       logger.InfoLevel,
    Output:      logger.ConsoleOutput,
    Format:      "text",
}
```

### Вывод в файл

```go
config := logger.Config{
    ServiceName: "api-server",
    Level:       logger.DebugLevel,
    Output:      logger.FileOutput,
    FilePath:    "/var/log/api-server.log",
    Format:      "json",
}
```

### Вывод в консоль и файл

```go
config := logger.Config{
    ServiceIntegration: "api-server",
    Level:             logger.InfoLevel,
    Output:            logger.BothOutput,
    FilePath:          "/var/log/api-server.log",
    Format:            "text",
}
```

## Форматы вывода

### Текстовый формат

```
time="2024-01-15T10:30:00Z" level=info msg="Service started" service="my-service" func="main.main()" file="main.go:25"
```

### JSON формат

```json
{
  "caller": "main.go:25",
  "func": "main.main()",
  "level": "info",
  "msg": "Service started",
  "service": "my-service",
  "time": "2024-01-15T10:30:00Z"
}
```

## Интеграция с Echo

Для интеграции с Echo framework можно использовать middleware:

```go
import (
    "github.com/labstack/echo/v4"
    "exrate/reviews-backend/internal/logger"
)

func main() {
    // Создание логгера
    log, _ := logger.New(config)
    
    e := echo.New()
    
    // Middleware для логирования запросов
    e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            
            err := next(c)
            
            log.WithFields(map[string]interface{}{
                "method":     c.Request().Method,
                "uri":        c.Request().RequestURI,
                "status":     c.Response().Status,
                "latency_ms": time.Since(start).Milliseconds(),
            }).Info("HTTP request processed")
            
            return err
        }
    })
}
```

## Тестирование

Запуск тестов:

```bash
go test ./internal/logger
```

Запуск тестов с покрытием:

```bash
go test -cover ./internal/logger
```

## Зависимости

- `github.com/sirupsen/logrus` - основная библиотека логирования
- `github.com/stretchr/testify` - для тестирования (dev dependency) 