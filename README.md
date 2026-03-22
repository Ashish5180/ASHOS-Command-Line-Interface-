# ASHOS — Event-Driven Personal CLI Operating System

> A modular, developer-first, terminal-native personal operating system built in Go.  
> Designed with production-grade architecture: dependency injection, interface-driven design, concurrency safety, atomic persistence, and extensibility.

```
⚙️  ASHOS — Personal CLI Operating System
   Your digital cockpit.
   Control tasks, focus, system insights, and productivity intelligence.
```

---

## Vision

ASHOS is **not** a todo CLI.

It is a modular, event-driven, extensible **personal operating system** designed for developers who think in systems.

**It aims to demonstrate:**

- Clean Architecture principles
- Composition Root pattern
- Interface-first design
- Pluggable infrastructure
- Concurrency-safe state management
- Extensible module system
- Event-driven orchestration

---

## Table of Contents

- [Core Features](#core-features)
- [Architecture](#architecture)
- [Storage Engine Design](#storage-engine-design)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Engineering Practices](#engineering-practices)
- [Upcoming Features](#upcoming-features)
- [Tech Stack](#tech-stack)
- [Engineering Principles](#engineering-principles)
- [Project Maturity Roadmap](#project-maturity-roadmap)
- [Positioning Statement](#positioning-statement)

---

## Core Features

### 1. Task Management

Full CRUD lifecycle via CLI.

| Command | Description |
|---|---|
| `ash task add <title>` | Create a new task |
| `ash task list` | List all tasks with status |
| `ash task done <id>` | Mark a task as completed |
| `ash task delete <id>` | Permanently delete a task |

**Engineering Highlights:**

- Auto-increment ID generation (scans existing tasks for max ID)
- Input validation (empty title, max 200 chars)
- Invariant enforcement (no double-complete)
- Repository pattern with storage abstraction via `storage.Store`
- Event-driven notifications (broadcasts `TaskCreated`, `TaskCompleted`, etc.)
- SQLite persistence with ACID guarantees

---

### 2. Focus Mode

| Command | Description |
|---|---|
| `ash focus start` | Start a focus session |
| `Ctrl+C` | Gracefully stop and persist the session |

**Engineering Highlights:**

- Background goroutine with `time.Ticker` for second-level tracking
- Minute-level progress updates printed to stdout
- Signal-safe shutdown (`SIGINT`, `SIGTERM`)
- Persisted session history via storage engine
- Aggregated daily focus calculation in system status

---

### 3. System Status

| Command | Description |
|---|---|
| `ash status` | Show uptime, pending tasks, and today's focus time |

**Displays:**

- ASHOS process uptime
- Pending task count (aggregated from task service)
- Focus time today (summed from persisted sessions)

---

### 4. Interactive Dashboard (TUI)

Built with **Bubble Tea** (Elm architecture).
| Command | Description |
|---|---|
| `ash dash` | Launch the interactive dashboard |

**Features:**
- **Two-panel layout** — Tasks panel + System stats panel
- **Focus status indicator** — Shows `FOCUSING...` when active
- **Auto-refresh** — Updates every 5 seconds via `tea.Tick`

---

### 5. AI Brain (Intelligence Layer)
ASHOS includes a powerful RAG-based AI system that analyzes your productivity data.

| Command | Description |
|---|---|
| `ash ai ask <query>` | Context-aware Q&A with relevance scores |
| `ash ai daily-digest` | **NEW**: Advanced behavioral analysis and productivity scores |
| `ash ai standup` | Auto-generate professional daily standups from logs |
| `ash ai suggest` | Smart task recommendations based on historical patterns |

**Engineering Highlights:**
- **Vector Search (RAG)**: Uses `sqlite-vec` (vec0) for local vector embeddings and similarity search.
- **Content Deduplication**: MD5-based hashing prevents duplicate embedding of unchanged content.
- **Local Fallback**: Supports **Ollama** (`llama3.2`) for offline intelligence via the `--local` flag.
- **Conversation Memory**: Maintains state for multi-turn dialogues (last 10 turns).
- **Graceful Health Checks**: Automatically pings local LLM and falls back to cloud if unavailable.

---

## Architecture

### High-Level Overview

```
┌──────────────┐
│   main.go    │  ← Composition Root (DI wiring & Event Subscriptions)
└──────┬───────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│   cmd/       │────▶│   App struct  │  ← Dependency container
│  (Cobra CLI) │     │ (all services)│
└──────┬───────┘     └──────────────┘
       │
       ▼
┌──────────────────────────────────────────┐
│         internal/ (business logic)       │
│                                          │
│  ┌────────────┐  ┌──────────────┐        │
│  │   task/    │  │   focus/     │        │
│  │  Service   │  │  Manager     │        │
│  │  Repository│  │  Session     │        │
│  └─────┬──────┘  └──────┬───────┘        │
│        │                │                │
│        ▼                ▼                │
│  ┌─────────────────────────────┐         │
│  │       storage.Store         │         │ ← Interface
│  │  (SQLite engine / Memory)   │         │
│  └─────────────────────────────┘         │
│                                          │
│  ┌────────────┐  ┌──────────────┐        │
│  │  system/   │  │    ai/       │        │
│  │  Service   │  │  Brain (RAG) │        │
│  └────────────┘  └──────────────┘        │
│                                          │
│      📡  internal/core/event             │ ← Event Bus
└──────────────────────────────────────────┘
       │
       ▼
┌──────────────┐
│   data/      │  ← SQLite database
│  ashos.db    │
└──────────────┘
```

### Layered Architecture

```
CLI (cmd/)  →  Service (internal/*/service.go)  →  Repository  →  Storage Engine
   ↑                    ↑                              ↑               ↑
Presentation      Business Logic                Data Contract    Infrastructure
```

### Dependency Flow

```
main.go → storage.Store → task.Repository → task.Service ─┐
                        → focus.Manager ───────────────────┤
                        → system.Service ──────────────────┤
                                                           ▼
                                                    cmd.App (container)
                                                           │
                                                    cmd.RegisterCommands()
```

### Composition Root

All dependencies are wired in `main.go` — no module constructs its own dependencies.

**This allows:**

- Swapping storage engines (JSON → SQLite) by changing **one line**
- Injecting test implementations
- Future plugin injection
- Event system integration

---

## Storage Engine Design

### `Store` Interface

The universal data contract. Every module talks **only** to this interface.

```go
type Store interface {
    Save(ctx context.Context, collection, key string, value any) error
    Get(ctx context.Context, collection, key string, dest any) error
    Delete(ctx context.Context, collection, key string) error
    List(ctx context.Context, collection string, dest any) error
    Close() error
}
```

### Implementations

| Engine | Use Case | Description |
|---|---|---|
| `sqliteStore` | Production | Relational persistence via SQLite (CGO-free) |

### Guarantees

- **Atomic file writes** — Write to temp file, then rename (crash-safe)
- **Concurrency protection** — `sync.RWMutex` on all operations
- **Typed errors** — `ErrNotFound`, `ErrCorruptData`, `ErrReadFailed`, `ErrWriteFailed`
- **Context-ready API** — All methods accept `context.Context` for future timeouts/tracing

---

## Project Structure

```
.
├── main.go                     # Composition root — all DI wiring happens here
├── go.mod                      # Module: ashos
├── ash                         # Compiled binary
│
├── cmd/                        # CLI layer (Cobra commands)
│   ├── root.go                 # Root "ash" command + App struct + RegisterCommands
│   ├── task.go                 # ash task {add,list,done,delete}
│   ├── status.go               # ash status
│   ├── focus.go                # ash focus start
│   ├── dash.go                 # ash dash (TUI launcher)
│   ├── note.go                 # (placeholder — future)
│   └── sys.go                  # (placeholder — future)
│
├── internal/                   # Private application code
│   ├── task/                   # Task module (Service + Repository)
│   │   ├── model.go            # Task struct
│   │   ├── service.go          # Business logic (Validation + Event Publishing)
│   │   ├── repository.go       # Repository interface
│   │   └── store_repo.go       # StoreRepository adapter
│   │
│   ├── focus/                  # Focus tracking module
│   │   ├── manager.go          # Focus manager (concurrency + events)
│   │   └── model.go            # SessionRecord structs
│   │
│   ├── core/                   # Shared system internals
│   │   └── event/              # 📡 Event Bus system (Pub/Sub)
│   │       ├── bus.go          # Async event engine
│   │       └── types.go        # Shared event types
│   │
│   ├── storage/                # Universal Storage engine
│   │   ├── storage.go          # Store interface
│   │   ├── sqlite.go           # SQLite engine (ACID, CGO-free)
│   │   └── error.go            # Typed storage errors
│   │
│   ├── system/                 # System service (Uptime, Aggregates)
│   └── ui/                     # Bubble Tea TUI components
│
├── data/                       # Runtime persistence
│   └── ashos.db                # SQLite 3 database file
└── ash                         # Final compiled binary
```

---

## Getting Started

### Prerequisites

- **Go 1.24+** (module uses `go 1.24.5`)
- A terminal that supports ANSI colors (for the TUI dashboard)

### Build & Run

```bash
# Clone the repository
git clone <repo-url> && cd Personal-os

# Build the binary
go build -o ash .

# Verify
./ash --help
```

### Quick Start

```bash
# Add some tasks
./ash task add "Read Clean Architecture chapter 5"
./ash task add "Refactor storage layer"
./ash task add "Write unit tests"

# List tasks
./ash task list

# Complete a task
./ash task done 1

# Check system status
./ash status

# Launch the interactive dashboard
./ash dash

# Start a focus session (Ctrl+C to stop)
./ash focus start
```

---

## Usage

### Task Commands

```bash
# Add a task (multi-word titles work with or without quotes)
./ash task add Buy groceries
./ash task add "Read 10 pages of Go book"

# List all tasks with pending count
./ash task list

# Mark task #3 as done
./ash task done 3

# Delete task #2 permanently
./ash task delete 2
```

### Focus Commands

```bash
# Start a focus session (runs until Ctrl+C)
./ash focus start
# → Prints minute-level progress: "⏱️ Focus Progress: 5m0s"
# → Ctrl+C gracefully saves the session
```

### Status & Dashboard

```bash
# Quick text-based status
./ash status

# Full interactive TUI dashboard (auto-refreshes every 5s)
./ash dash
# → Press 'r' to refresh manually, 'q' to quit
```

---

## Engineering Practices

### 1. Dependency Injection (Composition Root Pattern)

All dependencies are created and wired in `main.go` — the **single composition root**. No module constructs its own dependencies. Swapping the storage engine from JSON to SQLite requires changing **only** `main.go`.

```go
store, _ := storage.NewSQLiteStore("data/ashos.db")
taskRepo  := task.NewStoreRepository(store)
taskSvc   := task.NewService(taskRepo)
```

### 2. Interface-Driven Design

The `storage.Store` interface is the universal data contract. All modules depend on the **interface**, never on concrete implementations. The primary engine is `sqliteStore`.

### 3. Repository Pattern

The task module uses a **Repository interface** to decouple business logic from storage:

```go
type Repository interface {
    GetAll() ([]Task, error)
    SaveAll([]Task) error
}
```

`StoreRepository` acts as an **adapter** bridging `task.Repository` ↔ `storage.Store`. The service layer never knows how or where data is stored.

### 4. Layered Architecture

Each layer has a single responsibility:

| Layer | Location | Responsibility |
|---|---|---|
| Presentation | `cmd/` | Parse args, format output |
| Business Logic | `internal/*/service.go` | Validation, rules, orchestration |
| Data Contract | `internal/*/repository.go` | Data access interface |
| Infrastructure | `internal/storage/` | Serialization, file I/O, concurrency |

### 5. Typed Error Handling

Custom error types enable callers to programmatically distinguish error categories:

```go
ErrNotFound     // key doesn't exist → return empty list, don't crash
ErrStorageInit  // engine failed to start
ErrReadFailed   // I/O read error
ErrWriteFailed  // I/O write error
ErrCorruptData  // deserialization failure
```

Uses `errors.As()` for type-safe error matching — no string comparison.

### 6. ACID Transactions
SQLite provides Atomic, Consistent, Isolated, and Durable guarantees. Even if the process crashes mid-write, the database remains in a valid state.

### 7. Concurrency Safety

- SQLite engine uses internal locking and connection pooling for safety.
- Focus session state is protected by its own `sync.Mutex`
- The in-memory `TaskManager` model is also mutex-guarded

### 8. Graceful Signal Handling

Focus sessions catch `SIGINT`/`SIGTERM` via Go's `os/signal` package:

1. Session is persisted to storage
2. Stop signal is sent to the background goroutine
3. Goroutine cleans up and exits

### 9. Go Project Layout Conventions

| Directory | Purpose |
|---|---|
| `internal/` | Private packages — cannot be imported externally |
| `cmd/` | CLI command definitions (Cobra pattern) |
| `pkg/` | Public utilities (future shared rendering) |
| `data/` | Runtime data directory (auto-created, `.gitignore`-able) |

### 10. Elm Architecture for TUI (Bubble Tea)

The dashboard follows the **Model-Update-View** pattern:

- **`Init()`** — Fetch initial data + start tick timer
- **`Update()`** — Handle key events, window resize, data messages
- **`View()`** — Pure function rendering the entire UI from model state

No side effects in `View()`. All I/O happens through **commands** (`tea.Cmd`).

### 11. Responsive Terminal UI

- Side-by-side panels when width ≥ 60 columns
- Stacked layout when width < 60 columns
- Fallback width of 80 when terminal size is unavailable

### 12. Context-Ready Storage API

All `Store` methods accept `context.Context` — ready for timeouts, cancellation, and tracing when the storage engine evolves to network-backed stores.

---

## Upcoming Features

> The next evolution — from CLI tool to event-driven personal OS.

### 1. Event Bus System (Implemented)

ASHOS uses an **internal event-driven architecture** to decouple modules.

**Why:** Loose coupling, extensibility, cross-module reactions without direct dependencies.

```go
type EventBus struct {
    subscribers map[reflect.Type][]Handler
    mu          sync.RWMutex
}
```

**Events Supported:** `TaskCreated`, `TaskCompleted`, `TaskDeleted`, `FocusStarted`, `FocusEnded`.

**Use Case:** When a `TaskCompleted` event is published, the system can automatically trigger UI updates, productivity scoring, or notifications without the Task module knowing about them.

---

### 2. SQLite Storage Engine

Drop-in replacement via the `storage.Store` interface:

```go
storage.NewSQLiteStore("data.db")
```

**Benefits:** ACID guarantees, structured queries, indexing, faster reads.

**Architecture Rule:** Swapping JSON → SQLite requires modifying **only** `main.go`. Zero business logic changes.

---

### 3. Focus Analytics Engine (Implemented)

Transform focus tracking into **behavioral intelligence** with visual charts and scoring.

- **Trend Analysis**: ASCII line graphs via `asciigraph`.
- **Daily Heatmap**: Visual bar charts showing focus intensity per day.
- **Deep Work Score**: Algorithmically calculated productivity rating (0-100).
- **Consistency Tracker**: Streak detection for focus habits.

**Commands:**
```bash
ash focus analytics
```

**Features:**
- 🔥 Primary color highlights for deep work (1h+)
- 🏆 Gold highlights for elite focus (4h+)
- ⏱️ Accurate to-the-minute tracking


---

### 4. Sprint Mode

Track productivity by sprint:

```bash
ash sprint start "Feature X"
ash sprint end
```

**Tracks:** Tasks completed, focus time, session density, productivity score.

```
Sprint Report:
  Tasks Completed: 4
  Focus Time:      2h 34m
  Deep Work Score:  87
```

`SprintManager` aggregates `TaskService`, `FocusManager`, and `EventBus`.

---

### 5. Plugin System

Long-term vision — pluggable modules:

```go
type Plugin interface {
    Name() string
    Init(app *App) error
}
```

**Goals:** Load dynamic modules, register new CLI commands, subscribe to events, extend ASHOS without core modification.

```bash
ash plugin install <module>
```

---

### 6. System Monitoring Module

```bash
ash sys
```

**Planned:** CPU usage, memory usage, disk usage, OS info, Go runtime stats, ASHOS internal stats.  
**Purpose:** Reinforce the "Personal OS" identity.

---

### 7. Scriptable DSL (Experimental)

Allow automation via `.ash` scripts:

```bash
# morning.ash
focus 25
task add "Review PR"
status
```

```bash
ash run morning.ash
```

**Introduces:** Mini interpreter, command parser abstraction, reusable CLI engine.

---

### 8. Testing Strategy Expansion

| Test Type | Scope |
|---|---|
| Unit tests | Service logic (using `memoryStore`) |
| Integration tests | Storage engine I/O |
| Concurrency tests | Mutex & goroutine safety |
| Event bus tests | Publish/subscribe correctness |

**Using:** `MemoryStore`, table-driven tests, context cancellation simulation.

---

### 9. Dashboard Upgrade

- Global status bar with live clock
- Animated focus spinner
- Keyboard legend footer
- Dynamic refresh via event bus instead of polling

---

## Tech Stack

| Component | Tool | Purpose |
|---|---|---|
| CLI Framework | [Cobra](https://github.com/spf13/cobra) | Command tree, args parsing, help generation |
| TUI Framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Interactive terminal UI (Elm architecture) |
| TUI Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Declarative terminal styles, colors, borders |
| Charts | [Asciigraph](https://github.com/guptarohit/asciigraph) | ASCII line charts for trend visualization |
| Storage | SQLite (`modernc.org/sqlite`) | ACID-compliant relational persistence (CGO-free) |
| Event System | Internal Event Bus | Asynchronous Pub/Sub (stdlib + reflection) |
| Concurrency | `sync`, `time`, `os/signal` (stdlib) | Mutexes, tickers, signal handling |

**Zero external runtime dependencies beyond Cobra + Charm libraries.** Everything else is Go stdlib.

---

## Engineering Principles

| # | Principle | Applied Where |
|---|---|---|
| 1 | Composition Root Pattern | `main.go` — single DI wiring point |
| 2 | Interface-Driven Design | `storage.Store`, `task.Repository` |
| 3 | Repository Pattern | `task.Service` ↔ `task.Repository` ↔ `storage.Store` |
| 4 | Layered Architecture | `cmd/` → Service → Repository → Storage |
| 5 | Event-Driven Extension | Planned `EventBus` for cross-module reactions |
| 6 | Infrastructure Swappability | SQLite support via interface |
| 7 | Concurrency Safety | `sql.DB` connection pooling & mutexes |
| 8 | ACID Persistence | SQLite's native transaction safety |
| 9 | Context Propagation | `context.Context` on all `Store` methods |
| 10 | Single Responsibility per Layer | Each layer owns one concern |

---

## Project Maturity Roadmap

### Phase 1 — Foundation (Current)

- [x] Stable CLI with Cobra
- [x] JSON storage engine with atomic writes
- [x] Interactive TUI dashboard (Bubble Tea)
- [x] Focus session persistence
- [x] Composition root DI pattern
- [x] Interface-driven storage abstraction
- [x] Concurrency-safe state management

### Phase 2 — Intelligence (Completed)

- [x] SQLite storage engine
- [x] Event bus system
- [x] Focus analytics engine
- [x] RAG-based AI Brain (OpenAI + Ollama)
- [x] Vector deduplication & schema versioning
- [x] Behavioral Daily Digest

### Phase 3 — Productivity (In Progress)

- [x] Manual Re-ingestion tool
- [x] Sprint mode (Partial)
- [ ] Plugin system
- [ ] System monitoring module (`ash sys`)
- [ ] DSL scripting (`ash run`)

### Phase 4 — Scale

- [ ] Remote sync (optional)
- [ ] Distributed store support
- [ ] Observability hooks
- [ ] Metrics export

---

## Positioning Statement

**ASHOS demonstrates:**

- **Clean Architecture** in Go
- **Production-grade persistence** patterns
- **Event-driven design** for extensibility
- **Concurrency-safe** system programming
- **Extensible module architecture**
- **Developer-focused** tooling mindset

> *A developer-centric, extensible personal operating system built with production-grade architectural principles.*

---

## License

Personal project. Intended for learning, experimentation, and architectural demonstration.
