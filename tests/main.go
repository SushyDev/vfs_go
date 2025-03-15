package main

import (
	"github.com/sushydev/vfs_go"
)

func main() {
	fs, err := filesystem.New("./vfs.db")
	if err != nil {
		panic(err)
	}

	root, err := fs.Root()
	if err != nil {
		panic(err)
	}

	err = fs.MkDir(root.GetId(), "dir")
	if err != nil {
		panic(err)
	}

	err = fs.Touch(root.GetId(), "file.txt")
	if err != nil {
		panic(err)
	}

	_, err = fs.Open(1)
	if err != nil {
		panic(err)
	}
}
