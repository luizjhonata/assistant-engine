# assistant-engine

Personal reminder CLI that creates Windows Task Scheduler tasks to send Microsoft Teams notifications via webhooks.

## Prerequisites

- Go 1.21+
- WSL2 with access to `powershell.exe`
- Microsoft Teams incoming webhook URL

## Installation

```bash
go install github.com/luizjhonata/assistant-engine/cmd/assistant-engine@latest
```

Or build from source:

```bash
git clone git@github.com:luizjhonata/assistant-engine.git
cd assistant-engine
go build -o assistant-engine ./cmd/assistant-engine/
```

## Configuration

Create `~/.assistant-engine/config.json`:

```json
{
  "webhook_url": "https://your-org.webhook.office.com/webhookb2/...",
  "default_time": "09:00",
  "default_delay_hours": 24
}
```

### Setting up a Teams Webhook

1. In Microsoft Teams, go to the channel where you want notifications
2. Click the **...** menu → **Connectors** (or **Workflows**)
3. Search for **Incoming Webhook** and configure it
4. Copy the webhook URL and paste it in `config.json`

## Usage

### Add a reminder

```bash
# Default: +24h at 09:00
assistant-engine add "Review catalog PR"

# With specific date and time
assistant-engine add "Deploy to production" -d 2026-04-25 -t 14:00

# With a detailed message
assistant-engine add "Sprint review" -m "Prepare demo for the catalog feature" -d 2026-04-28 -t 10:00
```

### List reminders

```bash
assistant-engine list
```

### Remove a reminder

```bash
assistant-engine remove <id>
```

### Snooze a reminder

```bash
# Postpone to a new date/time
assistant-engine snooze <id> -d 2026-04-30 -t 09:00
```

### Clear all reminders

```bash
assistant-engine clear
```

## How It Works

1. The CLI creates a reminder and stores it in `~/.assistant-engine/reminders.json`
2. A Windows Task Scheduler task is registered via `powershell.exe`
3. At the scheduled time, the task runs a PowerShell script that sends an HTTP POST to the Teams webhook
4. The task auto-deletes 1 hour after execution

## Architecture

```
cmd/assistant-engine/     CLI entry point (Cobra)
internal/
  domain/                 Reminder entity, repository and scheduler interfaces
  application/            Commands (add, remove, snooze, clear) and queries (list)
  infrastructure/
    config/               Config file loader
    persistence/          JSON file-based reminder repository
    scheduler/            Windows Task Scheduler via PowerShell
```

## License

MIT
