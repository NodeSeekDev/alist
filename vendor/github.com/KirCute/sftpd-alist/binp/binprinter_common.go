package binp

import "encoding/binary"

// Printer type. Don't touch the internals.
type Printer struct {
	w []byte
}

// Create a new printer with empty output.
func Out() *Printer {
	return &Printer{[]byte{}}
}

// Create a new printer with output prefixed with the given byte slice.
func OutWith(b []byte) *Printer {
	return &Printer{b}
}

// Create a new printer with an empty slice with the capacity given below.
func OutCap(initialcap int) *Printer {
	return &Printer{make([]byte, 0, initialcap)}
}

// Output a byte.
func (p *Printer) Byte(d byte) *Printer {
	p.w = append(p.w, d)
	return p
}

// Output a byte, synonym for .Byte.
func (p *Printer) B8(d byte) *Printer {
	p.w = append(p.w, d)
	return p
}

var z16 = make([]byte, 16)

// Align to boundary
func (p *Printer) Align(n int) *Printer {
	r := len(p.w) % n
	if r == 0 {
		return p
	}
	r = n - r
	for r > 0 {
		cur := r
		if cur > 16 {
			cur = 16
		}
		p.w = append(p.w, z16[:cur]...)
		r -= cur
	}

	return p
}

// Skip (zero-fill) some bytes.
func (p *Printer) Skip(n int) *Printer {
	for n > 0 {
		cur := n
		if cur > 16 {
			cur = 16
		}
		p.w = append(p.w, z16[:cur]...)
		n -= cur
	}

	return p
}

// Output a raw byte slice with no length prefix.
func (p *Printer) Bytes(d []byte) *Printer {
	p.w = append(p.w, d...)
	return p
}

// Output a raw string with no length prefix.
func (p *Printer) String(d string) *Printer {
	p.w = append(p.w, []byte(d)...)
	return p
}

// Output a string terminated by a null-byte
func (p *Printer) String0(d string) *Printer {
	return p.String(d).Byte(0)
}

// Get the output as a byte slice.
func (p *Printer) Out() []byte {
	return p.w
}

// Start counting bytes for the length field in question.
func (p *Printer) LenStart(l *Len) *Printer {
	l.start = len(p.w)
	return p
}

// Call LenDone for all the arguments
func (p *Printer) LensDone(ls ...*Len) *Printer {
	for _, l := range ls {
		p.LenDone(l)
	}
	return p
}

// Fill fields associated with this length with the current offset.
func (p *Printer) LenDone(l *Len) *Printer {
	plen := len(p.w) - l.start
	for _, ls := range l.ls {
		switch ls.size {
		case 2:
			binary.BigEndian.PutUint16(p.w[ls.offset:], uint16(plen))
		case 4:
			binary.BigEndian.PutUint32(p.w[ls.offset:], uint32(plen))
		}
	}
	return p
}

// Type for handling length fields.
type Len struct {
	ls    []ls
	start int
}

type ls struct {
	offset uint32
	size   uint32
}
