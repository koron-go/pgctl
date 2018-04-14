package pgctl

import (
	"fmt"
	"os"
	"sync"
)

// Server represents PostgreSQL instance.
type Server struct {
	m   sync.Mutex
	r   bool
	dir string
	io  *InitDBOptions
	so  *StartOptions
}

// NewServer creates an instance of PostgreSQL.
func NewServer(dataDir string) *Server {
	return &Server{dir: dataDir}
}

// InitDBOptions sets InitDBOptions
func (srv *Server) InitDBOptions(io *InitDBOptions) error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if srv.r {
		return ErrAlreadyRunning
	}
	srv.io = io
	return nil
}

// StartOptions sets StartOptions
func (srv *Server) StartOptions(so *StartOptions) error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if srv.r {
		return ErrAlreadyRunning
	}
	srv.so = so
	return nil
}

// Start starts PostgreSQL in background.
func (srv *Server) Start() error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if srv.r {
		return ErrAlreadyRunning
	}
	if srv.io == nil {
		srv.io = &InitDBOptions{}
	}
	if srv.so == nil {
		srv.so = &StartOptions{}
	}
	if _, err := os.Stat(srv.dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := InitDB(srv.dir, srv.io); err != nil {
			return err
		}
	}
	if err := Start(srv.dir, srv.so); err != nil {
		return err
	}
	srv.r = true
	return nil
}

// Stop stops an instance of PostgreSQL.
func (srv *Server) Stop() error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if !srv.r {
		return ErrNotRunning
	}
	if err := Stop(srv.dir); err != nil {
		return err
	}
	srv.r = false
	return nil
}

// Name returns data source name if server is running.
// Otherwise returns empty string.
func (srv Server) Name() (string) {
	srv.m.Lock()
	defer srv.m.Unlock()
	if !srv.r {
		return ""
	}
	return fmt.Sprintf("postgres://%[1]s@%[2]s:%[3]s/%[1]s", srv.io.user(), srv.so.host(), srv.so.portString())
}
