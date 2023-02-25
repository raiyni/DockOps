package main

import (
	"errors"
	"flag"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/raiyni/compose-ops/pkg/config"
	"github.com/raiyni/compose-ops/pkg/git"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var metaData *config.Metadata
var store *config.Store
var dataDir *string

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	configFile := flag.String("config", "config.yml", "config file path")
	dataDir = flag.String("data", "data", "data dir path")

	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if _, err := os.Stat(*dataDir); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	metaData = config.LoadMetaData(*dataDir)
	store = config.NewSource(*configFile, "main")

	initialClone()
}

func initialClone() {
	for _, s := range store.Services {
		g := git.NewGitClient(s, *dataDir)

		savedHash := metaData.Hash(s.Name)
		remoteHash, err := g.PullMostRecent(context.TODO(), savedHash)
		if err != nil {
			log.Error().Msgf("unable to clone: %s", err)
		}

		if savedHash == remoteHash && err == nil {
			log.Debug().Msgf("hash up to date: %s", s.Name)
		} else if savedHash != remoteHash {
			log.Debug().Msgf("new hash found for %s: %s", s.Name, remoteHash)
			metaData.SetHash(s.Name, remoteHash)
		}
	}

	err := metaData.Save()
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
