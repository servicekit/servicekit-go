package invoker

import (
	"time"
)

// Delay defin a interface that has a method: GetDelay returns delay seconds
type Delay interface {
	GetDelay() time.Duration
}

// commonDelay is an implementation base on Delay
type commonDelay struct {
	delay time.Duration
}

// NewDelay returns a commonDelay
func NewDelay(delay time.Duration) Delay {
	return &commonDelay{
		delay: delay,
	}
}

// GetDelay returns delay seconds
func (c *commonDelay) GetDelay() time.Duration {
	return c.delay
}

// fibDelay is an implementation base on Delay
// delay seconds will be increased by Fibonacci algorithm
type fibDelay struct {
	delay time.Duration
	next  func() time.Duration
}

// NewFibDelay return a fibDelay
func NewFibDelay(delay time.Duration) Delay {
	first, second := 1, 2

	return &fibDelay{
		delay: delay,
		next: func() time.Duration {
			ret := first
			first, second = second, first+second
			return time.Duration(ret)
		},
	}
}

// GetDelay returns delay seconds
func (f *fibDelay) GetDelay() time.Duration {
	return f.delay * f.next()
}
