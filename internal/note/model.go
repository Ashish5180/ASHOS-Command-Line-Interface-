package note

import "time"

type Note struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
