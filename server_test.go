package pgctl

import (
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
)

func TestDB(t *testing.T) {
	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmpdir := t.TempDir()
	dir := filepath.Join(tmpdir, "data")

	srv := NewServer(dir)
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to stat DB: %s", err)
	}
	t.Cleanup(func() {
		srv.Stop()
	})

	rx := regexp.MustCompile(`^postgres://postgres@([^:]*):([0-9]*)/postgres\?sslmode=disable$`)
	m := rx.FindStringSubmatch(srv.Name())
	if len(m) != 3 {
		t.Fatalf("not match: target=%q match=%+v", srv.Name(), m)
	}
	if m[1] != "localhost" {
		t.Errorf("unexpected host: %s", m[1])
	}
	if n, err := strconv.Atoi(m[2]); err != nil || n != int(srv.Port()) {
		t.Errorf("unexpected port: got=%s want=%d err=%s", m[2], srv.Port(), err)
	}
}
