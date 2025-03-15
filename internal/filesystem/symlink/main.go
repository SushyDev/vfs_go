package symlink

import (
	"github.com/sushydev/vfs_go/interfaces"
	database_interfaces "github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Symlink struct {
	entity database_interfaces.Symlink
}

var _ interfaces.Symlink = &Symlink{}

func New(entity database_interfaces.Symlink) (*Symlink, error) {
	return &Symlink{
		entity: entity,
	}, nil
}

func (symlink *Symlink) GetId() uint64 {
	return uint64(symlink.entity.GetId())
}

func (symlink *Symlink) GetSourceNodeId() uint64 {
	return uint64(symlink.entity.GetSourceNodeId())
}

func (symlink *Symlink) GetTargetNodeId() uint64 {
	return uint64(symlink.entity.GetTargetNodeId())
}

func (symlink *Symlink) SetSourceNodeId(sourceNodeId uint64) {
	symlink.entity.SetSourceNodeId(int64(sourceNodeId))
}

func (symlink *Symlink) SetTargetNodeId(targetNodeId uint64) {
	symlink.entity.SetTargetNodeId(int64(targetNodeId))
}

func (symlink *Symlink) GetEntity() database_interfaces.Symlink {
	return symlink.entity
}
