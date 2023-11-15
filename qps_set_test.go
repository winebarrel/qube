package qube_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/qube"
)

func Test_QPSSet_Even(t *testing.T) {
	assert := assert.New(t)

	dps := []qube.DataPoint{
		// 1 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 14, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 15, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 16, 0, time.UTC), Duration: 1 * time.Millisecond},
		// 2 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 17, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 17, 1, time.UTC), Duration: 1 * time.Millisecond},
		// 6 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 1, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 2, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 3, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 4, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 5, time.UTC), Duration: 1 * time.Millisecond},
		// 7 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 1, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 2, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 3, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 4, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 5, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 19, 6, time.UTC), Duration: 1 * time.Millisecond},
	}

	qpsSet := qube.NewQPSSet(dps)
	minQPS, maxQPS, medianQPS := qpsSet.Stats()

	assert.Equal(float64(1), minQPS)
	assert.Equal(float64(7), maxQPS)
	assert.Equal(float64(4), medianQPS)
}

func Test_QPSSet_Odd(t *testing.T) {
	assert := assert.New(t)

	dps := []qube.DataPoint{
		// 1 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 14, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 15, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 16, 0, time.UTC), Duration: 1 * time.Millisecond},
		// 2 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 17, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 17, 1, time.UTC), Duration: 1 * time.Millisecond},
		// 6 qps
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 0, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 1, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 2, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 3, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 4, time.UTC), Duration: 1 * time.Millisecond},
		{Time: time.Date(2023, 10, 11, 12, 13, 18, 5, time.UTC), Duration: 1 * time.Millisecond},
	}

	qpsSet := qube.NewQPSSet(dps)
	minQPS, maxQPS, medianQPS := qpsSet.Stats()

	assert.Equal(float64(1), minQPS)
	assert.Equal(float64(6), maxQPS)
	assert.Equal(float64(1), medianQPS)
}
