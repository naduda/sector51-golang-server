package main

import (
	"flag"
	"log"

	"github.com/naduda/sector51-golang/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.json", "path to config file")
}

func main() {
	flag.Parse()

	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
