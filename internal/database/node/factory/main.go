package factory

import (
	"database/sql"

	"github.com/sushydev/vfs_go/internal/database/node"
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

type Factory struct {
	db *sql.DB
}

func New(db *sql.DB) *Factory {
	return &Factory{db: db}
}

func (factory *Factory) New(row interfaces.RowScanner) (interfaces.Node, error) {
	var id int64
	var name string
	var parentId int64
	var path string
	var mode int64
	var uid int
	var gid int
	var modTime string
	var createTime string
	var accessTime string

	err := row.Scan(
		&id,
		&name,
		&parentId,
		&path,
		&mode,
		&uid,
		&gid,
		&modTime,
		&createTime,
		&accessTime,
	)
	if err != nil {
		return nil, err
	}

	return node.New(
		id,
		name,
		parentId,
		path,
		mode,
		uid,
		gid,
		modTime,
		createTime,
		accessTime,
	)
}
