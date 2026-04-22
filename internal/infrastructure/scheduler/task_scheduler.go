package scheduler

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf16"

	"github.com/luizjhonata/assistant-engine/internal/domain"
	"github.com/luizjhonata/assistant-engine/internal/infrastructure/config"
)

const taskPrefix = "AssistantEngine_"

type WindowsTaskScheduler struct {
	webhookURL  string
	webhookType config.WebhookType
	mentionID   string
	mentionName string
}

func NewWindowsTaskScheduler(cfg config.Config) *WindowsTaskScheduler {
	return &WindowsTaskScheduler{
		webhookURL:  cfg.WebhookURL,
		webhookType: cfg.WebhookType,
		mentionID:   cfg.MentionID,
		mentionName: cfg.MentionName,
	}
}

func (s *WindowsTaskScheduler) Schedule(reminder domain.Reminder) error {
	taskName := taskPrefix + reminder.ID
	payload := s.buildPayload(reminder)

	innerScript := fmt.Sprintf(
		"$body = [System.Text.Encoding]::UTF8.GetBytes('%s'); Invoke-RestMethod -Uri '%s' -Method Post -ContentType 'application/json; charset=utf-8' -Body $body",
		payload, s.webhookURL,
	)
	encodedCmd := encodeForPowerShell(innerScript)

	scheduleTime := reminder.ScheduledAt.Format("15:04")
	scheduleDate := reminder.ScheduledAt.Format("01/02/2006")

	createCmd := fmt.Sprintf(
		`$action = New-ScheduledTaskAction -Execute 'powershell.exe' -Argument '-NoProfile -WindowStyle Hidden -EncodedCommand %s'; `+
			`$trigger = New-ScheduledTaskTrigger -Once -At '%s %s'; `+
			`$trigger.EndBoundary = (Get-Date '%s %s').AddHours(1).ToString('yyyy-MM-ddTHH:mm:ss'); `+
			`$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries; `+
			`Register-ScheduledTask -TaskName '%s' -Action $action -Trigger $trigger -Settings $settings -Force`,
		encodedCmd, scheduleDate, scheduleTime, scheduleDate, scheduleTime, taskName,
	)

	return runPowerShell(createCmd)
}

func (s *WindowsTaskScheduler) buildPayload(reminder domain.Reminder) string {
	title := escapeForPowerShellString(reminder.Title)
	message := escapeForPowerShellString(reminder.Message)

	if message == "" {
		message = title
	}

	if s.webhookType == config.WebhookClassic {
		return fmt.Sprintf(
			`{"title":"%s","text":"%s"}`,
			title, message,
		)
	}

	mentionTag := ""
	mentionEntities := ""
	if s.mentionID != "" && s.mentionName != "" {
		mentionTag = fmt.Sprintf(" <at>%s</at>", s.mentionName)
		mentionEntities = fmt.Sprintf(
			`,"msteams":{"entities":[{"type":"mention","text":"<at>%s</at>","mentioned":{"id":"%s","name":"%s"}}]}`,
			s.mentionName, s.mentionID, s.mentionName,
		)
	}

	return fmt.Sprintf(
		`{"type":"message","attachments":[{"contentType":"application/vnd.microsoft.card.adaptive","contentUrl":null,"content":{"$schema":"http://adaptivecards.io/schemas/adaptive-card.json","type":"AdaptiveCard","version":"1.4","body":[{"type":"TextBlock","text":"%s","weight":"Bolder","size":"Medium"},{"type":"TextBlock","text":"%s%s","wrap":true}]%s}}]}`,
		title, message, mentionTag, mentionEntities,
	)
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

func encodeForPowerShell(script string) string {
	runes := utf16.Encode([]rune(script))
	bytes := make([]byte, len(runes)*2)
	for i, r := range runes {
		bytes[i*2] = byte(r)
		bytes[i*2+1] = byte(r >> 8)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func runPowerShell(command string) error {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("powershell command failed: %w\noutput: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func escapeForPowerShellString(s string) string {
	s = strings.ReplaceAll(s, `'`, `''`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}
