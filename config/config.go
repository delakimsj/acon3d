package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	User struct {
		UseAuthHeader bool `yaml:"use_auth_header"`
	} `yaml:"User"`
}

func (c *Config) GetConfig(path string) {
	var f *os.File
	var err error

	f, err = os.Open(path)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(c)
	if err != nil {
		panic(err)
	}
}
