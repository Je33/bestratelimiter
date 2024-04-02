package model

import "time"

type Limit struct {
	firstAttempt time.Time
	lastAttempt  time.Time
	count        int
}

func NewLimit() *Limit {
	return &Limit{
		firstAttempt: time.Time{},
		lastAttempt:  time.Time{},
		count:        0,
	}
}

func (l *Limit) Reset() {
	l.firstAttempt = time.Now()
	l.lastAttempt = time.Now()
	l.count = 0
}

func (l *Limit) Increment() {
	l.count++
	l.lastAttempt = time.Now()
	if l.firstAttempt.IsZero() {
		l.firstAttempt = l.lastAttempt
	}
}

func (l *Limit) GetCount() int {
	return l.count
}

func (l *Limit) GetFirstAttempt() time.Time {
	return l.firstAttempt
}

func (l *Limit) GetLastAttempt() time.Time {
	return l.lastAttempt
}

func (l *Limit) SetCount(count int) {
	l.count = count
}

func (l *Limit) SetFirstAttempt(t time.Time) {
	l.firstAttempt = t
}

func (l *Limit) SetLastAttempt(t time.Time) {
	l.lastAttempt = t
}
