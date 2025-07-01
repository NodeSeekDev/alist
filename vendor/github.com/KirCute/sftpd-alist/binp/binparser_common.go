package binp

import "errors"

// Parser type. Don't touch the internals.
type Parser struct {
	r   []byte
	off int
}

// Create a new parser with the given buffer. Panics on error.
func NewParser(b []byte) *Parser {
	return &Parser{b, 0}
}

// Parse a byte from the buffer.
func (p *Parser) Byte(d *byte) *Parser {
	if p == nil || len(p.r) < p.off+1 {
		return nil
	}
	*d = p.r[p.off]
	p.off++
	return p
}

// Parse a byte from the buffer, synonym for .Byte.
func (p *Parser) B8(d *byte) *Parser {
	if p == nil || len(p.r) < p.off+1 {
		return nil
	}
	*d = p.r[p.off]
	p.off++
	return p
}

// Parse n bytes from the buffer and copy to a []byte pointer that is allocated.
func (p *Parser) Bytes(n int, d *[]byte) *Parser {
	if p == nil || n > len(p.r[p.off:]) {
		return nil
	}
	*d = make([]byte, n)
	copy(*d, p.r[p.off:])
	p.off += n
	return p
}

// Parse n bytes from the buffer and copy to the supplied []byte.
func (p *Parser) BytesCopy(n int, d []byte) *Parser {
	if p == nil || n > len(p.r[p.off:]) {
		return nil
	}
	copy(d, p.r[p.off:p.off+n])
	p.off += n
	return p
}

// Parse n bytes from the buffer to a []byte pointer that refers to the parser internal buffer.
func (p *Parser) BytesPeek(n int, d *[]byte) *Parser {
	if p == nil || n > len(p.r[p.off:]) {
		return nil
	}
	*d = p.r[p.off : p.off+n]
	p.off += n
	return p
}

// Ensure the input is aligned possibly skipping bytes.
func (p *Parser) Align(n int) *Parser {
	if p == nil {
		return nil
	}
	r := p.off % n
	if r != 0 {
		p.off += n - r
	}
	return p
}

// Skip bytes.
func (p *Parser) Skip(n int) *Parser {
	p.off += n
	return p
}

// Parse a null terminated string.
func (p *Parser) String0(d *string) *Parser {
	if p == nil {
		return nil
	}
	for i, ch := range p.r[p.off:] {
		if ch == 0 {
			p.String(i, d)
			p.off++
			return p
		}
	}
	panic("String0: null byte not found")
}

// Check that we are at the end of input.
func (p *Parser) End() error {
	if p == nil || p.off != len(p.r) {
		return eparse
	}
	return nil
}

var eparse = errors.New("binparser invalid input")

// Check that we are at the end of input.
func (p *Parser) AtEnd() bool {
	return p.off == len(p.r)
}

// Peek the rest of input as raw bytes.
func (p *Parser) PeekRest(d *[]byte) *Parser {
	*d = p.r[p.off:]
	return p
}

// Parse n bytes from the buffer to a string pointer.
func (p *Parser) String(n int, d *string) *Parser {
	if p == nil || n > len(p.r[p.off:]) {
		return nil
	}
	bs := p.r[p.off : p.off+n]
	*d = string(bs)
	p.off += n
	return p
}
