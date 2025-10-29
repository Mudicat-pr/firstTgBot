package config

import (
	"log"
	"os"

	"github.com/Mudicat-pr/firstTgBot/pkg/e"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AdminID  int64  `yaml:"admin_id"`
	TokenStr string `yaml:"bot_token"`
}

func ReadConfig() (cfg *Config, err error) {
	defer func() { err = e.WrapIfErr("can't read config", err) }()
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("I cant found config: %v", err)
		return nil, err
	}
	cfg = &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func MustToken() (token string) {
	cfg, err := ReadConfig()
	if err != nil {
		log.Fatalf("Empty token: %v", err)
		return ""
	}
	token = cfg.TokenStr
	if token == "" {
		log.Fatal("Empty token")
		return ""
	}
	return token
}
