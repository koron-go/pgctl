package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koron-go/pgctl"
)

var dataDir string

func main() {
	flag.StringVar(&dataDir, "datadir", "", "database data directory")
	flag.Parse()
	if dataDir == "" {
		log.Fatal("require -datadir")
	}
	if flag.NArg() < 1 {
		log.Fatal("require sub-commands: initdb, start, status or stop")
	}
	sc := flag.Arg(0)
	switch sc {
	case "initdb":
		err := pgctl.InitDB(dataDir, nil)
		if err != nil {
			log.Fatalf("InitDB failed: %s", err)
		}
	case "start":
		err := pgctl.Start(dataDir, nil)
		if err != nil {
			log.Fatalf("InitDB failed: %s", err)
		}
	case "status":
		err := pgctl.Status(dataDir)
		if err != nil {
			log.Fatalf("Status failed: %s", err)
		}
		fmt.Printf("database is running: %s\n", dataDir)
	case "stop":
		err := pgctl.Stop(dataDir)
		if err != nil {
			log.Fatalf("Stop failed: %s", err)
		}
	default:
		log.Fatalf("unknown sub-command: %s", sc)
	}
}
