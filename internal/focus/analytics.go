package focus

import (
	"ashos/internal/storage"
	"context"
	"errors"
	"math"
	"time"
)

// AnalyticsReport holds the processed productivity data.
type AnalyticsReport struct {
	DailyStats    []DailyStat
	TotalDuration time.Duration
	Streak        int
	DeepWorkScore int
}

// DailyStat holds aggregated duration for a specific day.
type DailyStat struct {
	Date     time.Time
	Duration time.Duration
}

// GetAnalyticsReport processes focus history to generate insights.
func (m *Manager) GetAnalyticsReport(days int) (AnalyticsReport, error) {
	report := AnalyticsReport{
		DailyStats: make([]DailyStat, days),
	}

	now := time.Now()
	// Initialize daily stats for the last 'days' days
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		report.DailyStats[days-1-i] = DailyStat{Date: date}
	}

	var sessions []SessionRecord
	err := m.store.Get(context.Background(), "focus_data", "sessions", &sessions)
	if err != nil {
		var notFound storage.ErrNotFound
		if errors.As(err, &notFound) {
			// No sessions yet, return empty report with zero stats
			return report, nil
		}
		return AnalyticsReport{}, err
	}

	// Aggregate durations by day
	for _, sess := range sessions {
		for i := range report.DailyStats {
			if sess.StartTime.Year() == report.DailyStats[i].Date.Year() &&
				sess.StartTime.YearDay() == report.DailyStats[i].Date.YearDay() {
				report.DailyStats[i].Duration += sess.Duration
				report.TotalDuration += sess.Duration
			}
		}
	}

	// Calculate Streak
	report.Streak = calculateStreak(sessions)

	// Calculate Deep Work Score (0-100)
	// Base: Avg focus time per day (Target 4h/day = 100%)
	avgSeconds := report.TotalDuration.Seconds() / float64(days)
	targetSeconds := 4.0 * 3600.0
	baseScore := (avgSeconds / targetSeconds) * 80.0 // 80 points for duration

	// Consistency bonus: Streak points
	streakBonus := math.Min(float64(report.Streak)*4.0, 20.0) // Max 20 points for streak

	report.DeepWorkScore = int(math.Min(baseScore+streakBonus, 100))

	return report, nil
}

func calculateStreak(sessions []SessionRecord) int {
	if len(sessions) == 0 {
		return 0
	}

	// Unique days with at least 15 mins of focus
	daysSet := make(map[string]bool)
	for _, s := range sessions {
		if s.Duration >= 15*time.Minute {
			daysSet[s.StartTime.Format("2006-01-02")] = true
		}
	}

	streak := 0
	current := time.Now()

	for {
		key := current.Format("2006-01-02")
		if daysSet[key] {
			streak++
			current = current.AddDate(0, 0, -1)
		} else {
			// If it's today and we haven't hit 15m yet, check yesterday
			if key == time.Now().Format("2006-01-02") {
				current = current.AddDate(0, 0, -1)
				continue
			}
			break
		}
	}

	return streak
}
