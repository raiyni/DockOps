package config

import (
	"testing"
)

func TestConfigBinding(t *testing.T) {
	s := NewSource("../../config.yml", "test1")

	if s.auths == nil {
		t.Error("auths should exist")
	}

	if len(s.auths) == 0 {
		t.Error("auths shouldn't be empty")
	}

	if len(s.Services) == 0 {
		t.Error("services shouldn't be empty")
	}

	if s.Services[0].Url != "https://github.com/raiyni/compose-ops-example.git" {
		t.Error("services[0] url should be https://github.com/raiyni/compose-ops-example.git")
	}

	if s.Services[1].AuthObj.Username != "testy" {
		t.Error("services[1] should have an auth")
	}
}
