package pgctl

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPgctl(t *testing.T) {
	const port = 5453

	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmp, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	dir := filepath.Join(tmp, "data")

	// before InitDB
	err = Start(dir, &StartOptions{Port: port})
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

	err = Start(dir, &StartOptions{Port: port})
	if err != nil {
		t.Fatalf("Start() failed: %s", err)
	}
	time.Sleep(3 * time.Second)

	// after Start()
	err = Status(dir)
	if err != nil {
		t.Errorf("after Start: Status() failed: %s", err)
	}
	err = Start(dir, &StartOptions{Port: port})
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
