package v1

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Project struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type JiraHttpReqField struct {
	Project     `json:",inline"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func TestProbe(t *testing.T) {
	httPProbe := HTTPProbe{}
	d1, _ := json.Marshal(httPProbe)
	fmt.Println(string(d1))
	// dataProject := Project{
	// 	Key:   "name",
	// 	Value: "zhangsan",
	// }
	dd := JiraHttpReqField{
		Project:     Project{},
		Summary:     "asdf",
		Description: "asdf",
	}
	d2, _ := json.Marshal(dd)
	fmt.Println(string(d2))
}
