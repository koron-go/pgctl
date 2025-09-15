/*
Package tpg provides utilities to write tests with pgctl package.
*/
package tpg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/koron-go/pgctl"
)

// Server is a PostgreSQL instance for test.
type Server struct {
	tb   testing.TB
	dir  string
	psrv *pgctl.Server
	pn   uint16
}

// New creates an independent instance of PostgreSQL server and starts it.
func New(tb testing.TB) *Server {
	tb.Helper()
	dir, err := os.MkdirTemp("", "tpg-")
	if err != nil {
		tb.Fatal("failed to create dir for PostgreSQL server:", err)
	}
	psrv := pgctl.NewServer(filepath.Join(dir, "data"))
	err = psrv.Start()
	if err != nil {
		os.RemoveAll(dir)
		tb.Fatal("failed to start PostgreSQL server:", err)
	}
	return &Server{tb: tb, dir: dir, psrv: psrv, pn: psrv.Port()}
}

// Close closes a server and removes all related resources.
func (s *Server) Close() {
	s.tb.Helper()
	if err := s.psrv.Stop(); err != nil {
		if err != pgctl.ErrNotRunning {
			s.tb.Error("failed to stop PostgreSQL server:", err)
		}
	}
	if err := os.RemoveAll(s.dir); err != nil {
		s.tb.Errorf("failed to remove %q: %s", s.dir, err)
	}
}

// Dir gets assigned working directory.
func (s *Server) Dir() string {
	return s.dir
}

// Port gets assigned port number.
func (s *Server) Port() uint16 {
	return s.pn
}

// Name returns data source name.
func (s *Server) Name() string {
	return s.psrv.Name()
}
