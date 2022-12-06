package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
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

func TestHttpProbe(t *testing.T) {
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

func TestProbe(t *testing.T) {
	// httPProbe := Probe{
	// 	Spec: ProbeSpec{
	// 		Module: Module{
	// 			Timeout: Duration(1000),
	// 		},
	// 	},
	// }
	// d, _ := yaml.Marshal(httPProbe)
	// //d, _ := json.Marshal(httPProbe)
	// fmt.Println(string(d))
	str := "{\"metadata\":{\"creationTimestamp\":null},\"spec\":{\"targets\":null,\"module\":{\"prober\":\"\",\"timeout\":\"10s\",\"http\":{\"basic_auth\":{\"username\":\"\"},\"authorization\":{},\"oauth2\":{\"client_id\":\"\",\"client_secret\":\"\",\"client_secret_file\":\"\",\"token_url\":\"\",\"tls_config\":{}},\"tls_config\":{}},\"tcp\":{\"tls_config\":{}},\"icmp\":{},\"dns\":{\"tls_config\":{},\"validate_answer_rrs\":{},\"validate_authority_rrs\":{},\"validate_additional_rrs\":{}},\"grpc\":{\"tls_config\":{}}},\"pause\":false},\"status\":{}}"
	httPProbe2 := Probe{}
	err := json.Unmarshal([]byte(str), &httPProbe2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(httPProbe2.Spec.Module.Timeout)
	httPProbe2.Spec.Module.Timeout = httPProbe2.Spec.Module.Timeout * 2
	d2, _ := json.Marshal(&httPProbe2)
	fmt.Println(string(d2))
	d3, _ := yaml.Marshal(&httPProbe2)
	fmt.Println(string(d3))
	fmt.Println(httPProbe2)
	// httPProbe3 := Probe{}
	// _ = yaml.Unmarshal(d3, &httPProbe3)
	// fmt.Println(string(d3))
}

func TestDuration(t *testing.T) {
	d := Duration(time.Second * 5)
	jsonStr, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsonStr)
	fmt.Println(string(jsonStr))

	yamlStr, err := yaml.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(yamlStr))

	yamlStr2, err := yaml.Marshal(time.Duration(200))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(yamlStr2))

	var d1 Duration
	b := []byte(`10s`)
	yaml.Unmarshal(b, &d1)
	fmt.Println(d1)
}
