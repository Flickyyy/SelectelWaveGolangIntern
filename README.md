# loglint

Go-линтер для проверки лог-записей, совместимый с golangci-lint.

Анализирует вызовы `log/slog` и `go.uber.org/zap`, проверяя сообщения на соответствие правилам оформления и безопасности.

## Правила

| # | Правило | Пример нарушения |
|---|---------|------------------|
| 1 | Сообщение начинается со строчной буквы | `slog.Info("Starting server")` |
| 2 | Только английский язык | `slog.Info("запуск сервера")` |
| 3 | Без спецсимволов и эмодзи | `slog.Info("started!🚀")` |
| 4 | Без чувствительных данных | `slog.Info("token: " + token)` |

Для правил 1 и 3 реализованы `SuggestedFixes` — автоматические исправления.

## Поддерживаемые логгеры

- `log/slog` — пакетные функции (`slog.Info`, `slog.Error`, ...) и методы `*slog.Logger`, включая `Context`-варианты
- `go.uber.org/zap` — методы `*zap.Logger` и `*zap.SugaredLogger` (`Infow`, `Errorf`, ...)

## Сборка и запуск

```bash
# сборка
make build

# запуск на проекте
./bin/loglint ./...

# или через go run
go run ./cmd/loglint ./...
```

## Тесты

```bash
make test
# или
go test -v -race ./...
```

## Интеграция с golangci-lint

### Вариант 1: module plugin (golangci-lint v1.57+)

Создать `.custom-gcl.yml`:

```yaml
version: v1.64.0
plugins:
  - module: 'github.com/Flickyyy/SelectelWaveGolangIntern'
    import: 'github.com/Flickyyy/SelectelWaveGolangIntern/plugin'
    version: latest
```

Собрать кастомный golangci-lint:

```bash
golangci-lint custom
```

### Вариант 2: Go plugin (.so)

```bash
go build -buildmode=plugin -o loglint.so ./plugin
```

Указать путь к `.so` в конфигурации golangci-lint.

## Структура проекта

```
loglint.go          — определение анализатора, обнаружение лог-вызовов
rules.go            — реализация правил проверки
loglint_test.go     — тесты через analysistest
cmd/loglint/        — standalone бинарник (singlechecker)
plugin/             — плагин для golangci-lint
testdata/           — тестовые данные (slog, zap)
```

## Стек

- Go 1.24+
- `golang.org/x/tools/go/analysis` — фреймворк для статического анализа
- `analysistest` — тестирование анализатора

## Примеры

Запуск на тестовом файле:

```bash
$ cat example.go
package main

import "log/slog"

func main() {
    slog.Info("Starting server")
    password := "secret"
    slog.Info("password: " + password)
}

$ ./bin/loglint ./...
example.go:6:15: log message should start with a lowercase letter
example.go:8:15: log message may contain sensitive data (keyword: "password")
```

## CI

GitHub Actions: lint -> test -> build. Конфигурация в `.github/workflows/ci.yml`.

## Использование AI

AI-инструменты использовались для генерации заглушек (stub) в тестовых данных для zap и для первичной вычитки документации.
