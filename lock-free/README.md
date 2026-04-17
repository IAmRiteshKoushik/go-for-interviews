# Checklist for Lock-Free Data Structures

- [] Shared state should be minimal
- [] Published nodes should become immutable
- [] Every operation needs a clear linearization point
- [] Every CAS path failure must preserve invariants
- [] Avoid reclaimatoin problems by not reusing nodes early
- [] Test under contention
- [] Benchmark against mutexes
