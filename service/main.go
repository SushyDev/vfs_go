package service

import (
	"fmt"
	"syscall"

	"github.com/sushydev/vfs_go"
	"github.com/sushydev/vfs_go/interfaces"
)

func GetRoot(fileSystem *filesystem.FileSystem) (interfaces.Node, error) {
	node, err := fileSystem.Root()
	if err != nil {
		return nil, err
	}

	if !node.GetMode().IsDir() {
		return nil, syscall.ENOTDIR
	}

	return node, nil
}

func GetFile(fileSystem *filesystem.FileSystem, id uint64) (interfaces.Node, error) {
	node, err := fileSystem.Open(id)
	if err != nil {
		return nil, err
	}

	if !node.GetMode().IsRegular() {
		return nil, syscall.EISDIR
	}

	return node, nil
}

func GetDirectory(fileSystem *filesystem.FileSystem, id uint64) (interfaces.Node, error) {
	node, err := fileSystem.Open(id)
	if err != nil {
		return nil, err
	}

	if !node.GetMode().IsDir() {
		return nil, syscall.ENOTDIR
	}

	return node, nil
}

func FindFile(fileSystem *filesystem.FileSystem, name string) (interfaces.Node, error) {
	node, err := fileSystem.Find(name)
	if err != nil {
		return nil, err
	}

	if !node.GetMode().IsRegular() {
		return nil, syscall.EISDIR
	}

	return node, nil
}

func FindDirectory(fileSystem *filesystem.FileSystem, name string) (interfaces.Node, error) {
	node, err := fileSystem.Find(name)
	if err != nil {
		return nil, err
	}

	if !node.GetMode().IsDir() {
		return nil, syscall.ENOTDIR
	}

	return node, nil
}

func FindOrCreateFile(fileSystem *filesystem.FileSystem, parentId uint64, name string) (interfaces.Node, error) {
	existingNode, err := fileSystem.Lookup(parentId, name)
	switch err {
	case nil:
		if !existingNode.GetMode().IsRegular() {
			return nil, fmt.Errorf("node %s is not a file", name)
		}

		return existingNode, nil
	case syscall.ENOENT:
		err := fileSystem.Touch(parentId, name)
		if err != nil {
			return nil, err
		}

		node, err := fileSystem.Lookup(parentId, name)
		if err != nil {
			return nil, err
		}

		if !node.GetMode().IsRegular() {
			return nil, syscall.EISDIR
		}

		return node, nil
	default:
		return nil, err
	}
}

func FindOrCreateDirectory(fileSystem *filesystem.FileSystem, parentId uint64, name string) (interfaces.Node, error) {
	existingNode, err := fileSystem.Lookup(parentId, name)
	switch err {
	case nil:
		if !existingNode.GetMode().IsDir() {
			return nil, fmt.Errorf("node %s is not a directory", name)
		}

		return existingNode, nil
	case syscall.ENOENT:
		err := fileSystem.MkDir(parentId, name)
		if err != nil {
			return nil, err
		}

		node, err := fileSystem.Lookup(parentId, name)
		if err != nil {
			return nil, err
		}

		if !node.GetMode().IsDir() {
			return nil, syscall.ENOTDIR
		}

		return node, nil
	default:
		return nil, err
	}
}
