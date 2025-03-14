package main

import (
	"github.com/sushydev/vfs_go/database"
	"github.com/sushydev/vfs_go/filesystem"
)

func main() {
	database, err := database.New()
	if err != nil {
		panic(err)
	}

	fs := filesystem.New(database)

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
