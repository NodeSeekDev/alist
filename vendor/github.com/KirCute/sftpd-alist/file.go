package sftpd

import (
	"io"
	"os"
	"time"
)

type Attr struct {
	Flags        uint32
	Size         uint64
	Uid, Gid     uint32
	User, Group  string
	Mode         os.FileMode
	ATime, MTime time.Time
	Extended     []string
}

type NamedAttr struct {
	Name string
	Attr
}

const (
	ATTR_SIZE    = ssh_FILEXFER_ATTR_SIZE
	ATTR_UIDGID  = ssh_FILEXFER_ATTR_UIDGID
	ATTR_MODE    = ssh_FILEXFER_ATTR_PERMISSIONS
	ATTR_TIME    = ssh_FILEXFER_ATTR_ACMODTIME
	MODE_REGULAR = os.FileMode(0)
	MODE_DIR     = os.ModeDir
)

type Dir interface {
	io.Closer
	Readdir(count int) ([]NamedAttr, error)
}

type File interface {
	FileTransfer
	FStat() (*Attr, error)
	FSetStat(*Attr) error
}

type FileSystem interface {
	OpenFile(name string, flags uint32, attr *Attr) (File, error)
	OpenDir(name string) (Dir, error)
	Remove(name string) error
	Rename(old, new string, flags uint32) error
	Mkdir(name string, attr *Attr) error
	Rmdir(name string) error
	Stat(name string, islstat bool) (*Attr, error)
	SetStat(name string, attr *Attr) error
	ReadLink(path string) (string, error)
	CreateLink(path string, target string, flags uint32) error
	RealPath(path string) (string, error)
}

// FillFrom fills an Attr from a os.FileInfo
func (a *Attr) FillFrom(fi os.FileInfo) {
	*a = Attr{}
	a.Flags = ATTR_SIZE | ATTR_MODE
	a.Size = uint64(fi.Size())
	a.Mode = fi.Mode()
	a.MTime = fi.ModTime()
}

func fileModeToSftp(m os.FileMode) uint32 {
	var raw = uint32(m.Perm())
	switch {
	case m.IsDir():
		raw |= 0040000
	case m.IsRegular():
		raw |= 0100000
	}
	return raw
}

func sftpToFileMode(raw uint32) os.FileMode {
	var m = os.FileMode(raw & 0777)
	switch {
	case raw&0040000 != 0:
		m |= os.ModeDir
	case raw&0100000 != 0:
		// regular
	}
	return m
}

// FileSystemExtensionFileList is a convenience extension to allow to return file listing
// without requiring to implement the methods Open/Readdir for your custom afero.File
// From: github.com/fclairamb/ftpserverlib
type FileSystemExtensionFileList interface {
	// ReadDir reads the directory named by name and return a list of directory entries.
	ReadDir(name string) ([]NamedAttr, error)
}

// FileSystemExtentionFileTransfer is a convenience extension to allow to transfer files
// without requiring to implement the methods Create/Open/OpenFile for your custom afero.File.
// From: github.com/fclairamb/ftpserverlib
type FileSystemExtentionFileTransfer interface {
	// GetHandle return an handle to upload or download a file based on flags:
	// os.O_RDONLY indicates a download
	// os.O_WRONLY indicates an upload and can be combined with os.O_APPEND (resume) or
	// os.O_CREATE (upload to new file/truncate)
	//
	// offset is the argument of a previous REST command, if any, or 0
	GetHandle(name string, flags uint32, attr *Attr, offset uint64) (FileTransfer, error)
}

// FileTransfer defines the inferface for file transfers.
// From: github.com/fclairamb/ftpserverlib
type FileTransfer interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}
