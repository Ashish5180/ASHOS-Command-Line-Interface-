package event

import "time"

// --- Task Events ---

type TaskCreated struct {
	ID    int
	Title string
}

type TaskCompleted struct {
	ID    int
	Title string
}

type TaskDeleted struct {
	ID int
}

// --- Focus Events ---

type FocusStarted struct {
	StartTime time.Time
}

type FocusEnded struct {
	StartTime time.Time
	Duration  time.Duration
	Summary   string
}

// --- Note Events ---

type NoteCreated struct {
	ID        int
	Content   string
	CreatedAt time.Time
}

// --- Sprint Events ---

type SprintEnded struct {
	ID        int
	Title     string
	Summary   string
	TasksDone int
	FocusTime time.Duration
	CreatedAt time.Time
}

// --- System Events ---

type SystemStarted struct {
	Time time.Time
}

type SystemShutdown struct {
	Uptime time.Duration
}
