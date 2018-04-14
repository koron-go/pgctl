package pgctl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmp, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	dir := filepath.Join(tmp, "db")

	db := NewDB(dir)
	if err := db.Start(); err != nil {
		t.Fatal("failed to stat DB", err)
	}
	defer db.Stop()

	n, err := db.Name()
	if err != nil {
		t.Fatal("db.Name() failed", err)
	}
	if n != "postgres://postgres@127.0.0.1:5432/postgres" {
		t.Error("db.Name() returns unexpected:", n)
	}
}
