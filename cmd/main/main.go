package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"spire-reader/app/config"
	"spire-reader/app/server"
)

var (
	configPath string
	profile    string
)

func init() {
	baseDir, _ := os.Getwd()
	flag.StringVar(&configPath, "config", filepath.Join(baseDir, "configs", "spire-reader.toml"), "Path to config files")
	flag.StringVar(&profile, "profile", "local", "Profile where app runs, accept two options docker/local")
}

func main() {
	flag.Parse()
	config := configs.NewConfig()

	if err := configs.ParseConfig(configPath, profile, config); err != nil {
		log.Fatalf("%v", err)
	}
	if err := server.Start(config); err != nil {
		log.Fatalf("%v", err)
	}
}
