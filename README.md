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

There are two methods to create a webhook in Microsoft Teams. **Workflows** is the recommended approach by Microsoft — the classic Incoming Webhook connector is being deprecated.

#### Method 1: Workflows (recommended)

This method uses Power Automate to create a webhook endpoint that posts Adaptive Cards to a channel.

1. Open **Microsoft Teams** and go to the channel where you want notifications
2. Click the **+** (add a tab) or go to **...** → **Workflows**
3. Search for the template **"Post to a channel when a webhook request is received"**
4. Select it and follow the setup wizard:
   - Name the workflow (e.g., "Assistant Engine Notifications")
   - Confirm the Team and Channel where messages will be posted
   - Click **Add workflow**
5. After creation, Teams shows the **webhook URL** — copy it
6. The URL looks like: `https://*.logic.azure.com:443/workflows/...`
7. Paste it in your `config.json` and set `webhook_type` to `"workflow"`:

```json
{
  "webhook_url": "https://prod-XX.westus.logic.azure.com:443/workflows/...",
  "webhook_type": "workflow",
  "default_time": "09:00",
  "default_delay_hours": 24
}
```

#### Method 2: Incoming Webhook (legacy)

This method uses the classic Office 365 connector. It still works in some tenants but Microsoft is phasing it out.

1. Open **Microsoft Teams** and go to the channel where you want notifications
2. Click the **...** menu next to the channel name → **Connectors** (or **Manage channel** → **Connectors**)
3. Search for **"Incoming Webhook"** and click **Configure**
4. Give it a name (e.g., "Assistant Engine") and optionally upload an icon
5. Click **Create** — Teams generates a webhook URL
6. The URL looks like: `https://your-org.webhook.office.com/webhookb2/...`
7. Paste it in your `config.json` and set `webhook_type` to `"classic"`:

```json
{
  "webhook_url": "https://your-org.webhook.office.com/webhookb2/...",
  "webhook_type": "classic",
  "default_time": "09:00",
  "default_delay_hours": 24
}
```

> **Note:** If you don't see the Connectors option, your organization may have disabled it. Use Method 1 instead.

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
