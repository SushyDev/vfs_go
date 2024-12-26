package vfs

import (
	"fmt"
	"github.com/sushydev/vfs_go/database"
	"github.com/sushydev/vfs_go/node"
	"github.com/sushydev/vfs_go/service"
	"log"
	"sync"
)

type FileSystemInterface interface {
	GetRoot() *node.Directory

	// --- Directory

	FindOrCreateDirectory(name string, parent *node.Directory) (*node.Directory, error)
	FindDirectory(name string, parent *node.Directory) (*node.Directory, error)
	CreateDirectory(name string, parent *node.Directory) (*node.Directory, error)
	DeleteDirectory(directory *node.Directory) error
	UpdateDirectory(directory *node.Directory, name string, parent *node.Directory) (*node.Directory, error)
	GetDirectory(identifier uint64) (*node.Directory, error)

	// --- File

	FindOrCreateFile(name string, parent *node.Directory, contentType string, data string) (*node.File, error)
	FindFile(name string, parent *node.Directory) (*node.File, error)
	CreateFile(name string, parent *node.Directory, contentType string, data string) (*node.File, error)
	DeleteFile(file *node.File) error
	UpdateFile(file *node.File, name string, parent *node.Directory, contentType string, data string) (*node.File, error)
	GetFile(identifier uint64) (*node.File, error)

	// --- Node

	GetChildNodes(parent *node.Directory) ([]*node.Node, error)
	FindChildNode(name string, parent *node.Directory) (*node.Node, error)
	GetNode(identifier uint64) (*node.Node, error)
}

var _ FileSystemInterface = &FileSystem{}

type FileSystem struct {
	root *node.Directory

	nodeService      *service.NodeService
	directoryService *service.DirectoryService
	fileService      *service.FileService

	mu sync.RWMutex
}

func NewFileSystem(name string, file string) (*FileSystem, error) {
	database, err := database.New(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to create index\n%w", err)
	}

	nodeService := service.NewNodeService()
	directoryService := service.NewDirectoryService(database, nodeService)
	fileService := service.NewFileService(database, nodeService)

	fileSystem := &FileSystem{
		nodeService:      nodeService,
		directoryService: directoryService,
		fileService:      fileService,
	}

	root, err := fileSystem.FindOrCreateDirectory(name, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to get root directory\n%w", err)
	}

	fileSystem.root = root

	return fileSystem, nil
}

func (fileSystem *FileSystem) GetRoot() *node.Directory {
	fileSystem.mu.Lock()
	defer fileSystem.mu.Unlock()

	return fileSystem.root
}

// --- Directory

func (fileSystem *FileSystem) FindOrCreateDirectory(name string, parent *node.Directory) (*node.Directory, error) {
	directory, err := fileSystem.FindDirectory(name, parent)
	if err != nil {
		return nil, err
	}

	if directory != nil {
		return directory, nil
	}

	directory, err = fileSystem.CreateDirectory(name, parent)
	if err != nil {
		return nil, err
	}

	return directory, nil
}

func (fileSystem *FileSystem) FindDirectory(name string, parent *node.Directory) (*node.Directory, error) {
	fileSystem.mu.Lock()
	defer fileSystem.mu.Unlock()

	return fileSystem.directoryService.FindDirectory(name, parent)
}

func (fileSystem *FileSystem) CreateDirectory(name string, parent *node.Directory) (*node.Directory, error) {
	fileSystem.mu.Lock()

	nodeId, err := fileSystem.directoryService.CreateDirectory(name, parent)
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory\n%w", err)
	}

	if nodeId == nil {
		return nil, fmt.Errorf("Failed to create directory\n")
	}

	fileSystem.mu.Unlock()

	return fileSystem.GetDirectory(*nodeId)
}

func (fileSystem *FileSystem) DeleteDirectory(directory *node.Directory) error {
	fileSystem.mu.Lock()
	defer fileSystem.mu.Unlock()

	// get all child nodes
	childNodes, err := fileSystem.directoryService.GetChildNodes(directory)
	if err != nil {
		return fmt.Errorf("Failed to get child nodes\n%w", err)
	}

	// delete all child nodes
	for _, childNode := range childNodes {
		switch childNode.GetType() {
		case node.DirectoryNode:
			directory, err := fileSystem.GetDirectory(childNode.GetIdentifier())
			if err != nil {
				return fmt.Errorf("Failed to get directory\n%w", err)
			}

			err = fileSystem.DeleteDirectory(directory)
			if err != nil {
				return fmt.Errorf("Failed to delete directory\n%w", err)
			}
		case node.FileNode:
			file, err := fileSystem.GetFile(childNode.GetIdentifier())
			if err != nil {
				return fmt.Errorf("Failed to get file\n%w", err)
			}

			err = fileSystem.DeleteFile(file)
			if err != nil {
				return fmt.Errorf("Failed to delete file\n%w", err)
			}
		}
	}

	return fileSystem.directoryService.DeleteDirectory(directory.GetIdentifier())
}

func (fileSystem *FileSystem) UpdateDirectory(directory *node.Directory, name string, parent *node.Directory) (*node.Directory, error) {
	fileSystem.mu.Lock()

	node := directory.GetNode()
	if node == nil {
		return nil, fmt.Errorf("Node is nil")
	}

	err := fileSystem.directoryService.UpdateDirectory(node.GetIdentifier(), name, parent)
	if err != nil {
		return nil, fmt.Errorf("Failed to update directory\n%w", err)
	}

	fileSystem.mu.Unlock()

	return fileSystem.GetDirectory(node.GetIdentifier())
}

func (fileSystem *FileSystem) GetDirectory(identifier uint64) (*node.Directory, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.directoryService.GetDirectory(identifier)
}

// --- File

func (fileSystem *FileSystem) FindOrCreateFile(name string, parent *node.Directory, contentType string, data string) (*node.File, error) {
	file, err := fileSystem.FindFile(name, parent)
	if err != nil {
		log.Printf("Failed to find file %s\n", name)
		return nil, err
	}

	if file != nil {
		return file, nil
	}

	file, err = fileSystem.CreateFile(name, parent, contentType, data)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fileSystem *FileSystem) CreateFile(name string, parent *node.Directory, contentType string, data string) (*node.File, error) {
	fileSystem.mu.Lock()

	identifier, err := fileSystem.fileService.CreateFile(name, parent, contentType, data)
	if err != nil {
		return nil, fmt.Errorf("Failed to register file\n%w", err)
	}

	fileSystem.mu.Unlock()

	return fileSystem.GetFile(*identifier)
}

func (fileSystem *FileSystem) DeleteFile(file *node.File) error {
	fileSystem.mu.Lock()
	defer fileSystem.mu.Unlock()

	return fileSystem.fileService.DeleteFile(file.GetNode().GetIdentifier())
}

func (fileSystem *FileSystem) UpdateFile(file *node.File, name string, parent *node.Directory, contentType string, data string) (*node.File, error) {
	fileSystem.mu.Lock()

	node := file.GetNode()
	if node == nil {
		return nil, fmt.Errorf("Node is nil")
	}

	err := fileSystem.fileService.UpdateFile(
		node.GetIdentifier(),
		name,
		parent,
		contentType,
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to update file\n%w", err)
	}

	fileSystem.mu.Unlock()

	return fileSystem.GetFile(node.GetIdentifier())
}

func (fileSystem *FileSystem) GetFile(identifier uint64) (*node.File, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.fileService.GetFile(identifier)
}

func (fileSystem *FileSystem) FindFile(name string, parent *node.Directory) (*node.File, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.fileService.FindFile(name, parent)
}

func (fileSystem *FileSystem) GetFiles(parent *node.Directory) ([]*node.File, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.fileService.GetFiles(parent)
}

// --- Node

func (fileSystem *FileSystem) GetChildNodes(parent *node.Directory) ([]*node.Node, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.directoryService.GetChildNodes(parent)
}

func (fileSystem *FileSystem) GetNode(identifier uint64) (*node.Node, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	// return fileSystem.nodeService.GetNode(identifier)
	return nil, nil
}

func (fileSystem *FileSystem) FindChildNode(name string, parent *node.Directory) (*node.Node, error) {
	fileSystem.mu.RLock()
	defer fileSystem.mu.RUnlock()

	return fileSystem.directoryService.FindChildNode(name, parent)
}
