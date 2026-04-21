# assistant-engine

Personal reminder CLI that creates Windows Task Scheduler tasks to send Microsoft Teams notifications via webhooks.

## Prerequisites

- Go 1.21+
- WSL2 with access to `powershell.exe`
- Microsoft Teams with permission to create channels and workflows

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
  "webhook_url": "https://...",
  "webhook_type": "workflow",
  "mention_id": "your.email@company.com",
  "mention_name": "YourFirstName",
  "default_time": "09:00",
  "default_delay_hours": 24
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `webhook_url` | yes | Teams webhook URL (see setup guide below) |
| `webhook_type` | no | `"workflow"` (default) or `"classic"` |
| `mention_id` | no | Your Teams email/UPN for @mention notifications |
| `mention_name` | no | Display name used in the @mention tag |
| `default_time` | no | Default reminder time if not specified (default: `"09:00"`) |
| `default_delay_hours` | no | Default delay when no date is given (default: `24`) |

### Notifications and @mentions

By default, messages posted to a Teams channel do **not** generate push notifications or activity feed alerts — they appear silently in the channel. To ensure you are actively notified when a reminder fires, configure `mention_id` and `mention_name`. This adds an @mention to every notification, which triggers the Teams activity feed alert just like a direct mention from another person.

If you omit these fields, reminders will still be posted to the channel, but you will need to check the channel manually.

## Setting up a Teams Channel

Teams requires every channel to belong to a team — you cannot create a standalone channel. For personal reminders, the recommended approach is:

1. Open **Microsoft Teams** and go to any team you belong to
2. Click **"+"** or **"..."** → **Add channel** (PT-BR: **Adicionar canal**)
3. Set the channel name (e.g., "Assistant Engine")
4. Set **Privacy** to **Private** — only people you explicitly add will see it
5. Set **Type** to **Standard** (PT-BR: **Publicações**), not Threads — notifications are simpler with standard posts
6. Create the channel without adding anyone else

Since you are the only member, no one else on your team will see the channel or its messages.

## Setting up a Teams Webhook

There are two methods. **Workflows** is the recommended approach — the classic Incoming Webhook connector is being deprecated by Microsoft.

### Method 1: Workflows (recommended)

This method uses Power Automate to create a webhook endpoint that posts Adaptive Cards to a channel.

1. Go to the channel you created above
2. Click **"..."** → **Workflows** (PT-BR: **Fluxos de trabalho**)
3. Search for **"webhook"** in the template search bar
4. Select **"Post to a channel when a webhook request is received"** (PT-BR: **"Enviar alertas de webhook para um canal"**)
5. Follow the setup wizard:
   - Name the workflow (e.g., "Assistant Engine Notifications")
   - Confirm the Team and Channel
   - Click **Add workflow** (PT-BR: **Adicionar fluxo de trabalho**)
6. After creation, Teams displays the **webhook URL** — copy it
7. Paste it in your `config.json` with `"webhook_type": "workflow"`

> **Note:** The URL format varies by tenant. Common formats include:
> - `https://*.logic.azure.com:443/workflows/...`
> - `https://*.powerplatform.com:443/powerautomate/automations/direct/workflows/...`
>
> Both work the same way — the format depends on your organization's Power Platform configuration.

### Method 2: Incoming Webhook (legacy)

This method uses the classic Office 365 connector. It still works in some tenants but is being phased out.

1. Go to the channel you created above
2. Click **"..."** → **Connectors** (or **Manage channel** → **Connectors**)
3. Search for **"Incoming Webhook"** and click **Configure**
4. Give it a name (e.g., "Assistant Engine") and optionally upload an icon
5. Click **Create** — Teams generates a webhook URL
6. The URL looks like: `https://your-org.webhook.office.com/webhookb2/...`
7. Paste it in your `config.json` with `"webhook_type": "classic"`

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
2. A Windows Task Scheduler task is registered via `powershell.exe` from WSL
3. At the scheduled time, the task runs a PowerShell script that sends an HTTP POST to the Teams webhook
4. The task expires 1 hour after its scheduled time (one-time trigger with `EndBoundary`)

## Architecture

```
cmd/assistant-engine/     CLI entry point (Cobra)
internal/
  domain/                 Reminder entity, repository and scheduler interfaces
  application/            Commands (add, remove, snooze, clear) and queries (list)
  infrastructure/
    config/               Config file loader (~/.assistant-engine/config.json)
    persistence/          JSON file-based reminder repository (~/.assistant-engine/reminders.json)
    scheduler/            Windows Task Scheduler via PowerShell (EncodedCommand)
```

## Troubleshooting

### Webhook URL contains `&` characters

The scheduler uses PowerShell's `-EncodedCommand` (Base64 UTF-16LE) to avoid issues with special characters in webhook URLs. This is handled automatically — no manual escaping is needed.

### Notifications are not appearing

- Verify the webhook is working: `assistant-engine add "Test" -d <tomorrow> -t <soon>` and check the channel
- If messages appear in the channel but you don't get notified, make sure `mention_id` and `mention_name` are configured in `config.json`
- Check that the email in `mention_id` matches your Teams account (UPN format: `user@domain.com`)

### Teams UI is in another language

The template and menu names vary by language. Search for **"webhook"** in the Workflows/Connectors search bar — the keyword works regardless of the UI language.

## License

MIT
