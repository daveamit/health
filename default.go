package health

import (
	"net/http"
	"strings"
)

var defaultHealth = newHealth()

func EnsureService(name string, namespace string) {
	name = safeNameReplacer.Replace(name)
	namespace = safeNameReplacer.Replace(namespace)
	defaultHealth.EnsureService(name, namespace)
}
func ServiceUp(name string, namespace string) {
	name = safeNameReplacer.Replace(name)
	namespace = safeNameReplacer.Replace(namespace)
	defaultHealth.ServiceUp(name, namespace)
}
func ServiceDown(name string, namespace string) {
	name = safeNameReplacer.Replace(name)
	namespace = safeNameReplacer.Replace(namespace)
	defaultHealth.ServiceDown(name, namespace)
}
func PrometheusScrapHandler() http.Handler {
	return defaultHealth.PrometheusScrapHandler()
}
func HealthCheckHandler() http.Handler {
	return defaultHealth.HealthCheckHandler()
}

func SetSafeNameReplacer(r Replacer) {
	// nul replacer essentially means NOP replacer
	if r == nil {
		r = NewCharacterToUnderscoreReplacer()
	}
	safeNameReplacer = r
}
func GetSafeNameReplacer() Replacer {
	return safeNameReplacer
}

type Replacer interface {
	Replace(string) string
}

var safeNameReplacer = NewCharacterToUnderscoreReplacer()

// NewCharacterToUnderscoreReplacer will retun a replacer which will replace given chars to _
func NewCharacterToUnderscoreReplacer(c ...string) Replacer {
	if len(c) > 0 {
		chars := make([]string, len(c)*2)

		for i := range c {
			chars[i] = c[i]
			chars[i+1] = "_"
		}
		return strings.NewReplacer(chars...)
	}
	return strings.NewReplacer()
}

func ClearItems() {
	defaultHealth.ClearItems()
}
