package types

import "time"

type ActivityLogEntry struct {
    Description string
    LogTime time.Time
}

type IActivityLog interface {
    NewLogEntry(logEntry ActivityLogEntry)
}

type ActivityLog struct {
    MaxSize int
    Entries []ActivityLogEntry
}

func (al *ActivityLog) NewLogEntry(logEntry ActivityLogEntry) {
    // Ensure a default max size if not set
    if al.MaxSize == 0 {
        defaultMaxSize := 3
        al.MaxSize = defaultMaxSize
    }

    // Append the new log entry
    al.Entries = append(al.Entries, logEntry)

    // Trim the log if it exceeds the max size
    if len(al.Entries) > al.MaxSize {
        al.Entries = al.Entries[len(al.Entries)-al.MaxSize:]
    }
}