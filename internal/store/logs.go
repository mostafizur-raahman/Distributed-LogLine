package store

import "time"

type LogRow struct {
	ID        int64
	Level     string
	Message   string
	Service   string
	Timestamp time.Time
	Data      map[string]any
	CreatedAt time.Time
}
