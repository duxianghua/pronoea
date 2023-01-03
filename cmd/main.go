package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/duxianghua/pronoea/internal"
	"github.com/duxianghua/pronoea/internal/config"
	"github.com/duxianghua/pronoea/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	APP_NAME = "watchman"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	pflag.Bool("debug", false, "sets log level to debug")
	pflag.String("config", "./config.yaml", "config file path")
	pflag.String("web.addr", "0.0.0.0:8081", "web service listen address")
	pflag.String("init", "false", "disable init")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("log level is debug")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetConfigFile(viper.GetString("config"))
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("ncounter an error in loading configuration file: %w", err))
	}
	viper.SetEnvPrefix("WM")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&config.Options); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
	log.Info().Interface("Options", config.Options).Msg("current profile")

	// if viper.GetBool("init") {
	// 	stopCh := make(chan struct{})
	// 	defer close(stopCh)
	// 	prom := metrcs.NewMetricsStore("https://thanos.infra.homepartners.com/")
	// 	//prom.Register()
	// 	prom.Sync(time.Minute, stopCh)

	log.Info().Msg("starting controllers manager")
	controllers.StartMgr()

	log.Info().Msg("starting web server")
	internal.Service().Run(viper.GetString("web.addr"))
}
