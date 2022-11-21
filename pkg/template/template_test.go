package template

import (
	// "html/template"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
)

// func TestTemplate(t *testing.T) {
// 	type Inventory struct {
// 		Name  string
// 		Count uint
// 	}
// 	sweaters := Inventory{"wool", 17}
// 	tmpl, err := template.New("test").Parse("Hello {{ .Name }}")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = tmpl.Execute(os.Stdout, sweaters)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func TestLoadTemplate(t *testing.T) {
	tmpl, _ := template.FromGlobs("../../templates/email.tmpl")
	jsonFile, err := os.Open("./alerts.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var message webhook.Message
	json.Unmarshal(byteValue, &message)
	fmt.Println(message.Data)
	if err != nil {
		panic(err)
	}
	fmt.Println(tmpl)
	str, err := tmpl.ExecuteHTMLString(
		`{{ template "email.default.html" . }}`, message.Data)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
