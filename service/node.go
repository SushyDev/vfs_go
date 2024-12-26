package service

import (
	"database/sql"
	"sync"

	vfs_node "github.com/sushydev/vfs_go/node"
)

type NodeService struct {
	mu sync.RWMutex
}

func NewNodeService() *NodeService {
	return &NodeService{}
}

func GetParentIdentifier(parent *vfs_node.Directory) (sql.NullInt64, error) {
	var parentIdentifier sql.NullInt64

	if parent == nil {
		return parentIdentifier, nil
	}

	node := parent.GetNode()

	if node == nil {
		return parentIdentifier, serviceError("Parent node is nil", nil)
	}

	parentIdentifier.Scan(node.GetIdentifier())

	return parentIdentifier, nil
}

func (service *NodeService) CreateNode(tx *sql.Tx, name string, parent *vfs_node.Directory, nodeType vfs_node.NodeType) (*uint64, error) {
	existingNodeIdentifier, err := service.FindNode(tx, name, parent)
	if err != nil {
		return nil, serviceError("Failed to find node", err)
	}

	if existingNodeIdentifier != nil {
		return nil, serviceError("Node already exists", nil)
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	query := `
        INSERT INTO nodes (name, parent_id, type)
        VALUES (?, ?, ?)
        RETURNING id
    `

	parentIdentifier, err := GetParentIdentifier(parent)
	if err != nil {
		return nil, serviceError("Failed to get parent identifier", err)
	}

	row := tx.QueryRow(query, name, parentIdentifier, nodeType.String())

	var identifier uint64
	err = row.Scan(&identifier)
	if err != nil {
		return nil, serviceError("Failed to scan node", err)
	}

	return &identifier, nil
}

func (service *NodeService) UpdateNode(tx *sql.Tx, identifier uint64, name string, parent *vfs_node.Directory) error {
	existingNodeIdentifier, err := service.FindNode(tx, name, parent)
	if err != nil {
		return serviceError("Failed to find node", err)
	}

	if existingNodeIdentifier != nil && *existingNodeIdentifier != identifier {
		return serviceError("Node with name already exists", nil)
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	parentIdentifier, err := GetParentIdentifier(parent)
	if err != nil {
		return serviceError("Failed to get parent identifier", err)
	}

	query := `
        UPDATE nodes SET name = ?, parent_id = ?
        WHERE id = ?
    `

	result, err := tx.Exec(query, name, parentIdentifier, identifier)
	if err != nil {
		return serviceError("Failed to update node", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceError("Failed to get rows affected", err)
	}

	if rowsAffected != 1 {
		return serviceError("Failed to update node", nil)
	}

	return nil
}

func (service *NodeService) DeleteNode(tx *sql.Tx, identifier uint64) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	query := `
        DELETE FROM nodes
        WHERE id = ?
    `

	result, err := tx.Exec(query, identifier)
	if err != nil {
		return serviceError("Failed to delete node", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceError("Failed to get rows affected", err)
	}

	if rowsAffected == 0 {
		return serviceError("No node found with ID", nil)
	}

	return nil
}

func (service *NodeService) FindNode(tx *sql.Tx, name string, parent *vfs_node.Directory) (*uint64, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	parentIdentifier, err := GetParentIdentifier(parent)
	if err != nil {
		return nil, serviceError("Failed to get parent identifier", err)
	}

	query := `
        SELECT id
        FROM nodes
        WHERE name = ? AND parent_id = ?
    `

	row := tx.QueryRow(query, name, parentIdentifier)

	var identifier uint64
	err = row.Scan(&identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, serviceError("Failed to scan node", err)
	}

	return &identifier, nil
}

// --- Helpers

type row interface {
	Scan(dest ...interface{}) error
}

func getNodeFromRow(row row) (*vfs_node.Node, error) {
	var identifier uint64
	var name string
	var parentIdentifier sql.NullInt64
	var nodeTypeStr string

	err := row.Scan(&identifier, &name, &parentIdentifier, &nodeTypeStr)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, sql.ErrNoRows
		default:
			return nil, serviceError("Failed to scan node", err)
		}
	}

	var parentIdentifierPtr *uint64
	if parentIdentifier.Valid {
		parentIdentifierPtr = new(uint64)
		*parentIdentifierPtr = uint64(parentIdentifier.Int64)
	}

	nodeType := vfs_node.NodeTypeFromString(nodeTypeStr)

	node := vfs_node.NewNode(identifier, name, parentIdentifierPtr, nodeType)

	return node, nil
}
