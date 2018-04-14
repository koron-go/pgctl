package pgctl

import (
	"fmt"
	"os"
	"sync"
)

// DB represents PostgreSQL instance.
type DB struct {
	m   sync.Mutex
	r   bool
	dir string
	io  *InitDBOptions
	so  *StartOptions
}

// NewDB creates an instance of PostgreSQL.
func NewDB(dataDir string) *DB {
	return &DB{dir: dataDir}
}

// InitDBOptions sets InitDBOptions
func (db *DB) InitDBOptions(io *InitDBOptions) error {
	db.m.Lock()
	defer db.m.Unlock()
	if db.r {
		return ErrAlreadyRunning
	}
	db.io = io
	return nil
}

// StartOptions sets StartOptions
func (db *DB) StartOptions(so *StartOptions) error {
	db.m.Lock()
	defer db.m.Unlock()
	if db.r {
		return ErrAlreadyRunning
	}
	db.so = so
	return nil
}

// Start starts PostgreSQL in background.
func (db *DB) Start() error {
	db.m.Lock()
	defer db.m.Unlock()
	if db.r {
		return ErrAlreadyRunning
	}
	if db.io == nil {
		db.io = &InitDBOptions{}
	}
	if db.so == nil {
		db.so = &StartOptions{}
	}
	if _, err := os.Stat(db.dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := InitDB(db.dir, db.io); err != nil {
			return err
		}
	}
	if err := Start(db.dir, db.so); err != nil {
		return err
	}
	db.r = true
	return nil
}

// Stop stops an instance of PostgreSQL.
func (db *DB) Stop() error {
	db.m.Lock()
	defer db.m.Unlock()
	if !db.r {
		return ErrNotRunning
	}
	if err := Stop(db.dir); err != nil {
		return err
	}
	db.r = false
	return nil
}

// Name gets data source name.
func (db DB) Name() (string, error) {
	db.m.Lock()
	defer db.m.Unlock()
	if !db.r {
		return "", ErrNotRunning
	}
	return fmt.Sprintf("postgres://%[1]s@%[2]s:%[3]s/%[1]s", db.io.user(), db.so.host(), db.so.portString()), nil
}
