package pkg

import (
	"bytes"
	"strings"
	"time"
)

// LogEntry represents a single log entry with a timestamp.
type LogEntry struct {
	// Timestamp is the time when the log entry was created (from Docker).
	Timestamp time.Time `json:"timestamp"`
	// Content is the actual log message.
	Content string `json:"content"`
}

// LogsWriter is a custom writer that processes log data and sends it to a channel.
type LogsWriter struct {
	channel chan<- LogEntry
	buffer  bytes.Buffer
}

// NewLogsWriter creates a new LogsWriter that sends log entries to the provided channel.
func NewLogsWriter(channel chan<- LogEntry) *LogsWriter {
	return &LogsWriter{channel: channel}
}

// parseLine parses a single line of log data into a LogEntry.
func parseLine(line string) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, "\n")

	var content string
	var timestamp time.Time
	var err error

	lineEntries := strings.Split(line, " ")
	if len(lineEntries) == 1 {
		timestamp = time.Now().UTC()
		content = "\n"
	} else if len(lineEntries) >= 2 {
		timestamp, err = time.Parse(time.RFC3339Nano, lineEntries[0])
		// TODO: do we really want to return an error here, or just use the current time?
		if err != nil {
			return nil, err
		}
		content = strings.Join(lineEntries[1:], " ")
	} else {
		timestamp = time.Now().UTC()
		content = line
	}

	return &LogEntry{
		Timestamp: timestamp,
		Content:   content,
	}, nil
}

func (w *LogsWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)

	for {
		line, err := w.buffer.ReadString('\n')
		if err != nil {
			break
		}

		entry, err := parseLine(line)
		if err != nil {
			return 0, err
		}
		w.channel <- *entry
	}

	return len(data), nil
}

// FlushRemaining flushes any remaining data in the buffer as a log entry.
func (w *LogsWriter) FlushRemaining() error {
	if w.buffer.Len() == 0 {
		return nil
	}

	entry, err := parseLine(w.buffer.String())
	if err == nil {
		return err
	}

	w.channel <- *entry
	return nil
}
