package main

import (
	"flag"
	"fmt"
	"github.com/benschw/anvil-mgr/anvilrepo"
	"github.com/elazarl/go-bindata-assetfs"
	"log"
	"os"
)

func main() {
	var bind string
	var path string

	flag.StringVar(&bind, "bind", "localhost:8080", "Address to bind to")
	flag.StringVar(&path, "repo-path", "./repo", "anvil repository path root")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command> \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// pull desired command/operation from args
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Command argument required")
	}
	cmd := flag.Arg(0)

	if _, err := os.Stat(path); err != nil {
		os.Mkdir(path, 0755)
	}

	// Configure Server
	s, err := anvilrepo.NewAnvilRepoService(bind, path, &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "dist"})
	if err != nil {
		log.Fatal(err)
	}
	// Run Main App
	switch cmd {
	case "serve":

		// Start Server
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		log.Fatalf("Unknown Command: %s", cmd)
	}

}
