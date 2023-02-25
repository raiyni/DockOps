package config

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	"github.com/rs/zerolog/log"
)

type metadata struct {
	dir  string
	data *config.Config
}

func LoadMetaData(dataDir string) *metadata {
	c := config.NewEmpty("data")
	c.WithOptions(config.ParseEnv)
	c.AddDriver(json.Driver)

	err := c.LoadFiles(fmt.Sprintf("%s/metadata.json", dataDir))
	if err != nil {
		log.Error().Msgf("failed to load metadata file: %s", err.Error())
	}

	return &metadata{
		dir:  dataDir,
		data: c,
	}
}

func (m *metadata) Hash(key string) string {
	return m.data.String(fmt.Sprintf("hash.%s", key))
}

func (m *metadata) Save() error {
	buf := new(bytes.Buffer)

	config.JSONMarshalIndent = "    "
	_, err := m.data.DumpTo(buf, config.JSON)
	ioutil.WriteFile(fmt.Sprintf("%s/metadata.json", m.dir), buf.Bytes(), 0755)

	return err
}
