package filesystem

import (
	"syscall"
	"io/fs"
	"fmt"

	"github.com/sushydev/vfs_go/internal/database"
	"github.com/sushydev/vfs_go/interfaces"
	node_repository "github.com/sushydev/vfs_go/internal/filesystem/node/repository"
	node_content_repository "github.com/sushydev/vfs_go/internal/filesystem/node_content/repository"
	symlink_repository "github.com/sushydev/vfs_go/internal/filesystem/symlink/repository"
)

type FileSystem struct {
	database       *database.Database
	nodeRepository *node_repository.Repository
	nodeContentRepository *node_content_repository.Repository
	symlinkRepository *symlink_repository.Repository
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

func getPath(parentNode interfaces.Node, name string) string {
	if parentNode.GetPath() == "/" {
		return "/" + name
	}

	return parentNode.GetPath() + "/" + name
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

	if parentNode == nil {
		return nil, syscall.ENOENT
	}

	if !parentNode.GetMode().IsDir() {
		return nil, syscall.ENOTDIR
	}

	return f.nodeRepository.GetChildren(parentNode)
}

func (f *FileSystem) Lookup(parentId uint64, name string) (interfaces.Node, error) {
	parentNode, err := f.nodeRepository.Get(parentId)
	if err != nil {
		return nil, err
	}

	if parentNode == nil {
		return nil, syscall.ENOENT
	}

	if !parentNode.GetMode().IsDir() {
		return nil, syscall.ENOTDIR
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

	path := getPath(parentNode, name)

	return f.database.InsertNode(name, parentNode.GetEntity(), path, uint32(fs.ModeDir), 0, 0, 0, "", "")
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

	path := getPath(parentNode, name)

	return f.database.InsertNode(name, parentNode.GetEntity(), path, 0, 0, 0, 0, "", "")
}

func (f *FileSystem) WriteFile(id uint64, content []byte) (int, error) {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return 0, err
	}

	if node.GetMode() != fs.FileMode(0) {
		return 0, fmt.Errorf("node %s is not a file", node.GetName())
	}

	nodeContent, err := f.nodeContentRepository.GetByNode(node)
	nodeContent.SetContent(content)

	err = f.database.SaveNodeContent(nodeContent.GetEntity())
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

	if node.GetMode() != fs.FileMode(0) {
		return nil, fmt.Errorf("node %s is not a file", node.GetName())
	}

	nodeContent, err := f.nodeContentRepository.GetByNode(node)
	if err != nil {
		return nil, err
	}

	return nodeContent.GetContent(), nil
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

func (f *FileSystem) Move(id uint64, name string, newParentId uint64) error {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return err
	}

	if !node.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	parentNode, err := f.nodeRepository.Get(newParentId)
	if err != nil {
		return err
	}

	if !parentNode.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	path := getPath(parentNode, name)

	node.SetName(name)
	node.SetPath(path)
	node.SetParentId(newParentId)

	err = f.database.SaveNode(node.GetEntity())
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) Rename(id uint64, newName string, newParentId uint64) error {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return err
	}

	parentNode, err := f.nodeRepository.Get(newParentId)
	if err != nil {
		return err
	}

	if !parentNode.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	path := getPath(parentNode, newName)

	node.SetName(newName)
	node.SetPath(path)
	node.SetParentId(newParentId)

	err = f.database.SaveNode(node.GetEntity())
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystem) Link(id uint64, name string, parentId uint64) error {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return err
	}

	if node.GetMode().IsDir() {
		return syscall.EISDIR
	}

	parentNode, err := f.nodeRepository.Get(parentId)
	if err != nil {
		return err
	}

	if !parentNode.GetMode().IsDir() {
		return syscall.ENOTDIR
	}

	path := getPath(parentNode, name)

	err = f.database.InsertNode(name, parentNode.GetEntity(), path, uint32(fs.ModeSymlink), 0, 0, 0, "", "")
	if err != nil {
		return err
	}

	sourceNode, err := f.nodeRepository.GetByParentAndName(parentNode, name)
	if err != nil {
		return err
	}

	return f.database.InsertSymlink(sourceNode.GetEntity(), node.GetEntity())
}

func (f *FileSystem) ReadLink(id uint64) (string, error) {
	node, err := f.nodeRepository.Get(id)
	if err != nil {
		return "", err
	}

	if node.GetMode() != fs.ModeSymlink {
		return "", syscall.EINVAL
	}

	symlinkEntity, err := f.database.GetSymlinkBySourceNode(node.GetEntity())
	if err != nil {
		return "", err
	}

	symlink, err := f.symlinkRepository.GetByEntity(symlinkEntity)
	if err != nil {
		return "", err
	}

	targetNode, err := f.nodeRepository.Get(symlink.GetTargetNodeId())
	if err != nil {
		return "", err
	}

	return targetNode.GetPath(), nil
}

func (f *FileSystem) Save(node interfaces.Node) error {
	return f.database.SaveNode(node.GetEntity())
}
