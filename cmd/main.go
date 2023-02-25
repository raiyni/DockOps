package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raiyni/compose-ops/pkg/config"
	"github.com/rs/zerolog"
)

func init() {

}

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	configFile := flag.String("config", "config.yml", "config file path")
	dataDir := flag.String("data", "data", "data dir path")

	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if _, err := os.Stat(*dataDir); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	metaData := config.LoadMetaData(*dataDir)
	store := config.NewSource(*configFile, "main")

	fmt.Println(store.Services)
	fmt.Println(metaData)
}
