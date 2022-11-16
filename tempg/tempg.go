/*
Package tempg provides temporary PostgreSQL instance.
*/
package tempg

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/koron-go/pgctl"
)

// Server is a PostgreSQL instance for test.
type Server struct {
	dir  string
	port uint16
	pg   *pgctl.Server
}

// New creates an independent instance of PostgreSQL server and starts it.
func New() (*Server, error) {
	dir, err := os.MkdirTemp("", "tempg-")
	if err != nil {
		return nil, err
	}
	pg := pgctl.NewServer(filepath.Join(dir, "data"))
	port, err := start(pg)
	if err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	return &Server{
		dir:  dir,
		port: port,
		pg:   pg,
	}, nil
}

// defaultPort is default port listening by PostgreSQL.
var defaultPort uint16 = 25432

var lastIndex uint16

const numPorts uint16 = 1024

var mu sync.Mutex

func newPort() uint16 {
	mu.Lock()
	newPort := defaultPort + lastIndex%numPorts
	lastIndex++
	mu.Unlock()
	return newPort
}

func start(pg *pgctl.Server) (uint16, error) {
	// select unused port
	var err error
	for i := uint16(0); i < 3; i++ {
		port := newPort()
		pg.StartOptions(&pgctl.StartOptions{Port: port})
		err = pg.Start()
		if err == nil {
			return port, nil
		}
	}
	return 0, err
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
	return s.port
}

// Name returns data source name.
func (s *Server) Name() string {
	return s.pg.Name()
}
