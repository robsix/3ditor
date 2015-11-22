package golog

import(
	`time`
)

// Logs to stdout, no history is kept so all calls to get() and getById() will return errors
func newConsoleLog(printToStdOut bool, lineSpacing int) Log {
	put := func(le LogEntry){}
	getById := func(logId string) (LogEntry, error) {return LogEntry{}, &consoleLogNoStorageError{}}
	get := func(before time.Time, level level, limit int) ([]LogEntry, error) {return nil, &consoleLogNoStorageError{}}
	return NewLog(put, getById, get, printToStdOut, lineSpacing)
}

type consoleLogNoStorageError struct{}
func (e *consoleLogNoStorageError) Error() string {return `ConsoleLog only prints to stdout, it does not store any log data`}

func NewConsoleLog(lineSpacing int) Log {
	return newConsoleLog(true, lineSpacing)
}

func NewDevNullLog() Log {
	return newConsoleLog(false, 0)
}
