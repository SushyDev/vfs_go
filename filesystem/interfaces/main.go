package interfaces

import (
	"io/fs"

	database_interfaces "github.com/sushydev/vfs_go/database/interfaces"
)

type FileSystem interface {
	Root() (Node, error)
	Open(identifier uint64) (Node, error)
	ReadDir(identifier uint64) ([]Node, error)
	Lookup(identifier uint64, name string) (Node, error)
	MkDir(parentId uint64, name string) error
	GetNodeByParentAndName(parentId uint64, name string) (Node, error)
}

type Node interface {
	GetId() uint64
	GetName() string
	GetParentId() uint64
	GetPath() string
	GetContent() []byte
	GetSize() int64
	GetMode() fs.FileMode
	GetUid() int
	GetGid() int
	GetModTime() string
	GetCreateTime() string
	GetAccessTime() string

	SetName(name string)
	SetParentId(parentId uint64)
	SetPath(path string)
	SetContent(content []byte)
	SetSize(size int64)
	SetMode(mode uint32)
	SetUid(uid int)
	SetGid(gid int)
	SetModTime(modTime string)
	SetCreateTime(createTime string)
	SetAccessTime(accessTime string)

	GetEntity() database_interfaces.Node
}
