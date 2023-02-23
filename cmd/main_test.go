package main

import (
	"context"
	"flag"
	"testing"

	"github.com/raiyni/compose-ops/pkg/config"
	"github.com/raiyni/compose-ops/pkg/git"
	"github.com/rs/zerolog/log"
)

var debug = flag.Bool("debug", false, "sets log level to debug")
var configFile = flag.String("config", "config.yml", "config file path")

func TestLatestCommit(t *testing.T) {
	store := config.NewSource(*configFile, "main")

	g := git.NewGitClient(store.Services[0])
	hash, err := g.LatestCommit(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if hash == "" {
		t.Error("hash should not be empty")
	} else {
		log.Info().Str("hash", hash).Msg("")
	}
}
