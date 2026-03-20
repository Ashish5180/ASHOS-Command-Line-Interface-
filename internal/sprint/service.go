package sprint

import (
	"ashos/internal/core/event"
	"ashos/internal/storage"
	"context"
	"time"
)

type Sprint struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Summary   string        `json:"summary"`
	TasksDone int           `json:"tasks_done"`
	FocusTime time.Duration `json:"focus_time"`
	CreatedAt time.Time     `json:"created_at"`
}

type Service struct {
	store storage.Store
	bus   *event.EventBus
}

func NewService(store storage.Store, bus *event.EventBus) *Service {
	return &Service{
		store: store,
		bus:   bus,
	}
}

func (s *Service) EndSprint(ctx context.Context, title, summary string, tasksDone int, focusTime time.Duration) error {
	var sprints []Sprint
	_ = s.store.Get(ctx, "productivity", "sprints", &sprints)

	newID := 1
	for _, spr := range sprints {
		if spr.ID >= newID {
			newID = spr.ID + 1
		}
	}

	spr := Sprint{
		ID:        newID,
		Title:     title,
		Summary:   summary,
		TasksDone: tasksDone,
		FocusTime: focusTime,
		CreatedAt: time.Now(),
	}

	sprints = append(sprints, spr)
	if err := s.store.Save(ctx, "productivity", "sprints", sprints); err != nil {
		return err
	}

	s.bus.Publish(event.SprintEnded{
		ID:        spr.ID,
		Title:     spr.Title,
		Summary:   spr.Summary,
		TasksDone: spr.TasksDone,
		FocusTime: spr.FocusTime,
		CreatedAt: spr.CreatedAt,
	})

	return nil
}

func (s *Service) ListSprints(ctx context.Context) ([]Sprint, error) {
	var sprints []Sprint
	err := s.store.Get(ctx, "productivity", "sprints", &sprints)
	return sprints, err
}
