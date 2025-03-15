package node

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Node struct {
	id          int64
	name        string
	parentId    int64
	path        string
	mode        int64
	uid         int
	gid         int
	modTime     string
	createTime  string
	accessTime  string
}

var _ interfaces.Node = &Node{}

func New(
	id int64,
	name string,
	parentId int64,
	path string,
	mode int64,
	uid int,
	gid int,
	modTime string,
	createTime string,
	accessTime string,
) (*Node, error) {
	return &Node{
		id:          id,
		name:        name,
		parentId:    parentId,
		path:        path,
		mode:        mode,
		uid:         uid,
		gid:         gid,
		modTime:     modTime,
		createTime:  createTime,
		accessTime:  accessTime,
	}, nil
}

func (node *Node) GetId() int64 {
	return node.id
}

func (node *Node) GetName() string {
	return node.name
}

func (node *Node) GetParentId() int64 {
	return node.parentId
}

func (node *Node) GetPath() string {
	return node.path
}

func (node *Node) GetMode() int64 {
	return node.mode
}

func (node *Node) GetUid() int {
	return node.uid
}

func (node *Node) GetGid() int {
	return node.gid
}

func (node *Node) GetModTime() string {
	return node.modTime
}

func (node *Node) GetCreateTime() string {
	return node.createTime
}

func (node *Node) GetAccessTime() string {
	return node.accessTime
}

func (node *Node) SetName(name string) {
	node.name = name
}

func (node *Node) SetParentId(parentId int64) {
	node.parentId = parentId
}

func (node *Node) SetPath(path string) {
	node.path = path
}

func (node *Node) SetMode(mode int64) {
	node.mode = mode
}

func (node *Node) SetUid(uid int) {
	node.uid = uid
}

func (node *Node) SetGid(gid int) {
	node.gid = gid
}

func (node *Node) SetModTime(modTime string) {
	node.modTime = modTime
}

func (node *Node) SetCreateTime(createTime string) {
	node.createTime = createTime
}

func (node *Node) SetAccessTime(accessTime string) {
	node.accessTime = accessTime
}
