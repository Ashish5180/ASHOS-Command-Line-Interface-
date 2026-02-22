package task

import (
	"fmt"
	"strings"
	"time"
)

// Service is the ONLY gateway to task operations.
// CLI commands call Service → Service validates → Service persists via Repository.
// CLI kabhi directly Repository ko touch nahi karega.
type Service struct {
	repo Repository
}

// NewService creates a task service with injected repository.
// Repository kya hai (JSON? Memory?) — Service ko pata nahi, na hi pata hona chahiye.
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// ---------------------------------------------------------------------------
// Business Methods — Validation + Logic + Persistence
// ---------------------------------------------------------------------------

// AddTask validates input, generates ID, and persists a new task.
//
// Kyon Service mein hai? Kyunki:
// 1. Title validation (empty/too-long) = business rule
// 2. ID generation = invariant (unique hona chahiye)
// 3. Persistence = repo ke through, direct nahi
func (s *Service) AddTask(title string) (int, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, fmt.Errorf("task title cannot be empty")
	}
	if len(title) > 200 {
		return 0, fmt.Errorf("task title too long (max 200 characters)")
	}

	tasks, err := s.repo.GetAll()
	if err != nil {
		return 0, fmt.Errorf("failed to load tasks: %w", err)
	}

	// Auto-generate next ID — scan existing tasks for max ID
	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	task := Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	tasks = append(tasks, task)
	if err := s.repo.SaveAll(tasks); err != nil {
		return 0, fmt.Errorf("failed to save task: %w", err)
	}

	return nextID, nil
}

// ListTasks returns all tasks (completed + pending).
// Read-only operation — koi mutation nahi hoti.
func (s *Service) ListTasks() ([]Task, error) {
	return s.repo.GetAll()
}

// CompleteTask marks a task as done.
//
// Invariants protected:
// 1. Task exist karna chahiye (invalid ID → error)
// 2. Already completed task dobara complete nahi ho sakta
func (s *Service) CompleteTask(id int) error {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			if tasks[i].Completed {
				return fmt.Errorf("task #%d is already completed", id)
			}
			tasks[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task #%d not found", id)
	}

	return s.repo.SaveAll(tasks)
}

// DeleteTask removes a task permanently.
//
// Protection: Invalid ID pe error — silent fail nahi hoga.
func (s *Service) DeleteTask(id int) error {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	filtered := make([]Task, 0, len(tasks))
	found := false
	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue // skip = delete
		}
		filtered = append(filtered, t)
	}

	if !found {
		return fmt.Errorf("task #%d not found", id)
	}

	return s.repo.SaveAll(filtered)
}

// PendingCount returns the number of incomplete tasks.
func (s *Service) PendingCount() (int, error) {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, t := range tasks {
		if !t.Completed {
			count++
		}
	}
	return count, nil
}
