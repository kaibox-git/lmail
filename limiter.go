package lmail

import (
	"fmt"
	"time"
)

type Limiter struct {
	queue []time.Time
	max   int
	dif   time.Duration
}

func (l *Limiter) Add() {
	l.queue = append(l.queue, time.Now())
	if len(l.queue) > l.max {
		l.queue = l.queue[1 : l.max+1]
	}
}
func (l *Limiter) IsOn() bool {
	ln := len(l.queue)
	if ln >= l.max {
		return l.queue[ln-1].Sub(l.queue[0]) < l.dif && time.Since(l.queue[0]) < l.dif
	}
	return false
}

func (l *Limiter) Check() {
	for i := range l.queue {
		fmt.Printf("%d. %v\n", i, l.queue[i])
	}
	fmt.Println("")
}

func NewLimiter(max int, dif time.Duration) *Limiter {
	return &Limiter{
		queue: make([]time.Time, 0, 2),
		max:   max,
		dif:   dif,
	}
}
