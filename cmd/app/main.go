package main

import (
	"1._file-sync/internal/app"
	"flag"
)

const configPath = "config/config.yaml"

func main() {
	cfgPath := flag.String("config", configPath, "config path")
	flag.Parse()
	app.Run(*cfgPath)
}
