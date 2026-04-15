# 📖 ASHOS Professional Command Guide

Welcome to the **ASHOS Official Command Interface (OCI)** documentation. This guide provides a comprehensive breakdown of all available directives within your terminal-native personal operating system.

---

## 🏗️ 1. Task Orchestration
*Manage your productivity lifecycle with interface-driven CRUD operations.*

| Command | Action | AI Impact |
|:---|:---|:---|
| `ash task add <title>` | Initializes a new task | Immediate Vector Ingestion |
| `ash task list` | Renders a prioritized list | Context retrieval |
| `ash task done <id>` | Signs off a completion event | Productivity Scoring |
| `ash task delete <id>` | Purges task from persistence | Memory cleaning |

### Examples
```bash
$ ash task add "Implement Dependency Injection in storage"
✅ Task added. [Event: TaskCreated]

$ ash task list
ID  Status   Title
[1] [ ]      Implement Dependency Injection in storage
```

---

## ⏱️ 2. Deep Work & Focus
*Quantify your concentration with behavioral tracking.*

| Command | Action | Analysis |
|:---|:---|:---|
| `ash focus start` | Launches a focus window | Real-time tracking |
| `ash focus analytics` | Renders productivity heatmaps | Trend Forecasting |

### Focus Intelligence
When you finish a session (`Ctrl+C`), the AI analyzes your focus summary to calculate your **Deep Work Density**.
```bash
$ ash focus analytics
📊 ASHOS 7-DAY FOCUS DENSITY
Mon  | ▓▓▓▓▓▓▓▓ 4h
Tue  | ▓▓▓▓ 2h
Insights: "Focus peaks detected on Mondays. High density sessions (>1h) increasing."
```

---

## 🧠 3. Intelligence Layer (AI Brain)
*Interact with your Personal OS Brain for data-driven insights.*

| Command | Action | Mode |
|:---|:---|:---|
| `ash ai ask "<query>"` | RAG-based context Q&A | Cloud / Local (--local) |
| `ash ai weekly-sync` | **[NEW]** Global database sync & personality update | Cloud Sync |
| `ash ai daily-digest` | Strategic behavioral analysis | Comprehensive |
| `ash ai standup` | Generates professional Slack updates | Clipboard-linked |
| `ash ai reingest` | Force-refreshes Vector memory | Maintenance |

### 🚀 NEW: Weekly Sync Logic
The `weekly-sync` command is the heartbeat of your AI OS. It fetches all recent notes and tasks from **Appwrite Cloud**, analyzes the week's patterns, and **updates your personality context**.

```bash
$ ash ai weekly-sync
☁️  Appwrite Sync Enabled
📂 Starting Weekly Sync (From: 2026-04-08)
🔍 Fetching weekly notes... [3 found]
🔍 Fetching weekly tasks... [12 found]
🤖 Updating 'About Me' personality context...
✨ AI Summary Generated: "Ashish is rapidly scaling ASHOS with a focus on API reliability..."
🎉 Weekly Sync Complete!
```

---

## ✍️ 4. Knowledge Management (Notes)
*Capture insights and ephemeral data points.*

| Command | Action | RAG Tier |
|:---|:---|:---|
| `ash note add "<text>"` | Archiving new journal entry | Semantic Search Ready |
| `ash note list` | Sequential history review | Temporal View |

---

## 🌊 5. Project Sprints
*Group work into high-impact execution cycles.*

| Command | Action | Impact |
|:---|:---|:---|
| `ash sprint end` | Strategic closure of work cycle | Sprint Report Generation |
| `ash sprint list` | Milestone review | Historical Archives |

---

## 🛠️ 6. System Diagnostics
*Process-level insights of your CLI OS.*

| Command | Action | Metric |
|:---|:---|:---|
| `ash status` | Brief diagnostics | Uptime, Task Latency |
| `ash dash` | Full Interactive TUI | System Dashboard |

---

## 🪟 7. Slack Command Center
*Automated remote reporting.*

| Command | Action | Integration |
|:---|:---|:---|
| `ash slack connect` | Auth link via Bot Token | Webhook / API |
| `ash slack send` | Manual broadcast | Channel-targeted |

---

### Pro-Tips:
- Use `-h` after any command for real-time flag documentation.
- All AI operations support `--local` to prioritize **Ollama** privacy.
- For the best experience, run `ash ai weekly-sync` every Sunday to calibrate your OS personality.
