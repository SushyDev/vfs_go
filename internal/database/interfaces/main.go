package interfaces

type Entity interface {
	GetId() int64
}

type Node interface{
	Entity

	GetName() string
	GetParentId() int64
	GetPath() string
	GetMode() int64
	GetUid() int
	GetGid() int
	GetModTime() string
	GetCreateTime() string
	GetAccessTime() string

	SetName(string)
	SetParentId(int64)
	SetPath(string)
	SetMode(int64)
	SetUid(int)
	SetGid(int)
	SetModTime(string)
	SetCreateTime(string)
	SetAccessTime(string)
}

type NodeRelationship interface {
	GetNodeId() int64
	SetNodeId(int64)
}

type NodeContent interface {
	Entity
	NodeRelationship

	GetContent() []byte
	SetContent([]byte)
}

type NodeAttribute interface {
	Entity
	NodeRelationship

	GetKey() string
	GetValue() string
}

type Symlink interface {
	Entity

	GetSourceNodeId() int64
	GetTargetNodeId() int64

	SetSourceNodeId(int64)
	SetTargetNodeId(int64)
}

type Database interface {
	GetNode(id int64) (Node, error)
	SaveNode(Node) error
	GetSymlink(id int64) (Symlink, error)
	SaveSymlink(Symlink) error
}

type RowScanner interface {
	Scan(dest ...any) error
}
