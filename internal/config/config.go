package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Conf is to store our service configuration

var (
	Cfg *Confing
)

func init() {
	if Cfg == nil {
		Cfg = &Confing{
			Database: DBConfig{
				Host: "127.0.0.1",
			},
		}
	}
}

type Confing struct {
	Database DBConfig    `yaml:"database"`
	Email    EmailConfig `yaml:"email"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DATABASE_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Database string `yaml:"database" env:"DB_DATABASE"`
}

type Scenarios struct {
	PrometheusRemoteWriteRul string `yaml:"k6_prometheus_remote_write_url" env:"K6_PROMETHEUS_RW_SERVER_URL"`
}

type EmailConfig struct {
	From     string   `yaml:"from" env:"EMAIL_FROM"`
	Host     string   `yaml:"host" env:"EMAIL_SMARTHOST"`
	Port     int      `yaml:"port" env:"EMAIL_PORT"`
	Username string   `yaml:"username" env:"EMAIL_USERNAME"`
	Password string   `yaml:"password" env:"EMAIL_PASSWORD"`
	Html     string   `yaml:"html" env:"EMAIL_HTML"`
	Subject  string   `yaml:"subject" env:"EMAIL_SUBJECT"`
	Bcc      string   `yaml:"bcc" env:"EMAIL_BCC"`
	Cc       []string `yaml:"cc" env:"EMAIL_CC"`
}

// // NewConf return new Conf instance from env
// func NewConf() (*Confing, error) {
// 	var conf Confing

// 	err := cleanenv.ReadEnv(&conf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &conf, nil
// }

func InitConfig(path string) error {
	v := viper.New()

	if len(path) > 0 {
		v.SetConfigFile(path)
		err := v.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("ncounter an error in loading configuration file: %w", err))
		}
	}
	v.BindEnv("test.host")
	v.SetEnvPrefix("WM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.Unmarshal(&Cfg); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
	return nil
}
