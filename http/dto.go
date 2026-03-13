package http

import (
	"encoding/json"
	"errors"
	"time"
)

/*
Данная дтошка будет ждать JSON следующего формата:

	{
		"complete": true
	}
*/
type CompleteTaskDTO struct {
	Complete bool
}

// Ниже - data transfer object - объект, который не хранит информацию,
// а которая ТОЛЬКО принимает входящий HTTP-запрос
type TaskDTO struct {
	Title       string
	Description string
}

func (t TaskDTO) ValidateForCreate() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}

	if t.Description == "" {
		return errors.New("description is empty")
	}

	return nil
}

// Для отправки ошибки в HTTP ответе
type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e *ErrorDTO) ToString() string {
	// функция MarshalIndent() делает красивые отступы
	b, err := json.MarshalIndent(e, "", "    ")

	// Тут не стоит обрабатывать ошибку, потому что даже если она будет,
	// тогда произошло что-то фатально (ошибка в базовой структуре),
	// поэтому просто если она будет кинем панику
	if err != nil {
		panic(err)
	}

	return string(b)
}
