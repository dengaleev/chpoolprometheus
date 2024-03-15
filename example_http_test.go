package chpoolprometheus_test

import (
	"context"
	"log"
	"net/http"

	"github.com/ClickHouse/ch-go/chpool"
	"github.com/dengaleev/chpoolprometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExampleCollector() {
	ctx := context.Background()

	pool, err := chpool.Dial(ctx, chpool.Options{})
	if err != nil {
		log.Fatalf("can't init Clickhouse pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping: %v", err)
	}

	reg := prometheus.NewRegistry()
	cllctr := chpoolprometheus.NewCollector(pool, map[string]string{})
	reg.MustRegister(cllctr)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	err = http.ListenAndServe(":4815", nil)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
