package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var configFilePath = filepath.Join(GetConfDir(), "config.json")

type Config struct {
	ClientID          string
	ClientSecret      string
	BackupLinksSecret string
	DownloadDir       string
}

func (t *Kanna) loadConfiguration() error {
	if err := os.MkdirAll(GetConfDir(), os.ModePerm); err != nil {
		return err
	}

	t.Config = &Config{}

	if confBytes, err := ioutil.ReadFile(configFilePath); err == nil {
		err = json.Unmarshal(confBytes, t.Config)
		if err != nil {
			return err
		}
	}
	err := t.Config.sanitiseConfigurations()
	if err != nil {
		return err
	}

	return t.saveConfiguration()
}

func (t *Kanna) saveConfiguration() error {
	confBytes, err := json.MarshalIndent(t.Config, "", "\t")
	if err != nil {
		return err
	}

	if err = os.MkdirAll(GetConfDir(), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath, confBytes, os.ModePerm)
}

func (c *Config) sanitiseConfigurations() error {
	if c.DownloadDir == "" {
		downloadDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		c.DownloadDir = filepath.Join(downloadDir, "Downloads")
	}
	return nil
}

func GetConfDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, err = os.UserHomeDir()
		if err != nil {
			configDir = ""
		}
	}

	return filepath.Join(configDir, "kanna")
}
