package application

import (
	"fmt"
	"time"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type AddReminderCommand struct {
	repository domain.ReminderRepository
	scheduler  domain.TaskScheduler
}

func NewAddReminderCommand(repository domain.ReminderRepository, scheduler domain.TaskScheduler) *AddReminderCommand {
	return &AddReminderCommand{repository: repository, scheduler: scheduler}
}

type AddReminderInput struct {
	Title       string
	Message     string
	ScheduledAt time.Time
}

func (c *AddReminderCommand) Execute(input AddReminderInput) (domain.Reminder, error) {
	reminder, err := domain.NewReminder(input.Title, input.Message, input.ScheduledAt)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("creating reminder: %w", err)
	}

	if err := c.repository.Save(reminder); err != nil {
		return domain.Reminder{}, fmt.Errorf("saving reminder: %w", err)
	}

	if err := c.scheduler.Schedule(reminder); err != nil {
		_ = c.repository.DeleteByID(reminder.ID)
		return domain.Reminder{}, fmt.Errorf("scheduling reminder: %w", err)
	}

	return reminder, nil
}
