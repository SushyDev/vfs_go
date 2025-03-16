package database

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

// todo return last inserted id
func (d *Database) InsertNode(
	name string,
	parent interfaces.Node,
	path string,
	mode uint32,
	uid int,
	gid int,
	modTime int,
	createTime string,
	accessTime string,
) error {
	var parentId int64

	if parent != nil {
		parentId = int64(parent.GetId())
	} else {
		parentId = 0
	}

	parsedMode := int64(mode)

	_, err := d.db.Exec(`
		INSERT INTO nodes (name, parent_id, path, mode, uid, gid, mod_time, create_time, access_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, name, parentId, path, parsedMode, uid, gid, modTime, createTime, accessTime)

	return err
}

func (database *Database) GetNode(id int64) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE id = ?
	`, id)

	return database.nodeFactory.New(row)
}

func (database *Database) GetNodeByName(name string) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE name = ?
	`, name)

	return database.nodeFactory.New(row)
}

func (database *Database) GetNodesByParent(parent interfaces.Node) ([]interfaces.Node, error) {
	rows, err := database.db.Query(`
		SELECT id, name, parent_id, path, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE parent_id = ?
	`, parent.GetId())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []interfaces.Node
	for rows.Next() {
		file, err := database.nodeFactory.New(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, file)
	}

	return nodes, nil
}

func (database *Database) GetNodeByParentAndName(parent interfaces.Node, name string) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE parent_id = ? AND name = ?
	`, parent.GetId(), name)

	return database.nodeFactory.New(row)
}

func (d *Database) DeleteNode(node interfaces.Node) error {
	_, err := d.db.Exec(`
		DELETE FROM nodes
		WHERE id = ?
	`, node.GetId())
	return err
}

func (database *Database) SaveNode(node interfaces.Node) error {
	_, err := database.db.Exec(`
		UPDATE nodes
		SET name = ?, parent_id = ?, path = ?, mode = ?, uid = ?, gid = ?, mod_time = ?, create_time = ?, access_time = ?
		WHERE id = ?
	`,
		node.GetName(),
		node.GetParentId(),
		node.GetPath(),
		node.GetMode(),
		node.GetUid(),
		node.GetGid(),
		node.GetModTime(),
		node.GetCreateTime(),
		node.GetAccessTime(),
		node.GetId(),
	)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
