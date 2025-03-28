package runtime_data

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

type Config struct {
	MaxAge  time.Duration
	DataDir string

	key    []byte
	loaded bool
}

var config Config

func GetConfig() *Config {
	return &config
}

func LoadConfig(path string) error {
	if config.loaded {
		return fmt.Errorf("config already loaded")
	}

	cfg, err := loadOrCreateConfig(path)
	if err != nil {
		return err
	}
	dataDir := cfg.Section("storage").Key("dataDir").MustString("images")
	maxAgeStr := cfg.Section("storage").Key("maxFileAge").MustString("24h")
	if dataDir[len(dataDir)-1] != '/' {
		dataDir += "/"
	}

	if config.MaxAge, err = time.ParseDuration(maxAgeStr); err != nil {
		return err
	}

	config.DataDir = "data/" + dataDir
	if err = ensureDataDir(config.DataDir); err != nil {
		return err
	}

	config.loaded = true
	return nil
}

func ensureDataDir(dataDir string) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	return nil
}

func loadOrCreateConfig(path string) (*ini.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := ini.Empty()
		cfg.Section("storage").Key("dataDir").SetValue("data")
		cfg.Section("storage").Key("maxFileAge").SetValue("24h")

		err := cfg.SaveTo(path)
		if err != nil {
			return nil, err
		}
		log.Printf("Created default config at %s", path)
	}

	return ini.Load(path)
}

func (cfg *Config) InitKey() error {
	if len(cfg.key) > 0 {
		return fmt.Errorf("key already initialized")
	}
	cfg.key = make([]byte, 32)
	if _, err := rand.Read(cfg.key); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) GetKey() []byte {
	return cfg.key
}
