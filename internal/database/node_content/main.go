package node_content

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type NodeContent struct {
	id      int64
	nodeId  int64
	content []byte
}

var _ interfaces.NodeContent = &NodeContent{}

func New(
	id int64,
	nodeId int64,
	content []byte,
) (*NodeContent, error) {
	return &NodeContent{
		id:      id,
		nodeId:  nodeId,
		content: content,
	}, nil
}

func (nodeContent *NodeContent) GetId() int64 {
	return nodeContent.id
}

func (nodeContent *NodeContent) GetNodeId() int64 {
	return nodeContent.nodeId
}

func (nodeContent *NodeContent) GetContent() []byte {
	return nodeContent.content
}

func (nodeContent *NodeContent) SetNodeId(nodeId int64) {
	nodeContent.nodeId = nodeId
}

func (nodeContent *NodeContent) SetContent(content []byte) {
	nodeContent.content = content
}

func (nodeContent *NodeContent) GetEntity() interfaces.NodeContent {
	return nodeContent
}
