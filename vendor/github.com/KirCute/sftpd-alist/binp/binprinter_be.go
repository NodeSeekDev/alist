package binp

import "encoding/binary"

// Output 2 bigendian bytes.
func (p *Printer) B16(d uint16) *Printer {
	var bytes [2]byte
	binary.BigEndian.PutUint16(bytes[:], d)
	p.w = append(p.w, bytes[0], bytes[1])
	return p
}

// Output 4 bigendian bytes.
func (p *Printer) B32(d uint32) *Printer {
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[:], d)
	p.w = append(p.w, bytes[0], bytes[1], bytes[2], bytes[3])
	return p
}

// Output 4 bigendian bytes.
func (p *Printer) B64(d uint64) *Printer {
	var bytes [8]byte
	binary.BigEndian.PutUint64(bytes[:], d)
	p.w = append(p.w, bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7])
	return p
}

// Output a string with a 4 byte bigendian length prefix and no trailing null.
func (p *Printer) B32String(d string) *Printer {
	return p.B32(uint32(len(d))).String(d)
}

// Output bytes with a 4 byte bigendian length prefix and no trailing null.
func (p *Printer) B32Bytes(d []byte) *Printer {
	return p.B32(uint32(len(d))).Bytes(d)
}

// Output a string with a 2 byte bigendian length prefix and no trailing null.
func (p *Printer) B16String(d string) *Printer {
	if len(d) > 0xffff {
		panic("binprinter: string too long")
	}
	return p.B16(uint16(len(d))).String(d)
}

// Output a string with a 1 byte bigendian length prefix and no trailing null.
func (p *Printer) B8String(d string) *Printer {
	if len(d) > 0xff {
		panic("binprinter: string too long")
	}
	return p.Byte(byte(len(d))).String(d)
}

// Add a 16 bit field at the current location that will be filled with the length.
func (p *Printer) LenB16(l *Len) *Printer {
	l.ls = append(l.ls, ls{uint32(len(p.w)), 2})
	return p.B16(0)
}

// Add a 32 bit field at the current location that will be filled with the length.
func (p *Printer) LenB32(l *Len) *Printer {
	l.ls = append(l.ls, ls{uint32(len(p.w)), 4})
	return p.B32(0)
}
