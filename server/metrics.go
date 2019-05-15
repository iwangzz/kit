package main

import (
	"time"
	"github.com/go-kit/kit/metrics"
)

type metricMiddleware struct {
	Service
	RequestCount       	metrics.Counter
	RequestLatency      metrics.Histogram
}

func MetricMiddleware(count metrics.Counter, latency metrics.Histogram) ServiceMiddleware {
	return func(next Service) Service {
		return &metricMiddleware{next, count, latency}
	}
}

func (m *metricMiddleware) Add (a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Add"}
		m.RequestCount.With(lvs...).Add(1)
		m.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.Service.Add(a, b)
}

func (m *metricMiddleware) Subtract (a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Subtract"}
		m.RequestCount.With(lvs...).Add(1)
		m.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.Service.Subtract(a, b)
}

func (m *metricMiddleware) Multiply (a, b int) int {
	defer func(begin time.Time) {
		lvs := []string{"method", "Multiply"}
		m.RequestCount.With(lvs...).Add(1)
		m.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.Service.Multiply(a, b)
}

func (m *metricMiddleware) Divide (a, b int) (int, error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Divide"}
		m.RequestCount.With(lvs...).Add(1)
		m.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.Service.Divide(a, b)
}

func (m *metricMiddleware) HealthCheck() bool {
	defer func(begin time.Time) {
		lvs := []string{"method", "HealthCheck"}
		m.RequestCount.With(lvs...).Add(1)
		m.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return m.Service.HealthCheck()
}
