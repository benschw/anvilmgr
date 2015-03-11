package main

import (
	"flag"
	"fmt"
	"git.bvops.net/scm/auto/anvilmgr.git/anvilrepo"
	"github.com/benschw/opin-go/config"
	"github.com/elazarl/go-bindata-assetfs"
	"log"
	"os"
)

func main() {

	var cfg struct {
		Bind     string
		RepoPath string
		LogPath  string
	}

	var cfgPath string
	var bind string
	var path string

	flag.StringVar(&cfgPath, "config", "/etc/anvilmgr/anvilmgr.yaml", "Path to Config File")
	// flag.StringVar(&bind, "bind", "", "Address to bind to")
	// flag.StringVar(&path, "repo-path", "", "anvil repository path root")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command> \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// load config from file if config path exists
	if _, err := os.Stat(cfgPath); err == nil {
		if err := config.Bind(cfgPath, &cfg); err != nil {
			log.Fatal(err)
		}
	}

	// load config from flags if they are set
	if bind != "" {
		cfg.Bind = bind
	}
	if path != "" {
		cfg.RepoPath = path
	}

	// pull desired command/operation from args
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Command argument required")
	}
	cmd := flag.Arg(0)

	f, err := os.OpenFile(cfg.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// create repo path if it doesn't exist
	if _, err := os.Stat(path); err != nil {
		os.Mkdir(path, 0755)
	}

	// Configure Server
	s, err := anvilrepo.NewAnvilRepoService(
		cfg.Bind,
		cfg.RepoPath,
		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "dist"},
	)
	if err != nil {
		log.Fatal(err)
	}
	// Run Main App
	switch cmd {
	case "serve":
		log.Println("Starting anvilmgr server")
		if flag.NArg() == 0 || cfg.Bind == "" || cfg.RepoPath == "" {
			flag.Usage()
			log.Fatal("bind and repo-path required")
		}

		// Start Server
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		log.Fatalf("Unknown Command: %s", cmd)
	}

}
