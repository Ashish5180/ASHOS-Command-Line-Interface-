package note

import (
	"ashos/internal/core/event"
	"ashos/internal/storage"
	"context"
	"time"
)

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

func (s *Service) AddNote(ctx context.Context, content string) error {
	var notes []Note
	_ = s.store.Get(ctx, "journal", "entries", &notes)

	newID := 1
	for _, n := range notes {
		if n.ID >= newID {
			newID = n.ID + 1
		}
	}

	note := Note{
		ID:        newID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	notes = append(notes, note)
	if err := s.store.Save(ctx, "journal", "entries", notes); err != nil {
		return err
	}

	s.bus.Publish(event.NoteCreated{
		ID:        note.ID,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
	})

	return nil
}

func (s *Service) ListNotes(ctx context.Context) ([]Note, error) {
	var notes []Note
	err := s.store.Get(ctx, "journal", "entries", &notes)
	return notes, err
}
