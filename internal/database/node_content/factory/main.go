package factory

import (
	"database/sql"

	"github.com/sushydev/vfs_go/internal/database/node_content"
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Factory struct {
	db *sql.DB
}

func New(db *sql.DB) *Factory {
	return &Factory{db: db}
}

func (factory *Factory) New(row interfaces.RowScanner) (interfaces.NodeContent, error) {
	var id int64
	var nodeId int64
	var content []byte

	err := row.Scan(
		&id,
		&nodeId,
		&content,
	)
	if err != nil {
		return nil, err
	}

	return node_content.New(
		id,
		nodeId,
		content,
	)
}
