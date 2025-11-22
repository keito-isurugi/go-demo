package model

import "time"

type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Note      string    `json:"note"`
    DoneFlag  bool      `json:"done_flag"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
