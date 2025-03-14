package filesystem

import (
	"syscall"
	"io/fs"

	"github.com/sushydev/vfs_go/database"
	"github.com/sushydev/vfs_go/filesystem/interfaces"
	node_repository "github.com/sushydev/vfs_go/internal/filesystem/node/repository"
)

type FileSystem struct {
	database       *database.Database
	nodeRepository *node_repository.Repository
}

var _ interfaces.FileSystem = &FileSystem{}

func New(path string) (*FileSystem, error) {
	database, err := database.New(path)
	if err != nil {
		return nil, err
	}

	return &FileSystem{
		database:       database,
		nodeRepository: node_repository.New(database),
	}, nil
}

func (f *FileSystem) Root() (interfaces.Node, error) {
	return f.nodeRepository.Get(0)
}

func (f *FileSystem) Open(id uint64) (interfaces.Node, error) {
	return f.nodeRepository.Get(id)
}

func (f *FileSystem) Find(name string) (interfaces.Node, error) {
	return f.nodeRepository.GetByName(name)
}

func (f *FileSystem) ReadDir(id uint64) ([]interfaces.Node, error) {
	parentNode, err := f.nodeRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return f.nodeRepository.GetChildren(parentNode)
}

func (f *FileSystem) Lookup(id uint64, name string) (interfaces.Node, error) {
	parentNode, err := f.nodeRepository.Get(id)
	if err != nil {
		return nil, err
	}

	return f.nodeRepository.GetByParentAndName(parentNode, name)
}

func (f *FileSystem) MkDir(parentId uint64, name string) error {
	parentNode, err := f.nodeRepository.Get(parentId)
	if err != nil {
		return err
	}

	if !parentNode.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	path := parentNode.GetPath() + "/" + name

	return f.database.InsertNode(name, parentNode.GetEntity(), path, nil, 0, uint32(fs.ModeDir), 0, 0, 0, "", "")
}

// TODO RmDir -f flag
func (f *FileSystem) RmDir(id uint64) error {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return err
	}

	children, err := f.nodeRepository.GetChildren(node)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		return syscall.ENOTEMPTY
	}

	err = f.database.DeleteNode(node.GetEntity())
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) Touch(parentId uint64, name string) error {
	parentNode, err := f.nodeRepository.Get(parentId)
	if err != nil {
		return err
	}

	if !parentNode.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	path := parentNode.GetPath() + "/" + name

	return f.database.InsertNode(name, parentNode.GetEntity(), path, nil, 0, 0, 0, 0, 0, "", "")
}

func (f *FileSystem) WriteFile(id uint64, content []byte) (int, error) {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return 0, err
	}

	if node.GetMode().IsDir() {
		return 0, syscall.EISDIR
	}

	node.SetContent(content)

	err = f.database.SaveNode(node.GetEntity())
	if err != nil {
		return 0, err
	}

	return len(content), nil
}

func (f *FileSystem) ReadFile(id uint64) ([]byte, error) {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return nil, err
	}

	if node.GetMode().IsDir() {
		return nil, syscall.EISDIR
	}

	return node.GetContent(), nil
}

func (f *FileSystem) RemoveFile(id uint64) error {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return err
	}

	if node.GetMode().IsDir() {
		return syscall.EISDIR
	}

	err = f.database.DeleteNode(node.GetEntity())
	if err != nil {
		return err
	}

	return nil
}
