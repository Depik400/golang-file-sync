package main

import (
	"flag"
	"golang-file-sync/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	cfgPath := flag.String("config", configPath, "config path")
	flag.Parse()
	app.Run(*cfgPath)
}
