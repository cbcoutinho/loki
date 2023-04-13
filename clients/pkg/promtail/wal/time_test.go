package wal

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	// allow 10% delta of the minimum testing interval
	delta = time.Millisecond * 25
)

func TestBackoffTimer(t *testing.T) {
	var min = time.Millisecond * 300
	var max = time.Second
	timer := newBackoffTimer(min, max)

	now := time.Now()
	<-timer.C()
	require.WithinDuration(t, now.Add(min), time.Now(), delta, "expected backing off timer to fire in the minimum")

	// backoff, and expect it will take twice the time
	now = time.Now()
	timer.backoff()
	<-timer.C()
	require.WithinDuration(t, now.Add(min*2), time.Now(), delta, "expected backing off timer to fire in the twice the minimum")

	// backoff capped, backoff will actually be 1200ms, but capped at 1000
	now = time.Now()
	timer.backoff()
	<-timer.C()
	require.WithinDuration(t, now.Add(max), time.Now(), delta, "expected backing off timer to fire in the max")
}
