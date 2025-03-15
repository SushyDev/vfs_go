package database

import (
	"database/sql"

	"github.com/sushydev/vfs_go/internal/database/interfaces"
	node_factory "github.com/sushydev/vfs_go/internal/database/node/factory"
	node_content_factory "github.com/sushydev/vfs_go/internal/database/node_content/factory"
	symlink_factory "github.com/sushydev/vfs_go/internal/database/symlink/factory"

	_ "modernc.org/sqlite"
)

var schema = `
-- Main nodes table that stores both regular files and directories
CREATE TABLE IF NOT EXISTS nodes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,                                             -- Base name of the file/directory
	parent_id INTEGER,                                              -- Parent directory ID (NULL for root)
	path TEXT NOT NULL UNIQUE,                                      -- Full path for easy lookup
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

-- Index for faster name lookups
CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name);

-- Index for faster mode lookups
CREATE INDEX IF NOT EXISTS idx_nodes_type ON nodes(mode);

-- Index for faster parent directory lookups
CREATE INDEX IF NOT EXISTS idx_nodes_parent ON nodes(parent_id);

-- Insert the root directory
INSERT OR IGNORE INTO nodes (id, name, parent_id, path, mode, mod_time, create_time, access_time)
VALUES (0, 'root', -1, '/', 2147483648, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

---- File contents table that stores file content ----

-- File contents table that stores file content
CREATE TABLE IF NOT EXISTS node_contents (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,                                    -- Node ID
	content BLOB NOT NULL,                                       -- File content
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE -- Ensure node exists
);

-- Index for faster content lookups
CREATE INDEX IF NOT EXISTS idx_contents_node ON node_contents(node_id);

---- File metadata table that stores extended metadata ----

-- Node attributes table that stores extended attributes
CREATE TABLE IF NOT EXISTS node_attributes (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	node_id INTEGER NOT NULL,                                    -- Node ID
	key TEXT NOT NULL,                                           -- Attribute key
	value TEXT NOT NULL,                                         -- Attribute value
	FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE -- Ensure node exists
);

-- Index for faster attribute lookups
CREATE INDEX IF NOT EXISTS idx_attributes_node ON node_attributes(node_id);

---- Symlink table that stores symbolic links ----

-- Symlink table that stores symbolic links
CREATE TABLE IF NOT EXISTS symlinks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	source_node_id INTEGER NOT NULL,                                     -- Source node ID
	target_node_id INTEGER NOT NULL,                                     -- Target node ID
	FOREIGN KEY (source_node_id) REFERENCES nodes(id) ON DELETE CASCADE, -- Ensure source node exists
	FOREIGN KEY (target_node_id) REFERENCES nodes(id) ON DELETE CASCADE  -- Ensure target node exists
);

-- Index for faster symlink source lookups
CREATE INDEX IF NOT EXISTS idx_symlinks_source ON symlinks(source_node_id);

-- Index for faster symlink target lookups
CREATE INDEX IF NOT EXISTS idx_symlinks_target ON symlinks(target_node_id);
`

type Database struct {
	db          *sql.DB
	nodeFactory *node_factory.Factory
	nodeContentFactory *node_content_factory.Factory
	symlinkFactory *symlink_factory.Factory
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
