package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/ilyakaznacheev/cleanenv"
)

func TestEnvConfig(t *testing.T) {
	os.Setenv("EMAIL_PORT", "123")
	config, err := NewConf()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}

func TestYamlConfig(t *testing.T) {
	// os.Setenv("EMAIL_PORT", "123")
	var conf Confing
	err := cleanenv.ReadConfig("./config.yaml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)
}

func TestNewConfig(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	cfg, err := InitConfig("./config2.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
