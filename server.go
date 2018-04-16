package pgctl

import (
	"os"
	"sync"
)

// Server represents PostgreSQL instance.
type Server struct {
	m   sync.Mutex
	r   bool
	dir string
	io  InitDBOptions
	so  StartOptions
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
	if io != nil {
		srv.io = *io
	} else {
		srv.io = InitDBOptions{}
	}
	return nil
}

// StartOptions sets StartOptions
func (srv *Server) StartOptions(so *StartOptions) error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if srv.r {
		return ErrAlreadyRunning
	}
	if so != nil {
		srv.so = *so
	} else {
		srv.so = StartOptions{}
	}
	return nil
}

// Start starts PostgreSQL in background.
func (srv *Server) Start() error {
	srv.m.Lock()
	defer srv.m.Unlock()
	if srv.r {
		return ErrAlreadyRunning
	}
	if _, err := os.Stat(srv.dir); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := InitDB(srv.dir, &srv.io); err != nil {
			return err
		}
	}
	if err := Start(srv.dir, &srv.so); err != nil {
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

// IsRunning checks PostgreSQL server is running or not.
func (srv *Server) IsRunning() bool {
	srv.m.Lock()
	defer srv.m.Unlock()
	return srv.r
}

// Name returns data source name if server is running.
// Otherwise returns empty string.
func (srv *Server) Name() string {
	srv.m.Lock()
	defer srv.m.Unlock()
	if !srv.r {
		return ""
	}
	return Name(&srv.io, &srv.so)
}
