package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/koron-go/pgctl"
)

var (
	dataDir string
	port    uint
)

func main() {
	flag.StringVar(&dataDir, "datadir", "", "database data directory")
	flag.UintVar(&port, "port", 0, "port number (start sub-command)")
	flag.Parse()
	if dataDir == "" {
		log.Fatal("require -datadir")
	}
	if port > 65535 {
		log.Fatal("port must be smaller than 65536")
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
		so := &pgctl.StartOptions{}
		if port > 0 {
			so.Port = uint16(port)
		}
		err := pgctl.Start(dataDir, so)
		if err != nil {
			log.Fatalf("Start failed: %s", err)
		}
		fmt.Println(pgctl.Name(&pgctl.InitDBOptions{}, so))
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
