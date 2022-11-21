package metrcs

import (
	"fmt"
	"testing"

	"github.com/prometheus/common/model"
)

func TestXxx(t *testing.T) {
	ps := PromStore{
		Address: "https://thanos.infra.homepartners.com/",
	}
	ps.init()
	data, err := ps.Query("kube_pod_labels")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	metrics := data.(model.Vector)
	for _, i := range metrics {
		if i.Metric["label_project"] != "" {
			fmt.Println(i.Metric["label_project"])
		}
	}
}
