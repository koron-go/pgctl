package pgctl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPgctl(t *testing.T) {
	if !IsAvailable() {
		t.Skip("can't find pg_ctl")
	}

	tmp, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	dir := filepath.Join(tmp, "data")

	// before InitDB
	err = Start(dir, nil)
	if err != ErrStartDatabase {
		t.Error("unexpected Start() result", err)
	}
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}

	err = InitDB(dir, nil)
	if err != nil {
		t.Fatal("InitDB() failed", err)
	}

	// after InitDB() and before Start()
	err = InitDB(dir, nil)
	if err != ErrAlreadyExists {
		t.Error("InitDB() failed", err)
	}
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}

	err = Start(dir, nil)
	if err != nil {
		t.Fatal("Start() failed", err)
	}
	time.Sleep(3 * time.Second)

	// after Start()
	err = Status(dir)
	if err != nil {
		t.Error("Status() failed", err)
	}
	err = Start(dir, nil)
	if err != ErrAlreadyRunning {
		t.Error("Start() failed", err)
	}

	err = Stop(dir)
	if err != nil {
		t.Fatal("Stop() failed", err)
	}

	// after Stop()
	err = Stop(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
	err = Status(dir)
	if err != ErrNotRunning {
		t.Error("unexpected Stop() result", err)
	}
}
