package domain

type TaskScheduler interface {
	Schedule(reminder Reminder) error
	Unschedule(reminderID string) error
	UnscheduleAll() error
}
