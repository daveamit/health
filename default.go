package health

import "net/http"

var defaultHealth = newHealth()

func EnsureService(name string, namespace string) {
	defaultHealth.EnsureService(name, namespace)
}
func ServiceUp(name string, namespace string) {
	defaultHealth.ServiceUp(name, namespace)
}
func ServiceDown(name string, namespace string) {
	defaultHealth.ServiceDown(name, namespace)
}
func PrometheusScrapHandler() http.Handler {
	return defaultHealth.PrometheusScrapHandler()
}
func HealthCheckHandler() http.Handler {
	return defaultHealth.HealthCheckHandler()
}
