package roaring

import (
	"fmt"
	"sort"
)

// common to rle32.go and rle16.go

// rleVerbose controls whether p() prints show up.
// The testing package sets this based on
// testing.Verbose().
var rleVerbose bool

// p is a shorthand for fmt.Printf with beginning and
// trailing newlines. p() makes it easy
// to add diagnostic print statements.
func p(format string, args ...interface{}) {
	if rleVerbose {
		fmt.Printf("\n"+format+"\n", args...)
	}
}

// MaxUint32 is only used internally for the endx
// value when UpperLimit32 is stored; users should
// only ever store up to UpperLimit32.
const MaxUint32 = 4294967295

// UpperLimit32 is the largest
// integer we can store in an RunContainer32. As
// we need to reserve one value for the open
// interval endpoint endx, this is MaxUint32 - 1.
const UpperLimit32 = MaxUint32 - 1

// MaxUint16 is only used internally for the endx
// value when UpperLimit16 is stored; users should
// only ever store up to UpperLimit16.
const MaxUint16 = 65535

// UpperLimit16 is the largest
// integer we can store in an RunContainer16. As
// we need to reserve one value for the open
// interval endpoint endx, this is MaxUint16 - 1.
const UpperLimit16 = MaxUint16 - 1

// searchOptions allows us to accelerate runContainer32.search with
// prior knowledge of (mostly lower) bounds. This is used by Union
// and Intersect.
type searchOptions struct {

	// start here instead of at 0
	StartIndex int

	// upper bound instead of len(rc.iv);
	// EndxIndex == 0 means ignore the bound and use
	// EndxIndex == n ==len(rc.iv) which is also
	// naturally the default for search()
	// when opt = nil.
	EndxIndex int
}

// And finds the intersection of rc and b.
func (rc *runContainer32) And(b *Bitmap) *Bitmap {
	out := NewBitmap()
	for _, p := range rc.iv {
		for i := p.start; i < p.endx; i++ {
			if b.Contains(i) {
				out.Add(i)
			}
		}
	}
	return out
}

// Xor returns the exclusive-or of rc and b.
func (rc *runContainer32) Xor(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			if out.Contains(v) {
				out.RemoveRange(uint64(v), uint64(v+1))
			} else {
				out.Add(v)
			}
		}
	}
	return out
}

// Or returns the union of rc and b.
func (rc *runContainer32) Or(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			out.Add(v)
		}
	}
	return out
}

func showHash(name string, h map[int]bool) {
	hv := []int{}
	for k := range h {
		hv = append(hv, k)
	}
	sort.Sort(sort.IntSlice(hv))
	stringH := ""
	for i := range hv {
		stringH += fmt.Sprintf("%v, ", hv[i])
	}

	p("%s is (len %v): %s", name, len(h), stringH)
}

// trial is used in the randomized testing of runContainers
type trial struct {
	n           int
	percentFill float64
	ntrial      int

	// only in the union test
	percentDelete float64
}

// And finds the intersection of rc and b.
func (rc *runContainer16) And(b *Bitmap) *Bitmap {
	out := NewBitmap()
	for _, p := range rc.iv {
		for i := p.start; i < p.endx; i++ {
			if b.Contains(uint32(i)) {
				out.Add(uint32(i))
			}
		}
	}
	return out
}

// Xor returns the exclusive-or of rc and b.
func (rc *runContainer16) Xor(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			w := uint32(v)
			if out.Contains(w) {
				out.RemoveRange(uint64(w), uint64(w+1))
			} else {
				out.Add(w)
			}
		}
	}
	return out
}

// Or returns the union of rc and b.
func (rc *runContainer16) Or(b *Bitmap) *Bitmap {
	out := b.Clone()
	for _, p := range rc.iv {
		for v := p.start; v < p.endx; v++ {
			out.Add(uint32(v))
		}
	}
	return out
}
