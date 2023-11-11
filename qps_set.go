package qube

import (
	"sort"
	"time"
)

type QPSSet []float64

func NewQPSSet(dps []DataPoint) QPSSet {
	if len(dps) == 0 {
		return nil
	}

	sort.Slice(dps, func(i, j int) bool {
		return dps[i].Time.Before(dps[j].Time)
	})

	baseTime := dps[0].Time
	countSet := []int{0}

	// Calculate number of queries per second
	for _, v := range dps {
		if baseTime.Add(1 * time.Second).Before(v.Time) {
			baseTime = baseTime.Add(1 * time.Second)
			countSet = append(countSet, 0)
		}

		countSet[len(countSet)-1]++
	}

	qpsSet := make([]float64, len(countSet))

	// Convert "countSet" to float64 array
	for i, v := range countSet {
		qpsSet[i] = float64(v)
	}

	return qpsSet
}

func (qpsSet QPSSet) Stats() (minQPS float64, maxQPS float64, medianQPS float64) {
	sort.Slice(qpsSet, func(i, j int) bool {
		return qpsSet[i] < qpsSet[j]
	})

	minQPS = qpsSet[0]
	maxQPS = qpsSet[len(qpsSet)-1]

	median := len(qpsSet) / 2
	medianNext := median + 1

	if len(qpsSet) == 1 {
		medianQPS = qpsSet[0]
	} else if len(qpsSet) == 2 {
		medianQPS = (qpsSet[0] + qpsSet[1]) / 2
	} else if len(qpsSet)%2 == 0 {
		medianQPS = (qpsSet[median] + qpsSet[medianNext]) / 2
	} else {
		medianQPS = qpsSet[medianNext]
	}

	return
}
