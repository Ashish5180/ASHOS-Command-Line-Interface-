# 🚀 ASHOS: The Event-Driven Personal AI Operating System

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Tech Stack](https://img.shields.io/badge/Architecture-Clean_Hexagonal-blue?style=for-the-badge)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
[![AI](https://img.shields.io/badge/Intelligence-RAG_Driven-orange?style=for-the-badge)](https://ollama.com)

> **ASHOS** is a modular, terminal-native personal operating system designed for developers who treat their productivity like a high-performance system. Built with **Dependency Injection**, **Event-Driven Orchestration**, and a **Recursive AI Brain**.

```text
    ⚙️  ASHOS — Digital Cockpit v2.0
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    Control tasks, focus, and intelligence 
    directly from your terminal.
```

---

## 🌟 Vision
ASHOS transcends the "task manager" category. It's a cohesive **Intelligence Engine** that learns from your behavior, synchronizes across the cloud, and evolves its personality based on your work patterns.

### Core Architectural Pillars
- **Clean Architecture**: Decoupled business logic from infrastructure.
- **Atomic Persistence**: ACID-compliant storage via SQLite.
- **Recursive RAG**: Vector-based memory that grows as you work.
- **Event-Driven Extension**: Modules communicate via a high-performance internal bus.

---

## 🔥 Key Intelligence Features

### 🧠 The AI Brain (RAG-Powered)
ASHOS doesn't just store data; it understands it. Using **Local Vector Embeddings** (via `sqlite-vec`), ASHOS performs semantic search across your entire history.
- **Context-Aware Q&A**: "What was I working on when I mentioned the database fix?"
- **Behavioral Digest**: Strategic analysis of your productivity trends.
- **Standup Automation**: Professional status updates generated from your logs.

### ☁️ Cloud-Native Synchronization
Using a **Direct API Bypass** architecture with **Appwrite**, ASHOS keeps your notes, tasks, and AI context perfectly in sync across multiple machines.
- **Recursive Scanning**: Handles nested JSON structures with 100% reliability.
- **Personality Calibration**: Weekly syncs update the `About Me` context, allowing the AI to "know" you better every week.

### ⏱️ Deep Work & Focus
ASHOS implements a high-precision focus engine that tracks **concentration intervals** and provides ASCII analytics to visualize your "Flow State" sessions.

---

## 🛠️ Project Structure

```text
├── cmd/                # CLI Layer (Cobra Implementation)
├── internal/           # Core Logic (The Private Engine)
│   ├── ai/             # Intelligence Layer (RAG & Memory)
│   ├── task/           # Task Orchestration
│   ├── focus/          # Concentration Tracking
│   ├── core/event/     # Internal Pub/Sub Bus
│   └── storage/        # Multi-engine Persistence (SQLite + Appwrite)
├── data/               # Local ACID Persistence
└── main.go             # Composition Root (The Heart)
```

---

## 🚀 Getting Started

### Installation
Ensure you have **Go 1.24+** installed.

```bash
# Clone the OS
git clone https://github.com/Ashish5180/ASHOS-Command-Line-Interface-
cd ASHOS-Command-Line-Interface-

# Build the Engine
go build -o ash main.go

# Initialize
./ash status
```

---

## ⚡ Quick Start Directives
For a full list of commands, see the [📖 Command Guide](./COMMAND_GUIDE.md).

| Command | Action |
|:---|:---|
| `ash task add "Title"` | Log a milestone |
| `ash focus start` | Initiate a deep work cycle |
| `ash ai weekly-sync` | Synchronize OS Brain with Cloud |
| `ash dash` | Launch the Terminal Command Center |

---

## 💠 Tech Stack

| Component | Technology |
|:---|:---|
| **Core Runtime** | Go (Golang) |
| **Storage Engine** | SQLite (Pure Go) |
| **Cloud Backend** | Appwrite |
| **AI Processing** | OpenAI (Cloud) / Ollama (Local) |
| **Vector Search** | `sqlite-vec` (v0.1.0) |
| **TUI Interface** | Charm (Bubble Tea & Lipgloss) |

---

## 📈 Roadmap: The Path to Maturity
- [x] **Phase 1**: Foundation (CLI + SQLite)
- [x] **Phase 2**: Intelligence (RAG + Semantic Search)
- [x] **Phase 3**: Connectivity (Appwrite Cloud Sync + Slack)
- [ ] **Phase 4**: Automation (DSL Scripting + System Monitoring)

---

## 🤝 Positioning Statement
**ASHOS is built for the 1% developer.** It's about taking control of your cognitive load by offloading memory to an intelligent, terminal-native system. It demonstrates how to build enterprise-grade software in Go while maintaining the speed and simplicity of the CLI.

---

## 📄 License
Architectural demonstration project by **Ashish Yaduvanshi**. For educational and personal productivity use.
