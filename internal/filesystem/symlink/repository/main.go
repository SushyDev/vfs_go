package repository

import (
	"syscall"

	"github.com/sushydev/vfs_go/internal/database"
	"github.com/sushydev/vfs_go/interfaces"
	"github.com/sushydev/vfs_go/internal/filesystem/symlink"
	database_interfaces "github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Repository struct {
	database *database.Database
}

func New(database *database.Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) Get(id int64) (interfaces.Symlink, error) {
	entitiy, err := r.database.GetSymlink(id)
	if err != nil {
		return nil, syscall.ENOENT
	}

	return symlink.New(entitiy)
}

func (r *Repository) GetByEntity(entity database_interfaces.Symlink) (interfaces.Symlink, error) {
	return symlink.New(entity)
}
