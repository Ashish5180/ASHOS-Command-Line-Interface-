package system

import (
	"ashos/internal/focus"
	"ashos/internal/storage"
	"ashos/internal/task"
	"context"
	"time"
)

// Status represents the system's current state.
type Status struct {
	Uptime       time.Duration
	PendingTasks int
	FocusToday   time.Duration
}

// Service handles system status data gathering.
type Service struct {
	startTime   time.Time
	taskService *task.Service
	store       storage.Store
}

// NewService creates a new system status service.
func NewService(taskService *task.Service, store storage.Store) *Service {
	return &Service{
		startTime:   time.Now(),
		taskService: taskService,
		store:       store,
	}
}

// GetStatus returns the current dashboard data.
func (s *Service) GetStatus() Status {
	uptime := time.Since(s.startTime).Truncate(time.Second)

	pending := 0
	if s.taskService != nil {
		tasks, _ := s.taskService.ListTasks()
		for _, t := range tasks {
			if !t.Completed {
				pending++
			}
		}
	}

	focusToday := time.Duration(0)
	var sessions []focus.SessionRecord
	err := s.store.List(context.Background(), "focus_data", &sessions)
	if err == nil {
		now := time.Now()
		for _, sess := range sessions {
			// Check if session was today
			if sess.StartTime.Year() == now.Year() &&
				sess.StartTime.YearDay() == now.YearDay() {
				focusToday += sess.Duration
			}
		}
	} else {
		// Try individual key if List fails (depending on JSONStore implementation)
		_ = s.store.Get(context.Background(), "focus_data", "sessions", &sessions)
		now := time.Now()
		for _, sess := range sessions {
			if sess.StartTime.Year() == now.Year() &&
				sess.StartTime.YearDay() == now.YearDay() {
				focusToday += sess.Duration
			}
		}
	}

	return Status{
		Uptime:       uptime,
		PendingTasks: pending,
		FocusToday:   focusToday.Truncate(time.Second),
	}
}
