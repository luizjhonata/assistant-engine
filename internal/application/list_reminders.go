package application

import (
	"fmt"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type ListRemindersQuery struct {
	repository domain.ReminderRepository
}

func NewListRemindersQuery(repository domain.ReminderRepository) *ListRemindersQuery {
	return &ListRemindersQuery{repository: repository}
}

func (q *ListRemindersQuery) Execute() ([]domain.Reminder, error) {
	reminders, err := q.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listing reminders: %w", err)
	}
	return reminders, nil
}
