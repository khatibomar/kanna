package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

func New(envFile *os.File) (*Config, error) {

	b, err := ioutil.ReadAll(envFile)
	if err != nil {
		return &Config{}, err
	}
	var conf Config
	err = toml.Unmarshal(b, &conf)
	return &conf, err
}
