package database

import (
	"github.com/sushydev/vfs_go/internal/database/interfaces"
)

func (database *Database) GetNodeContent(id int64) (interfaces.NodeContent, error) {
	row := database.db.QueryRow("SELECT * FROM node_content WHERE id = ?", id)

	return database.nodeContentFactory.New(row)
}

func (database *Database) GetNodeContentByNode(node interfaces.Node) (interfaces.NodeContent, error) {
	row := database.db.QueryRow("SELECT * FROM node_content WHERE node_id = ?", node.GetId())

	return database.nodeContentFactory.New(row)
}

func (database *Database) SaveNodeContent(nodeContent interfaces.NodeContent) error {
	_, err := database.db.Exec(
		"INSERT INTO node_content (node_id, content) VALUES (?, ?)",
		nodeContent.GetNodeId(),
		nodeContent.GetContent(),
	)
	if err != nil {
		return err
	}

	return nil
}
