package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

//Logger 用于打印调试日志
type Logger interface {
	Errorf(format string, v ...interface{}) //错误日志打印
	Infof(format string, v ...interface{})  //进程日志打印
	Debugf(format string, v ...interface{}) //调试日志打印
	Print(args ...interface{})              //打印错误日志
	Printf(format string, v ...interface{}) //打印错误日志
}

//LogLevel 日志级别, 为调试/信息/错误
type LogLevel uint8

//日志级别
const (
	DebugLevel LogLevel = iota //调试
	InfoLevel                  //信息
	ErrorLevel                 //错误
)

type defaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

func newNilLogger() Logger {
	d := &defaultLogger{
		level:  ErrorLevel,
		logger: log.New(os.Stderr, "[etl]", log.Lmicroseconds|log.LstdFlags|log.Lshortfile),
	}
	return d
}

//NewDefaultLogger 生成一个日志打印Logger，level可以是DebugLevel，InfoLevel，ErrorLevel
func NewDefaultLogger(writer io.Writer, level LogLevel, prefix string) Logger {
	d := &defaultLogger{
		level:  level,
		logger: log.New(writer, prefix, log.Lmicroseconds|log.LstdFlags|log.Lshortfile),
	}
	return d
}

func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	if d.level <= ErrorLevel {
		d.logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (d *defaultLogger) Infof(format string, args ...interface{}) {
	if d.level <= InfoLevel {
		d.logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (d *defaultLogger) Debugf(format string, args ...interface{}) {
	if d.level <= DebugLevel {
		d.logger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (d *defaultLogger) Print(args ...interface{}) {
	d.logger.Output(2, fmt.Sprint(args...))
}

func (d *defaultLogger) Printf(format string, v ...interface{}) {
	d.logger.Output(2, fmt.Sprintf(format, v...))
}

var (
	lw   = loggerWrapper{l: newNilLogger()}
	_log = lw.logger()
)

type loggerWrapper struct {
	l  Logger
	mu sync.RWMutex
}

func (l *loggerWrapper) setLogger(logger Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l = logger
}

func (l *loggerWrapper) logger() Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.l
}

//SetLogger 设置一个符合Logger日志来打印调试信息
func SetLogger(logger Logger) {
	lw.setLogger(logger)
	_log = lw.logger()
}

func GetLogger() Logger {
	return _log
}
