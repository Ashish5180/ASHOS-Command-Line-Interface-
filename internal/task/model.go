package task

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a single task
type Task struct {
	ID        int
	Title     string
	Completed bool
	CreatedAt time.Time
}

// TaskManager handles all task operations
type TaskManager struct {
	tasks  map[int]*Task
	nextID int
	mu     sync.RWMutex // Thread safety ke liye
}

// NewTaskManager creates a new task manager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// AddTask - Naya task add karta hai with validation
func (tm *TaskManager) AddTask(title string) (int, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Input validation - khali title nahi ho sakta
	if title == "" {
		return 0, fmt.Errorf("title khali nahi ho sakta")
	}

	if len(title) > 200 {
		return 0, fmt.Errorf("title 200 characters se zyada nahi ho sakta")
	}

	// ID auto-generate hota hai
	id := tm.nextID
	tm.nextID++

	tm.tasks[id] = &Task{
		ID:        id,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	return id, nil
}

// ListTasks - Sabhi tasks return karta hai
func (tm *TaskManager) ListTasks() []Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, *task)
	}

	return tasks
}

// CompleteTask - Task ko complete mark karta hai
func (tm *TaskManager) CompleteTask(id int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("task ID %d nahi mila", id)
	}

	// Invariant check - already complete task ko dobara complete nahi kar sakte
	if task.Completed {
		return fmt.Errorf("task pehle se hi complete hai")
	}

	task.Completed = true
	return nil
}

// DeleteTask - Task ko delete karta hai
func (tm *TaskManager) DeleteTask(id int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	_, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("task ID %d nahi mila", id)
	}

	delete(tm.tasks, id)
	return nil
}

// PendingCount - Kitne pending tasks hain
func (tm *TaskManager) PendingCount() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	count := 0
	for _, task := range tm.tasks {
		if !task.Completed {
			count++
		}
	}

	return count
}
