## TODO Items:

### Metrics
- Max number of messages/events to handle per invocation to Handle()

#### No panic recovery in consumers

If a `Handler.Handle()` panics, it kills the goroutine. `defer waiter.Done()` in `compositeListener.Listen()` fires, but the panic propagates and crashes the process. The Java LMAX has `ExceptionHandler` on `BatchEventProcessor` that catches exceptions, invokes a callback, and optionally continues. Here, a single misbehaving handler takes down everything with no hook point.

#### `sharedSequencer.Reserve` — irrevocable reservation blocks forever on shutdown

`sequencer_shared.go:80`: `this.reservedSequence.Add(slots)` is irrevocable. If the buffer is full and `Close()` is called (consumers stop draining), the producer spins in the slow-path loop forever. There's no exit path — the wait strategy has no cancellation signal, and the loop has no check for closed state. Java's `MultiProducerSequencer.next()` uses a CAS loop that could be made interruptible; the `Add` approach here trades that for throughput.

This is the most dangerous issue in the codebase. A producer calling `Reserve` on a full buffer after `Close()` will hang the goroutine permanently.

#### Cache line layout is x86-centric

`defaultSequencer` is annotated as "64B total — one cache line" and `sharedSequencer` as "128B — two cache lines." On Apple Silicon (128B lines), both structs fit in a single cache line, which means the hot/cold field separation in `sharedSequencer` provides no isolation. On 32B platforms (ARM32, MIPS), `defaultSequencer` spans two cache lines and the hot fields straddle the boundary. The code is correct everywhere but the performance optimization only works as intended on x86.

#### `newSequence()` retry loop is unbounded

`sequence.go:15`: The allocation loop retries `new(atomicSequence)` until it gets a cache-line-aligned pointer, with no iteration limit. Go's allocator will almost always return aligned pointers for these sizes, but there's no guarantee. An allocation failure or pathological allocator behavior would spin forever. A bounded retry with an over-allocate-and-align fallback would be more robust.

#### Error sentinels as magic int64 values

`ErrReservationSize = -1` and `ErrCapacityUnavailable = -2` are bare int64 constants, not `error` types. A caller that forgets to check the return value will use -1 or -2 as a sequence number, silently indexing the ring buffer at `(-1 & mask)` or `(-2 & mask)` — corrupting data with no signal. Java throws `InsufficientCapacityException`. The current API makes misuse easy and silent.

#### Redundant barrier check for first handler group

In `listener.go:60`, for the first handler group, `upstreamBarrier` and `committedBarrier` are the same object (both set to `committedBarrier` in `config.go:34`). So the gated branch (`else if upperSequence = this.committedBarrier.Load(...)`) is unreachable for group 0 — it will always return the same result as the first check. Minor waste, not a bug.

#### Compared to Java LMAX — missing features

| Feature | Java LMAX | This port |
|---|---|---|
| **WorkerPool** (events partitioned across handlers) | Yes | No — fan-out only (every handler sees every event) |
| **Handler panic/exception recovery** | `ExceptionHandler` callback | Process crash |
| **Dependency DAGs** | Arbitrary via `SequenceBarrier` | Linear pipeline only |
| **Blocking Reserve timeout / cancellation** | Via `WaitStrategy` + `AlertException` | Blocks forever, no exit |
| **`SequenceReportingEventHandler`** (mid-batch progress) | Yes | Always publishes full batch |
| **EventTranslator** (structured publish API) | Yes | User manages ring buffer directly |

The biggest functional gap is **WorkerPool** — many real-world use cases need work distribution (one event to one handler), not just fan-out. The biggest robustness gap is the inability to interrupt a blocked `Reserve`.
