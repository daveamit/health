![Test](https://github.com/daveamit/health/actions/workflows/test.yml/badge.svg)

# Background
This package helps setup health check based on status of external dependencies. The idea is to add all external dependencies like database, queue connections etc, and based on their status, the health handler will return `500` or `200` response with details about all the dependencies.

# Installation

```shell
go get github.com/daveamit/health@latest
```

# Limitation
Current implementation maintains everything in default instance of `healthImpl` as this package is intended to be used as singleton instance. But if there are use cases where I need to expose API such that `health` interface and `newHealth` are make public and provide added api to setup `registry` and do various customization for prometheus

# Usage

```golang
package main

import "github.com/daveamit/health"

func main() {
    // Register a dep
    health.EnsureService("posgtres", "default")


    // postgresConn is a hypothetical function, assume that it return
    // a valid pg connection on sucess and err on failure.
    pg, err := postgresConn()

    if err != nil {
        // tell the package that `postgres` service is connected and usable.
        health.ServiceUp("postgres", "default")
    } else {
        health.ServiceDown("postgres", "default")
    }

    // so something with the pg connection
    ...
    ...
    ...

    // serve prom and health endpoints @ 8000 port
    http.Handle("/prom", health.PrometheusScrapHandler())
    http.Handle("/health", health.HealthCheckHandler())
    http.ListenAndServe(":8000", nil)
}
```
