package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/elias/axiom/engine"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGrey   = "\033[90m"
)

type LogMessage struct {
	Time     time.Time
	Tick     int64
	Color    string
	Level    string
	Filename string
	Message  string
}

type logger struct {
	queue    chan LogMessage
	filepath string
	tick     *engine.Tick
}

func (l *logger) run() {
	log.SetOutput(
		&lumberjack.Logger{
			Filename:   l.filepath,
			MaxSize:    30,
			MaxBackups: 3,
			MaxAge:     14,
			Compress:   true,
		},
	)
	log.SetFlags(0)

	for msg := range l.queue {
		timeTick := fmt.Sprintf("%s%s%s", colorGrey, msg.Time.Format("15:04:05.000"), colorReset)
		log.Printf("%s %v %s[%s] %s %s%s\n", timeTick, l.tick.Tick(), msg.Color, msg.Level, msg.Filename, msg.Message, colorReset)
	}
}

func Init(filepath string, tick *engine.Tick) {
	l := &logger{
		queue:    make(chan LogMessage, 10),
		filepath: filepath,
		tick:     tick,
	}

	go l.run()
	_logger = l
}

func (l *logger) log(level, color, format string, args ...any) {
	_, file, _, ok := runtime.Caller(2)
	fileName := "UNKNOWN"
	if ok {
		fileName = strings.ToUpper(filepath.Base(file))
	}

	msg := fmt.Sprintf(format, args...)
	l.queue <- LogMessage{
		Time:     time.Now(),
		Tick:     l.tick.Tick(),
		Level:    level,
		Filename: fileName,
		Message:  msg,
		Color:    color,
	}
}

var _logger *logger

func Info(format string, args ...any) {
	_logger.log("TELEMETRY", colorReset, format, args...)
}
func Error(format string, args ...any) {
	_logger.log("FAULT", colorRed, format, args...)
}
func Warn(format string, args ...any) {
	_logger.log("WARN", colorYellow, format, args...)
}
func Ok(format string, args ...any) {
	_logger.log("STABLE", colorGreen, format, args...)
}
