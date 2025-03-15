package node_content

import (
	"github.com/sushydev/vfs_go/interfaces"
	database_interfaces "github.com/sushydev/vfs_go/internal/database/interfaces"
)

type NodeContent struct {
	entity database_interfaces.NodeContent
}

var _ interfaces.NodeContent = &NodeContent{}

func New(entity database_interfaces.NodeContent) (*NodeContent, error) {
	return &NodeContent{
		entity: entity,
	}, nil
}

func (nodeContent *NodeContent) GetId() uint64 {
	return uint64(nodeContent.entity.GetId())
}

func (nodeContent *NodeContent) GetNodeId() uint64 {
	return uint64(nodeContent.entity.GetNodeId())
}

func (nodeContent *NodeContent) GetContent() []byte {
	return nodeContent.entity.GetContent()
}

func (nodeContent *NodeContent) SetNodeId(nodeId uint64) {
	nodeContent.entity.SetNodeId(int64(nodeId))
}

func (nodeContent *NodeContent) SetContent(content []byte) {
	nodeContent.entity.SetContent(content)
}

func (nodeContent *NodeContent) GetEntity() database_interfaces.NodeContent {
	return nodeContent.entity
}
