-- SQLite schema for a filesystem where everything is a file

PRAGMA foreign_keys = ON;  -- Enable foreign key support

CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,              -- File name
    parent_id INTEGER,               -- Parent directory ID (NULL for root)
    type TEXT NOT NULL,              -- 'regular', 'directory', or 'symlink'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    size INTEGER DEFAULT 0,          -- Size in bytes (0 for directories)
    permissions TEXT DEFAULT '644',  -- Unix-style permissions
    owner TEXT,                      -- Owner name/id
    target_id INTEGER,               -- For symlinks: the ID of the target file
    
    -- Self-referential foreign key for parent directories
    FOREIGN KEY (parent_id) REFERENCES files(id) ON DELETE CASCADE,
    
    -- Foreign key for symlink targets
    FOREIGN KEY (target_id) REFERENCES files(id) ON DELETE SET NULL,
    
    -- Enforce uniqueness of file names within a directory
    UNIQUE(parent_id, name)
);

-- Index for faster lookups by parent
CREATE INDEX idx_files_parent_id ON files(parent_id);

-- Index for faster lookups by type
CREATE INDEX idx_files_type ON files(type);

-- Index for faster lookups of symlink targets
CREATE INDEX idx_files_target_id ON files(target_id) WHERE target_id IS NOT NULL;

-- For file content storage (only for regular files)
CREATE TABLE file_contents (
    file_id INTEGER PRIMARY KEY,
    content BLOB,              -- Binary content of the file
    hash TEXT,                 -- Hash of the content for integrity checks
    
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);

-- Optional: Extended attributes table
CREATE TABLE file_attributes (
    file_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT,
    
    PRIMARY KEY (file_id, key),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);
