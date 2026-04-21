package domain

type ReminderRepository interface {
	Save(reminder Reminder) error
	FindByID(id string) (Reminder, error)
	FindAll() ([]Reminder, error)
	DeleteByID(id string) error
	DeleteAll() error
}
