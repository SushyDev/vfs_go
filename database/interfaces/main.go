package interfaces

type Node interface{
	GetId() uint64
	GetName() string
	GetParentId() uint64
	GetPath() string
	GetContent() []byte
	GetSize() int64
	GetMode() uint32
	GetUid() int
	GetGid() int
	GetModTime() string
	GetCreateTime() string
	GetAccessTime() string

	SetName(string)
	SetParentId(uint64)
	SetContent([]byte)
	SetSize(int64)
	SetMode(uint32)
	SetUid(int)
	SetGid(int)
	SetModTime(string)
	SetCreateTime(string)
	SetAccessTime(string)
}

type Database interface {
	GetNode(id uint64) (Node, error)
	SaveNode(Node) error
}

type RowScanner interface {
	Scan(dest ...any) error
}
