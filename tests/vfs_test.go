package tests

import (
	"testing"

	vfs "github.com/sushydev/vfs_go"
	"github.com/sushydev/vfs_go/node"
)

func TestNewFileSystem(t *testing.T) {
	fs, err := vfs.NewFileSystem()
	if err != nil {
		t.Fatalf("Failed to create FileSystem: %v", err)
	}
	if fs.GetRoot() == nil {
		t.Fatalf("Root directory not initialized")
	}
}

func TestCreateDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, err := fs.CreateDirectory("test_create_directory", root)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if dir == nil || dir.GetName() != "test_create_directory" {
		t.Fatalf("Directory not created correctly")
	}
}

func TestCreateExistingDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	_, err := fs.CreateDirectory("test_create_existing_directory", root)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	_, err = fs.CreateDirectory("test_create_existing_directory", root)
	if err == nil {
		t.Fatalf("Expected error when creating existing directory, got nil")
	}
}

func TestDeleteDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, err := fs.CreateDirectory("test_delete_directory", root)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	err = fs.DeleteDirectory(dir)
	if err != nil {
		t.Fatalf("Failed to delete directory: %v", err)
	}
}

func TestDeleteNonEmptyDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	dir, _ := fs.CreateDirectory("test_delete_non_empty_directory", root)
	fs.CreateFile("testFile", dir, "text/plain", "host")
	err := fs.DeleteDirectory(dir)
	if err != nil {
		t.Fatalf("Failed to delete non-empty directory: %v", err)
	}
}

func TestCreateFile(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	file, err := fs.CreateFile("test_create_file", root, "text/plain", "host")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	if file == nil || file.GetName() != "test_create_file" {
		t.Fatalf("File not created correctly")
	}
}

func TestCreateExistingFile(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	fs.CreateFile("test_create_existing_file", root, "text/plain", "host")
	_, err := fs.CreateFile("test_create_existing_file", root, "text/plain", "host")
	if err == nil {
		t.Fatalf("Expected error when creating existing file, got nil")
	}
}

func TestDeleteFile(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	file, _ := fs.CreateFile("test_delete_file", root, "text/plain", "host")
	err := fs.DeleteFile(file)
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}
}

func TestMoveFileToImpossibleLocation(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	file, _ := fs.CreateFile("test_move_file_to_impossible_location", root, "text/plain", "host")
	_, err := fs.UpdateFile(file, "test_move_file_to_impossible_location", &node.Directory{}, "text/plain", "host")
	if err == nil {
		t.Fatalf("Expected error when moving file to impossible location, got nil")
	}
}

func TestMoveDirectoryToImpossibleLocation(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	dir, _ := fs.CreateDirectory("test_move_directory_to_impossible_location", root)
	_, err := fs.UpdateDirectory(dir, "test_move_directory_to_impossible_location", &node.Directory{})
	if err == nil {
		t.Fatalf("Expected error when moving directory to impossible location, got nil")
	}
}

func TestMoveFile(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()
	file, _ := fs.CreateFile("test_move_file", root, "text/plain", "host")
	dir, _ := fs.CreateDirectory("test_move_file_dir", root)
	file, err := fs.UpdateFile(file, "test_move_file", dir, "text/plain", "host")
	if err != nil {
		t.Fatalf("Failed to move file: %v", err)
	}
}

func TestMoveDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, _ := fs.CreateDirectory("test_move_directory", root)
	dir2, _ := fs.CreateDirectory("test_move_directory_dir", root)

	dir, err := fs.UpdateDirectory(dir, "test_move_directory", dir2)
	if err != nil {
		t.Fatalf("Failed to move directory: %v", err)
	}
}

func TestFindDirectory(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, _ := fs.CreateDirectory("test_find_directory", root)

	foundDir, err := fs.FindDirectory("test_find_directory", root)
	if err != nil {
		t.Fatalf("Failed to find directory: %v", err)
	}
	if foundDir == nil || foundDir.GetName() != dir.GetName() {
		t.Fatalf("Directory not found correctly")
	}
}

func TestFindFile(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	file, _ := fs.CreateFile("test_find_file", root, "text/plain", "host")

	foundFile, err := fs.FindFile("test_find_file", root)
	if err != nil {
		t.Fatalf("Failed to find file: %v", err)
	}
	if foundFile == nil || foundFile.GetName() != file.GetName() {
		t.Fatalf("File not found correctly")
	}
}

func TestGetChildNode(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, _ := fs.CreateDirectory("test_get_child_node_dir", root)
	file, _ := fs.CreateFile("test_get_child_node_file", dir, "text/plain", "host")

	node, err := fs.GetChildNode("test_get_child_node_file", dir)
	if err != nil {
		t.Fatalf("Failed to get child node: %v", err)
	}

	if node == nil || node.GetName() != file.GetName() {
		t.Fatalf("Child node not found correctly")
	}
}

func TestGetChildNodes(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, _ := fs.CreateDirectory("test_get_child_nodes_dir", root)
	fs.CreateFile("test_get_child_nodes_file", dir, "text/plain", "host")

	nodes, err := fs.GetChildNodes(dir)
	if err != nil {
		t.Fatalf("Failed to get child nodes: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("Expected 1 child node, got %d", len(nodes))
	}
}

func TestGetFiles(t *testing.T) {
	fs, _ := vfs.NewFileSystem()
	root := fs.GetRoot()

	dir, _ := fs.CreateDirectory("test_get_files_dir", root)
	fs.CreateFile("test_get_files_file", dir, "text/plain", "host")

	files, err := fs.GetFiles(dir)
	if err != nil {
		t.Fatalf("Failed to get files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}
}
