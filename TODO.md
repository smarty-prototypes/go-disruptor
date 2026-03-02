### TODO Items:

#### Metrics
- Max number of messages/events to handle per invocation to Handle()

#### Create a compiler flag to enable use of storeRel / loadAcq semantics

```
import _ "unsafe" // for go:linkname

//go:linkname loadAcq32 internal/runtime/atomic.LoadAcq
func loadAcq32(ptr *uint32) uint32

//go:linkname storeRel32 internal/runtime/atomic.StoreRel
func storeRel32(ptr *uint32, val uint32)
```

Then make the following changes:
```
// OLD
storeRel32(&this.committedSlots[lower&mask], uint32(lower>>this.shift))
this.committedSlots[lower&mask].Load() != int32(lower>>this.shift)

// NEW
this.committedSlots[lower&mask].Store(int32(lower >> this.shift))
loadAcq32(&this.committedSlots[lower&mask]) != uint32(lower>>this.shift)
```
