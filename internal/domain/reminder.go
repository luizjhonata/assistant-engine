package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEmptyTitle       = errors.New("reminder title must not be empty")
	ErrPastSchedule     = errors.New("scheduled time must be in the future")
	ErrReminderNotFound = errors.New("reminder not found")
)

type ReminderStatus string

const (
	StatusActive  ReminderStatus = "active"
	StatusFired   ReminderStatus = "fired"
	StatusSnoozed ReminderStatus = "snoozed"
)

type Reminder struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Message     string         `json:"message"`
	ScheduledAt time.Time      `json:"scheduled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	Status      ReminderStatus `json:"status"`
}

func NewReminder(title, message string, scheduledAt time.Time) (Reminder, error) {
	if title == "" {
		return Reminder{}, ErrEmptyTitle
	}

	if !scheduledAt.After(time.Now()) {
		return Reminder{}, ErrPastSchedule
	}

	id, err := generateID()
	if err != nil {
		return Reminder{}, fmt.Errorf("generating reminder ID: %w", err)
	}

	return Reminder{
		ID:          id,
		Title:       title,
		Message:     message,
		ScheduledAt: scheduledAt,
		CreatedAt:   time.Now(),
		Status:      StatusActive,
	}, nil
}

func (r Reminder) Snooze(newTime time.Time) (Reminder, error) {
	if !newTime.After(time.Now()) {
		return r, ErrPastSchedule
	}

	return Reminder{
		ID:          r.ID,
		Title:       r.Title,
		Message:     r.Message,
		ScheduledAt: newTime,
		CreatedAt:   r.CreatedAt,
		Status:      StatusActive,
	}, nil
}

func (r Reminder) IsActive() bool {
	return r.Status == StatusActive
}

func generateID() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
