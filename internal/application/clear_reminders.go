package application

import (
	"fmt"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type ClearRemindersCommand struct {
	repository domain.ReminderRepository
	scheduler  domain.TaskScheduler
}

func NewClearRemindersCommand(repository domain.ReminderRepository, scheduler domain.TaskScheduler) *ClearRemindersCommand {
	return &ClearRemindersCommand{repository: repository, scheduler: scheduler}
}

func (c *ClearRemindersCommand) Execute() error {
	if err := c.scheduler.UnscheduleAll(); err != nil {
		return fmt.Errorf("unscheduling all reminders: %w", err)
	}

	if err := c.repository.DeleteAll(); err != nil {
		return fmt.Errorf("deleting all reminders: %w", err)
	}

	return nil
}
