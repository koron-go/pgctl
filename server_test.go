package pgctl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
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
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to stat DB: %s", err)
	}
	defer srv.Stop()

	n := srv.Name()
	if n != "postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable" {
		t.Errorf("srv.Name() returns unexpected: %s", n)
	}
}
