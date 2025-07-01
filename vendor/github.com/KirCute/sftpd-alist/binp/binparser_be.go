package binp

import "encoding/binary"

// Parse 4 bigendian bytes from the buffer.
func (p *Parser) B32(d *uint32) *Parser {
	if p == nil || len(p.r) < p.off+4 {
		return nil
	}
	*d = binary.BigEndian.Uint32(p.r[p.off:])
	p.off += 4
	return p
}

// Parse 8 bigendian bytes from the buffer.
func (p *Parser) B64(d *uint64) *Parser {
	if p == nil || len(p.r) < p.off+8 {
		return nil
	}
	*d = binary.BigEndian.Uint64(p.r[p.off:])
	p.off += 8
	return p
}

// Parse 2 bigendian bytes from the buffer.
func (p *Parser) B16(d *uint16) *Parser {
	if p == nil || len(p.r) < p.off+2 {
		return nil
	}
	*d = binary.BigEndian.Uint16(p.r[p.off:])
	p.off += 2
	return p
}

// Parse a string with a 4 byte bigendian length prefix to a []byte pointer.
func (p *Parser) B32Bytes(d *[]byte) *Parser {
	var v uint32
	return p.B32(&v).Bytes(int(v), d)
}

// Parse a string with a 4 byte bigendian length prefix to a string pointer.
func (p *Parser) B32String(d *string) *Parser {
	var v uint32
	return p.B32(&v).String(int(v), d)
}

// Parse a string with a 2 byte bigendian length prefix to a string pointer.
func (p *Parser) B16String(d *string) *Parser {
	var v uint16
	return p.B16(&v).String(int(v), d)
}

// Parse a string with a 1 byte length prefix to a string pointer.
func (p *Parser) B8String(d *string) *Parser {
	var v uint8
	return p.Byte(&v).String(int(v), d)
}
