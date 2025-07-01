// Variable sized pools of bytes using sync.Pools underneath.
package bytepool

import (
	"sync"
)

// Alloc a slice of bytes with the given size. Sizes till MaxSize are supported, for larger
// sizes nil will be returned. The slices are taken from sync.Pool caches sized to the
// nearest power of two and truncated to the size supplied. This function may be
// called concurrently.
func Alloc(size int) []byte {
	if size <= MaxSize {
		sz := int32(size)
		for i, v := range sizes {
			if sz <= v {
				bs := pools[i].Get().([]byte)
				return bs[0:size]
			}
		}
	}
	return nil
}

// Free a byte slice allocated by Alloc. The argument may be shorted from the original allocation,
// it is reset to len = cap. This function may be called concurrently.
func Free(bs []byte) {
	bs = bs[0:cap(bs)]
	l := int32(len(bs))
	for i, v := range sizes {
		if l == v {
			pools[i].Put(bs)
			return
		}
	}
}

const npools = 20

// MaxSize is the supported maximum size
const MaxSize = 1 << 23

var pools = [npools]sync.Pool{}
var sizes = [npools]int32{
	1 << 4,
	1 << 5,
	1 << 6,
	1 << 7,
	1 << 8,
	1 << 9,
	1 << 10,
	1 << 11,
	1 << 12,
	1 << 13,
	1 << 14,
	1 << 15,
	1 << 16,
	1 << 17,
	1 << 18,
	1 << 19,
	1 << 20,
	1 << 21,
	1 << 22,
	1 << 23,
}

func init() {
	for i := range pools {
		pools[i].New = creater(1 << (uint32(i) + 4)).create
	}
}

type creater int32

func (c creater) create() interface{} {
	return make([]byte, int(c))
}
