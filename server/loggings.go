package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	Service
	log 		log.Logger
}

func LoggingMiddleware(log log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &loggingMiddleware{next, log}
	}
}

func (l *loggingMiddleware) Add (a, b int) int {
	defer func(begin time.Time) {
		l.log.Log(
			"func", "Add",
			"a", a,
			"b", b,
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.Service.Add(a, b)
}

func (l *loggingMiddleware) Subtract (a, b int) int {
	defer func(begin time.Time) {
		l.log.Log(
			"func", "Subtract",
			"a", a,
			"b", b,
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.Service.Subtract(a, b)
}

func (l *loggingMiddleware) Multiply (a, b int) int {
	defer func(begin time.Time) {
		l.log.Log(
			"func", "Multiply",
			"a", a,
			"b", b,
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.Service.Multiply(a, b)
}

func (l *loggingMiddleware) Divide (a, b int) (int, error) {
	defer func(begin time.Time) {
		l.log.Log(
			"func", "Divide",
			"a", a,
			"b", b,
			"took", time.Since(begin),
		)
	}(time.Now())

	return l.Service.Divide(a, b)
}

func (l *loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		l.log.Log(
			"func", "HealthCheck",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = l.Service.HealthCheck()

	return result
}

