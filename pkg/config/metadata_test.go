package config

import "testing"

func TestMetadata(t *testing.T) {
	c := LoadMetaData("../../testdata")

	if c.data.String("hash.unicorn") == "" {
		t.Error("metadata should exist")
	}
}
