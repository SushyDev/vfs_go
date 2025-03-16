package repository

import (
	"github.com/sushydev/vfs_go/internal/database"
	"github.com/sushydev/vfs_go/interfaces"
	"github.com/sushydev/vfs_go/internal/filesystem/node"
)

type Repository struct {
	database *database.Database
}

func New(database *database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) Get(id uint64) (interfaces.Node, error) {
	entitiy, err := r.database.GetNode(int64(id))
	if err != nil {
		return nil, err
	}

	return node.New(entitiy)
}

func (r *Repository) GetByName(name string) (interfaces.Node, error) {
	entity, err := r.database.GetNodeByName(name)
	if err != nil {
		return nil, err
	}

	return node.New(entity)
}

func (r *Repository) GetByParentAndName(parent interfaces.Node, name string) (interfaces.Node, error) {
	entity, err := r.database.GetNodeByParentAndName(parent.GetEntity(), name)
	if err != nil {
		return nil, err
	}

	return node.New(entity)
}

func (r *Repository) GetChildren(parent interfaces.Node) ([]interfaces.Node, error) {
	entities, err := r.database.GetNodesByParent(parent.GetEntity())
	if err != nil {
		return nil, err
	}

	var nodes []interfaces.Node
	for _, entity := range entities {
		node, err := node.New(entity)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
