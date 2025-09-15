package pgctl

import (
	"path/filepath"
	"testing"
)

func TestPgctl(t *testing.T) {
	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmpdir := t.TempDir()
	dir := filepath.Join(tmpdir, "data")

	// before InitDB
	err := Start(dir, &StartOptions{})
	if err != ErrNotInitialized {
		t.Errorf("before InitDB: Start() failed unexpectedly: %s", err)
	}
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Errorf("before InitDB: Stop() failed unexpectedly: %s", err)
	}
	err = Status(dir)
	if err != ErrNotInitialized {
		t.Errorf("before InitDB: Status failed unexpectedly: %s", err)
	}

	err = InitDB(dir, nil)
	if err != nil {
		t.Fatalf("InitDB() failed: %s", err)
	}

	// after InitDB() and before Start()
	err = InitDB(dir, nil)
	if err != ErrAlreadyExists {
		t.Errorf("after InitDB: InitDB() failed unexpectedly: %s", err)
	}
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Errorf("after InitDB: Stop() failed unexpectedly: %s", err)
	}
	err = Status(dir)
	if err != ErrNotRunning {
		t.Errorf("after InitDB: Status() failed unexpectedly: %s", err)
	}

	err = Start(dir, &StartOptions{})
	if err != nil {
		t.Fatalf("Start() failed: %s", err)
	}

	// after Start()
	err = Status(dir)
	if err != nil {
		t.Errorf("after Start: Status() failed: %s", err)
	}
	err = Start(dir, &StartOptions{})
	if err != ErrAlreadyRunning {
		t.Errorf("after Start: Start() failed unexpectedly: %s", err)
	}

	err = Stop(dir)
	if err != nil {
		t.Fatalf("Stop() failed: %s", err)
	}

	// after Stop()
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Errorf("after Stop: Stop() failed unexpectedly: %s", err)
	}
	err = Status(dir)
	if err != ErrNotRunning {
		t.Errorf("after Stop: Status() failed unexpectedly: %s", err)
	}
}
