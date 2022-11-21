package models

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg/v10"
)

func TestServiceInfo(t *testing.T) {
	db := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "grafana",
		Password: "rB6imVLSGWJCyrB",
		Database: "grafana",
	})
	defer db.Close()

	// select all
	var serviceInfoList []ServiceInfo
	err := db.Model(&serviceInfoList).Where("service = 'hpa-desktop-tool'").Select()
	if err != nil {
		panic(err)
	}
	fmt.Println(serviceInfoList)
}
