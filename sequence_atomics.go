//go:build !relaxed_atomics

package disruptor

import "sync/atomic"

func int32Load(addr *atomic.Int32) int32    { return addr.Load() }
func int32Store(addr *atomic.Int32, val int32) { addr.Store(val) }
