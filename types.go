package health

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type health interface {
	EnsureService(name string, namespace string)
	ServiceUp(name string, namespace string)
	ServiceDown(name string, namespace string)
	PrometheusScrapHandler() http.Handler
	HealthCheckHandler() http.Handler
	ClearItems()
}

func newHealth() health {
	return &healthImpl{}
}

type serviceState byte

const (
	UndefinedServiceState serviceState = iota
	RunningServiceState
	StoppedServiceState
)

func (s serviceState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
func (s serviceState) String() string {
	switch s {
	case RunningServiceState:
		return "Up"
	case StoppedServiceState:
		return "Down"
	default:
		return "undefined"
	}
}

type service struct {
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"`
	State     serviceState     `json:"state"`
	Gauge     prometheus.Gauge `json:"-"`
}

type healthImpl struct {
	items []service
}

func (h *healthImpl) setServiceState(name string, namespace string, state serviceState) {
	for index := range h.items {
		if h.items[index].Name == name && h.items[index].Namespace == namespace {
			h.items[index].State = state
			switch state {
			case RunningServiceState:
				h.items[index].Gauge.Set(1)
			default:
				h.items[index].Gauge.Set(0)
			}
			return
		}
	}
}

func (h *healthImpl) ServiceUp(name string, namespace string) {
	h.setServiceState(name, namespace, RunningServiceState)
}
func (h *healthImpl) ServiceDown(name string, namespace string) {
	h.setServiceState(name, namespace, StoppedServiceState)
}

func (h *healthImpl) EnsureService(name string, namespace string) {
	for _, i := range h.items {
		if i.Name == name && i.Namespace == namespace {
			return
		}
	}
	gague := promauto.NewGauge(prometheus.GaugeOpts{Namespace: namespace, Name: name, Help: name + " up status"})
	h.items = append(h.items, service{Name: name, Namespace: namespace, Gauge: gague})
}

func (h *healthImpl) PrometheusScrapHandler() http.Handler {

	return promhttp.Handler()
}
func (h *healthImpl) HealthCheckHandler() http.Handler {
	return h
}

func (h *healthImpl) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(h.items)
	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rsp, "Failed to marshal service status: ", err)
		return
	}

	running := true
	for _, i := range h.items {
		if i.State != RunningServiceState {
			running = false
			rsp.WriteHeader(http.StatusInternalServerError)
			break
		}
	}

	if running {
		rsp.WriteHeader(http.StatusOK)
	}

	rsp.Write(bytes)
}

func (h *healthImpl) ClearItems() {
	for _, h := range h.items {
		prometheus.Unregister(h.Gauge)
	}
	h.items = make([]service, 0)
}
