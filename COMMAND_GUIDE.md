# 📖 ASHOS Command Guide

A simple, straightforward reference for all commands in **ASHOS**.

---

## 🚀 Installation & Global Access

Access `ash` directly from any terminal directory without `./`:

```bash
# Option 1: Global Go Install
go install

# Option 2: System Symlink
sudo ln -sf $(pwd)/ash /usr/local/bin/ash
```

---

## ⚡ Quick Cheat Sheet

```bash
# Tasks
ash task add "Title"       # Add a new task
ash task list              # View all tasks
ash task done <id>         # Mark task as complete
ash task delete <id>       # Delete a task

# AI Brain (RAG)
ash ai ask "question"      # Ask AI about your tasks, notes, & context
ash ai ingest-about-me     # Embed data/about_me.json into AI memory
ash ai suggest             # Get smart suggestion for next task
ash ai standup             # Auto-generate daily standup
ash ai daily-digest        # Generate daily productivity report
ash ai weekly-sync         # Sync with Appwrite & update context
ash ai reingest --all      # Refresh vector database

# Focus Tracking
ash focus start            # Start focus timer (Ctrl+C to stop)
ash focus analytics        # Show 7-day focus chart & streak score

# Notes & Journals
ash note add "text"        # Add journal entry
ash note list              # List past notes

# Sprints
ash sprint end             # Log sprint review
ash sprint list            # View past sprints

# Slack Integration
ash slack connect -t TOKEN -c CHANNEL  # Connect Slack bot
ash slack send "message"               # Post to Slack channel
ash slack list                         # View recent channel messages

# Dashboard & System
ash status                 # View quick status
ash dash                   # Launch interactive TUI dashboard
```

---

## 📝 1. Task Management (`ash task`)

Manage your daily tasks and to-dos.

| Command | Usage | Description |
|:---|:---|:---|
| `add` | `ash task add <title>` | Create a new task (auto-ingested into AI memory) |
| `list` | `ash task list` | Display all pending and completed tasks |
| `done` | `ash task done <id>` | Mark a specific task as completed by ID |
| `delete` | `ash task delete <id>` | Permanently delete a task by ID |

**Examples:**
```bash
ash task add "Implement auth middleware"
ash task list
ash task done 1
ash task delete 2
```

---

## 🧠 2. AI Brain & RAG (`ash ai`)

Interact with your vector-embedding-powered AI assistant.

| Command | Usage | Description |
|:---|:---|:---|
| `ask` | `ash ai ask "<query>"` | Semantic Q&A across your tasks, notes, and profile |
| `ingest-about-me` | `ash ai ingest-about-me` | Embeds `data/about_me.json` into vector database |
| `suggest` | `ash ai suggest` | AI analyzes your habits & recommends what to do next |
| `standup` | `ash ai standup` | Generates Slack-ready daily standup report |
| `daily-digest` | `ash ai daily-digest` | Detailed productivity intelligence breakdown |
| `weekly-sync` | `ash ai weekly-sync` | Syncs data from cloud and updates AI personality |
| `reingest` | `ash ai reingest --all` | Re-embeds all records into the vector database |

**Flags:**
- `--local` or `-l`: Use local Ollama model instead of Cloud NVIDIA API.

**Examples:**
```bash
ash ai ask "What was I working on yesterday?"
ash ai ask "Tell me about Ashish's projects" --local
ash ai ingest-about-me
```

---

## ⏱️ 3. Deep Work & Focus (`ash focus`)

Track focus sessions and visualize work density.

| Command | Usage | Description |
|:---|:---|:---|
| `start` | `ash focus start` | Starts a live focus timer. Press `Ctrl+C` when done |
| `analytics` | `ash focus analytics` | Renders 7-day focus intensity chart and streak score |

**Examples:**
```bash
ash focus start
ash focus analytics
```

---

## ✍️ 4. Notes & Journaling (`ash note`)

Capture quick notes and daily logs.

| Command | Usage | Description |
|:---|:---|:---|
| `add` | `ash note add "<text>"` | Save a note (automatically indexed by RAG AI) |
| `list` | `ash note list` | View all saved notes in chronological order |

**Examples:**
```bash
ash note add "Fixed memory leak in buffer pool using sync.Pool"
ash note list
```

---

## 🌊 5. Sprint Tracking (`ash sprint`)

Group work into execution cycles.

| Command | Usage | Description |
|:---|:---|:---|
| `end` | `ash sprint end` | Summarize and archive current work sprint |
| `list` | `ash sprint list` | View historical sprint logs |

---

## 🪟 6. Slack Integration (`ash slack`)

Connect ASHOS to your Slack workspace.

| Command | Usage | Description |
|:---|:---|:---|
| `connect` | `ash slack connect -t <token> -c <channel>` | Link Slack Bot token & Channel ID |
| `send` | `ash slack send "<message>"` | Broadcast a message to the linked channel |
| `list` | `ash slack list` | Fetch recent channel messages |

---

## 🛠️ 7. System Dashboard (`ash status` / `ash dash`)

| Command | Usage | Description |
|:---|:---|:---|
| `status` | `ash status` | Quick terminal summary (uptime, tasks, focus time) |
| `dash` | `ash dash` | Interactive Bubble Tea TUI dashboard (Press `q` to exit) |
