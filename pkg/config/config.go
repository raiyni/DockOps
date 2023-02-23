package config

import (
	"fmt"
	"os"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/rs/zerolog/log"
)

type Auth struct {
	Username    string
	Password    string
	KeyPath     string
	KeyPassword string
}

type Service struct {
	Name        string
	Url         string
	Ref         string `example:"refs/heads/master"`
	Interval    int    `default:"5"`
	Auth        string
	AuthObj     Auth
	Hash        string
	Path        string
	ComposeFile string `default:"docker-compose.yml"`
}

type store struct {
	auths    map[string]Auth
	Services []Service
}

func New() *store {
	return NewSource("config.yml", "main")
}

func NewSource(source, name string) *store {
	c := config.NewEmpty(name)
	c.WithOptions(config.ParseEnv, config.ParseDefault)
	c.AddDriver(yaml.Driver)

	err := c.LoadFiles(source)
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	auths := makeAuths(c)
	services := makeServices(c, auths)

	s := &store{
		Services: services,
		auths:    auths,
	}

	return s
}

func makeAuths(c *config.Config) map[string]Auth {
	ca := c.SubDataMap("auth")
	if c == nil {
		return nil
	}

	auths := make(map[string]Auth)
	for _, k := range ca.Keys() {
		a := &Auth{}

		key := "auth." + k
		err := c.MapStruct(key, a)
		if err != nil {
			log.Warn().Msgf("cannot decode %s", key)
		} else {
			auths[k] = *a
			log.Debug().Msgf("decoded %s as %s", key, *a)
		}
	}

	return auths
}

func makeServices(c *config.Config, auths map[string]Auth) []Service {
	ss := c.Get("services").([]any)
	if ss == nil {
		return nil
	}

	services := make([]Service, len(ss))
	for i := range ss {
		s := &Service{}

		key := fmt.Sprintf("services.%d", i)
		sub := c.SubDataMap(key)
		name := sub.Keys()[0]

		err := c.MapStruct(fmt.Sprintf("%s.%s", key, name), s)
		if err != nil {
			log.Warn().Msgf("cannot decode %s", key)
		} else {
			s.Name = name
			if s.Path == "" {
				s.Path = name
			}

			if s.Auth != "" {
				if auths == nil {
					log.Error().Msg("no auths defined")
				} else {
					a, ok := auths[s.Auth]
					if !ok {
						log.Error().Msgf("unable to find Auth %s", s.Auth)
					} else {
						s.AuthObj = a
					}
				}
			}

			services[i] = *s
			log.Debug().Msgf("decoded %s as {%s, %s, %d, %s}", key, s.Url, s.Ref, s.Interval, s.AuthObj)
		}
	}

	return services
}
