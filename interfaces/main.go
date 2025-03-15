package interfaces

import (
	"io/fs"

	database_interfaces "github.com/sushydev/vfs_go/internal/database/interfaces"
)

type FileSystem interface {
	Open(id uint64) (Node, error)
	ReadDir(parentId uint64) ([]Node, error)
	Lookup(parentId uint64, name string) (Node, error)
	MkDir(parentId uint64, name string) error
}

type Entry interface {
	GetId() uint64
}

type Node interface {
	Entry

	GetName() string
	GetParentId() uint64
	GetPath() string
	GetMode() fs.FileMode
	GetUid() int
	GetGid() int
	GetModTime() string
	GetCreateTime() string
	GetAccessTime() string

	SetName(name string)
	SetParentId(parentId uint64)
	SetPath(path string)
	SetMode(mode uint32)
	SetUid(uid int)
	SetGid(gid int)
	SetModTime(modTime string)
	SetCreateTime(createTime string)
	SetAccessTime(accessTime string)

	GetEntity() database_interfaces.Node
}

type NodeContent interface {
	Entry

	GetNodeId() uint64
	GetContent() []byte

	SetNodeId(nodeId uint64)
	SetContent(content []byte)

	GetEntity() database_interfaces.NodeContent
}

type Symlink interface {
	Entry

	GetSourceNodeId() uint64
	GetTargetNodeId() uint64

	SetSourceNodeId(nodeId uint64)
	SetTargetNodeId(nodeId uint64)

	GetEntity() database_interfaces.Symlink
}
