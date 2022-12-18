package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/duxianghua/pronoea/internal/config"
)

func TestInitConfig(t *testing.T) {
	//os.Setenv("WM_DATABASE_HOST", "asdf")
	os.Setenv("WM_TEST_HOST", "WM_TEST_HOST")
	config.InitConfig("config.yaml")
	//config.InitConfig("")
	fmt.Println(config.Cfg)
}
