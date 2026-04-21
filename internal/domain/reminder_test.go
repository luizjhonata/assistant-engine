package domain_test

import (
	"testing"
	"time"

	"github.com/luizjhonata/assistant-engine/internal/domain"
)

func TestNewReminder(t *testing.T) {
	t.Run("creates reminder with valid inputs", func(t *testing.T) {
		scheduledAt := time.Now().Add(24 * time.Hour)

		reminder, err := domain.NewReminder("Review PR", "Check the catalog PR", scheduledAt)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reminder.ID == "" {
			t.Error("expected non-empty ID")
		}
		if reminder.Title != "Review PR" {
			t.Errorf("expected title 'Review PR', got %q", reminder.Title)
		}
		if reminder.Message != "Check the catalog PR" {
			t.Errorf("expected message 'Check the catalog PR', got %q", reminder.Message)
		}
		if reminder.Status != domain.StatusActive {
			t.Errorf("expected status 'active', got %q", reminder.Status)
		}
		if reminder.CreatedAt.IsZero() {
			t.Error("expected non-zero CreatedAt")
		}
	})

	t.Run("rejects empty title", func(t *testing.T) {
		scheduledAt := time.Now().Add(24 * time.Hour)

		_, err := domain.NewReminder("", "some message", scheduledAt)

		if err != domain.ErrEmptyTitle {
			t.Errorf("expected ErrEmptyTitle, got %v", err)
		}
	})

	t.Run("rejects past schedule time", func(t *testing.T) {
		pastTime := time.Now().Add(-1 * time.Hour)

		_, err := domain.NewReminder("Test", "message", pastTime)

		if err != domain.ErrPastSchedule {
			t.Errorf("expected ErrPastSchedule, got %v", err)
		}
	})

	t.Run("allows empty message", func(t *testing.T) {
		scheduledAt := time.Now().Add(1 * time.Hour)

		reminder, err := domain.NewReminder("Quick note", "", scheduledAt)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reminder.Message != "" {
			t.Errorf("expected empty message, got %q", reminder.Message)
		}
	})
}

func TestReminderSnooze(t *testing.T) {
	t.Run("snoozes to a future time", func(t *testing.T) {
		original, err := domain.NewReminder("Meeting", "standup", time.Now().Add(1*time.Hour))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		newTime := time.Now().Add(48 * time.Hour)
		snoozed, err := original.Snooze(newTime)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if snoozed.ID != original.ID {
			t.Error("snooze should preserve the original ID")
		}
		if !snoozed.ScheduledAt.Equal(newTime) {
			t.Errorf("expected scheduled time %v, got %v", newTime, snoozed.ScheduledAt)
		}
		if snoozed.Status != domain.StatusActive {
			t.Errorf("expected status 'active' after snooze, got %q", snoozed.Status)
		}
	})

	t.Run("rejects snooze to past time", func(t *testing.T) {
		original, err := domain.NewReminder("Meeting", "standup", time.Now().Add(1*time.Hour))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = original.Snooze(time.Now().Add(-1 * time.Hour))

		if err != domain.ErrPastSchedule {
			t.Errorf("expected ErrPastSchedule, got %v", err)
		}
	})
}

func TestReminderIsActive(t *testing.T) {
	reminder, err := domain.NewReminder("Test", "msg", time.Now().Add(1*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reminder.IsActive() {
		t.Error("new reminder should be active")
	}
}
