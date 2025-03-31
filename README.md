[![Version](https://img.shields.io/badge/Version-1.1.0-blue)](https://github.com/Your-RoGr/gogger/tree/master)
[![Latest Release](https://img.shields.io/github/v/release/Your-RoGr/gogger)](https://github.com/Your-RoGr/gogger/releases)
[![codecov.io](https://codecov.io/gh/Your-RoGr/gogger/branch/master/graph/badge.svg?branch=master)](https://codecov.io/gh/Your-RoGr/gogger?branch=master)
![License](https://img.shields.io/github/license/Your-RoGr/gogger)
![Downloads](https://img.shields.io/github/downloads/Your-RoGr/gogger/total)
[![Go Report Card](https://goreportcard.com/badge/Your-RoGr/gogger)](https://goreportcard.com/report/github.com/Your-RoGr/gogger)
![GitHub Stars](https://img.shields.io/github/stars/Your-RoGr/gogger?style=social)

# Gogger

[English](README.md) | [Русский](README.ru.md)

This repository is dedicated to a logger program implemented in the Go programming language.

- [Features](#features)
- [Technologies](#technologies)
- [Usage](#usage)
- [Installation](#Installation)

## Features

Logging levels are implemented as a structure similar to an `enum` and are called `LogLevel`:

| Level    | Representation              |
| -------- | --------------------------- |
| DEBUG    | `gogger.DEBUG`     |
| INFO     | `gogger.INFO`      |
| WARNING  | `gogger.WARNING`   |
| ERROR    | `gogger.ERROR`     |

The following setup functions are available:

| Function           | Arguments                                                                   | Fields in Gogger struct                                   | Description                                                                                                                          |
|--------------------|-----------------------------------------------------------------------------|----------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| SetLogLevel        | level `LogLevel`                                                            | -                                                        | Sets the logging level for both console and file                                                                                     |
| SetLogLevelConsole | level `LogLevel`                                                            | INFO                                                     | Sets the logging level for console                                                                                                   |
| SetLogLevelFile    | level `LogLevel`                                                            | WARNING                                                  | Sets the logging level for file                                                                                                      |
| SetLogFormat       | format `string`                                                             | "[%timestamp%] [%level%] %message%"                      | Sets the log output format                                                                                                           |
| SetUseConsoleLog   | console `bool`                                                              | true                                                     | Sets the flag for using console output (true - enable)                                                                               |
| SetUseFileLog      | file `bool`                                                                 | true                                                     | Sets the flag for using file output (true - enable)                                                                                  |
| SetClearAll        | clearAll `bool`                                                             | false                                                    | When true, deletes all log files in the directory with the same name when creating a Gogger object or when calling SetFilename      |
| SetFilename        | filename `string`, pathFolder `string` = "logs", maxEntries `int` = 1000000 | pathFolder `string` = "logs", maxEntries `int` = 1000000 | Sets a new filename                                                                                                                  |
| SetMaxEntries      | maxEntries `int`                                                            | maxEntries `int` = 1000000                               | Sets the number of entries in one file                                                                                               |
| SetMaxFiles        | maxFiles `int`                                                              | maxFiles `int` = 5                                       | Sets the maximum number of files                                                                                                     |

## Technologies

Gogger uses the following technologies:

- [Go](https://golang.org/) - Programming language

## Usage

| Function             | Arguments                    | Description                                   |
| -------------------- | ---------------------------- | --------------------------------------------- |
| NewGogger            | filename `string`, pathFolder `string` = "logs", maxEntries `int` = 1000000, maxFiles `int` = 5 | Creates a new instance of the Gogger class |
| Log                  | level `LogLevel`, message `string` | Writes a log with the specified logging level |
| Debug                | debugMessage `string`         | Writes a log with DEBUG logging level         |
| Info                 | infoMessage `string`          | Writes an informational log                   |
| Warning              | warningMessage `string`       | Writes a warning                              |
| Error                | errorMessage `string`         | Writes an error log                           |

Example usage in a Go program:

```go
package main

import (
	"gogger/gogger"
)

func main() {
    logger, err := gogger.NewGogger("logfile.txt", "logs", 8, 5)
    if err != nil {
        // Error handling
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

Console output:

```sh
[30-09-2020 21:59:05] [WARNING] console Warning message
[30-09-2020 21:59:05] [INFO] console Info message
[30-09-2020 21:59:05] [WARNING] console Warning message
[30-09-2020 21:59:05] [ERROR] console Error message
```

As you can see from the example, the log with DEBUG level is not displayed in the console.

## Installation

To install the package, use the command:

```bash
go get github.com/Your-RoGr/gogger
```

## License

gogger is MIT-Licensed
