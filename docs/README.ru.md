# goschedviz — Визуализация работы планировщика Go

[![Build Status](https://github.com/JustSkiv/goschedviz/workflows/build/badge.svg)](https://github.com/JustSkiv/goschedviz/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/JustSkiv/goschedviz.svg)](https://pkg.go.dev/github.com/JustSkiv/goschedviz)
[![Release](https://img.shields.io/github/release/JustSkiv/goschedviz.svg?style=flat-square)](https://github.com/JustSkiv/goschedviz/releases)
[![Telegram](https://img.shields.io/badge/Telegram-@ntuzov-blue?logo=telegram&logoColor=white)](https://t.me/ntuzov)

[![Go Report Card](https://goreportcard.com/badge/github.com/JustSkiv/goschedviz)](https://goreportcard.com/report/github.com/JustSkiv/goschedviz)
[![codecov](https://codecov.io/gh/JustSkiv/goschedviz/branch/main/graph/badge.svg)](https://app.codecov.io/gh/JustSkiv/goschedviz)
[![Coverage Status](https://coveralls.io/repos/github/JustSkiv/goschedviz/badge.svg?branch=main)](https://coveralls.io/github/JustSkiv/goschedviz?branch=main)

[English version](../README.md)

Инструмент для визуализации работы планировщика Go в терминале. Помогает понять поведение планировщика Go через
отображение метрик в реальном времени.

![Демонстрация работы](../demo.gif)

⚠️ **Важно**: Этот инструмент предназначен только для образовательных целей. Он разработан для помощи в понимании работы
планировщика Go и не должен использоваться в продакшен-окружении или критически важных проектах. В нём могут быть ошибки
и он не оптимизирован для производительности.

## Возможности

- Мониторинг метрик планировщика Go в реальном времени с использованием [GODEBUG schedtrace](https://pkg.go.dev/github.com/maximecaron/gotraining/topics/profiling/godebug/schedtrace)
- Мониторинг количества горутин через runtime метрики
- Консольный интерфейс с несколькими виджетами:
    - Таблица текущих значений планировщика
    - Диаграммы локальных очередей (LRQ)
    - Индикаторы для GRQ, горутин, потоков и простаивающих процессоров
    - Два графика истории (линейная и логарифмическая шкалы)
    - Цветовая легенда метрик
- Поддержка мониторинга любой Go-программы

## Установка

### Вариант 1: Из исходного кода

Клонируйте и соберите проект:

```bash
git clone https://github.com/JustSkiv/goschedviz
cd goschedviz
make build
```

Исполняемый файл будет создан в директории `bin`.

### Вариант 2: Используя go install

```bash
go install github.com/JustSkiv/goschedviz/cmd/goschedviz@latest
```

Это установит исполняемый файл `goschedviz` в директорию `$GOPATH/bin`. Убедитесь, что эта директория добавлена в PATH.

## Использование

### Базовое использование

```bash
goschedviz -target=/path/to/your/program.go -period=1000
```

Где:
- `-target`: Путь к Go-программе для мониторинга
- `-period`: Период GODEBUG schedtrace в миллисекундах (по умолчанию: 1000)

### Добавление метрик горутин в вашу программу

Для включения мониторинга количества горутин добавьте reporter метрик в вашу программу:

```go
package main

import (
    "time"
    "github.com/JustSkiv/goschedviz/pkg/metrics"
)

func main() {
    // Инициализация reporter'а метрик
    reporter := metrics.NewReporter(time.Second)
    reporter.Start()
    defer reporter.Stop()

    // Ваша логика программы
    ...
}
```

Reporter автоматически будет отправлять метрики о количестве горутин, которые будут отображаться в goschedviz.

### Управление

- `q` или `Ctrl+C`: Выход из программы
- Поддерживается изменение размера терминала

## Пример

1. Создайте простую тестовую программу (example.go):

```go
package main

import "time"

func main() {
	// Создаем нагрузку на планировщик
	for i := 0; i < 1000; i++ {
		go func() {
			time.Sleep(time.Second)
		}()
	}
	time.Sleep(10 * time.Second)
}
```

2. Запустите визуализацию:

```bash
goschedviz -target=example.go
```

Или попробуйте готовый пример:

```bash
# Простой пример с интенсивной нагрузкой на CPU, GOMAXPROCS=2 и большим количеством горутин
goschedviz -target=examples/simple/main.go
```

Приветствуются новые интересные демонстрационные примеры. Особенно ценны примеры, показывающие различное поведение
планировщика (подробности в [Contributing](CONTRIBUTING.ru.md)).

Что могут продемонстрировать хорошие примеры:

- Нагрузку на CPU против I/O операций
- Различные конфигурации GOMAXPROCS
- Приложения с сетевой нагрузкой
- Операции с интенсивным использованием памяти
- Специфические паттерны работы планировщика или краевые случаи

Это помогает другим лучше понять поведение планировщика Go в различных сценариях.

## Понимание вывода

Интерфейс отображает несколько ключевых метрик:

- **Таблица текущих значений**: Показывает текущее состояние планировщика, включая GOMAXPROCS, количество потоков и т.д.
- **Столбцы локальных очередей**: Визуализирует длину очереди для каждого P (процессора)
- **Индикаторы метрик**: 
  * GRQ - длина глобальной очереди выполнения
  * Количество активных горутин
  * Количество системных потоков
  * Количество простаивающих процессоров
- **Графики истории**:
  * Линейная шкала для точного отслеживания значений
  * Логарифмическая шкала для лучшей визуализации больших диапазонов
- **Легенда**: Цветовая кодировка метрик на графиках:
  * GRQ - Глобальная очередь (зеленый)
  * LRQ - Сумма локальных очередей (пурпурный)
  * THR - Системные потоки (красный)
  * IDL - Простаивающие процессоры (желтый)
  * GRT - Горутины (голубой)

## Как это работает

Инструмент:

1. Запускает вашу Go-программу с включенным GODEBUG=schedtrace
2. Анализирует вывод трассировки планировщика в реальном времени
3. Собирает дополнительные runtime метрики (количество горутин)
4. Визуализирует все метрики через терминальный интерфейс

## Требования

- Go 1.23 или новее
- UNIX-подобная операционная система (Linux, macOS)
- Терминал с поддержкой цветов

## Разработка

```bash
# Сборка проекта
make build

# Запуск тестов
make test

# Очистка артефактов сборки
make clean
```

## Ресурсы автора

- [YouTube канал](https://www.youtube.com/@nikolay_tuzov) - Туториалы по Go
- [@ntuzov](https://t.me/ntuzov) - Основной Telegram-канал: гайды, новости, анонсы и многое другое
- [@golang_digest](https://t.me/golang_digest) - Полезные материалы и ресурсы по Go

## Участие в разработке

Ваше участие приветствуется! Неважно, исправляете ли вы баги, улучшаете документацию или добавляете новые функции —
ваша помощь ценна для проекта.

Если вы новичок в open source или Go-разработке, этот проект отлично подойдет для старта. Рекомендую сначала
ознакомиться с [руководством по контрибьюту](CONTRIBUTING.md).

Не стесняйтесь задавать вопросы — это помогает вам учиться и развиваться ❤️

## Цитирование

Если вы используете goschedviz в своем проекте, исследовании или учебных материалах, пожалуйста, укажите ссылку:

```
Этот проект использует goschedviz (https://github.com/JustSkiv/goschedviz) от Николая Тузова
```

## Лицензия

MIT License — подробности в файле [LICENSE](LICENSE).

