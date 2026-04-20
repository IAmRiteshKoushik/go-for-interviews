# Go Concurrency — Interview Question Bank

> **Target:** Senior Go Developer | Live Coding Round
> **Topics:** Go Concurrency Patterns · sync Package · context Package · Worker Pool Pattern

---

## Low Difficulty

### Q1. Goroutines vs threads — explain the difference
**Prompt:** What is a goroutine? How does it differ from an OS thread? What 
does the Go runtime scheduler do that makes goroutines cheap?

**Tags:** `goroutines` `runtime` `scheduling`

---

### Q2. What does `go func(){}()` actually do?
**Prompt:** Walk me through what happens when you prefix a function call with 
`go`. What is the lifecycle of that goroutine? When does it terminate?

**Tags:** `goroutines` `lifecycle`

---

### Q3. Buffered vs unbuffered channels
**Prompt:** Explain the difference between `make(chan int)` and 
`make(chan int, 5)`. When would you choose one over the other? What happens if 
a sender sends to a full buffered channel?

**Tags:** `channels` `blocking`

---

### Q4. What is a WaitGroup and how do you use it?
**Prompt:** Explain `sync.WaitGroup`. Write a snippet that launches 5 goroutines, 
each printing their index, and blocks the main goroutine until all are done.

**Tags:** `sync` `goroutines` `WaitGroup`

---

### Q5. What is a race condition? How do you detect one in Go?
**Prompt:** Define a race condition. What flag do you pass to `go test` or 
`go run` to detect it? Show a minimal example of a data race.

**Tags:** `race conditions` `sync` `tooling`

---

### Q6. Closing a channel — rules and pitfalls
**Prompt:** When should you close a channel? What happens if you send on a 
closed channel? What happens if you receive from a closed channel? Who should 
be responsible for closing?

**Tags:** `channels` `panic` `ownership`

---

## Medium Difficulty

### Q7. Implement a fan-out / fan-in pipeline
**Prompt:** You have a stream of integers on an input channel. Fan-out to 3 
worker goroutines that square each number, then fan-in results into a single 
output channel. Write the full implementation.

**Tags:** `channels` `goroutines` `fan-out` `fan-in` `pipelines`

---

### Q8. Build a bounded worker pool
**Prompt:** Implement a worker pool that processes a slice of jobs 
(type `func() error`) with a maximum of N concurrent workers. Collect all 
errors. The pool must shut down cleanly after all jobs are done.

**Tags:** `worker-pool` `goroutines` `channels` `sync` `error handling`

---

### Q9. Use a Mutex to protect a shared counter
**Prompt:** Write a struct `SafeCounter` with an `Inc()` and `Value()` method 
that is safe for concurrent use. 1000 goroutines each call `Inc()` once. Assert 
the final value is 1000.

**Tags:** `sync` `Mutex` `race conditions`

---

### Q10. Context cancellation — propagate a cancel signal
**Prompt:** You have a function `doWork(ctx context.Context)` that runs an 
infinite loop. Wire up `context.WithCancel` in the caller so that pressing 
Ctrl+C (os.Signal) cancels the context and `doWork` exits cleanly.

**Tags:** `context` `cancellation` `goroutines` `signals`

---

### Q11. Context timeout vs deadline
**Prompt:** What is the difference between `context.WithTimeout` and 
`context.WithDeadline`? Write a function that calls a slow external API but 
returns an error if it takes longer than 2 seconds.

**Tags:** `context` `timeout` `deadline`

---

### Q12. Select statement with a done channel
**Prompt:** Implement a generator that sends fibonacci numbers on a channel 
indefinitely. The caller must be able to stop it at any time using a done 
channel passed via context. Use a `select` with a timeout of 500ms per value.

**Tags:** `channels` `select` `context` `goroutines`

---

### Q13. sync.Once — implement a thread-safe singleton
**Prompt:** Implement a config loader that reads from disk. The load must 
happen exactly once even if 100 goroutines call `GetConfig()` concurrently. 
Use `sync.Once`.

**Tags:** `sync` `Once` `goroutines`

---

### Q14. Detect and fix a goroutine leak
**Prompt:** The following code starts a goroutine for every HTTP request that 
it never cleans up. Identify the leak, explain why it occurs, and refactor to 
prevent it. (Caller provides broken code.)

**Tags:** `goroutines` `leaks` `context` `lifecycle`

---

## High Difficulty

### Q15. Worker pool with dynamic concurrency and graceful shutdown
**Prompt:** Build a worker pool where: (1) concurrency is configurable at 
runtime via a channel, (2) in-flight jobs complete before shutdown, (3) new 
jobs are rejected post-shutdown with a typed error, and (4) all errors are 
aggregated and returned. No goroutine leaks.

**Tags:** `worker-pool` `channels` `context` `sync` `error handling` `shutdown`

---

### Q16. Rate-limited worker pool with backpressure
**Prompt:** Implement a worker pool that processes at most R requests per 
second (token bucket or ticker-based). If the pool is saturated, callers should 
block — not drop jobs. Support context cancellation for waiting callers.

**Tags:** `worker-pool` `rate limiting` `channels` `context` `backpressure`

---

### Q17. errgroup — concurrent subtask orchestration
**Prompt:** You need to call 4 microservices in parallel. If any one fails 
the whole operation should cancel immediately and return the first error. 
Use `golang.org/x/sync/errgroup`. Explain how errgroup compares to manual 
WaitGroup + channel error collection.

**Tags:** `errgroup` `context` `goroutines` `error handling` `worker-pool`

---

### Q18. Implement a semaphore using a buffered channel
**Prompt:** Without using any external library, implement a `Semaphore` type 
with `Acquire(ctx)` and `Release()` methods. Acquire should block if the 
semaphore is at capacity and respect context cancellation. Write a test that 
proves it.

**Tags:** `sync` `channels` `context` `semaphore` `concurrency patterns`

---

### Q19. Pipeline with context propagation and partial failure
**Prompt:** Build a 3-stage pipeline: fetch → transform → persist. Each stage 
runs concurrently. If any stage returns an error mid-stream, the entire 
pipeline must cancel cleanly without goroutine leaks. Use context throughout.

**Tags:** `pipelines` `context` `channels` `goroutines` `error handling` `cancellation`

---

### Q20. sync.Map vs Mutex-guarded map — when and why
**Prompt:** Explain the internal design of `sync.Map`. When does it outperform 
a `sync.RWMutex`-guarded map and when does it not? Write a benchmark that 
demonstrates the tradeoff under high read contention.

**Tags:** `sync` `sync.Map` `Mutex` `benchmarking` `concurrency patterns`

---

## Recommended Study Order

Given the topic emphasis on worker pools, context propagation, and the 
sync package:

1. **Worker pool progression** → Q8 → Q15 → Q16
2. **Context propagation** → Q10 → Q11 → Q19
3. **sync package depth** → Q9 → Q13 → Q20
4. **Channel & select fluency** → Q7 → Q12
5. **Foundations** → Q1 → Q2 → Q3 → Q4 → Q5 → Q6

---

*Generated for FAANG-track Senior Go Developer live coding round prep.*
