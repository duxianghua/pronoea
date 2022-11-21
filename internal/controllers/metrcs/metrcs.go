package metrcs

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	httpapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
)

const (
	CONTAINER_CPU_USAGE_EXPRESS = "ceil(avg_over_time(label_replace(sum(rate(container_cpu_usage_seconds_total{container_name!='POD',container!=''}[5m])) by (container, pod, namespace) *on(pod) group_left(owner_name, owner_kind) kube_pod_owner{owner_kind='ReplicaSet'}, 'deployment', '$1', 'owner_name', '(.*)-(.*)$')[1d:])*1000)"
)

var (
	SERVICE_ALERTS_SUM = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "service_state_score",
		Help: "Current service alerts state"},
		[]string{"service_name", "service_env", "service_platform"},
	)
)

type MetricsStore struct {
	Address string
	API     promv1.API
}

func NewMetricsStore(address string) *MetricsStore {
	store := MetricsStore{
		Address: address}
	store.init()
	return &store
}

func (p *MetricsStore) init() {
	prometheus.MustRegister(SERVICE_ALERTS_SUM)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 30 * time.Second,
	}
	// create config
	config := httpapi.Config{
		Address:      p.Address,
		RoundTripper: tr,
	}
	client, err := httpapi.NewClient(config)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	p.API = promv1.NewAPI(client)
}

func (p *MetricsStore) Query(exp string) (model.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	value, _, err := p.API.Query(ctx, exp, time.Now())
	return value, err
}

func (p *MetricsStore) Register() {
	data, _ := p.Query("count by (label_project, env, platform) (kube_pod_labels{label_project!=\"\"})")
	metrics := data.(model.Vector)
	for _, i := range metrics {
		if i.Metric["label_project"] != "" {
			SERVICE_ALERTS_SUM.With(
				prometheus.Labels{
					"service_name":     string(i.Metric["label_project"]),
					"service_env":      string(i.Metric["env"]),
					"service_platform": string(i.Metric["platform"]),
				},
			).Set(0)
		}
	}
}

func (p *MetricsStore) Score() {
	log.Info().Msg("start score process")
	p.Register()
	data, _ := p.Query("count by (alertname, project, severity, env, platform) (ALERTS{alertstate=\"firing\", project!=\"\"})")
	alerts := data.(model.Vector)
	for _, i := range alerts {
		log.Info().Interface("alert", i)
		if i.Metric["project"] != "" {
			log.Info().Interface("alert", i)
			name := i.Metric["project"]
			severity := i.Metric["severity"]
			v := 0.0
			switch severity {
			case "p1":
				v = 8
			case "p2":
				v = 4
			case "p3":
				v = 2
			case "p4":
				v = 1
			default:
				v = 0
			}
			SERVICE_ALERTS_SUM.With(
				prometheus.Labels{
					"service_name":     string(name),
					"service_env":      string(i.Metric["env"]),
					"service_platform": string(i.Metric["platform"]),
				},
			).Add(v)
		}
	}
}

func (p *MetricsStore) Sync(interval time.Duration, stopCh <-chan struct{}) {
	go func() {
		for {
			log.Info().Msg("start sync process")
			p.Score()
			select {
			case <-stopCh:
				log.Info().Msg("exits sync process")
				return
			case <-time.After(interval):
			}
		}
	}()
}
