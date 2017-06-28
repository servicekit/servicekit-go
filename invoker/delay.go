package invoker

import (
	"time"
)

type Delay interface {
	GetDelay() time.Duration
}

type commonDelay struct {
	delay time.Duration
}

func NewDelay(delay time.Duration) Delay {
	return &commonDelay{
		delay: delay,
	}
}

func (c *commonDelay) GetDelay() time.Duration {
	return c.delay
}

type fibDelay struct {
	delay time.Duration
	next  func() time.Duration
}

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

func (f *fibDelay) GetDelay() time.Duration {
	return f.delay * f.next()
}
