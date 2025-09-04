package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

var (
	ErrMissingConfig = errors.New(".gatorconfig.json is missing")
)

type Config struct {
	DbUrl    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func Read() (*Config, error) {
	fp, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.Username = username
	return cfg.write()
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, configFileName)
	return path, nil
}

func (cfg Config) write() error {
	fp, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(cfg)
}
