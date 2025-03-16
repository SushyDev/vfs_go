package database

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

func (database *Database) InsertNodeContent(node interfaces.Node, content []byte) error {
	_, err := database.db.Exec(
		"INSERT INTO node_contents (node_id, content) VALUES (?, ?)",
		node.GetId(),
		content,
	)
	if err != nil {
		return err
	}

	return nil
}

func (database *Database) GetNodeContent(id int64) (interfaces.NodeContent, error) {
	row := database.db.QueryRow("SELECT id, node_id, content FROM node_contents WHERE id = ?", id)

	return database.nodeContentFactory.New(row)
}

func (database *Database) GetNodeContentByNode(node interfaces.Node) (interfaces.NodeContent, error) {
	row := database.db.QueryRow("SELECT id, node_id, content FROM node_contents WHERE node_id = ?", node.GetId())

	return database.nodeContentFactory.New(row)
}

func (database *Database) SaveNodeContent(nodeContent interfaces.NodeContent) error {
	_, err := database.db.Exec(
		"UPDATE node_contents SET content = ? WHERE id = ?",
		nodeContent.GetContent(),
		nodeContent.GetId(),
	)
	if err != nil {
		return err
	}

	return nil
}
