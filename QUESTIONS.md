Here are the most common concurrency questions for senior Go engineers in live coding rounds:

---

## Goroutine & Channel Fundamentals

1. **Fan-out / Fan-in** — Distribute work across N goroutines, collect results into one channel
2. **Pipeline pattern** — Chain stages where each stage reads from upstream and writes to downstream channel
3. **Worker pool** — Fixed N workers processing jobs from a shared channel; handle graceful shutdown
4. **Timeout on channel receive** — Use `select` + `time.After` to avoid blocking forever
5. **Done channel / cancellation** — Propagate cancellation to goroutines using `context.WithCancel`

---

## Synchronization Primitives

6. **WaitGroup usage** — Common trap: calling `wg.Add` inside the goroutine instead of before launch
7. **Mutex vs RWMutex** — Implement a thread-safe cache; know when read-heavy workloads benefit from `RWMutex`
8. **`sync.Once`** — Implement a singleton / lazy-initialized resource
9. **`sync.Map` vs `map + RWMutex`** — Trade-offs; when does `sync.Map` win?
10. **Semaphore pattern** — Rate-limit goroutine concurrency using a buffered channel as a semaphore

---

## Classic Problems

11. **Bounded concurrency** — Process 1000 URLs with at most 10 concurrent HTTP calls
12. **Merge N channels** — Write `merge(channels ...chan T) chan T`
13. **Or-done channel** — Wrap a channel so it respects a done signal
14. **Rate limiter** — Implement token bucket or leaky bucket using tickers
15. **Publish-subscribe broker** — Multiple subscribers receive from a single event stream

---

## Error Handling in Concurrent Code

16. **`errgroup`** — Use `golang.org/x/sync/errgroup` to run N goroutines and collect the first error
17. **Error propagation via channel** — Return `(result, error)` pairs through a result struct on a channel
18. **Panic recovery in goroutines** — Panics don't propagate across goroutine boundaries; how do you handle them?

---

## Common Traps (High signal for senior level)

19. **Loop variable capture** — Classic `go func() { use(i) }()` bug; fixed with `i := i` or Go 1.22+ semantics
20. **Goroutine leak** — Identify and fix a goroutine blocked on a channel nobody reads from
21. **Deadlock diagnosis** — Two goroutines each waiting for the other; spot it, explain the runtime's detection
22. **Closing a closed channel** — Panics; pattern: only the sender closes, or use `sync.Once`
23. **Unbuffered vs buffered channel semantics** — When does a send block? Explain with a concrete example

---

## Advanced (Distinguishes senior from mid-level)

24. **`context` propagation** — Thread `ctx` through a call chain; cancel all downstream work on timeout
25. **Atomic operations** — Use `sync/atomic` for a lock-free counter; when to prefer it over mutex
26. **Memory model / happens-before** — Why can a goroutine see a stale value even without a data race?
27. **Select with default** — Non-blocking channel operations; implement a try-send/try-receive
28. **Backpressure** — What happens when producers outpace consumers? How do you design for it?

---
