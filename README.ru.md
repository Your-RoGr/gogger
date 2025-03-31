[![Version](https://img.shields.io/badge/Version-1.1.0-blue)](https://github.com/Your-RoGr/gogger/tree/master)
[![Latest Release](https://img.shields.io/github/v/release/Your-RoGr/gogger)](https://github.com/Your-RoGr/gogger/releases)
[![codecov.io](https://codecov.io/gh/Your-RoGr/gogger/branch/master/graph/badge.svg?branch=master)](https://codecov.io/gh/Your-RoGr/gogger?branch=master)
![License](https://img.shields.io/github/license/Your-RoGr/gogger)
![Downloads](https://img.shields.io/github/downloads/Your-RoGr/gogger/total)
[![Go Report Card](https://goreportcard.com/badge/Your-RoGr/gogger)](https://goreportcard.com/report/github.com/Your-RoGr/gogger)
![GitHub Stars](https://img.shields.io/github/stars/Your-RoGr/gogger?style=social)

# Gogger

[English](README.md) | [Русский](README.ru.md)

Этот репозиторий посвящен программе логгера, реализованной на языке программирования Golang.

- [Особенности](#Особенности)
- [Технологии](#Технологии)
- [Использование](#Использование)
- [Установка](#Установка)

## Особенности

Уровни логирования реализованы в виде структуры, аналогичной `enum`, и называются `LogLevel`:

| Уровень   | Представление               |
| --------- | --------------------------- |
| DEBUG     | `gogger.DEBUG`     |
| INFO      | `gogger.INFO`      |
| WARNING   | `gogger.WARNING`   |
| ERROR     | `gogger.ERROR`     |

Доступны следующие функции установки:

| Функция            | Аргументы                                                                   | Поля в struct Gogger                                     | Описание                                                                                                                             |
|--------------------|-----------------------------------------------------------------------------|----------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| SetLogLevel        | level `LogLevel`                                                            | -                                                        | Устанавливает уровень логирования для консоли и файла                                                                                |
| SetLogLevelConsole | level `LogLevel`                                                            | INFO                                                     | Устанавливает уровень логирования для консоли                                                                                        |
| SetLogLevelFile    | level `LogLevel`                                                            | WARNING                                                  | Устанавливает уровень логирования для файла                                                                                          |
| SetLogFormat       | format `string`                                                             | "[%timestamp%] [%level%] %message%"                      | Устанавливает формат вывода логов                                                                                                    |
| SetUseConsoleLog   | console `bool`                                                              | true                                                     | Устанавливает флаг использования вывода в консоль (true - включить)                                                                  |
| SetUseFileLog      | file `bool`                                                                 | true                                                     | Устанавливает флаг использования вывода в файлы (true - включить)                                                                    |
| SetClearAll        | clearAll `bool`                                                             | false                                                    | При true удаляет все файлы логов в директории с таким же наименованием при создании объекта класса Gogger или при вызове SetFilename |
| SetFilename        | filename `string`, pathFolder `string` = "logs", maxEntries `int` = 1000000 | pathFolder `string` = "logs", maxEntries `int` = 1000000 | Устанавливает новое название файлов                                                                                                  |
| SetMaxEntries      | maxEntries `int`                                                            | maxEntries `int` = 1000000                               | Устанавливает количество записей в одном файле                                                                                       |
| SetMaxFiles        | maxFiles `int`                                                              | maxFiles `int` = 5                                       | Устанавливает максимальное количество файлов                                                                                         |

## Технологии

Gogger использует следующие технологии:

- [Go](https://golang.org/) - Язык программирования

## Использование

| Функция              | Аргументы                    | Описание                                      |
| --------------------- | ---------------------------- | --------------------------------------------- |
| NewGogger            | filename `string`, pathFolder `string` = "logs", maxEntries `int` = 1000000, maxFiles `int` = 5 | Создает новый экземпляр класса Gogger    |
| Log                  | level `LogLevel`, message `string` | Записывает лог с указанным уровнем логирования |
| Debug                | debugMessage `string`         | Записывает лог с уровнем логирования DEBUG      |
| Info                 | infoMessage `string`          | Записывает информационный лог                   |
| Warning              | warningMessage `string`       | Записывает предупреждение                      |
| Error                | errorMessage `string`         | Записывает лог об ошибке                        |

Пример использования в программе на Go:

```go
package main

import (
	"gogger/gogger"
)

func main() {
    logger, err := gogger.NewGogger("logfile.txt", "logs", 8, 5)
    if err != nil {
        // Обработка ошибки
        return
    }
    defer logger.Close()
    logger.Log(gogger.WARNING, "console Warning message")
    logger.Debug("console Debug message")
    logger.Info("console Info message")
    logger.Warning("console Warning message")
    logger.Error("console Error message")
}
```

Вывод в консоли:

```sh
[30-09-2020 21:59:05] [WARNING] console Warning message
[30-09-2020 21:59:05] [INFO] console Info message
[30-09-2020 21:59:05] [WARNING] console Warning message
[30-09-2020 21:59:05] [ERROR] console Error message
```

Как видно из примера, лог с уровнем DEBUG не отобразился в консоли.

## Установка

Для установки пакета используйте команду:

```bash
go get github.com/Your-RoGr/gogger
```

## Лицензия

gogger is MIT-Licensed
