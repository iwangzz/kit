package main

import "errors"

type Service interface {
	Add(a, b int) int
	Subtract(a, b int) int
	Multiply(a, b int) int
	Divide(a, b int) (int, error)
	HealthCheck() bool
}

type CalcService struct {}

func(c *CalcService) Add(a, b int) int {
	return a + b
}

func (c *CalcService) Subtract(a, b int) int {
	return a - b
}

func (c *CalcService) Multiply(a, b int) int {
	return a * b
}

func (c *CalcService) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("cant be zero")
	}

	return a / b, nil
}

func (c *CalcService) HealthCheck() bool {
	return true
}

//middleware
type ServiceMiddleware func(Service) Service