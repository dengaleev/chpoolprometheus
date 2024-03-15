// Package chpoolprometheus defines the prometheus Collector for chpool.
package chpoolprometheus

import (
	"time"

	"github.com/ClickHouse/ch-go/chpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Stat defines the chpool.Stat interface.
type Stat interface {
	AcquireCount() int64
	AcquireDuration() time.Duration
	AcquiredResources() int32
	CanceledAcquireCount() int64
	ConstructingResources() int32
	EmptyAcquireCount() int64
	IdleResources() int32
	MaxResources() int32
	TotalResources() int32
}

type statFunc func() Stat

// Collector is implements prometheus.Collector interface that collect metrics produced by chpool.
type Collector struct {
	statFn statFunc

	acquires         *prometheus.Desc
	acquiresDuration *prometheus.Desc
	canceledAcquires *prometheus.Desc
	emptyAcquires    *prometheus.Desc

	acquiredConnections     *prometheus.Desc
	constructingConnections *prometheus.Desc
	idleConnections         *prometheus.Desc
	totalConnections        *prometheus.Desc
	maxConnections          *prometheus.Desc
}

// NewCollector creates a new Collector for the pool.
func NewCollector(pool *chpool.Pool, labels prometheus.Labels) *Collector {
	buildFQName := func(name string) string {
		return prometheus.BuildFQName("ch", "pool", name)
	}

	return &Collector{
		statFn: func() Stat {
			return pool.Stat()
		},

		acquires: prometheus.NewDesc(
			buildFQName("acquires_total"),
			"Cumulative count of successful acquires from the pool.",
			nil,
			labels,
		),
		acquiresDuration: prometheus.NewDesc(
			buildFQName("acquire_duration_nanoseconds"),
			"Total duration of all successful acquires from the pool.",
			nil,
			labels,
		),
		canceledAcquires: prometheus.NewDesc(
			buildFQName("canceled_acquires_total"),
			"Cumulative count of acquires from the pool that were canceled by a context.",
			nil,
			labels,
		),
		emptyAcquires: prometheus.NewDesc(
			buildFQName("empty_acquires_total"),
			"Cumulative count of successful acquires from the pool that waited for a connection to be released or constructed because the pool was empty.",
			nil,
			labels,
		),

		acquiredConnections: prometheus.NewDesc(
			buildFQName("acquired_connections"),
			"The number of currently acquired connections in the pool.",
			nil,
			labels,
		),
		constructingConnections: prometheus.NewDesc(
			buildFQName("constructing_connections"),
			"The number of connections with construction in progress in the pool.",
			nil,
			labels,
		),
		idleConnections: prometheus.NewDesc(
			buildFQName("idle_connections"),
			"The number of currently idle connections in the pool.",
			nil,
			labels,
		),
		totalConnections: prometheus.NewDesc(
			buildFQName("total_connections"),
			"Total number of connections currently in the pool. The value is the sum of constructing, acquired, and idle connections.",
			nil,
			labels,
		),
		maxConnections: prometheus.NewDesc(
			buildFQName("max_connections"),
			"The maximum size of the pool.",
			nil,
			labels,
		),
	}
}

// Describe implements the prometheus.Collector.Describe method.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements the prometheus.Collector.Collect method.
func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	stat := c.statFn()

	metrics <- prometheus.MustNewConstMetric(
		c.acquires,
		prometheus.CounterValue,
		float64(stat.AcquireCount()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.acquiresDuration,
		prometheus.CounterValue,
		float64(stat.AcquireDuration().Nanoseconds()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.canceledAcquires,
		prometheus.CounterValue,
		float64(stat.CanceledAcquireCount()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.emptyAcquires,
		prometheus.CounterValue,
		float64(stat.EmptyAcquireCount()),
	)

	metrics <- prometheus.MustNewConstMetric(
		c.acquiredConnections,
		prometheus.GaugeValue,
		float64(stat.AcquiredResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.constructingConnections,
		prometheus.GaugeValue,
		float64(stat.ConstructingResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.idleConnections,
		prometheus.GaugeValue,
		float64(stat.IdleResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.totalConnections,
		prometheus.GaugeValue,
		float64(stat.TotalResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		c.maxConnections,
		prometheus.GaugeValue,
		float64(stat.MaxResources()),
	)
}
