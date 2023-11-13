package main

import (
	"fmt"
	"gogger/gogger"
)

func main() {
	// Создаем экземпляр Gogger
	newGogger, err := gogger.NewGogger("example.log", "logs", 1000, 5)
	if err != nil {
		panic(err)
	}
	defer newGogger.Close()
	// Устанавливаем уровень логирования для консоли и файла
	newGogger.SetLogLevelConsole(gogger.DEBUG)

	// Устанавливаем формат лога
	err = newGogger.SetLogFormat("[%timestamp%] [%level%] %message%")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Устанавливаем использование консоли и файла для логирования
	newGogger.SetUseConsoleLog(true)
	newGogger.SetUseFileLog(true)

	// Устанавливаем опцию очистки всех файлов
	newGogger.SetClearAll(false)

	// Логирование различных сообщений
	newGogger.Debug("This is a debug message.")
	newGogger.Info("This is an info message.")
	newGogger.Warning("This is a warning message.")
	newGogger.Error("This is an error message.")

	// Устанавливаем новый файл, путь к папке и максимальное количество записей
	err = newGogger.SetFilename("newlog.log", "newlogs", 500)
	if err != nil {
		panic(err)
	}

	// Устанавливаем максимальное количество записей
	newGogger.SetMaxEntries(500)

	// Устанавливаем максимальное количество файлов
	newGogger.SetMaxFiles(3)

	// Логирование после изменения настроек
	newGogger.Info("This is a new info message.")
	newGogger.Error("This is a new error message.")
}
