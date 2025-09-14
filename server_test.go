package pgctl

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	const port = 5454

	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmp, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	dir := filepath.Join(tmp, "data")

	srv := NewServer(dir)
	srv.StartOptions(&StartOptions{Port: port})
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to stat DB: %s", err)
	}
	defer srv.Stop()

	var (
		got  = srv.Name()
		want = fmt.Sprintf("postgres://postgres@127.0.0.1:%d/postgres?sslmode=disable", port)
	)
	if got != want {
		t.Errorf("srv.Name() returns unexpected: %s", got)
	}
}
