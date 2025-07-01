package sftpd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/KirCute/sftpd-alist/binp"
	"io"
	"os"
	"time"

	"github.com/taruti/bytepool"
	"golang.org/x/crypto/ssh"
)

var sftpSubSystem = []byte{0, 0, 0, 4, 115, 102, 116, 112}

// IsSftpRequest checks whether a given ssh.Request is for sftp.
func IsSftpRequest(req *ssh.Request) bool {
	return req.Type == "subsystem" && bytes.Equal(sftpSubSystem, req.Payload)
}

var initReply = []byte{0, 0, 0, 5, ssh_FXP_VERSION, 0, 0, 0, 3}

type DebugLogger func(s string, v ...interface{})

// ServeChannel serves a ssh.Channel with the given FileSystem.
func ServeChannel(c ssh.Channel, fs FileSystem, debugf DebugLogger) error {
	defer func() { _ = c.Close() }()
	var h handles
	h.init()
	brd := bufio.NewReaderSize(c, 64*1024)
	var e error
	var plen int
	var op byte
	var bs []byte
	var id uint32
	for {
		if e != nil {
			debugf("Sending error: %v\n", e)
			e = writeErr(c, id, e, debugf)
			if e != nil {
				return e
			}
		}
		_ = discard(brd, plen)
		plen, op, e = readPacketHeader(brd)
		if e != nil {
			return e
		}
		plen--
		debugf("CR op=%v data len=%d\n", ssh_fxp(op), plen)
		if plen < 2 {
			return errors.New("Packet too short")
		}
		// Feeding too large values to peek is ok, it just errors.
		bs, e = brd.Peek(plen)
		if e != nil {
			return e
		}
		debugf("Data %X\n", bs)
		p := binp.NewParser(bs)
		switch op {
		case ssh_FXP_INIT:
			debugf("Init %v\n", initReply)
			e = wrc(c, initReply)
		case ssh_FXP_OPEN:
			var path string
			var flags uint32
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&path).B32(&flags), &a).End()
			if e != nil {
				return e
			}
			debugf("Open id=%v path=%s flags=%v\n", id, path, flags)
			if h.nfiles() >= maxFiles {
				e = errTooManyFiles
				continue
			}
			handle := h.newFile(&FileOpenArgs{path, flags, &a})
			debugf("Open ret: handle=%s\n", handle)
			e = writeHandle(c, id, handle)
		case ssh_FXP_CLOSE:
			var handle string
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			debugf("Close id=%v handle=%s\n", id, handle)
			h.closeHandle(handle)
			e = writeErr(c, id, nil, debugf)
		case ssh_FXP_READ:
			var handle string
			var offset uint64
			var length uint32
			var n int
			e = p.B32(&id).B32String(&handle).B64(&offset).B32(&length).End()
			if e != nil {
				return e
			}
			debugf("Read id=%v handle=%s offset=%v length=%v\n", id, handle, offset, length)
			f := h.getFile(handle)
			if f == nil {
				return errInvalidHandle
			}
			if length > 64*1024 {
				length = 64 * 1024
			}
			bs := bytepool.Alloc(int(length))
			if reader, ok := h.fr[handle]; ok {
				n, e = reader.ReadAt(bs, int64(offset))
			} else {
				if ft, ok := fs.(FileSystemExtentionFileTransfer); ok {
					var t FileTransfer
					t, e = ft.GetHandle(f.name, f.flags, f.attr, offset)
					if e == nil {
						br := &BufferedReader{r: t, cur: int64(offset)}
						n, e = br.ReadAt(bs, int64(offset))
						h.fr[handle] = br
					}
				} else {
					var file File
					file, e = fs.OpenFile(f.name, f.flags, f.attr)
					if e == nil {
						if offset != 0 {
							_, e = file.Seek(int64(offset), io.SeekStart)
						}
						if e == nil {
							br := &BufferedReader{r: file, cur: int64(offset)}
							n, e = br.ReadAt(bs, int64(offset))
							h.fr[handle] = br
						}
					}
				}
			}
			// Handle go readers that return io.EOF and bytes at the same time.
			if e == io.EOF && n > 0 {
				debugf("Read EOF, n=%d\n", n)
				e = nil
			}
			if e != nil {
				bytepool.Free(bs)
				continue
			}
			bs = bs[0:n]
			e = wrc(c, binp.Out().B32(1+4+4+uint32(len(bs))).Byte(ssh_FXP_DATA).B32(id).B32(uint32(len(bs))).Out())
			if e == nil {
				e = wrc(c, bs)
			}
			bytepool.Free(bs)
		case ssh_FXP_WRITE:
			var handle string
			var offset uint64
			var length uint32
			p.B32(&id).B32String(&handle).B64(&offset).B32(&length)
			debugf("Write id=%v handle=%s offset=%v length=%v\n", id, handle, offset, length)
			f := h.getFile(handle)
			if f == nil {
				return errInvalidHandle
			}
			var bs []byte
			e = p.BytesPeek(int(length), &bs).End()
			if e != nil {
				return e
			}
			if writer, ok := h.fw[handle]; ok {
				_, e = writer.WriteAt(bs, int64(offset))
			} else {
				if ft, ok := fs.(FileSystemExtentionFileTransfer); ok {
					var t FileTransfer
					t, e = ft.GetHandle(f.name, f.flags, f.attr, offset)
					if e == nil {
						aw := &AutoSeekWriter{w: t, cur: int64(offset)}
						_, e = aw.WriteAt(bs, int64(offset))
						h.fw[handle] = aw
					}
				} else {
					var file File
					file, e = fs.OpenFile(f.name, f.flags, f.attr)
					if e == nil {
						if offset != 0 {
							_, e = file.Seek(int64(offset), io.SeekStart)
						}
						if e == nil {
							aw := &AutoSeekWriter{w: file, cur: int64(offset)}
							_, e = aw.WriteAt(bs, int64(offset))
							h.fw[handle] = aw
						}
					}
				}
			}
			e = writeErr(c, id, e, debugf)
		case ssh_FXP_LSTAT, ssh_FXP_STAT:
			var path string
			var a *Attr
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("Stat/Lstat id=%d path=%s\n", id, path)
			a, e = fs.Stat(path, op == ssh_FXP_LSTAT)
			debugf("Stat/Lstat ret: %v %v\n", a, e)
			e = writeAttr(c, id, a, e, debugf)
		case ssh_FXP_FSTAT:
			var handle string
			var a *Attr
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			debugf("Fstat id=%d handle=%s\n", id, handle)
			f := h.getFile(handle)
			if f == nil {
				return errInvalidHandle
			}
			a, e = fs.Stat(f.name, false)
			debugf("Fstat ret: %v %v\n", a, e)
			e = writeAttr(c, id, a, e, debugf)
		case ssh_FXP_SETSTAT:
			var path string
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&path), &a).End()
			if e != nil {
				return e
			}
			debugf("SetStat id=%d path=%s\n", id, path)
			e = writeErr(c, id, fs.SetStat(path, &a), debugf)
		case ssh_FXP_FSETSTAT:
			var handle string
			var a Attr
			e = parseAttr(p.B32(&id).B32String(&handle), &a).End()
			if e != nil {
				return e
			}
			debugf("FSetStat id=%d handle=%s\n", id, handle)
			f := h.getFile(handle)
			if f == nil {
				return errInvalidHandle
			}
			e = writeErr(c, id, fs.SetStat(f.name, &a), debugf)
		case ssh_FXP_OPENDIR:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("Opendir id=%d path=%s\n", id, path)
			if e != nil {
				continue
			}
			handle := h.newDir(path)
			debugf("Opendir ret: handle=%s\n", handle)
			e = writeHandle(c, id, handle)
		case ssh_FXP_READDIR:
			var handle string
			e = p.B32(&id).B32String(&handle).End()
			if e != nil {
				return e
			}
			debugf("Readdir id=%d handle=%s\n", id, handle)
			f := h.getDir(handle)
			if f == "" {
				return errInvalidHandle
			}
			var fis []NamedAttr
			if dr, ok := h.dr[handle]; ok {
				fis, e = dr.Readdir(1024)
			} else {
				if frd, ok := fs.(FileSystemExtensionFileList); ok {
					var allFis []NamedAttr
					allFis, e = frd.ReadDir(f)
					if e == nil {
						dirReader := &DirReader{
							attrs: allFis,
							pos:   0,
						}
						h.dr[handle] = dirReader
						fis, e = dirReader.Readdir(1024)
					}
				} else {
					var dir Dir
					dir, e = fs.OpenDir(f)
					if e == nil {
						h.dr[handle] = dir
						fis, e = dir.Readdir(1024)
					}
				}
			}
			debugf("Readdir ret: %v => %v\n", fis, e)
			if e != nil {
				continue
			}
			var l binp.Len
			o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_NAME).B32(id).B32(uint32(len(fis)))
			for _, fi := range fis {
				n := fi.Name
				o.B32String(n).B32String(readdirLongName(&fi)).B32(fi.Flags)
				if fi.Flags&ATTR_SIZE != 0 {
					o.B64(uint64(fi.Size))
				}
				if fi.Flags&ATTR_UIDGID != 0 {
					o.B32(fi.Uid).B32(fi.Gid)
				}
				if fi.Flags&ATTR_MODE != 0 {
					o.B32(fileModeToSftp(fi.Mode))
				}
				if fi.Flags&ATTR_TIME != 0 {
					outTimes(o, &fi.Attr)
				}
			}
			o.LenDone(&l)
			e = wrc(c, o.Out())
		case ssh_FXP_REMOVE:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("Remove id=%d path=%s\n", id, path)
			e = writeErr(c, id, fs.Remove(path), debugf)
		case ssh_FXP_MKDIR:
			var path string
			var a Attr
			p = p.B32(&id).B32String(&path)
			e = parseAttr(p, &a).End()
			if e != nil {
				return e
			}
			debugf("Mkdir id=%d path=%s\n", id, path)
			e = writeErr(c, id, fs.Mkdir(path, &a), debugf)
		case ssh_FXP_RMDIR:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("Rmdir id=%d path=%s\n", id, path)
			e = writeErr(c, id, fs.Rmdir(path), debugf)
		case ssh_FXP_REALPATH:
			var path, newpath string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("RealPath id=%d path=%s\n", id, path)
			newpath, e = fs.RealPath(path)
			debugf("RealPath ret %s => %v\n", newpath, e)
			e = writeNameOnly(c, id, newpath, e, debugf)
		case ssh_FXP_RENAME:
			var oldName, newName string
			var flags uint32
			e = p.B32(&id).B32String(&oldName).B32String(&newName).B32(&flags).End()
			if e != nil {
				return e
			}
			debugf("Rename id=%d oldName=%s newName=%s flags=%x\n", id, oldName, newName, flags)
			e = writeErr(c, id, fs.Rename(oldName, newName, flags), debugf)
		case ssh_FXP_READLINK:
			var path string
			e = p.B32(&id).B32String(&path).End()
			if e != nil {
				return e
			}
			debugf("ReadLink id=%d path=%s\n", id, path)
			path, e = fs.ReadLink(path)
			debugf("ReadLink ret %s\n", path)
			e = writeNameOnly(c, id, path, e, debugf)
		case ssh_FXP_SYMLINK:
			e = writeErrCode(c, id, ssh_FX_OP_UNSUPPORTED, debugf)
		}
		if e != nil {
			debugf("Fatal error: %v\n", e)
			return e
		}
	}
}

var errInvalidHandle = errors.New("Client supplied an invalid handle")
var errTooManyFiles = errors.New("Too many files")

const maxFiles = 0x100

func readPacketHeader(rd *bufio.Reader) (int, byte, error) {
	bs := make([]byte, 5)
	_, e := io.ReadFull(rd, bs)
	if e != nil {
		return 0, 0, e
	}
	return int(binary.BigEndian.Uint32(bs)), bs[4], nil
}

func parseAttr(p *binp.Parser, a *Attr) *binp.Parser {
	p = p.B32(&a.Flags)
	if a.Flags&ssh_FILEXFER_ATTR_SIZE != 0 {
		p = p.B64(&a.Size)
	}
	if a.Flags&ssh_FILEXFER_ATTR_UIDGID != 0 {
		p = p.B32(&a.Uid).B32(&a.Gid)
	}
	if a.Flags&ssh_FILEXFER_ATTR_PERMISSIONS != 0 {
		var mode uint32
		p = p.B32(&mode)
		a.Mode = sftpToFileMode(mode)
	}
	if a.Flags&ssh_FILEXFER_ATTR_ACMODTIME != 0 {
		p = inTimes(p, a)
	}
	if a.Flags&ssh_FILEXFER_ATTR_EXTENDED != 0 {
		var count uint32
		p = p.B32(&count)
		if count > 0xFF {
			return nil
		}
		ss := make([]string, 2*int(count))
		for i := 0; i < int(count); i++ {
			var k, v string
			p = p.B32String(&k).B32String(&v)
			ss[2*i+0] = k
			ss[2*i+1] = v
		}
		a.Extended = ss
	}
	return p
}

func writeAttr(c ssh.Channel, id uint32, a *Attr, e error, debugf DebugLogger) error {
	if e != nil {
		return writeErr(c, id, e, debugf)
	}
	var l binp.Len
	o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_ATTRS).B32(id).B32(a.Flags)
	if a.Flags&ssh_FILEXFER_ATTR_SIZE != 0 {
		o = o.B64(a.Size)
	}
	if a.Flags&ssh_FILEXFER_ATTR_UIDGID != 0 {
		o = o.B32(a.Uid).B32(a.Gid)
	}
	if a.Flags&ssh_FILEXFER_ATTR_PERMISSIONS != 0 {
		o = o.B32(fileModeToSftp(a.Mode))
	}
	if a.Flags&ssh_FILEXFER_ATTR_ACMODTIME != 0 {
		outTimes(o, a)
	}
	if a.Flags&ssh_FILEXFER_ATTR_EXTENDED != 0 {
		count := uint32(len(a.Extended) / 2)
		o = o.B32(count)
		for _, s := range a.Extended {
			o = o.B32String(s)
		}
	}
	o.LenDone(&l)
	return wrc(c, o.Out())
}

func writeNameOnly(c ssh.Channel, id uint32, path string, e error, debugf DebugLogger) error {
	if e != nil {
		return writeErr(c, id, e, debugf)
	}
	var l binp.Len
	o := binp.Out().LenB32(&l).LenStart(&l).Byte(ssh_FXP_NAME).B32(id).B32(1)
	o.B32String(path).B32String(path).B32(0)
	o.LenDone(&l)
	return wrc(c, o.Out())
}

var failTmpl = []byte{0, 0, 0, 1 + 4 + 4 + 4 + 4, ssh_FXP_STATUS, 0, 0, 0, 0, 0, 0, 0, ssh_FX_FAILURE, 0, 0, 0, 0, 0, 0, 0, 0}

func writeFail(c ssh.Channel, id uint32) error {
	bs := make([]byte, len(failTmpl))
	copy(bs, failTmpl)
	binary.BigEndian.PutUint32(bs[5:], id)
	return wrc(c, bs)
}

func writeErrCode(c ssh.Channel, id uint32, code ssh_fx, debugf DebugLogger) error {
	bs := make([]byte, len(failTmpl))
	copy(bs, failTmpl)
	binary.BigEndian.PutUint32(bs[5:], id)
	debugf("Sending sftp error code %v\n", code)
	bs[12] = byte(code)
	return wrc(c, bs)
}

func writeErr(c ssh.Channel, id uint32, err error, debugf DebugLogger) error {
	var code ssh_fx
	switch {
	case err == nil:
		code = ssh_FX_OK
	case err == io.EOF:
		code = ssh_FX_EOF
	case os.IsPermission(err):
		code = ssh_FX_PERMISSION_DENIED
	case os.IsNotExist(err):
		code = ssh_FX_NO_SUCH_FILE
	default:
		code = ssh_FX_FAILURE
	}
	return writeErrCode(c, id, code, debugf)
}

func writeHandle(c ssh.Channel, id uint32, handle string) error {
	return wrc(c, binp.OutCap(4+9+len(handle)).B32(uint32(9+len(handle))).B8(ssh_FXP_HANDLE).B32(id).B32String(handle).Out())
}

func wrc(c ssh.Channel, bs []byte) error {
	_, e := c.Write(bs)
	return e
}

func discard(brd *bufio.Reader, n int) error {
	if n == 0 {
		return nil
	}
	m, e := io.Copy(io.Discard, &io.LimitedReader{R: brd, N: int64(n)})
	if int(m) == n && e == io.EOF {
		e = nil
	}
	return e
}

func outTimes(o *binp.Printer, a *Attr) {
	o.B32(uint32(a.ATime.Unix())).B32(uint32(a.MTime.Unix()))
}
func inTimes(p *binp.Parser, a *Attr) *binp.Parser {
	var at, mt uint32
	p = p.B32(&at).B32(&mt)
	a.ATime = time.Unix(int64(at), 0)
	a.MTime = time.Unix(int64(mt), 0)
	return p
}
