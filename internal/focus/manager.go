package focus

import (
	"ashos/internal/storage"
	"context"
	"fmt"
	"sync"
	"time"
)

// Session represents a focus session state.
type Session struct {
	StartTime time.Time
	Duration  time.Duration
	IsActive  bool
	mu        sync.Mutex
}

// Manager handles focus mode logic using Go concurrency.
type Manager struct {
	session *Session
	stopCh  chan bool
	store   storage.Store
}

// NewManager initializes the focus manager.
func NewManager(store storage.Store) *Manager {
	return &Manager{
		session: &Session{},
		stopCh:  make(chan bool),
		store:   store,
	}
}

// StartSession initiates a focus session in a background goroutine.
func (m *Manager) StartSession() {
	m.session.mu.Lock()
	if m.session.IsActive {
		m.session.mu.Unlock()
		return
	}
	m.session.StartTime = time.Now()
	m.session.IsActive = true
	m.session.mu.Unlock()

	// 🧠 TECHNIQUE: Goroutine & Ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		fmt.Println("🧠 Focus Session Started! Stay sharp.")

		for {
			select {
			case <-ticker.C:
				m.session.mu.Lock()
				m.session.Duration = time.Since(m.session.StartTime)
				if int(m.session.Duration.Seconds())%60 == 0 && m.session.Duration.Seconds() > 0 {
					fmt.Printf("⏱️ Focus Progress: %v\n", m.session.Duration.Truncate(time.Minute))
				}
				m.session.mu.Unlock()
			case <-m.stopCh:
				fmt.Println("\n🛑 Focus Session Stopped.")
				return
			}
		}
	}()
}

// StopSession sends a stop signal and saves the session.
func (m *Manager) StopSession() {
	m.session.mu.Lock()
	if !m.session.IsActive {
		m.session.mu.Unlock()
		return
	}
	m.session.IsActive = false
	m.session.mu.Unlock()

	m.saveSession()
	m.stopCh <- true
}

func (m *Manager) saveSession() {
	m.session.mu.Lock()
	record := SessionRecord{
		StartTime: m.session.StartTime,
		Duration:  m.session.Duration,
	}
	m.session.mu.Unlock()

	var sessions []SessionRecord
	_ = m.store.Get(context.Background(), "focus_data", "sessions", &sessions)
	sessions = append(sessions, record)
	_ = m.store.Save(context.Background(), "focus_data", "sessions", sessions)
}

// GetStats returns current session duration.
func (m *Manager) GetStats() time.Duration {
	m.session.mu.Lock()
	defer m.session.mu.Unlock()
	return m.session.Duration
}
