package database

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

func (d *Database) InsertSymlink(
	sourceNode interfaces.Node,
	targetNode interfaces.Node,
) error {
	_, err := d.db.Exec(`
		INSERT INTO symlinks (source_node_id, target_node_id)
		VALUES (?, ?)
	`, sourceNode.GetId(), targetNode.GetId())

	return err
}

func (database *Database) GetSymlink(id int64) (interfaces.Symlink, error) {
	row := database.db.QueryRow(`
		SELECT id, source_node_id, target_node_id
		FROM symlinks
		WHERE id = ?
	`, id)

	return database.symlinkFactory.New(row)
}

func (database *Database) GetSymlinkBySourceNode(sourceNode interfaces.Node) (interfaces.Symlink, error) {
	row := database.db.QueryRow(`
		SELECT id, source_node_id, target_node_id
		FROM symlinks
		WHERE source_node_id = ?
	`, sourceNode.GetId())

	return database.symlinkFactory.New(row)
}

func (database *Database) SaveSymlink(entity interfaces.Symlink) error {
	_, err := database.db.Exec(`
		UPDATE symlinks
		SET source_node_id = ?, target_node_id = ?
		WHERE id = ?
	`, entity.GetSourceNodeId(), entity.GetTargetNodeId(), entity.GetId())

	return err
}

func (database *Database) DeleteSymlink(symlink interfaces.Symlink) error {
	_, err := database.db.Exec(`
		DELETE FROM symlinks
		WHERE id = ?
	`, symlink.GetId())

	return err
}
