package main

import (
	"flag"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raiyni/compose-ops/pkg/config"
	"github.com/rs/zerolog"
)

func init() {

}

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	configFile := flag.String("config", "config.yml", "config file path")

	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	store := config.NewSource(*configFile, "main")

	fmt.Println(store.Services)
}
