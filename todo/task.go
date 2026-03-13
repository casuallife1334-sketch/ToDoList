package todo

import "time"

type Task struct {
	Title     string
	Text      string
	Completed bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(title string, text string) Task {
	return Task{
		Title: title,
		Text:  text,

		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	t.Completed = true

	completeTime := time.Now()
	t.CompletedAt = &completeTime
}

func (t *Task) Uncomplete() {
	t.Completed = false

	t.CompletedAt = nil
}
