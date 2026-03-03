package ui

import (
	"ashos/internal/focus"
	"ashos/internal/system"
	"ashos/internal/task"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	taskService   *task.Service
	systemService *system.Service
	focusManager  *focus.Manager

	tasks  []task.Task
	status system.Status
	err    error
	width  int
	height int
}

func NewModel(ts *task.Service, ss *system.Service, fm *focus.Manager) model {
	return model{
		taskService:   ts,
		systemService: ss,
		focusManager:  fm,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchTasks,
		m.fetchStatus,
		tick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			return m, tea.Batch(m.fetchTasks, m.fetchStatus)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case taskMsg:
		m.tasks = []task.Task(msg)
	case statusMsg:
		m.status = system.Status(msg)
	case tickMsg:
		return m, tea.Batch(m.fetchStatus, tick())
	case error:
		m.err = msg
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress 'q' to quit.", m.err)
	}

	// Dynamic Width Calculation
	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80 // Fallback
	}

	// Header
	header := TitleStyle.Render(" ASHOS DASHBOARD ")

	// Left Panel: Tasks
	var taskList string
	if len(m.tasks) == 0 {
		taskList = GrayStyle.Render("No pending tasks.")
	} else {
		var sb strings.Builder
		count := 0
		for i := range m.tasks {
			t := m.tasks[i]
			if !t.Completed {
				sb.WriteString(fmt.Sprintf("• %s\n", t.Title))
				count++
			}
			if count >= 8 {
				sb.WriteString(GrayStyle.Render(fmt.Sprintf("...and others")))
				break
			}
		}
		taskList = sb.String()
	}

	// Dynamic Styles based on terminal size
	taskBoxWidth := (contentWidth / 2) - 4
	if contentWidth < 60 {
		taskBoxWidth = contentWidth - 4
	}

	taskBox := MakeBox("TASKS", taskList, BoxStyle.Width(taskBoxWidth))

	// Right Panel: System Stats
	stats := fmt.Sprintf(
		"Uptime: %s\nTasks:  %d Pending\nFocus:  %s Today",
		HighlightStyle.Render(m.status.Uptime.String()),
		m.status.PendingTasks,
		HighlightStyle.Render(m.status.FocusToday.String()),
	)

	statsBoxWidth := (contentWidth / 2) - 4
	if contentWidth < 60 {
		statsBoxWidth = contentWidth - 4
	}
	statsBox := MakeBox("SYSTEM", stats, StatsBoxStyle.Width(statsBoxWidth))

	// Layout Logic
	var topRow string
	if contentWidth < 60 {
		topRow = lipgloss.JoinVertical(lipgloss.Left, taskBox, statsBox)
	} else {
		topRow = lipgloss.JoinHorizontal(lipgloss.Top, taskBox, statsBox)
	}

	// Bottom Row: Focus
	focusStatus := "Idle"
	if m.focusManager != nil && m.focusManager.IsActive() {
		focusStatus = HighlightStyle.Render("FOCUSING...")
	}
	focusBox := MakeBox("FOCUS", "Status: "+focusStatus, FocusBoxStyle.Width(contentWidth-4))

	// Footer
	footer := StatusStyle.Render("\n[q] Quit • [r] Refresh (Auto-refreshes every 5s)")

	return lipgloss.JoinVertical(lipgloss.Left, header, topRow, focusBox, footer)
}

// ---------------------------------------------------------------------------
// Messages & Commands
// ---------------------------------------------------------------------------

type taskMsg []task.Task
type statusMsg system.Status
type tickMsg time.Time

func (m model) fetchTasks() tea.Msg {
	if m.taskService == nil {
		return nil
	}
	tasks, err := m.taskService.ListTasks()
	if err != nil {
		return err
	}
	return taskMsg(tasks)
}

func (m model) fetchStatus() tea.Msg {
	if m.systemService == nil {
		return nil
	}
	return statusMsg(m.systemService.GetStatus())
}

func tick() tea.Cmd {
	return tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func RunDashboard(ts *task.Service, ss *system.Service, fm *focus.Manager) error {
	p := tea.NewProgram(NewModel(ts, ss, fm), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
