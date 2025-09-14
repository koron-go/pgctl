package tempg_test

import (
	"os"
	"sync"
	"testing"

	"github.com/koron-go/pgctl"
	"github.com/koron-go/pgctl/tempg"
)

func TestServers(t *testing.T) {
	if !pgctl.IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	s1, err := tempg.New()
	if err != nil {
		t.Fatalf("failed to start #1 server: %s", err)
	}
	defer s1.Close()

	s2, err := tempg.New()
	if err != nil {
		t.Fatalf("failed to start #2 server: %s", err)
	}
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
	_, err = os.Stat(dir)
	if err == nil {
		t.Errorf("data dir is remained: %s", dir)
	} else if !os.IsNotExist(err) {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestParallelServers(t *testing.T) {
	if !pgctl.IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	t.Log("starting")
	var wgStart sync.WaitGroup
	var wgClose sync.WaitGroup
	for i := 0; i < 16; i++ {
		wgStart.Add(1)
		wgClose.Add(1)
		go func(n int) {
			s, err := tempg.New()
			if err != nil {
				t.Errorf("failed to start #%d server: %s", n, err)
				wgStart.Done()
				wgClose.Done()
				return
			}
			wgStart.Done()
			wgStart.Wait()
			err = s.Close()
			if err != nil {
				t.Errorf("failed to close #%d server: %s", n, err)
			}
			wgClose.Done()
		}(i)
	}
	wgStart.Wait()
	t.Log("started")
	wgClose.Wait()
	t.Log("closed")
}
