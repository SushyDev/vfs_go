package symlink

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Symlink struct {
	id       int64
	sourceNodeId int64
	targetNodeId int64
}

var _ interfaces.Symlink = &Symlink{}

func New(id int64, sourceNodeId int64, targetNodeId int64) (*Symlink, error) {
	return &Symlink{
		id:       id,
		sourceNodeId: sourceNodeId,
		targetNodeId: targetNodeId,
	}, nil
}

func (symlink *Symlink) GetId() int64 {
	return symlink.id
}

func (symlink *Symlink) GetSourceNodeId() int64 {
	return symlink.sourceNodeId
}

func (symlink *Symlink) GetTargetNodeId() int64 {
	return symlink.targetNodeId
}

func (symlink *Symlink) SetSourceNodeId(sourceNodeId int64) {
	symlink.sourceNodeId = sourceNodeId
}

func (symlink *Symlink) SetTargetNodeId(targetNodeId int64) {
	symlink.targetNodeId = targetNodeId
}
