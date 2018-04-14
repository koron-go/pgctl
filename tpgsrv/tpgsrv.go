package tpgsrv

import (
	"io/ioutil"
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

// New creates an instance of PostgreSQL server and starts it.
func New(tb testing.TB) *Server {
	tb.Helper()
	dir, err := ioutil.TempDir("", "tpgsrv-")
	if err != nil {
		tb.Fatal("failed to create dir for PostgreSQL server:", err)
	}
	psrv := pgctl.NewServer(filepath.Join(dir, "data"))
	pn, err := start(psrv)
	if err != nil {
		os.RemoveAll(dir)
		tb.Fatal("failed to start PostgreSQL server:", err)
	}
	return &Server{tb: tb, dir: dir, psrv: psrv, pn: pn}
}

// DefaultPort is default port listening by PostgreSQL.
var DefaultPort uint16 = 15432

var lastIndex uint16

const numPorts uint16 = 10

func start(psrv *pgctl.Server) (uint16, error) {
	// select unused port
	var err error
	for i := uint16(0); i < 10; i++ {
		port := DefaultPort + (lastIndex+i)%numPorts
		psrv.StartOptions(&pgctl.StartOptions{Port: port})
		err = psrv.Start()
		if err == nil {
			lastIndex = i + 1
			return port, nil
		}
	}
	return 0, err
}

// Close closes a server and removes all related resources.
func (s *Server) Close() {
	s.tb.Helper()
	if err := s.psrv.Stop(); err != nil {
		s.tb.Error("failed to stop PostgreSQL server:", err)
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
