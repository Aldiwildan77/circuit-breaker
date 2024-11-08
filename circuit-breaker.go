package main

import (
	"errors"
	"sync"
	"time"
)

const (
	DefaultMaxFailures  = 3
	DefaultResetTimeout = 5 * time.Second

	EmptyFailure = 0
)

var (
	ErrCircuitOpen = errors.New("circuit is open")
)

type circuitBreaker struct {
	failures     int
	maxFailures  int
	isOpen       bool
	lock         sync.Mutex
	resetTimeout time.Duration
	lastAttempt  time.Time
}

type CircuitBreaker interface {
	IsAllowed() bool
	Call(f func() error) error
	GetFailures() int
	SetResetTimeout(resetTimeout time.Duration)
}

type Option func(*circuitBreaker)

func WithMaxFailures(maxFailures int) Option {
	return func(cb *circuitBreaker) {
		cb.maxFailures = maxFailures
	}
}

func WithResetTimeout(resetTimeout time.Duration) Option {
	return func(cb *circuitBreaker) {
		cb.resetTimeout = resetTimeout
	}
}

func NewCircuitBreaker(options ...Option) CircuitBreaker {
	cb := &circuitBreaker{
		maxFailures:  DefaultMaxFailures,
		resetTimeout: DefaultResetTimeout,
	}

	for _, option := range options {
		option(cb)
	}

	return cb
}

func (cb *circuitBreaker) IsAllowed() bool {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	if !cb.isOpen {
		return true
	}

	if time.Since(cb.lastAttempt) > cb.resetTimeout {
		cb.isOpen = false
		cb.failures = EmptyFailure
		return true
	}

	return false
}

func (cb *circuitBreaker) Call(f func() error) error {
	if !cb.IsAllowed() {
		return ErrCircuitOpen
	}

	err := f()

	cb.lock.Lock()
	defer cb.lock.Unlock()

	if err != nil {
		cb.failures++
		if cb.failures >= cb.maxFailures {
			cb.isOpen = true
			cb.lastAttempt = time.Now()
		}
	}

	return err
}

func (cb *circuitBreaker) GetFailures() int {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	return cb.failures
}

func (cb *circuitBreaker) SetResetTimeout(resetTimeout time.Duration) {
	cb.lock.Lock()
	defer cb.lock.Unlock()

	cb.resetTimeout = resetTimeout
}
