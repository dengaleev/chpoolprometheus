package chpoolprometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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

type Stater interface {
	Stat() Stat
}

type Collector struct {
	stater Stater

	acquires         prometheus.Desc
	acquiresDuration prometheus.Desc
	canceledAcquires prometheus.Desc
	emptyAcquires    prometheus.Desc

	acquiredConnections     prometheus.Desc
	constructingConnections prometheus.Desc
	idleConnections         prometheus.Desc
	totalConnections        prometheus.Desc
	maxConnections          prometheus.Desc
}

func NewCollector(stater Stater, labels prometheus.Labels) *Collector {
	return &Collector{
		stater: stater,

		acquires: *prometheus.NewDesc(
			"acquires_total",
			"Cumulative count of successful acquires from the pool.",
			nil,
			labels,
		),
		acquiresDuration: *prometheus.NewDesc(
			"acquire_duration_nanoseconds",
			"Total duration of all successful acquires from the pool.",
			nil,
			labels,
		),
		canceledAcquires: *prometheus.NewDesc(
			"canceled_acquires_total",
			"Cumulative count of acquires from the pool that were canceled by a context.",
			nil,
			labels,
		),
		emptyAcquires: *prometheus.NewDesc(
			"empty_acquires_total",
			"Cumulative count of successful acquires from the pool that waited for a resource to be released or constructed because the pool was empty.",
			nil,
			labels,
		),

		acquiredConnections: *prometheus.NewDesc(
			"acquired_connections",
			"The number of currently acquired resources in the pool.",
			nil,
			labels,
		),
		constructingConnections: *prometheus.NewDesc(
			"constructing_connections",
			"The number of resources with construction in progress in the pool.",
			nil,
			labels,
		),
		idleConnections: *prometheus.NewDesc(
			"idle_connections",
			"The number of currently idle resources in the pool.",
			nil,
			labels,
		),
		totalConnections: *prometheus.NewDesc(
			"total_connections",
			"Total number of resources currently in the pool. The value is the sum of ConstructingResources, AcquiredResources, and IdleResources.",
			nil,
			labels,
		),
		maxConnections: *prometheus.NewDesc(
			"max_connections",
			"The maximum size of the pool.",
			nil,
			labels,
		),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	stat := c.stater.Stat()

	metrics <- prometheus.MustNewConstMetric(
		&c.acquires,
		prometheus.CounterValue,
		float64(stat.AcquireCount()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.acquiresDuration,
		prometheus.CounterValue,
		float64(stat.AcquireDuration().Nanoseconds()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.canceledAcquires,
		prometheus.CounterValue,
		float64(stat.CanceledAcquireCount()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.emptyAcquires,
		prometheus.CounterValue,
		float64(stat.EmptyAcquireCount()),
	)

	metrics <- prometheus.MustNewConstMetric(
		&c.acquiredConnections,
		prometheus.GaugeValue,
		float64(stat.AcquiredResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.constructingConnections,
		prometheus.GaugeValue,
		float64(stat.ConstructingResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.idleConnections,
		prometheus.GaugeValue,
		float64(stat.IdleResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.totalConnections,
		prometheus.GaugeValue,
		float64(stat.TotalResources()),
	)
	metrics <- prometheus.MustNewConstMetric(
		&c.maxConnections,
		prometheus.GaugeValue,
		float64(stat.MaxResources()),
	)
}
