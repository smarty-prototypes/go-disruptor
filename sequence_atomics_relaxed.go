//go:build relaxed_atomics

package disruptor

import (
	"sync/atomic"
	"unsafe"
)

//go:linkname loadAcq32 internal/runtime/atomic.LoadAcq
func loadAcq32(ptr *uint32) uint32

//go:linkname storeRel32 internal/runtime/atomic.StoreRel
func storeRel32(ptr *uint32, val uint32)

func int32Load(addr *atomic.Int32) int32 {
	return int32(loadAcq32((*uint32)(unsafe.Pointer(addr))))
}
func int32Store(addr *atomic.Int32, val int32) {
	storeRel32((*uint32)(unsafe.Pointer(addr)), uint32(val))
}
