package gogger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// LogLevel тип перечисления для уровней логирования
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Gogger структура для логирования
type Gogger struct {
	filename          string
	fileStream        *os.File
	logLevelConsole   LogLevel
	logLevelFile      LogLevel
	logFormat         string
	pathFolder        string
	console           bool
	file              bool
	clearAll          bool
	logQueueFiles     []string
	maxEntries        int
	maxEntriesCounter int
	maxFiles          int
	logFileNumber     int
}

// NewGogger создает новый экземпляр Gogger
func NewGogger(filename, pathFolder string, maxEntries, maxFiles int) (*Gogger, error) {

	if !isValidFilename(filename) || !isValidPathFolder(pathFolder) || maxEntries <= 0 {
		return nil, fmt.Errorf("Invalid filename, path folder, or max_entries")
	}

	l := &Gogger{
		filename:          filename,
		pathFolder:        pathFolder,
		maxEntries:        maxEntries,
		maxEntriesCounter: maxEntries,
		maxFiles:          maxFiles,
	}

	if err := l.createFolder(); err != nil {
		return nil, err
	}

	l.addCurrentFiles()
	l.deleteAllFiles()
	l.openFile()

	return l, nil
}

// Close закрывает файловый поток при уничтожении Gogger
func (l *Gogger) Close() {
	if l.fileStream != nil {
		err := l.fileStream.Close()
		if err != nil {
			return
		}
	}
}

// Log записывает сообщение с уровнем логирования
func (l *Gogger) Log(level LogLevel, message string) {
	if l.file || l.console {
		timestamp := getFormattedTimestamp()
		levelString := getLogLevelString(level)
		formattedMessage := l.logFormat

		formattedMessage = replacePlaceholder(formattedMessage, "%timestamp%", timestamp)
		formattedMessage = replacePlaceholder(formattedMessage, "%level%", levelString)
		formattedMessage = replacePlaceholder(formattedMessage, "%message%", message)

		if l.file && level >= l.logLevelFile {
			l.writeLogsFile(formattedMessage)
		}
		if l.console && level >= l.logLevelConsole {
			l.writeLogsToConsole(formattedMessage)
		}

	} else {
		fmt.Println("No log input in use")
	}
}

// Debug записывает сообщение отладки
func (l *Gogger) Debug(debugMessage string) {
	l.Log(DEBUG, debugMessage)
}

// Info записывает информационное сообщение
func (l *Gogger) Info(infoMessage string) {
	l.Log(INFO, infoMessage)
}

// Warning записывает предупреждение
func (l *Gogger) Warning(warningMessage string) {
	l.Log(WARNING, warningMessage)
}

// Error записывает сообщение об ошибке
func (l *Gogger) Error(errorMessage string) {
	l.Log(ERROR, errorMessage)
}

// SetLogLevel устанавливает уровень логирования для консоли и файла
func (l *Gogger) SetLogLevel(level LogLevel) {
	l.logLevelConsole = level
	l.logLevelFile = level
}

// SetLogLevelConsole устанавливает уровень логирования для консоли
func (l *Gogger) SetLogLevelConsole(level LogLevel) {
	l.logLevelConsole = level
}

// SetLogLevelFile устанавливает уровень логирования для файла
func (l *Gogger) SetLogLevelFile(level LogLevel) {
	l.logLevelFile = level
}

// SetLogFormat устанавливает формат лога
func (l *Gogger) SetLogFormat(format string) error {
	var requiredElements = []string{"%timestamp%", "%level%", "%message%"}

	isValidFormat := false
	for _, element := range requiredElements {
		if strings.Contains(format, element) {
			isValidFormat = true
			break
		}
	}

	if isValidFormat {
		l.logFormat = format
		return nil
	}

	return fmt.Errorf("Invalid log format. The format must contain at least one of the following elements: %timestamp%, %level%, %message%")
}

// SetUseConsoleLog устанавливает использование консоли для логирования
func (l *Gogger) SetUseConsoleLog(console bool) {
	l.console = console
}

// SetUseFileLog устанавливает использование файла для логирования
func (l *Gogger) SetUseFileLog(file bool) {
	l.file = file
}

// SetClearAll устанавливает опцию очистки всех файлов
func (l *Gogger) SetClearAll(clearAll bool) {
	l.clearAll = clearAll
}

// SetFilename устанавливает имя файла, путь к папке и максимальное количество записей
func (l *Gogger) SetFilename(filename, pathFolder string, maxEntries int) error {
	if !isValidFilename(filename) || !isValidPathFolder(pathFolder) || maxEntries <= 0 {
		return fmt.Errorf("Invalid filename, path folder, or max_entries")
	}

	l.filename = filename
	l.maxEntries = maxEntries
	l.maxEntriesCounter = maxEntries
	l.pathFolder = pathFolder

	l.logFileNumber = 0
	l.logQueueFiles = nil

	err := l.createFolder()
	if err != nil {
		return err
	}
	l.addCurrentFiles()
	l.deleteAllFiles()
	l.openFile()

	return nil
}

// SetMaxEntries устанавливает максимальное количество записей
func (l *Gogger) SetMaxEntries(maxEntries int) {
	if l.maxEntries-l.maxEntriesCounter < maxEntries {
		if l.logFileNumber+1 == l.maxFiles {
			l.logFileNumber = 0
		} else {
			l.logFileNumber++
		}
	}

	l.maxEntries = maxEntries
	l.maxEntriesCounter = maxEntries
}

// SetMaxFiles устанавливает максимальное количество файлов
func (l *Gogger) SetMaxFiles(maxFiles int) {
	l.maxFiles = maxFiles
}

func (l *Gogger) openFile() {
	for {
		if len(l.logQueueFiles) > 0 {
			for l.maxFiles < len(l.logQueueFiles)+1 && l.maxFiles != 0 {
				l.deleteFirstFile()
			}
		}

		if l.maxEntries > l.getCountOfLines() || l.maxEntries == 0 {
			break
		}

		l.logQueueFiles = append(l.logQueueFiles, fmt.Sprintf("#%d%s", l.logFileNumber, l.filename))

		if l.logFileNumber+1 == l.maxFiles {
			l.logFileNumber = 0
		} else {
			l.logFileNumber++
		}
	}

	var filePath string
	if l.pathFolder == "" {
		filePath = fmt.Sprintf("#%d%s", l.logFileNumber, l.filename)
	} else {
		filePath = fmt.Sprintf("%s\\#%d%s", l.pathFolder, l.logFileNumber, l.filename)
	}
	l.logQueueFiles = append(l.logQueueFiles, filePath)

	var err error
	if l.fileStream, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		fmt.Printf("Error opening file: %v\n", err)
	}
}

func (l *Gogger) deleteFirstFile() {
	if len(l.logQueueFiles) > 0 {
		filePath := l.logQueueFiles[0]
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Error deleting file: %v\n", err)
		}
		l.logQueueFiles = l.logQueueFiles[1:]
	}
}

func (l *Gogger) deleteAllFiles() {
	if l.clearAll {
		for len(l.logQueueFiles) > 0 {
			l.deleteFirstFile()
		}
	}
}

func (l *Gogger) writeLogsToFile(formattedMessage string) {
	if l.fileStream == nil {
		l.openFile()
	}

	if _, err := io.WriteString(l.fileStream, formattedMessage+"\n"); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func (l *Gogger) writeLogsFile(formattedMessage string) {
	if l.maxEntriesCounter > 0 || l.maxEntries == 0 {
		l.maxEntriesCounter--
		l.writeLogsToFile(formattedMessage)
	} else {
		l.maxEntriesCounter = l.maxEntries
		if l.fileStream != nil {
			err := l.fileStream.Close()
			if err != nil {
				return
			}
		}

		if l.logFileNumber+1 == l.maxFiles {
			l.logFileNumber = 0
		} else {
			l.logFileNumber++
		}

		l.openFile()
		l.writeLogsToFile(formattedMessage)
	}
}

func (l *Gogger) writeLogsToConsole(formattedMessage string) {
	fmt.Println(formattedMessage)
}

func (l *Gogger) addCurrentFiles() {
	var folder string
	if l.pathFolder == "" {
		folder = "."
	} else {
		folder = l.pathFolder
	}

	files, err := filepath.Glob(filepath.Join(folder, fmt.Sprintf("*%s", l.filename)))
	if err != nil {
		fmt.Printf("Error reading files in directory: %v\n", err)
		return
	}

	for _, file := range files {
		l.logQueueFiles = append(l.logQueueFiles, file)
	}
}

func (l *Gogger) createFolder() error {
	if l.pathFolder != "" {
		if _, err := os.Stat(l.pathFolder); os.IsNotExist(err) {
			if err := os.Mkdir(l.pathFolder, 0755); err != nil {
				return fmt.Errorf("Failed to create folder: %v", err)
			}
		}
	}
	return nil
}

func (l *Gogger) getCountOfLines() int {

	filePath := fmt.Sprintf("#%d%s", l.logFileNumber, l.filename)

	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var lineCount int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	return lineCount
}

func isValidFilename(filename string) bool {
	pattern := regexp.MustCompile(`[a-zA-Z0-9_]+\.(txt|log)`)
	return pattern.MatchString(filename)
}

func isValidPathFolder(pathFolder string) bool {
	if pathFolder == "" {
		return true
	}
	pattern := regexp.MustCompile(`.+[a-zA-Z0-9_]`)
	return pattern.MatchString(pathFolder)
}

func getFormattedTimestamp() string {
	return time.Now().Format("02-01-2006 15:04:05")
}

func getLogLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return ""
	}
}

func replacePlaceholder(format, placeholder, value string) string {
	return strings.ReplaceAll(format, placeholder, value)
}
