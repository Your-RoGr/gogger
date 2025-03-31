package gogger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitGogger(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		filename   string
		pathFolder string
		maxEntries int
		maxFiles   int
		wantPanic  bool
	}{
		{
			name:       "Invalid filename",
			filename:   "test/log",
			pathFolder: tempDir,
			maxEntries: 100,
			maxFiles:   5,
			wantPanic:  true,
		},
		{
			name:       "Invalid max entries",
			filename:   "test.log",
			pathFolder: tempDir,
			maxEntries: 0,
			maxFiles:   5,
			wantPanic:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("InitGogger() panic = %v, wantPanic = %v", r, tt.wantPanic)
				}
				if r != nil {
					panicMsg := r.(string)
					if !strings.Contains(panicMsg, "invalid") {
						t.Errorf("Unexpected panic message: %s", panicMsg)
					}
				}
			}()

			InitGogger(tt.filename, tt.pathFolder, tt.maxEntries, tt.maxFiles)

			if tt.wantPanic {
				t.Errorf("InitGogger() did not panic as expected")
			}
		})
	}
}

func TestNewGogger(t *testing.T) {
	tempDir := t.TempDir()

	testCases := []struct {
		name       string
		filename   string
		pathFolder string
		maxEntries int
		maxFiles   int
		wantErr    bool
	}{
		{
			name:       "Valid input",
			filename:   "test.log",
			pathFolder: tempDir,
			maxEntries: 100,
			maxFiles:   5,
			wantErr:    false,
		},
		{
			name:       "Invalid filename",
			filename:   "test/log",
			pathFolder: tempDir,
			maxEntries: 100,
			maxFiles:   5,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGogger(tc.filename, tc.pathFolder, tc.maxEntries, tc.maxFiles)

			if tc.wantErr {
				if err == nil {
					t.Errorf("NewGogger() error = nil, wantErr %v", tc.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("NewGogger() error = %v, wantErr %v", err, tc.wantErr)
				}
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	Logger.SetLogLevel(DEBUG)
	defer Logger.Close()

	tests := []struct {
		level   LogLevel
		message string
	}{
		{DEBUG, "Debug message"},
		{INFO, "Info message"},
		{WARNING, "Warning message"},
		{ERROR, "Error message"},
	}

	for _, tt := range tests {
		Logger.Log(tt.level, tt.message)
	}

	logFile := filepath.Join(tempDir, "#0test.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) != 5 { // 4 log entries + 1 empty line
		t.Errorf("Expected 5 lines in log file, got %d", len(lines))
	}

	for i, tt := range tests {

		if !strings.Contains(lines[i], getLogLevelString(tt.level)) {
			t.Errorf("Line %d does not contain expected log level: %s", i, getLogLevelString(tt.level))
		}

		if !strings.Contains(lines[i], tt.message) {
			t.Errorf("Line %d does not contain expected message: %s", i, tt.message)
		}
	}
}

func TestSetLogLevel(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	// Logger.SetLogLevel
	Logger.SetLogLevel(DEBUG)

	if Logger.logLevelConsole != DEBUG {
		t.Errorf("Expected console log level DEBUG, got %v", Logger.logLevelConsole)
	}

	if Logger.logLevelFile != DEBUG {
		t.Errorf("Expected file log level DEBUG, got %v", Logger.logLevelFile)
	}

	Logger.SetLogLevel(INFO)

	if Logger.logLevelConsole != INFO {
		t.Errorf("Expected console log level INFO, got %v", Logger.logLevelConsole)
	}

	if Logger.logLevelFile != INFO {
		t.Errorf("Expected file log level INFO, got %v", Logger.logLevelFile)
	}

	Logger.SetLogLevel(ERROR)

	if Logger.logLevelConsole != ERROR {
		t.Errorf("Expected console log level ERROR, got %v", Logger.logLevelConsole)
	}

	if Logger.logLevelFile != ERROR {
		t.Errorf("Expected file log level ERROR, got %v", Logger.logLevelFile)
	}

	Logger.SetLogLevel(WARNING)

	if Logger.logLevelConsole != WARNING {
		t.Errorf("Expected console log level WARNING, got %v", Logger.logLevelConsole)
	}

	if Logger.logLevelFile != WARNING {
		t.Errorf("Expected file log level WARNING, got %v", Logger.logLevelFile)
	}

	// Logger.SetLogLevelConsole
	Logger.SetLogLevelConsole(ERROR)

	if Logger.logLevelConsole != ERROR {
		t.Errorf("Expected console log level ERROR, got %v", Logger.logLevelConsole)
	}

	Logger.SetLogLevelConsole(INFO)

	if Logger.logLevelConsole != INFO {
		t.Errorf("Expected console log level INFO, got %v", Logger.logLevelConsole)
	}
	Logger.SetLogLevelConsole(DEBUG)

	if Logger.logLevelConsole != DEBUG {
		t.Errorf("Expected console log level DEBUG, got %v", Logger.logLevelConsole)
	}

	Logger.SetLogLevelConsole(WARNING)

	if Logger.logLevelConsole != WARNING {
		t.Errorf("Expected console log level WARNING, got %v", Logger.logLevelConsole)
	}

	// Logger.SetLogLevelFile
	Logger.SetLogLevelFile(ERROR)

	if Logger.logLevelFile != ERROR {
		t.Errorf("Expected file log level ERROR, got %v", Logger.logLevelFile)
	}

	Logger.SetLogLevelFile(INFO)

	if Logger.logLevelFile != INFO {
		t.Errorf("Expected file log level INFO, got %v", Logger.logLevelFile)
	}
	Logger.SetLogLevelFile(DEBUG)

	if Logger.logLevelFile != DEBUG {
		t.Errorf("Expected file log level DEBUG, got %v", Logger.logLevelFile)
	}

	Logger.SetLogLevelFile(WARNING)

	if Logger.logLevelFile != WARNING {
		t.Errorf("Expected file log level WARNING, got %v", Logger.logLevelFile)
	}
}

func TestSetLogFormat(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	newFormat := "<%timestamp%> <%level%>: <%message%>"
	err := Logger.SetLogFormat(newFormat)
	if err != nil {
		t.Errorf("SetLogFormat returned unexpected error: %v", err)
	}

	if Logger.logFormat != newFormat {
		t.Errorf("Expected log format '%s', got '%s'", newFormat, Logger.logFormat)
	}

	invalidFormat := "Invalid format"
	err = Logger.SetLogFormat(invalidFormat)
	if err == nil {
		t.Error("SetLogFormat should return an error for invalid format")
	}
}

func TestSetUseConsoleLog(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	Logger.SetUseConsoleLog(false)
	if Logger.console {
		t.Error("Expected console logging to be disabled")
	}

	Logger.SetUseConsoleLog(true)
	if !Logger.console {
		t.Error("Expected console logging to be enabled")
	}
}

func TestSetUseFileLog(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	Logger.SetUseFileLog(false)
	if Logger.file {
		t.Error("Expected file logging to be disabled")
	}

	Logger.SetUseFileLog(true)
	if !Logger.file {
		t.Error("Expected file logging to be enabled")
	}
}

func TestSetFilename(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	newFilename := "new_test.log"
	newPathFolder := filepath.Join(tempDir, "new_folder")
	newMaxEntries := 200

	err := Logger.SetFilename(newFilename, newPathFolder, newMaxEntries)
	if err != nil {
		t.Errorf("SetFilename returned unexpected error: %v", err)
	}

	if Logger.filename != newFilename {
		t.Errorf("Expected filename '%s', got '%s'", newFilename, Logger.filename)
	}

	if Logger.pathFolder != newPathFolder {
		t.Errorf("Expected pathFolder '%s', got '%s'", newPathFolder, Logger.pathFolder)
	}

	if Logger.maxEntries != newMaxEntries {
		t.Errorf("Expected maxEntries %d, got %d", newMaxEntries, Logger.maxEntries)
	}

	if _, err := os.Stat(newPathFolder); os.IsNotExist(err) {
		t.Errorf("New folder '%s' was not created", newPathFolder)
	}
}

func TestSetMaxEntries(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	newMaxEntries := 200
	Logger.SetMaxEntries(newMaxEntries)

	if Logger.maxEntries != newMaxEntries {
		t.Errorf("Expected maxEntries %d, got %d", newMaxEntries, Logger.maxEntries)
	}

	if Logger.maxEntriesCounter != newMaxEntries {
		t.Errorf("Expected maxEntriesCounter %d, got %d", newMaxEntries, Logger.maxEntriesCounter)
	}
}

func TestSetMaxFiles(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	newMaxFiles := 10
	Logger.SetMaxFiles(newMaxFiles)

	if Logger.maxFiles != newMaxFiles {
		t.Errorf("Expected maxFiles %d, got %d", newMaxFiles, Logger.maxFiles)
	}
}

func TestFileRotation(t *testing.T) {
	tempDir := t.TempDir()
	maxEntries := 5
	maxFiles := 3
	InitGogger("test.log", tempDir, maxEntries, maxFiles)
	defer Logger.Close()

	// Write enough logs to cause file rotation
	for i := 0; i < maxEntries*maxFiles+1; i++ {
		Logger.Info(fmt.Sprintf("Log entry %d", i))
	}

	// Check if the correct number of log files were created
	files, err := filepath.Glob(filepath.Join(tempDir, "#*test.log"))
	if err != nil {
		t.Fatalf("Failed to list log files: %v", err)
	}

	if len(files) != maxFiles {
		t.Errorf("Expected %d log files, got %d", maxFiles, len(files))
	}

	for i := 0; i < maxEntries+1; i++ {
		Logger.Info(fmt.Sprintf("Log entry %d", i))
	}

	// Check if the oldest file was deleted
	oldestFile := filepath.Join(tempDir, "#0test.log")
	if _, err := os.Stat(oldestFile); !os.IsNotExist(err) {
		t.Errorf("Expected oldest file '%s' to be deleted", oldestFile)
	}
}

func TestLoggerClose(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)

	Logger.Info("Test log")
	Logger.Close()

	// Attempt to write after closing
	Logger.Info("This should not be logged")

	logFile := filepath.Join(tempDir, "#0test.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Test log") {
		t.Error("Log file does not contain the expected log entry")
	}

	if strings.Contains(string(content), "This should not be logged") {
		t.Error("Log file contains an entry that should not have been logged after closing")
	}
}

func TestLogLevelFiltering(t *testing.T) {
	tempDir := t.TempDir()
	InitGogger("test.log", tempDir, 100, 5)
	defer Logger.Close()

	Logger.SetLogLevelFile(WARNING)
	Logger.SetLogLevelConsole(ERROR)

	Logger.Debug("Debug message")
	Logger.Info("Info message")
	Logger.Warning("Warning message")
	Logger.Error("Error message")

	logFile := filepath.Join(tempDir, "#0test.log")
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	if strings.Contains(logContent, "Debug message") || strings.Contains(logContent, "Info message") {
		t.Error("Log file contains DEBUG or INFO messages when log level is set to WARNING")
	}

	if !strings.Contains(logContent, "Warning message") || !strings.Contains(logContent, "Error message") {
		t.Error("Log file does not contain WARNING or ERROR messages")
	}
}

func TestIsValidPathFolder(t *testing.T) {
	tests := []struct {
		name       string
		pathFolder string
		want       bool
	}{
		{"Empty string", "", true},
		{"Valid path with letters", "path/to/folder", true},
		{"Valid path with numbers", "folder123", true},
		{"Valid path with underscore", "my_folder", true},
		{"Valid path with mixed characters", "path/to/my_folder123", true},
		{"Invalid path with only special characters", "@#$%", false},
		{"Invalid path with Unicode characters", "路径/到/文件夹", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPathFolder(tt.pathFolder); got != tt.want {
				t.Errorf("isValidPathFolder(%q) = %v, want %v", tt.pathFolder, got, tt.want)
			}
		})
	}
}

func TestGetCountOfLines(t *testing.T) {
	// Setup
	logger, err := NewGogger("test.log", "", 10, 5)
	if err != nil {
		t.Fatalf("Failed to create Gogger instance: %v", err)
	}
	defer logger.Close()
	defer os.Remove("#0test.log")

	tests := []struct {
		name          string
		linesToWrite  int
		expectedCount int
		setupFunc     func()
		teardownFunc  func()
	}{
		{
			name:          "Empty file",
			linesToWrite:  0,
			expectedCount: 0,
		},
		{
			name:          "Single line",
			linesToWrite:  1,
			expectedCount: 1,
		},
		{
			name:          "Multiple lines",
			linesToWrite:  5,
			expectedCount: 5,
		},
		{
			name:          "Max lines",
			linesToWrite:  10,
			expectedCount: 10,
		},
		{
			name:          "Non-existent file",
			linesToWrite:  0,
			expectedCount: 0,
			setupFunc: func() {
				logger.filename = "nonexistent.log"
			},
			teardownFunc: func() {
				logger.filename = "test.log"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			// Write lines to the file
			for i := 0; i < tt.linesToWrite; i++ {
				logger.Info("Test log line")
			}

			// Get the count of lines
			count := logger.getCountOfLines()

			// Check the result
			if count != tt.expectedCount {
				t.Errorf("getCountOfLines() = %d, want %d", count, tt.expectedCount)
			}

			// Teardown
			if tt.teardownFunc != nil {
				tt.teardownFunc()
			}

			// Reset the file for the next test
			os.Truncate("#0test.log", 0)
			logger.maxEntriesCounter = logger.maxEntries
		})
	}
}

func TestOpenFile(t *testing.T) {
	testDir := t.TempDir()
	// Setup
	// testDir := "test_logs"
	// err := os.Mkdir(testDir, 0755)
	// if err != nil {
	// 	t.Fatalf("Failed to create test directory: %v", err)
	// }
	// defer os.RemoveAll(testDir)

	// Test cases
	testCases := []struct {
		name       string
		filename   string
		maxEntries int
		maxFiles   int
		setup      func(*Gogger)
		check      func(*testing.T, *Gogger)
	}{
		{
			name:       "New file creation",
			filename:   "test.log",
			maxEntries: 100,
			maxFiles:   3,
			setup:      func(l *Gogger) {},
			check: func(t *testing.T, l *Gogger) {
				if l.fileStream == nil {
					t.Error("File stream is nil")
				}
				if len(l.logQueueFiles) != 1 {
					t.Errorf("Expected 1 file in queue, got %d", len(l.logQueueFiles))
				}
			},
		},
		{
			name:       "Existing file opened",
			filename:   "existing.log",
			maxEntries: 100,
			maxFiles:   3,
			setup: func(l *Gogger) {
				filePath := filepath.Join(testDir, "#0existing.log")
				_, err := os.Create(filePath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			},
			check: func(t *testing.T, l *Gogger) {
				if l.fileStream == nil {
					t.Error("File stream is nil")
				}
				if len(l.logQueueFiles) != 1 {
					t.Errorf("Expected 1 file in queue, got %d", len(l.logQueueFiles))
				}
			},
		},
		{
			name:       "Max files reached",
			filename:   "max.log",
			maxEntries: 10,
			maxFiles:   2,
			setup: func(l *Gogger) {
				for i := 0; i < 3; i++ {
					filePath := filepath.Join(testDir, fmt.Sprintf("#%dmax.log", i))
					f, err := os.Create(filePath)
					if err != nil {
						t.Fatalf("Failed to create test file: %v", err)
					}
					for j := 0; j < 10; j++ {
						f.WriteString("Test log entry\n")
					}
					f.Close()
				}
			},
			check: func(t *testing.T, l *Gogger) {
				if l.fileStream == nil {
					t.Error("File stream is nil")
				}
				if len(l.logQueueFiles) != 0 {
					t.Errorf("Expected 0 files in queue, got %d", len(l.logQueueFiles))
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := &Gogger{
				filename:          tc.filename,
				pathFolder:        testDir,
				maxEntries:        tc.maxEntries,
				maxEntriesCounter: tc.maxEntries,
				maxFiles:          tc.maxFiles,
			}

			tc.setup(l)
			l.openFile()
			tc.check(t, l)

			// Cleanup
			if l.fileStream != nil {
				l.fileStream.Close()
			}
		})
	}
}

func TestDeleteAllFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gogger_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cases
	testCases := []struct {
		name     string
		setup    func(*Gogger)
		expected int
	}{
		{
			name: "No files to delete",
			setup: func(l *Gogger) {
				// No setup needed
			},
			expected: 0,
		},
		{
			name: "Delete single file",
			setup: func(l *Gogger) {
				createTestFile(t, tempDir, "#0test.log")
				l.logQueueFiles = append(l.logQueueFiles, filepath.Join(tempDir, "#0test.log"))
			},
			expected: 0,
		},
		{
			name: "Delete multiple files",
			setup: func(l *Gogger) {
				createTestFile(t, tempDir, "#0test.log")
				createTestFile(t, tempDir, "#1test.log")
				createTestFile(t, tempDir, "#2test.log")
				l.logQueueFiles = append(l.logQueueFiles,
					filepath.Join(tempDir, "#0test.log"),
					filepath.Join(tempDir, "#1test.log"),
					filepath.Join(tempDir, "#2test.log"),
				)
			},
			expected: 0,
		},
		{
			name: "Delete non-existent files",
			setup: func(l *Gogger) {
				l.logQueueFiles = append(l.logQueueFiles,
					filepath.Join(tempDir, "nonexistent1.log"),
					filepath.Join(tempDir, "nonexistent2.log"),
				)
			},
			expected: 0,
		},
		{
			name: "Delete mix of existing and non-existing files",
			setup: func(l *Gogger) {
				createTestFile(t, tempDir, "#0test.log")
				createTestFile(t, tempDir, "#1test.log")
				l.logQueueFiles = append(l.logQueueFiles,
					filepath.Join(tempDir, "#0test.log"),
					filepath.Join(tempDir, "nonexistent.log"),
					filepath.Join(tempDir, "#1test.log"),
				)
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new Gogger instance for each test case
			l := &Gogger{
				filename:   "test.log",
				pathFolder: tempDir,
			}

			// Setup the test case
			tc.setup(l)

			// Call the deleteAllFiles function
			l.deleteAllFiles()

			// Check if all files were deleted
			if len(l.logQueueFiles) != tc.expected {
				t.Errorf("Expected %d files remaining, but got %d", tc.expected, len(l.logQueueFiles))
			}

			// Check if the files were actually deleted from the file system
			for _, file := range l.logQueueFiles {
				if _, err := os.Stat(file); !os.IsNotExist(err) {
					t.Errorf("File %s was not deleted from the file system", file)
				}
			}
		})
	}
}

// Helper function to create test files
func createTestFile(t *testing.T, dir, filename string) {
	filepath := filepath.Join(dir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filepath, err)
	}
	defer file.Close()
}

// Extra tests
// func TestLoggerPerformance(t *testing.T) {
// 	tempDir := t.TempDir()
// 	InitGogger("perf.log", tempDir, 10000, 5)
// 	defer Logger.Close()

// 	numLogs := 10000
// 	start := time.Now()

// 	for i := 0; i < numLogs; i++ {
// 		Logger.Info(fmt.Sprintf("Performance test log entry %d", i))
// 	}

// 	duration := time.Since(start)
// 	logsPerSecond := float64(numLogs) / duration.Seconds()

// 	t.Logf("Logged %d entries in %v (%.2f logs/second)", numLogs, duration, logsPerSecond)

// 	if logsPerSecond < 10000 {
// 		t.Errorf("Logger performance below expected threshold: %.2f logs/second", logsPerSecond)
// 	}
// }

// func TestLoggerWithLargeMessages(t *testing.T) {
// 	tempDir := t.TempDir()
// 	InitGogger("large.log", tempDir, 100, 5)
// 	defer Logger.Close()

// 	largeMessage := strings.Repeat("A", 1024*1024) // 1MB message

// 	Logger.Info(largeMessage)

// 	logFile := filepath.Join(tempDir, "#0large.log")
// 	info, err := os.Stat(logFile)
// 	if err != nil {
// 		t.Fatalf("Failed to get log file info: %v", err)
// 	}

// 	if info.Size() < int64(len(largeMessage)) {
// 		t.Errorf("Log file size (%d) is smaller than the large message size (%d)", info.Size(), len(largeMessage))
// 	}
// }
