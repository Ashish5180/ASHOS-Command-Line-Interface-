package task

import (
	"context"
	"errors"

	"ashos/internal/storage"
)

// StoreRepository bridges task.Repository ↔ storage.Store.
//
// Kyon zaruri hai?
// task.Service ko sirf Repository interface chahiye (GetAll, SaveAll).
// Usse pata nahi hona chahiye ki data JSON file mein hai ya memory mein.
// Yeh struct wo bridge hai — Store engine ko Repository ki shakal deta hai.
type StoreRepository struct {
	store storage.Store
}

// NewStoreRepository creates a Repository backed by any storage.Store engine.
// main.go mein call hoga: task.NewStoreRepository(sqliteStore)
func NewStoreRepository(store storage.Store) Repository {
	return &StoreRepository{store: store}
}

// GetAll fetches all tasks from the store.
// Agar pehli baar run ho raha hai (file nahi hai), toh empty slice return karega — crash nahi karega.
func (r *StoreRepository) GetAll() ([]Task, error) {
	var tasks []Task
	err := r.store.Get(context.Background(), "task_data", "all_tasks", &tasks)
	if err != nil {
		// Pehli baar koi task nahi hoga — "not found" matlab empty list
		var notFound storage.ErrNotFound
		if errors.As(err, &notFound) {
			return []Task{}, nil
		}
		return nil, err
	}
	return tasks, nil
}

// SaveAll replaces the entire task list atomically in the store.
// Ek hi key ("all_tasks") ke under saari tasks stored hain.
// Isse delete + update + insert sab ek hi Save call mein ho jaata hai.
func (r *StoreRepository) SaveAll(tasks []Task) error {
	return r.store.Save(context.Background(), "task_data", "all_tasks", tasks)
}
