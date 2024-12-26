package service

import (
	"database/sql"
	"sync"

	vfs_node "github.com/sushydev/vfs_go/node"
)

type DirectoryService struct {
	db          *sql.DB
	nodeService *NodeService

	mu sync.RWMutex
}

func NewDirectoryService(db *sql.DB, nodeService *NodeService) *DirectoryService {
	return &DirectoryService{
		db:          db,
		nodeService: nodeService,
	}
}

func (service *DirectoryService) CreateDirectory(name string, parent *vfs_node.Directory) (*uint64, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	transaction, err := service.db.Begin()
	if err != nil {
		return nil, serviceError("Failed to begin transaction", err)
	}
	defer transaction.Rollback()

	identifier, err := service.nodeService.CreateNode(transaction, name, parent, vfs_node.DirectoryNode)
	if err != nil {
		return nil, serviceError("Failed to create node", err)
	}

	query := `
        INSERT INTO directories (node_id)
        VALUES (?)
    `

	result, err := transaction.Exec(query, identifier)
	if err != nil {
		return nil, serviceError("Failed to insert directory", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, serviceError("Failed to get rows affected", err)
	}

	if rowsAffected != 1 {
		return nil, serviceError("Failed to insert directory", err)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, serviceError("Failed to commit transaction", err)
	}

	return identifier, nil
}

func (service *DirectoryService) UpdateDirectory(identifier uint64, name string, parent *vfs_node.Directory) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	transaction, err := service.db.Begin()
	if err != nil {
		return serviceError("Failed to begin transaction", err)
	}
	defer transaction.Rollback()

	err = service.nodeService.UpdateNode(transaction, identifier, name, parent)
	if err != nil {
		return serviceError("Failed to update node", err)
	}

	query := `
        UPDATE directories SET node_id = ?
        WHERE node_id = ?
    `

	result, err := transaction.Exec(query, identifier, identifier)
	if err != nil {
		return serviceError("Failed to update directory", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceError("Failed to get rows affected", err)
	}

	if rowsAffected != 1 {
		return serviceError("Failed to update directory", err)
	}

	err = transaction.Commit()
	if err != nil {
		return serviceError("Failed to commit transaction", err)
	}

	return nil
}

func (service *DirectoryService) DeleteDirectory(identifier uint64) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	transaction, err := service.db.Begin()
	if err != nil {
		return serviceError("Failed to begin transaction", err)
	}
	defer transaction.Rollback()

	err = service.nodeService.DeleteNode(transaction, identifier)
	if err != nil {
		return serviceError("Failed to delete node", err)
	}

	err = transaction.Commit()
	if err != nil {
		return serviceError("Failed to commit transaction", err)
	}

	return nil
}

func (service *DirectoryService) GetDirectory(identifier uint64) (*vfs_node.Directory, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	query := `
        SELECT node_id, name, parent_id, type
        FROM DIRECTORIES
        LEFT JOIN nodes ON nodes.id = directories.node_id
        WHERE node_id = ? and type = ?
    `

	row := service.db.QueryRow(query, identifier, vfs_node.DirectoryNode.String())

	directory, err := getDirectoryFromRow(row)
	if err != nil {
		return nil, serviceError("Failed to get directory from row", err)
	}

	return directory, nil
}

func (service *DirectoryService) FindDirectory(name string, parent *vfs_node.Directory) (*vfs_node.Directory, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	var row *sql.Row

	if parent == nil {
		query := `
            SELECT nodes.id, name, parent_id, type
            FROM nodes
            LEFT JOIN directories ON directories.node_id = nodes.id
            WHERE name = ? and parent_id IS NULL and type = ?
        `

		row = service.db.QueryRow(query, name, vfs_node.DirectoryNode.String())
	} else {
		query := `
            SELECT nodes.id, name, parent_id, type
            FROM nodes
            LEFT JOIN directories ON directories.node_id = nodes.id
            WHERE name = ? and parent_id = ? and type = ?
        `

		row = service.db.QueryRow(query, name, parent.GetNode().GetIdentifier(), vfs_node.DirectoryNode.String())
	}

	directory, err := getDirectoryFromRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, serviceError("Failed to get directory from row", err)
	}

	return directory, nil
}

func (service *DirectoryService) FindChildNode(name string, parent *vfs_node.Directory) (*vfs_node.Node, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	query := `
        SELECT id, name, parent_id, type
        FROM nodes
        WHERE name = ? AND parent_id = ?
    `

	row := service.db.QueryRow(query, name, parent.GetNode().GetIdentifier())

	childNode, err := getNodeFromRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, serviceError("Failed to get node from row", err)
	}

	return childNode, nil
}

func (service *DirectoryService) GetChildNodes(parent *vfs_node.Directory) ([]*vfs_node.Node, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	query := `
        SELECT id, name, parent_id, type
        FROM nodes
        WHERE parent_id = ?
    `

	rows, err := service.db.Query(query, parent.GetIdentifier())
	if err != nil {
		return nil, serviceError("Failed to get child nodes", err)
	}
	defer rows.Close()

	nodes := make([]*vfs_node.Node, 0)

	for rows.Next() {
		node, err := getNodeFromRow(rows)
		if err != nil {
			return nil, serviceError("Failed to get node from row", err)
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// --- Helpers

func getDirectoryFromRow(row row) (*vfs_node.Directory, error) {
	node, err := getNodeFromRow(row)
	if err != nil {
		return nil, err
	}

	directory := vfs_node.NewDirectory(node)

	return directory, nil
}
