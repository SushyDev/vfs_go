package factory

import (
	"database/sql"

	"github.com/sushydev/vfs_go/internal/database/node"
	"github.com/sushydev/vfs_go/database/interfaces"
)

type Factory struct {
	db *sql.DB
}

func New(db *sql.DB) *Factory {
	return &Factory{db: db}
}

func (factory *Factory) NewNode(row interfaces.RowScanner) (interfaces.Node, error) {
	var id uint64
	var name string
	var parentId uint64
	var path string
	var content []byte
	var size int64
	var mode uint32
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
		&content,
		&size,
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
		content,
		size,
		mode,
		uid,
		gid,
		modTime,
		createTime,
		accessTime,
	)
}
