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

## 🧠 6. AI Brain (Intelligence Layer)
The intelligence layer that analyzes your tasks, focus sessions, notes, and sprints to provide deep insights.

| Command | Description |
|---|---|
| `ash ai ask <query>` | Ask about your data. **NEW**: Use `--local` for offline mode (Ollama). |
| `ash ai standup` | Interactively generate a standup and copy it to your clipboard (`pbcopy`). |
| `ash ai suggest` | Structured suggestions with reasoning and productivity tips. |
| `ash ai daily-digest` | Advanced behavioral analysis and productivity scores. |
| `ash ai reingest` | **NEW**: Force re-sync of all tasks, notes, etc. into vector store. |

### Example: Offline Mode (Local LLM)
```bash
$ ash ai ask "summary" --local
🚀 AI Brain: Switched to Local Mode (Ollama)
```

### Example: Retrieval Scores
```bash
$ ash ai ask "What was I doing last Tuesday?"
🧠 ASHOS AI BRAIN IS THINKING...

Retrieved Context (Relevance Scores):
  → "Database refactor task"     [94% match]
  → "SQLite engine focus session" [87% match]  

🤖 AI RESPONSE
Last Tuesday you were primarily working on the database refactor...
```

### Example: Daily Digest
```bash
$ ash ai daily-digest
📊 GENERATING DAILY INTELLIGENCE...

📊 ASHOS Daily Intelligence — March 21
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Productivity Score : 82/100 📈
Focus Time         : 3h 12m
Tasks Completed    : 4/10
AI Queries Made    : 5 (est)

💡 Insight: "Your work sessions are lengthening. Consider 
   shorter breaks to maintain this high baseline."
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Example: Manual Re-ingestion
```bash
$ ash ai reingest --type tasks
🧠 Re-ingesting TASKS...
🧠 AI Brain: Ingesting task 'Finalize OS Intelligence'...
✅ Re-ingestion complete!
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

## 🛠️ How to Check New Changes

To verify the new AI improvements, follow these steps:

1. **Rebuild the project**:
   ```bash
   go build -o ash main.go
   ```
   *If you encounter permission errors, try `go run main.go` instead:*
   ```bash
   go run main.go ai ask "Am I being productive today?"
   ```

2. **Test Multi-turn Memory**:
   Run two consecutive queries to see if the AI remembers the context:
   ```bash
   ./ash ai ask "What are my pending tasks?"
   ./ash ai ask "Which one should I do first?"
   ```

3. **Test Interactive Standup**:
   ```bash
   ./ash ai standup
   ```
   Answer `y` when prompted to copy to clipboard and verify it works.

4. **Test Smart Suggestions**:
   ```bash
   ./ash ai suggest
   ```

5. **Test Daily Intelligence Digest**:
   ```bash
   ./ash ai daily-digest
   ```

**Tip**: Always use `./ash --help` to see the most up-to-date command tree.
