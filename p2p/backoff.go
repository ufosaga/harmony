package p2p

import "time"

// Backoff is a backoff sleeper.
type Backoff interface {
	Reset()
	Sleep()
	Backoff()
}

// BackoffBase is the base implementation of backoff sleeper.  It handles the
// minimum/current/maximum backoff logic.  If maximum backoff is 0, there is no
// maximum.
type BackoffBase struct {
	Min, Cur, Max time.Duration
}

func NewBackoffBase(min, max time.Duration) *BackoffBase{
	return &BackoffBase{min, min, max}
}

// Reset the current sleep duration to its minimum value.
func (b *BackoffBase) Reset() {
	b.Cur = b.Min
}

// Sleep for the current duration, then adjust the duration.
func (b *BackoffBase) Sleep() {
	if b.Cur < b.Min {
		b.Cur = b.Min
	}
	time.Sleep(b.Cur)
	b.Backoff()
	if b.Max != 0 && b.Cur > b.Max {
		b.Cur = b.Max
	}
	if b.Cur < b.Min {
		b.Cur = b.Min
	}
}

// Adjust the duration.  Subtypes shall implement this.
func (b *BackoffBase) Backoff() {
	// default implementation does not backoff
}

// Exponential backoff.
type ExpBackoff struct {
	BackoffBase
	Factor float64
}

func NewExpBackoff(min, max time.Duration, factor float64) *ExpBackoff {
	return &ExpBackoff{*NewBackoffBase(min, max), factor}
}

func (b *ExpBackoff) Backoff() {
	b.Cur = time.Duration(float64(b.Cur) * b.Factor)
}
