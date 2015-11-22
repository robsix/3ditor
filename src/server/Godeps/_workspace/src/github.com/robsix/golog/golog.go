/*
A very simple logging utility
*/
package golog

import (
	"fmt"
	"github.com/robsix/3ditor/src/server/Godeps/_workspace/src/code.google.com/p/go-uuid/uuid"
	ct "github.com/robsix/3ditor/src/server/Godeps/_workspace/src/github.com/daviddengcn/go-colortext"
	"time"
)

type level string

const (
	ANY      = level(`ANY`)
	INFO     = level(`INFO`)
	WARNING  = level(`WARNING`)
	ERROR    = level(`ERROR`)
	CRITICAL = level(`CRITICAL`)
)

var (
	logEntriesConsoleChannel = make(chan LogEntry)
	consoleWorkerCreated     = false
)

type LogEntry struct {
	LogId   string    `json:"logId"`
	Time    time.Time `json:"time"`
	Level   level     `json:"level"`
	Message string    `json:"message"`
}

type Log interface {
	Info(a ...interface{}) LogEntry
	Warning(a ...interface{}) LogEntry
	Error(a ...interface{}) LogEntry
	Critical(a ...interface{}) LogEntry
	GetById(logId string) (LogEntry, error)
	Get(before time.Time, level level, limit int) ([]LogEntry, error)
}

type Put func(le LogEntry)
type GetById func(logId string) (LogEntry, error)
type Get func(before time.Time, level level, limit int) ([]LogEntry, error)

func NewLog(put Put, getById GetById, get Get, printToStdOut bool, lineSpacing int) Log {
	if !consoleWorkerCreated {
		go func() {
			for le := range logEntriesConsoleChannel {
				var levelPadding string
				switch le.Level {
				case INFO:
					levelPadding = `    `
					ct.Foreground(ct.Cyan, true)
				case WARNING:
					levelPadding = ` `
					ct.Foreground(ct.Yellow, true)
				case ERROR:
					levelPadding = `   `
					ct.Foreground(ct.Red, true)
				case CRITICAL:
					levelPadding = ``
					ct.ChangeColor(ct.Black, true, ct.Red, true)
				}
				fmt.Println(le.Time.Format(`15:04:05.00`), string(le.Level)+levelPadding, le.Message)
				ct.ResetColor()
				for i := 0; i < lineSpacing; i++ {
					fmt.Println(``)
				}
			}
		}()
		consoleWorkerCreated = true
	}
	return &log{
		put:                      put,
		getById:                  getById,
		get:                      get,
		printToStdOut:            printToStdOut,
		logEntriesConsoleChannel: logEntriesConsoleChannel,
	}
}

type log struct {
	put                      Put
	getById                  GetById
	get                      Get
	printToStdOut            bool
	logEntriesConsoleChannel chan LogEntry
}

func (l *log) log(level level, a ...interface{}) LogEntry {
	le := LogEntry{
		LogId:   uuid.New(),
		Time:    time.Now().UTC(),
		Level:   level,
		Message: fmt.Sprint(a...),
	}
	if l.printToStdOut {
		l.logEntriesConsoleChannel <- le
	}
	l.put(le)
	return le
}

func (l *log) Info(a ...interface{}) LogEntry {
	return l.log(INFO, a...)
}

func (l *log) Warning(a ...interface{}) LogEntry {
	return l.log(WARNING, a...)
}

func (l *log) Error(a ...interface{}) LogEntry {
	return l.log(ERROR, a...)
}

func (l *log) Critical(a ...interface{}) LogEntry {
	return l.log(CRITICAL, a...)
}

func (l *log) GetById(logId string) (LogEntry, error) {
	return l.getById(logId)
}

func (l *log) Get(before time.Time, level level, limit int) ([]LogEntry, error) {
	return l.get(before, level, limit)
}
