package state

import (
	"slices"
	"sync"
)

type GameLogger struct {
	logs []string
	mu   sync.Mutex
	size int
}

func NewGameLogger(size int) *GameLogger {
	return &GameLogger{size: size}
}

func (l *GameLogger) Print(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.logs) >= l.size {
		l.logs = l.logs[1:]
	}

	l.logs = append(l.logs, msg)
}

func (l *GameLogger) Println(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.logs) >= l.size {
		l.logs = l.logs[1:]
	}

	l.logs = append(l.logs, msg+"\n")
}

func (l *GameLogger) Read() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	return slices.Clone(l.logs)
}
