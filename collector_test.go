package chpoolprometheus_test

import (
	"testing"

	"github.com/dengaleev/chpoolprometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func TestInterface(*testing.T) {
	var _ prometheus.Collector = (*chpoolprometheus.Collector)(nil)
}
