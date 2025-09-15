/*
Package tempg provides temporary PostgreSQL instance.
*/
package tempg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/koron-go/pgctl"
)

// Server is a PostgreSQL instance for test.
type Server struct {
	dir string
	pg  *pgctl.Server
}

// New creates an independent instance of PostgreSQL server and starts it.
func New() (*Server, error) {
	dir, err := os.MkdirTemp("", "tempg-")
	if err != nil {
		return nil, err
	}
	pg := pgctl.NewServer(filepath.Join(dir, "data"))
	err = pg.Start()
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &Server{
		dir: dir,
		pg:  pg,
	}, nil
}

// Close closes a server and removes all related resources.
func (s *Server) Close() error {
	if err := s.pg.Stop(); err != nil {
		if err != pgctl.ErrNotRunning {
			return fmt.Errorf("failed to stop PostgreSQL server: %w", err)
		}
	}
	if err := os.RemoveAll(s.dir); err != nil {
		return fmt.Errorf("failed to remove %q : %w", s.dir, err)
	}
	return nil
}

// Dir gets assigned working directory.
func (s *Server) Dir() string {
	return s.dir
}

// Port gets assigned port number.
func (s *Server) Port() uint16 {
	return s.pg.Port()
}

// Name returns data source name.
func (s *Server) Name() string {
	return s.pg.Name()
}
