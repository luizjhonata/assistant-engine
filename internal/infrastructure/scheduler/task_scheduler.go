package scheduler

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/luizjhonata/assistant-engine/internal/domain"
	"github.com/luizjhonata/assistant-engine/internal/infrastructure/config"
)

const taskPrefix = "AssistantEngine_"

type WindowsTaskScheduler struct {
	webhookURL string
}

func NewWindowsTaskScheduler(cfg config.Config) *WindowsTaskScheduler {
	return &WindowsTaskScheduler{webhookURL: cfg.WebhookURL}
}

func (s *WindowsTaskScheduler) Schedule(reminder domain.Reminder) error {
	taskName := taskPrefix + reminder.ID

	title := escapeForJSON(reminder.Title)
	message := escapeForJSON(reminder.Message)

	payload := fmt.Sprintf(
		`{\"title\":\"%s\",\"text\":\"%s\"}`,
		title, message,
	)

	psAction := fmt.Sprintf(
		`Invoke-RestMethod -Uri '%s' -Method Post -ContentType 'application/json' -Body '%s'`,
		s.webhookURL, payload,
	)

	scheduleTime := reminder.ScheduledAt.Format("15:04")
	scheduleDate := reminder.ScheduledAt.Format("01/02/2006")

	createCmd := fmt.Sprintf(
		`$action = New-ScheduledTaskAction -Execute 'powershell.exe' -Argument '-NoProfile -WindowStyle Hidden -Command "%s"'; `+
			`$trigger = New-ScheduledTaskTrigger -Once -At '%s %s'; `+
			`$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -DeleteExpiredTaskAfterCompletion (New-TimeSpan -Hours 1); `+
			`$trigger.EndBoundary = (Get-Date '%s %s').AddHours(1).ToString('yyyy-MM-ddTHH:mm:ss'); `+
			`Register-ScheduledTask -TaskName '%s' -Action $action -Trigger $trigger -Settings $settings -Force`,
		psAction, scheduleDate, scheduleTime, scheduleDate, scheduleTime, taskName,
	)

	return runPowerShell(createCmd)
}

func (s *WindowsTaskScheduler) Unschedule(reminderID string) error {
	taskName := taskPrefix + reminderID

	cmd := fmt.Sprintf(
		`Unregister-ScheduledTask -TaskName '%s' -Confirm:$false -ErrorAction SilentlyContinue`,
		taskName,
	)

	return runPowerShell(cmd)
}

func (s *WindowsTaskScheduler) UnscheduleAll() error {
	cmd := fmt.Sprintf(
		`Get-ScheduledTask | Where-Object { $_.TaskName -like '%s*' } | Unregister-ScheduledTask -Confirm:$false -ErrorAction SilentlyContinue`,
		taskPrefix,
	)

	return runPowerShell(cmd)
}

func runPowerShell(command string) error {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powershell command failed: %w\noutput: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func escapeForJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}
