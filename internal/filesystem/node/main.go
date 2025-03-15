package node

import (
	"io/fs"

	"github.com/sushydev/vfs_go/interfaces"
	database_interfaces "github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Node struct {
	entity database_interfaces.Node
}

var _ interfaces.Node = &Node{}

func New(entity database_interfaces.Node) (*Node, error) {
	return &Node{
		entity: entity,
	}, nil
}

func (node *Node) GetId() uint64 {
	return uint64(node.entity.GetId())
}

func (node *Node) GetName() string {
	return node.entity.GetName()
}

func (node *Node) GetParentId() uint64 {
	return uint64(node.entity.GetParentId())
}

func (node *Node) GetPath() string {
	return node.entity.GetPath()
}

func (node *Node) GetMode() fs.FileMode {
	return fs.FileMode(uint32(node.entity.GetMode()))
}

func (node *Node) GetUid() int {
	return node.entity.GetUid()
}

func (node *Node) GetGid() int {
	return node.entity.GetGid()
}

func (node *Node) GetModTime() string {
	return node.entity.GetModTime()
}

func (node *Node) GetCreateTime() string {
	return node.entity.GetCreateTime()
}

func (node *Node) GetAccessTime() string {
	return node.entity.GetAccessTime()
}

func (node *Node) SetName(name string) {
	node.entity.SetName(name)
}

func (node *Node) SetParentId(parentId uint64) {
	node.entity.SetParentId(int64(parentId))
}

func (node *Node) SetPath(path string) {
	node.entity.SetPath(path)
}

func (node *Node) SetMode(mode uint32) {
	node.entity.SetMode(int64(mode))
}

func (node *Node) SetUid(uid int) {
	node.entity.SetUid(uid)
}

func (node *Node) SetGid(gid int) {
	node.entity.SetGid(gid)
}

func (node *Node) SetModTime(modTime string) {
	node.entity.SetModTime(modTime)
}

func (node *Node) SetCreateTime(createTime string) {
	node.entity.SetCreateTime(createTime)
}

func (node *Node) SetAccessTime(accessTime string) {
	node.entity.SetAccessTime(accessTime)
}

func (node *Node) GetEntity() database_interfaces.Node {
	return node.entity
}
