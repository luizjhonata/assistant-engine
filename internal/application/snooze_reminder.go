package application

import (
	"fmt"
	"time"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type SnoozeReminderCommand struct {
	repository domain.ReminderRepository
	scheduler  domain.TaskScheduler
}

func NewSnoozeReminderCommand(repository domain.ReminderRepository, scheduler domain.TaskScheduler) *SnoozeReminderCommand {
	return &SnoozeReminderCommand{repository: repository, scheduler: scheduler}
}

func (c *SnoozeReminderCommand) Execute(id string, newTime time.Time) (domain.Reminder, error) {
	existing, err := c.repository.FindByID(id)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("finding reminder: %w", err)
	}

	snoozed, err := existing.Snooze(newTime)
	if err != nil {
		return domain.Reminder{}, fmt.Errorf("snoozing reminder: %w", err)
	}

	if err := c.scheduler.Unschedule(id); err != nil {
		return domain.Reminder{}, fmt.Errorf("unscheduling old reminder: %w", err)
	}

	if err := c.repository.Save(snoozed); err != nil {
		return domain.Reminder{}, fmt.Errorf("saving snoozed reminder: %w", err)
	}

	if err := c.scheduler.Schedule(snoozed); err != nil {
		return domain.Reminder{}, fmt.Errorf("scheduling snoozed reminder: %w", err)
	}

	return snoozed, nil
}
