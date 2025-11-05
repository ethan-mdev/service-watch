package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/sse"
)

type Logger struct {
	eventFile   *os.File
	broadcaster *sse.Broadcaster
	mutex       sync.Mutex
}

func Start(logPath string, broadcaster *sse.Broadcaster) (*Logger, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		eventFile:   logFile,
		broadcaster: broadcaster,
	}, nil
}

func (l *Logger) Info(eventType string, data map[string]interface{}) {
	l.log("INFO", eventType, data)
}

func (l *Logger) Error(eventType string, data map[string]interface{}) {
	l.log("ERROR", eventType, data)
}

func (l *Logger) log(level, eventType string, data map[string]interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	event := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": level,
		"event": eventType,
		"data":  data,
	}

	// Write to JSONL file
	json.NewEncoder(l.eventFile).Encode(event)

	// Print to console for development visibility
	fmt.Printf("[%s] %s: %v\n", level, eventType, data)

	// Broadcast to SSE clients
	if l.broadcaster != nil {
		l.broadcaster.Broadcast(core.Event{
			Type: eventType,
			Data: data,
		})
	}
}

func (l *Logger) Close() error {
	return l.eventFile.Close()
}
