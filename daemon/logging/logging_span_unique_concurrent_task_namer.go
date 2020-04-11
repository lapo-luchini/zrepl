package logging

import (
	"fmt"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/willf/bitset"
)

type uniqueConcurrentTaskNamer struct {
	mtx    sync.Mutex
	active map[string]*bitset.BitSet

	bitvecLengthGauge *prometheus.GaugeVec
}

// bitvecLengthGauge may be nil
func newUniqueTaskNamer(bitvecLengthGauge *prometheus.GaugeVec) *uniqueConcurrentTaskNamer {
	return &uniqueConcurrentTaskNamer{
		active:            make(map[string]*bitset.BitSet),
		bitvecLengthGauge: bitvecLengthGauge,
	}
}

// appends `#%d` to `name` such that until `done` is called,
// it is guaranteed that `#%d` is not returned a second time for the same `name`
func (namer *uniqueConcurrentTaskNamer) UniqueConcurrentTaskName(name string) (uniqueName string, done func()) {
	if strings.Contains(name, "#") {
		panic(name)
	}
	namer.mtx.Lock()
	act, ok := namer.active[name]
	if !ok {
		act = bitset.New(64) // FIXME magic const
		namer.active[name] = act
	}
	id, ok := act.NextClear(0)
	if !ok {
		// if !ok, all bits are 1 and act.Len() returns the next bit
		id = act.Len()
		// FIXME unbounded growth without reclamation
	}
	act.Set(id)
	namer.mtx.Unlock()

	if namer.bitvecLengthGauge != nil {
		namer.bitvecLengthGauge.WithLabelValues(name).Set(float64(act.Len()))
	}

	return fmt.Sprintf("%s#%d", name, id), func() {
		namer.mtx.Lock()
		defer namer.mtx.Unlock()
		act, ok := namer.active[name]
		if !ok {
			panic("must be initialized upon entry")
		}
		act.Clear(id)
	}
}
