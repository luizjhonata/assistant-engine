package application

import (
	"fmt"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type RemoveReminderCommand struct {
	repository domain.ReminderRepository
	scheduler  domain.TaskScheduler
}

func NewRemoveReminderCommand(repository domain.ReminderRepository, scheduler domain.TaskScheduler) *RemoveReminderCommand {
	return &RemoveReminderCommand{repository: repository, scheduler: scheduler}
}

func (c *RemoveReminderCommand) Execute(id string) error {
	if _, err := c.repository.FindByID(id); err != nil {
		return fmt.Errorf("finding reminder: %w", err)
	}

	if err := c.scheduler.Unschedule(id); err != nil {
		return fmt.Errorf("unscheduling reminder: %w", err)
	}

	if err := c.repository.DeleteByID(id); err != nil {
		return fmt.Errorf("deleting reminder: %w", err)
	}

	return nil
}
