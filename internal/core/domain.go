package core

import "time"

type TaskDomain struct {
	ID          int
	Title       string
	Description string
	Completed   bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}

type TaskCompleteDomain struct {
	ID          int
	Completed   bool
	CompletedAt *time.Time
}
