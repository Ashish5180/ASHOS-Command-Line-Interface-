package focus

import "time"

// SessionRecord is used to persist focus session data.
type SessionRecord struct {
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"`
}

// Stats represents aggregated focus data.
type Stats struct {
	TotalToday time.Duration
}
