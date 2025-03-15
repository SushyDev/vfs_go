package factory

import (
	"database/sql"

	"github.com/sushydev/vfs_go/internal/database/symlink"
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Factory struct {
	db *sql.DB
}

func New(db *sql.DB) *Factory {
	return &Factory{db: db}
}

func (factory *Factory) New(row interfaces.RowScanner) (interfaces.Symlink, error) {
	var id int64
	var sourceNodeId int64
	var targetNodeId int64

	err := row.Scan(
		&id,
		&sourceNodeId,
		&targetNodeId,
	)
	if err != nil {
		return nil, err
	}

	return symlink.New(
		id,
		sourceNodeId,
		targetNodeId,
	)
}
