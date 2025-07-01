package sftpd

import (
	"io"
	"strconv"
)

type FileOpenArgs struct {
	name  string
	flags uint32
	attr  *Attr
}

type DirReader struct {
	attrs []NamedAttr
	pos   int
}

func (d *DirReader) Readdir(count int) ([]NamedAttr, error) {
	if d.pos >= len(d.attrs) {
		return nil, io.EOF
	}
	var nextPos int
	if d.pos+count > len(d.attrs) {
		nextPos = len(d.attrs)
	} else {
		nextPos = d.pos + count
	}
	ret := d.attrs[d.pos:nextPos]
	d.pos = nextPos
	return ret, nil
}

func (d *DirReader) Close() error {
	return nil
}

type BufferedReader struct {
	r   io.ReadSeekCloser
	cur int64
}

func (b *BufferedReader) ReadAt(p []byte, off int64) (int, error) {
	if off != b.cur {
		o, err := b.r.Seek(off, io.SeekStart)
		if err != nil {
			return 0, err
		}
		b.cur = o
	}
	num := 0
	var err error
	for {
		var n int
		n, err = b.r.Read(p[num:])
		num += n
		if err != nil || num >= len(p) {
			break
		}
	}
	b.cur += int64(num)
	return num, err
}

func (b *BufferedReader) Close() error {
	return b.r.Close()
}

type WriteSeekCloser interface {
	io.WriteSeeker
	io.Closer
}

type AutoSeekWriter struct {
	w   WriteSeekCloser
	cur int64
}

func (a *AutoSeekWriter) WriteAt(p []byte, off int64) (int, error) {
	if off != a.cur {
		o, err := a.w.Seek(off, io.SeekStart)
		if err != nil {
			return 0, err
		}
		a.cur = o
	}
	n, err := a.w.Write(p)
	a.cur += int64(n)
	return n, err
}

func (a *AutoSeekWriter) Close() error {
	return a.w.Close()
}

type WriteAtCloser interface {
	io.WriterAt
	io.Closer
}

type ReadAtCloser interface {
	io.ReaderAt
	io.Closer
}

type handles struct {
	f  map[string]*FileOpenArgs
	d  map[string]string
	fw map[string]WriteAtCloser
	fr map[string]ReadAtCloser
	dr map[string]Dir
	c  int64
}

func (h *handles) init() {
	h.f = map[string]*FileOpenArgs{}
	h.d = map[string]string{}
	h.fw = map[string]WriteAtCloser{}
	h.fr = map[string]ReadAtCloser{}
	h.dr = map[string]Dir{}
}

func (h *handles) closeHandle(k string) {
	if k == "" {
		return
	}
	if k[0] == 'f' {
		delete(h.f, k)
		if c, ok := h.fw[k]; ok {
			_ = c.Close()
			delete(h.fw, k)
		}
		if c, ok := h.fr[k]; ok {
			_ = c.Close()
			delete(h.fr, k)
		}
	} else if k[0] == 'd' {
		delete(h.d, k)
		if c, ok := h.dr[k]; ok {
			_ = c.Close()
			delete(h.dr, k)
		}
	}
}

func (h *handles) nfiles() int { return len(h.f) }

func (h *handles) ndir() int { return len(h.d) }

func (h *handles) newFile(f *FileOpenArgs) string {
	h.c++
	k := "f" + strconv.FormatInt(h.c, 16)
	h.f[k] = f
	return k
}
func (h *handles) newDir(f string) string {
	h.c++
	k := "d" + strconv.FormatInt(h.c, 16)
	h.d[k] = f
	return k
}
func (h *handles) getFile(n string) *FileOpenArgs {
	return h.f[n]
}
func (h *handles) getDir(n string) string {
	return h.d[n]
}
