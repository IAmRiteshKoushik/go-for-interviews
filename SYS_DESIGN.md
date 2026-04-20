# Senior Golang Interview — 2-Day Prep Roadmap

---

## Day 1 — Depth Pass (Core Theory + Design)

---

### 1. System Design / Architecture

| # | Topic | What to cover |
|---|-------|---------------|
| 1 | Distributed patterns | Sharding, consistent hashing, partition tolerance (CAP theorem in practice) |
| 2 | API design | REST vs gRPC trade-offs, idempotency keys, versioning strategies |
| 3 | Caching strategy | Write-through, write-behind, TTL invalidation, cache stampede prevention |
| 4 | Database trade-offs | ACID vs BASE, read replicas, index design, N+1 query patterns |
| 5 | Event-driven architecture | Kafka partitions, ordering guarantees, consumer groups, dead-letter queues |

**Key prep output:** Be able to design a system end-to-end on a whiteboard — e.g. "design a freight rate aggregator" or "design a real-time notification system". Know where each trade-off lives.

---

### 2. Component Architecture

| # | Topic | What to cover |
|---|-------|---------------|
| 1 | Interface design | Small interfaces, implicit satisfaction, designing for testability |
| 2 | Dependency injection | Manual DI, google/wire, uber/fx — pros/cons of each |
| 3 | Layered architecture | Handler → Service → Repository → Infrastructure separation |
| 4 | Plugin / strategy patterns | Go-idiomatic function types as first-class values vs interface dispatch |
| 5 | Package cohesion | When to split packages, when to flatten, preventing cyclic imports |

**Key prep output:** Be able to sketch a package structure for a non-trivial service and justify every boundary decision.

---

## Day 1 (cont.) → Day 2 Morning — Concurrency + Error Handling

---

### 3. Concurrency & Scalability

| # | Topic | What to cover |
|---|-------|---------------|
| 1 | Goroutine lifecycle | Goroutine leaks, stack growth, scheduler preemption, GOMAXPROCS tuning |
| 2 | Channel patterns | Fan-out/fan-in, pipeline, done channel, select semantics, nil channel tricks |
| 3 | sync primitives | Mutex vs RWMutex, WaitGroup, sync.Once, atomic ops (sync/atomic) |
| 4 | context.Context | Cancellation propagation, deadlines, value passing pitfalls, never store in structs |
| 5 | Worker pool sizing | CPU-bound vs I/O-bound heuristics, backpressure via buffered channel capacity |
| 6 | Data race detection | `go test -race`, common races (loop variable capture, concurrent map writes) |

**Key prep output:** Be able to implement a worker pool, a fan-out/fan-in pipeline, and a graceful shutdown loop from scratch in ~10 minutes.

**Your existing edge:** TokenManagerQueue (thundering herd / 401 refresh serialisation) is a production-grade concurrency answer. Have this story ready.

---

### 4. Error Handling & Edge Cases

| # | Topic | What to cover |
|---|-------|---------------|
| 1 | Error taxonomy | Sentinel errors vs typed errors vs opaque errors — when each is appropriate |
| 2 | Wrapping mechanics | `errors.Is` / `errors.As` / `%w` — unwrap chain behaviour |
| 3 | panic vs error | Only panic for true invariant violations; recover only at package/service boundaries |
| 4 | Partial failure | `golang.org/x/sync/errgroup` — first-error semantics vs aggregate error collection |
| 5 | Timeout edge cases | Context already cancelled on entry, zero-value timeouts, deadline already passed |
| 6 | Nil interface trap | Nil pointer inside a non-nil interface — the classic Go gotcha, must be able to explain and demonstrate |

**Key prep output:** Be able to write a function that fans out to N services concurrently, collects errors without leaking goroutines, and returns a meaningful aggregate error.

---

## Day 2 Afternoon — Production Readiness

---

### 5. Production Readiness

| # | Topic | What to cover |
|---|-------|---------------|
| 1 | Structured logging | `log/slog` / `zap` field-based logging, log levels, request-scoped trace IDs via context |
| 2 | Metrics & tracing | Prometheus counters/histograms/gauges, OpenTelemetry spans, cardinality discipline |
| 3 | Graceful shutdown | `os.Signal` handling, drain window, in-flight request tracking with WaitGroup |
| 4 | Health checks | Liveness vs readiness distinction, dependency probing (DB ping, queue lag) |
| 5 | Profiling | `net/http/pprof` endpoints, CPU/heap/goroutine profiles, reading flame graphs |
| 6 | Config management | Env vars, Viper, secret injection (never in source), zero-downtime config reloads |
| 7 | Testing pyramid | Table-driven unit tests, integration tests with `testcontainers-go`, contract tests |

**Key prep output:** Be able to describe exactly how a service you built handles a rolling deploy — signal received → stop accepting → drain → exit. No goroutine leaks, no dropped requests.

**Your existing edge:** TruckHai bootstrap gate (startup sequencing) and the SDUI mutation action handler with optimistic update + rollback are strong production-readiness stories. Have them ready as concrete examples.

---

## Meta-Strategy

For each topic, prepare **two levels of answer**:

- **Design answer** — whiteboard-style, component boxes, trade-offs articulated verbally.
- **Code answer** — something you could write in Go in under 10 minutes on a shared editor.

Senior interviews probe both in the same question. Starting with design and then dropping to code on request shows maturity.

### Time allocation (rough)

| Block | Topic | Hours |
|-------|-------|-------|
| Day 1 AM | System Design | 3h |
| Day 1 PM | Component Architecture | 2.5h |
| Day 1 Eve | Concurrency (theory) | 2h |
| Day 2 AM | Concurrency (coding) + Error Handling | 3.5h |
| Day 2 PM | Production Readiness + mock design session | 3h |

---

*Good luck — you have more production Go experience than most candidates. 
The goal is to translate what you've already built into interview-fluent 
language.*
