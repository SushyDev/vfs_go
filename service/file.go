package service

import (
	"database/sql"
	"sync"

	vfs_node "github.com/sushydev/vfs_go/node"
)

type FileService struct {
	db          *sql.DB
	nodeService *NodeService

	mu sync.RWMutex
}

func NewFileService(db *sql.DB, nodeService *NodeService) *FileService {
	return &FileService{
		db:          db,
		nodeService: nodeService,
	}
}

func (service *FileService) CreateFile(name string, parent *vfs_node.Directory, contentType string, data string) (*uint64, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	transaction, err := service.db.Begin()
	if err != nil {
		return nil, serviceError("Failed to begin transaction", err)
	}
	defer transaction.Rollback()

	identifier, err := service.nodeService.CreateNode(transaction, name, parent, vfs_node.FileNode)
	if err != nil {
		return nil, serviceError("Failed to create node", err)
	}

	query := `
        INSERT INTO files (node_id, content_type, data)
        VALUES (?, ?, ?) 
    `

	result, err := transaction.Exec(query, identifier, contentType, data)
	if err != nil {
		return nil, serviceError("Failed to insert file", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, serviceError("Failed to get rows affected", err)
	}

	if rowsAffected != 1 {
		return nil, serviceError("Failed to insert file", err)
	}

	err = transaction.Commit()
	if err != nil {
		return nil, serviceError("Failed to commit transaction", err)
	}

	return identifier, nil
}

func (service *FileService) UpdateFile(identifier uint64, name string, parent *vfs_node.Directory, contentType string, data string) error {
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
        UPDATE files SET content_type = ?, data = ?
        WHERE node_id = ? 
    `

	result, err := transaction.Exec(query, contentType, data, identifier)
	if err != nil {
		return serviceError("Failed to update file", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return serviceError("Failed to get rows affected", err)
	}

	if rowsAffected != 1 {
		return serviceError("Failed to update file", err)
	}

	err = transaction.Commit()
	if err != nil {
		return serviceError("Failed to commit transaction", err)
	}

	return nil
}

func (service *FileService) DeleteFile(identifier uint64) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	transaction, err := service.db.Begin()
	if err != nil {
		return serviceError("Failed to begin transaction", err)
	}

	query := `
        DELETE FROM files
        WHERE node_id = ?
    `

	_, err = transaction.Exec(query, identifier)
	if err != nil {
		return serviceError("Failed to delete file", err)
	}

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

func (service *FileService) GetFile(identifier uint64) (*vfs_node.File, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	query := `
        SELECT n.id, name, parent_id, type, f.content_type, f.data
        FROM nodes n
        LEFT JOIN files f ON n.id = f.node_id
        WHERE n.id = ? AND type = ?
    `

	row := service.db.QueryRow(query, identifier, vfs_node.FileNode.String())

	file, err := getFileFromRow(row)
	if err != nil {
		return nil, serviceError("Failed to get file from row", err)
	}

	return file, nil
}

func (service *FileService) ListFiles() ([]*vfs_node.File, error) {
	query := `
        SELECT id, name, parent_id, type, f.content_type, f.data
        FROM nodes
        LEFT JOIN files f ON nodes.id = f.node_id
        WHERE type = ?
    `

	rows, err := service.db.Query(query)
	if err != nil {
		return nil, serviceError("Failed to list files", err)
	}
	defer rows.Close()

	var files []*vfs_node.File

	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, serviceError("Error occurred during rows iteration", err)
		}

		file, err := getFileFromRow(rows)
		if err != nil {
			return nil, serviceError("Failed to get file from row", err)
		}

		files = append(files, file)
	}

	return files, nil
}

func (service *FileService) FindFile(name string, parent *vfs_node.Directory) (*vfs_node.File, error) {
	service.mu.RLock()
	defer service.mu.RUnlock()

	var row *sql.Row

	if parent == nil {
		query := `
            SELECT n.id, name, parent_id, type, f.content_type, f.data
            FROM nodes n
            LEFT JOIN files f ON n.id = f.node_id
            WHERE parent_id IS NULL AND name = ?
        `

		row = service.db.QueryRow(query, name)
	} else {
		query := `
            SELECT n.id, name, parent_id, type, f.content_type, f.data
            FROM nodes n
            LEFT JOIN files f ON n.id = f.node_id
            WHERE parent_id = ? AND name = ?
        `

		row = service.db.QueryRow(query, parent.GetNode().GetIdentifier(), name)
	}

	file, err := getFileFromRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, serviceError("Failed to get file from row", err)
	}

	return file, nil
}

func (service *FileService) GetFiles(parent *vfs_node.Directory) ([]*vfs_node.File, error) {
	var parentIdentifier sql.NullInt64
	if parent != nil {
		parentIdentifier.Scan(parent.GetNode().GetIdentifier())
	}

	query := `
        SELECT n.id, name, parent_id, type, f.content_type, f.data
        FROM nodes n
        LEFT JOIN files f ON n.id = f.node_id
        WHERE parent_id = ?
    `

	rows, err := service.db.Query(query, parentIdentifier)
	if err != nil {
		return nil, serviceError("Failed to find files by parent", err)
	}
	defer rows.Close()

	var files []*vfs_node.File

	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, serviceError("Error occurred during rows iteration", err)
		}

		file, err := getFileFromRow(rows)
		if err != nil {
			return nil, serviceError("Failed to get file from row", err)
		}

		files = append(files, file)
	}

	return files, nil
}

func getFileFromRow(row row) (*vfs_node.File, error) {
	var identifier uint64
	var name string
	var parentIdentifier sql.NullInt64
	var nodeTypeStr string
	var content_type string
	var data string

	err := row.Scan(&identifier, &name, &parentIdentifier, &nodeTypeStr, &content_type, &data)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, sql.ErrNoRows
		default:
			return nil, serviceError("Failed to scan file", err)
		}
	}

	var parent_identifier *uint64
	if parentIdentifier.Valid {
		parentIdentifierValue := uint64(parentIdentifier.Int64)
		parent_identifier = &parentIdentifierValue
	}

	nodeType := vfs_node.NodeTypeFromString(nodeTypeStr)

	node := vfs_node.NewNode(identifier, name, parent_identifier, nodeType)

	file := vfs_node.NewFile(node, content_type, data)

	return file, nil
}
