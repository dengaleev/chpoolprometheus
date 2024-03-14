package chpoolprometheus_test

import (
	"bytes"
	"testing"
	"time"

	_ "embed"

	"github.com/dengaleev/chpoolprometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInterface(*testing.T) {
	var _ prometheus.Collector = (*chpoolprometheus.Collector)(nil)
}

type exampleStat struct{}

func (exampleStat) AcquireCount() int64            { return 4 }
func (exampleStat) AcquireDuration() time.Duration { return 8 }
func (exampleStat) AcquiredResources() int32       { return 15 }
func (exampleStat) CanceledAcquireCount() int64    { return 16 }
func (exampleStat) ConstructingResources() int32   { return 23 }
func (exampleStat) EmptyAcquireCount() int64       { return 42 }
func (exampleStat) IdleResources() int32           { return 24 }
func (exampleStat) MaxResources() int32            { return 32 }
func (exampleStat) TotalResources() int32          { return 61 }

type exampleStater struct{}

func (exampleStater) Stat() chpoolprometheus.Stat { return exampleStat{} }

//go:embed testdata/test_collector.txt
var testCollectorExpected string

func TestCollector(t *testing.T) {
	collector := chpoolprometheus.NewCollector(exampleStater{}, map[string]string{
		"host": "example.com",
	})

	reg := prometheus.NewRegistry()
	reg.MustRegister(collector)

	assert.NoError(t, testutil.GatherAndCompare(reg, bytes.NewBufferString(testCollectorExpected)))
}
