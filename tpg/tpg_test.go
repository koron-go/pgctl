package tpg

import (
	"os"
	"testing"

	"github.com/koron-go/pgctl"
)

func TestServer(t *testing.T) {
	if !pgctl.IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	s1 := New(t)
	defer s1.Close()

	s2 := New(t)
	defer s2.Close()

	// two servers should run in parallel
	if s1.Dir() == s2.Dir() {
		t.Errorf("same dir: %s", s1.Dir())
	}
	if s1.Port() == s2.Port() {
		t.Errorf("same port: %d", s1.Port())
	}
	if s1.Name() == s2.Name() {
		t.Errorf("same name: %s", s1.Name())
	}

	// working dir should be removed after closed.
	dir := s1.Dir()
	s1.Close()
	_, err := os.Stat(dir)
	if err == nil {
		t.Errorf("data dir is remained: %s", dir)
	} else if !os.IsNotExist(err) {
		t.Errorf("unexpected error: %s", err)
	}
}
