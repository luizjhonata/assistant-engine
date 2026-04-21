package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/luizjhonata/assistant-engine/internal/application"
	"github.com/luizjhonata/assistant-engine/internal/domain"
)

type mockRepository struct {
	reminders map[string]domain.Reminder
	saveErr   error
	findErr   error
	deleteErr error
}

func newMockRepository() *mockRepository {
	return &mockRepository{reminders: make(map[string]domain.Reminder)}
}

func (m *mockRepository) Save(r domain.Reminder) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.reminders[r.ID] = r
	return nil
}

func (m *mockRepository) FindByID(id string) (domain.Reminder, error) {
	if m.findErr != nil {
		return domain.Reminder{}, m.findErr
	}
	r, ok := m.reminders[id]
	if !ok {
		return domain.Reminder{}, domain.ErrReminderNotFound
	}
	return r, nil
}

func (m *mockRepository) FindAll() ([]domain.Reminder, error) {
	result := make([]domain.Reminder, 0, len(m.reminders))
	for _, r := range m.reminders {
		result = append(result, r)
	}
	return result, nil
}

func (m *mockRepository) DeleteByID(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.reminders, id)
	return nil
}

func (m *mockRepository) DeleteAll() error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.reminders = make(map[string]domain.Reminder)
	return nil
}

type mockScheduler struct {
	scheduled   map[string]bool
	scheduleErr error
}

func newMockScheduler() *mockScheduler {
	return &mockScheduler{scheduled: make(map[string]bool)}
}

func (m *mockScheduler) Schedule(r domain.Reminder) error {
	if m.scheduleErr != nil {
		return m.scheduleErr
	}
	m.scheduled[r.ID] = true
	return nil
}

func (m *mockScheduler) Unschedule(id string) error {
	delete(m.scheduled, id)
	return nil
}

func (m *mockScheduler) UnscheduleAll() error {
	m.scheduled = make(map[string]bool)
	return nil
}

func TestAddReminderCommand(t *testing.T) {
	t.Run("adds and schedules a reminder", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		cmd := application.NewAddReminderCommand(repo, sched)

		input := application.AddReminderInput{
			Title:       "Deploy check",
			Message:     "Verify production deploy",
			ScheduledAt: time.Now().Add(24 * time.Hour),
		}

		reminder, err := cmd.Execute(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reminder.Title != "Deploy check" {
			t.Errorf("expected title 'Deploy check', got %q", reminder.Title)
		}
		if len(repo.reminders) != 1 {
			t.Errorf("expected 1 reminder in repo, got %d", len(repo.reminders))
		}
		if !sched.scheduled[reminder.ID] {
			t.Error("reminder should be scheduled")
		}
	})

	t.Run("rolls back repository on scheduler failure", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		sched.scheduleErr = errors.New("scheduler unavailable")
		cmd := application.NewAddReminderCommand(repo, sched)

		input := application.AddReminderInput{
			Title:       "Will fail",
			ScheduledAt: time.Now().Add(1 * time.Hour),
		}

		_, err := cmd.Execute(input)

		if err == nil {
			t.Fatal("expected error from scheduler failure")
		}
		if len(repo.reminders) != 0 {
			t.Error("reminder should have been rolled back from repository")
		}
	})

	t.Run("rejects invalid input", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		cmd := application.NewAddReminderCommand(repo, sched)

		input := application.AddReminderInput{
			Title:       "",
			ScheduledAt: time.Now().Add(1 * time.Hour),
		}

		_, err := cmd.Execute(input)

		if err == nil {
			t.Fatal("expected error for empty title")
		}
	})
}

func TestRemoveReminderCommand(t *testing.T) {
	t.Run("removes existing reminder", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		addCmd := application.NewAddReminderCommand(repo, sched)
		removeCmd := application.NewRemoveReminderCommand(repo, sched)

		reminder, _ := addCmd.Execute(application.AddReminderInput{
			Title:       "To remove",
			ScheduledAt: time.Now().Add(1 * time.Hour),
		})

		err := removeCmd.Execute(reminder.ID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.reminders) != 0 {
			t.Error("reminder should have been deleted")
		}
		if sched.scheduled[reminder.ID] {
			t.Error("reminder should have been unscheduled")
		}
	})

	t.Run("fails for non-existent reminder", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		cmd := application.NewRemoveReminderCommand(repo, sched)

		err := cmd.Execute("nonexistent")

		if err == nil {
			t.Fatal("expected error for non-existent reminder")
		}
	})
}

func TestSnoozeReminderCommand(t *testing.T) {
	t.Run("snoozes existing reminder", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		addCmd := application.NewAddReminderCommand(repo, sched)
		snoozeCmd := application.NewSnoozeReminderCommand(repo, sched)

		original, _ := addCmd.Execute(application.AddReminderInput{
			Title:       "Snoozable",
			ScheduledAt: time.Now().Add(1 * time.Hour),
		})

		newTime := time.Now().Add(48 * time.Hour)
		snoozed, err := snoozeCmd.Execute(original.ID, newTime)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !snoozed.ScheduledAt.Equal(newTime) {
			t.Errorf("expected new schedule time, got %v", snoozed.ScheduledAt)
		}
		if snoozed.ID != original.ID {
			t.Error("ID should be preserved after snooze")
		}
	})
}

func TestClearRemindersCommand(t *testing.T) {
	t.Run("clears all reminders", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		addCmd := application.NewAddReminderCommand(repo, sched)
		clearCmd := application.NewClearRemindersCommand(repo, sched)

		for i := 0; i < 3; i++ {
			_, _ = addCmd.Execute(application.AddReminderInput{
				Title:       "Reminder",
				ScheduledAt: time.Now().Add(time.Duration(i+1) * time.Hour),
			})
		}

		err := clearCmd.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.reminders) != 0 {
			t.Errorf("expected 0 reminders, got %d", len(repo.reminders))
		}
		if len(sched.scheduled) != 0 {
			t.Errorf("expected 0 scheduled tasks, got %d", len(sched.scheduled))
		}
	})
}

func TestListRemindersQuery(t *testing.T) {
	t.Run("returns all reminders", func(t *testing.T) {
		repo := newMockRepository()
		sched := newMockScheduler()
		addCmd := application.NewAddReminderCommand(repo, sched)
		listQuery := application.NewListRemindersQuery(repo)

		for i := 0; i < 2; i++ {
			_, _ = addCmd.Execute(application.AddReminderInput{
				Title:       "Reminder",
				ScheduledAt: time.Now().Add(time.Duration(i+1) * time.Hour),
			})
		}

		reminders, err := listQuery.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(reminders) != 2 {
			t.Errorf("expected 2 reminders, got %d", len(reminders))
		}
	})

	t.Run("returns empty list when no reminders", func(t *testing.T) {
		repo := newMockRepository()
		query := application.NewListRemindersQuery(repo)

		reminders, err := query.Execute()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(reminders) != 0 {
			t.Errorf("expected 0 reminders, got %d", len(reminders))
		}
	})
}
