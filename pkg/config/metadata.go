package config

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/rs/zerolog/log"
)

type Metadata struct {
	dir  string
	data *config.Config
}

func LoadMetaData(dataDir string) *Metadata {
	c := config.NewEmpty("data")
	c.WithOptions(config.ParseEnv)
	c.AddDriver(config.JSONDriver)

	err := c.LoadFiles(fmt.Sprintf("%s/metadata.json", dataDir))
	if err != nil {
		log.Error().Msgf("failed to load metadata file: %s", err.Error())
	}

	return &Metadata{
		dir:  dataDir,
		data: c,
	}
}

func (m *Metadata) Hash(key string) string {
	return m.data.String(fmt.Sprintf("hash.%s", key))
}

func (m *Metadata) SetHash(key, value string) error {
	return m.data.Set(fmt.Sprintf("hash.%s", key), value)
}

func (m *Metadata) Save() error {
	log.Debug().Msgf("saving metadata: %s", m.data.ToJSON())

	if m.data.IsEmpty() {
		m.data.Set("hash", make(map[string]string))
	}

	return m.data.DumpToFile(fmt.Sprintf("%s/metadata.json", m.dir), config.JSON)
}
