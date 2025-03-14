package database

import (
	"database/sql"

	"github.com/sushydev/vfs_go/database/interfaces"
	node_factory "github.com/sushydev/vfs_go/internal/database/node/factory"

	_ "modernc.org/sqlite"
)

var schema = `
-- Main nodes table that stores both regular files and directories
CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,                                             -- Base name of the file/directory
    parent_id INTEGER,                                              -- Parent directory ID (NULL for root)
    path TEXT NOT NULL UNIQUE,                                      -- Full path for easy lookup
    content BLOB,                                                   -- File content (NULL for directories)
    size INTEGER NOT NULL DEFAULT 0,                                -- File size in bytes
    mode INTEGER NOT NULL,                                          -- File mode bits (including directory bit)
    uid INTEGER NOT NULL DEFAULT 0,                                 -- Owner user ID
    gid INTEGER NOT NULL DEFAULT 0,                                 -- Owner group ID
    mod_time TIMESTAMP NOT NULL,                                    -- Last modification time
    create_time TIMESTAMP NOT NULL,                                 -- Creation time
    access_time TIMESTAMP NOT NULL,                                 -- Last access time
    FOREIGN KEY (parent_id) REFERENCES nodes(id) ON DELETE CASCADE, -- Ensure parent directory exists
    UNIQUE (parent_id, name)                                        -- Ensure unique names within a directory
);

-- Index for faster path lookups
CREATE INDEX IF NOT EXISTS idx_nodes_path ON nodes(path);

-- Index for faster parent directory lookups
CREATE INDEX IF NOT EXISTS idx_nodes_parent ON nodes(parent_id);

-- Insert the root directory
INSERT OR IGNORE INTO nodes (id, name, parent_id, path, mode, mod_time, create_time, access_time)
VALUES (0, 'root', 0, '/', 2147483648, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
`

type Database struct {
	db          *sql.DB
	nodeFactory *node_factory.Factory
}

var _ interfaces.Database = &Database{}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	
	return &Database{db: db}, nil
}

// todo return last inserted id
func (d *Database) InsertNode(
	name string,
	parent interfaces.Node,
	path string,
	content []byte,
	size int64,
	mode uint32, // todo ModeType and ModePerm should be separated
	uid int,
	gid int,
	modTime int,
	createTime string,
	accessTime string,
) error {
	var parentId int64

	if parent != nil {
		parentId = int64(parent.GetId())
	} else {
		parentId = 0
	}

	_, err := d.db.Exec(`
		INSERT INTO nodes (name, parent_id, path, content, size, mode, uid, gid, mod_time, create_time, access_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, name, parentId, path, content, size, mode, uid, gid, modTime, createTime, accessTime)

	return err
}

func (database *Database) GetNode(id uint64) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, content, size, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE id = ?
	`, id)

	return database.nodeFactory.NewNode(row)
}

func (database *Database) GetNodeByName(name string) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, content, size, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE name = ?
	`, name)

	return database.nodeFactory.NewNode(row)
}

func (database *Database) GetNodesByParent(parent interfaces.Node) ([]interfaces.Node, error) {
	rows, err := database.db.Query(`
		SELECT id, name, parent_id, path, content, size, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE parent_id = ?
	`, parent.GetId())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []interfaces.Node
	for rows.Next() {
		file, err := database.nodeFactory.NewNode(rows)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, file)
	}

	return nodes, nil
}

func (database *Database) GetNodeByParentAndName(parent interfaces.Node, name string) (interfaces.Node, error) {
	row := database.db.QueryRow(`
		SELECT id, name, parent_id, path, content, size, mode, uid, gid, mod_time, create_time, access_time
		FROM nodes
		WHERE parent_id = ? AND name = ?
	`, parent.GetId(), name)

	return database.nodeFactory.NewNode(row)
}

func (d *Database) DeleteNode(node interfaces.Node) error {
	_, err := d.db.Exec(`
		DELETE FROM nodes
		WHERE id = ?
	`, node.GetId())
	return err
}

func (database *Database) SaveNode(node interfaces.Node) error {
	_, err := database.db.Exec(`
		UPDATE nodes
		SET name = ?, parent_id = ?, path = ?, content = ?, size = ?, mode = ?, uid = ?, gid = ?, mod_time = ?, create_time = ?, access_time = ?
		WHERE id = ?
	`,
		node.GetName(),
		node.GetParentId(),
		node.GetPath(),
		node.GetContent(),
		node.GetSize(),
		node.GetMode(),
		node.GetUid(),
		node.GetGid(),
		node.GetModTime(),
		node.GetCreateTime(),
		node.GetAccessTime(),
		node.GetId(),
	)

	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}
