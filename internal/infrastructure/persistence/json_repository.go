package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/luizjhonata/assistant-engine/internal/domain"
	"github.com/luizjhonata/assistant-engine/internal/infrastructure/config"
)

const remindersFileName = "reminders.json"

type JSONReminderRepository struct {
	mu sync.Mutex
}

func NewJSONReminderRepository() *JSONReminderRepository {
	return &JSONReminderRepository{}
}

func (r *JSONReminderRepository) Save(reminder domain.Reminder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	reminders, err := r.loadAll()
	if err != nil {
		return err
	}

	found := false
	for i, existing := range reminders {
		if existing.ID == reminder.ID {
			reminders[i] = reminder
			found = true
			break
		}
	}
	if !found {
		reminders = append(reminders, reminder)
	}

	return r.writeAll(reminders)
}

func (r *JSONReminderRepository) FindByID(id string) (domain.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	reminders, err := r.loadAll()
	if err != nil {
		return domain.Reminder{}, err
	}

	for _, reminder := range reminders {
		if reminder.ID == id {
			return reminder, nil
		}
	}

	return domain.Reminder{}, domain.ErrReminderNotFound
}

func (r *JSONReminderRepository) FindAll() ([]domain.Reminder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.loadAll()
}

func (r *JSONReminderRepository) DeleteByID(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	reminders, err := r.loadAll()
	if err != nil {
		return err
	}

	filtered := make([]domain.Reminder, 0, len(reminders))
	for _, reminder := range reminders {
		if reminder.ID != id {
			filtered = append(filtered, reminder)
		}
	}

	return r.writeAll(filtered)
}

func (r *JSONReminderRepository) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.writeAll([]domain.Reminder{})
}

func (r *JSONReminderRepository) loadAll() ([]domain.Reminder, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Reminder{}, nil
		}
		return nil, fmt.Errorf("reading reminders file: %w", err)
	}

	if len(data) == 0 {
		return []domain.Reminder{}, nil
	}

	var reminders []domain.Reminder
	if err := json.Unmarshal(data, &reminders); err != nil {
		return nil, fmt.Errorf("parsing reminders file: %w", err)
	}

	return reminders, nil
}

func (r *JSONReminderRepository) writeAll(reminders []domain.Reminder) error {
	path, err := filePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(reminders, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling reminders: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing reminders file: %w", err)
	}

	return nil
}

func filePath() (string, error) {
	dir, err := config.DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, remindersFileName), nil
}
