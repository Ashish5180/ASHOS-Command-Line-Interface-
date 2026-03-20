# 📖 ASHOS Command Guide

Welcome to the **ASHOS Official Command Guide**. This document outlines every command available in the Personal CLI Operating System, its purpose, and example usage with responses.

---

## 🏗️ 1. Task Management
Handle your daily tasks with RAG-integrated CRUD operations.

| Command | Description |
|---|---|
| `ash task add <title>` | Create a new task. Automatically ingested into AI Brain. |
| `ash task list` | View all pending and completed tasks. |
| `ash task done <id>` | Mark a specific task as completed. |
| `ash task delete <id>` | Permanently remove a task from the database. |

### Example
```bash
$ ./ash task add "Refactor storage layer"
🧠 AI Brain: Ingesting task 'Refactor storage layer'...
✅ Task added.

$ ./ash task list
[1] [ ] Refactor storage layer
```

---

## ⏱️ 2. Focus Mode
Deep work tracking with behavioral analysis.

| Command | Description |
|---|---|
| `ash focus start` | Start a focus session. It will prompt for a summary on exit. |
| `ash focus analytics` | Show ASCII charts and productivity heatmap for the last 7 days. |

### Example (Session End)
```bash
$ ./ash focus start
🧠 Focus Session Started! Stay sharp.
Press Ctrl+C to end focus session.
^C
📝 What did you focus on? (Optional): Implementing event-driven RAG
🧠 AI Brain: Ingesting focus session (45m)...
🔋 [Event] Focus session complete: 45m tracked. 
Total Focus Time: 45m2s
```

---

## ✍️ 3. Journaling (Notes)
Capture personal insights, ideas, and journal entries.

| Command | Description |
|---|---|
| `ash note add <text>` | Add a new journal entry. Content is vector-embedded for RAG. |
| `ash note list` | List recent notes with timestamps. |

### Example
```bash
$ ./ash note add "Morning deep work is 2x more productive than evening."
🧠 AI Brain: Ingesting journal entry...
✍️  [Event] Personal insight archived.
📝 Note added to your brain.
```

---

## 🌊 4. Sprint Mode
Summarize large work cycles and project milestones.

| Command | Description |
|---|---|
| `ash sprint end` | Interactively close a sprint with summary, tasks, and time. |
| `ash sprint list` | Review historical sprint performance. |

### Example
```bash
$ ./ash sprint end
🏆 Sprint Title: Infrastructure Week
📝 Sprint Review Summary: Successfully migrated to event-driven RAG architecture.
✅ Tasks Completed (number): 12
⏱️ Total Focus Time (e.g., 2h30m): 15h

🧠 AI Brain: Ingesting sprint review 'Infrastructure Week'...
🌊 [Event] Sprint 'Infrastructure Week' closed. 
🌊 Sprint summarized and archived in the brain.
```

---

## 📊 5. System Status
Quick high-level diagnostics of your Personal OS.

| Command | Description |
|---|---|
| `ash status` | Show uptime, pending tasks, and today's total focus time. |

### Example
```bash
$ ./ash status
🚀 ASHOS STATUS
Uptime:        2h 15m
Pending Tasks: 4 
Focus Today:   3h 45m
```

---

## 🧠 6. AI Brain (RAG)
The intelligence layer that knows everything about your tasks, notes, and sprints.

| Command | Description |
|---|---|
| `ash ai ask <query>` | Ask anything about your data (tasks, sessions, notes). |
| `ash ai standup` | Automatically generate a daily standup based on recent activity. |
| `ash ai suggest` | AI analyzes your patterns to suggest the most logical next task. |

### Example
```bash
$ ./ash ai ask "What did I learn about my morning productivity?"
🧠 ASHOS AI BRAIN IS THINKING...
🤖 AI RESPONSE
Based on your journal entry from this morning, you noted that your 
DEEP WORK IS MOST EFFECTIVE BEFORE 11 AM, showing nearly 2x 
productivity compared to evening sessions. Focus your high-intensity 
tasks in this window.
```

---

## 🏎️ 7. Interactive Dashboard
The full terminal command center.

| Command | Description |
|---|---|
| `ash dash` | Launch the Bubble Tea TUI dashboard. |

### Keybindings
- `r`: Manual data refresh
- `q`: Exit dashboard
- `Ctrl+C`: Force quit

---

**Tip**: Always use `./ash --help` to see the most up-to-date command tree.
