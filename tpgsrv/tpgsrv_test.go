package tpgsrv

import (
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	s1 := New(t)
	defer s1.Close()

	s2 := New(t)
	defer s2.Close()

	// two servers should run in parallel
	if s1.Dir() == s2.Dir() {
		t.Error("same dir:", s1.Dir())
	}
	if s1.Port() == s2.Port() {
		t.Error("same port:", s1.Port())
	}
	if s1.Name() == s2.Name() {
		t.Error("same name:", s1.Name())
	}

	// working dir should be removed after closed.
	dir := s1.Dir()
	s1.Close()
	_, err := os.Stat(dir)
	if err == nil {
		t.Error("data dir is remained:", dir)
	} else if !os.IsNotExist(err) {
		t.Error("unexpected error:", err)
	}
}
