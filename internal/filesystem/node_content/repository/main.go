package repository

import (
	"syscall"

	"github.com/sushydev/vfs_go/internal/database"
	"github.com/sushydev/vfs_go/interfaces"
	"github.com/sushydev/vfs_go/internal/filesystem/node_content"
)

type Repository struct {
	database *database.Database
}

func New(database *database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) Get(id uint64) (interfaces.NodeContent, error) {
	entitiy, err := r.database.GetNodeContent(int64(id))
	if err != nil {
		return nil, syscall.ENOENT
	}

	return node_content.New(entitiy)
}

func (r *Repository) GetByNode(node interfaces.Node) (interfaces.NodeContent, error) {
	entity, err := r.database.GetNodeContentByNode(node.GetEntity())
	if err != nil {
		return nil, syscall.ENOENT
	}

	return node_content.New(entity)
}
