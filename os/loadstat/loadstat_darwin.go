// Copyright (c) 2015 Square, Inc

// Package loadstat implements metrics collection related to loadavg
package loadstat

import (
	"fmt"
	"github.com/square/inspect/metrics"
	"github.com/square/inspect/os/misc"
	"time"
)

/*
#include <stdlib.h>
*/
import "C"

// LoadStat represents load average metrics for 1/5/15 Minutes of
// current operating system.
// Caution: reflection is used to read this struct to discover names
// Do not add new types
type LoadStat struct {
	OneMinute     *metrics.Gauge
	FiveMinute    *metrics.Gauge
	FifteenMinute *metrics.Gauge
	m             *metrics.MetricContext
}

// New starts metrics collection every Step and registers with
// metricscontext
func New(m *metrics.MetricContext, Step time.Duration) *LoadStat {
	s := new(LoadStat)
	s.m = m
	// initialize all metrics and register them
	misc.InitializeMetrics(s, m, "loadstat", true)
	// collect once
	s.Collect()
	// collect metrics every Step
	ticker := time.NewTicker(Step)
	go func() {
		for _ = range ticker.C {
			s.Collect()
		}
	}()
	return s
}

// Collect populates Loadstat by using sysctl
func (s *LoadStat) Collect() {
	var loadavg [3]C.double

	C.getloadavg(&loadavg[0], 3)
	s.OneMinute.Set(misc.ParseFloat(fmt.Sprintf("%.2f", loadavg[0])))
	s.FiveMinute.Set(misc.ParseFloat(fmt.Sprintf("%.2f", loadavg[1])))
	s.FifteenMinute.Set(misc.ParseFloat(fmt.Sprintf("%.2f", loadavg[2])))
}
