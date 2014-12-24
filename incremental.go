// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0

package vlc

import (
	"sync"
)

// IncrementalInt is a concurency-safe incremental integer provider
type IncrementalInt struct {
	increment int
	lock      sync.Mutex
}

// Next returns with an integer that is exactly one higher as the previous call to Next() for this IncrementalInt
func (i *IncrementalInt) Next() int {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.increment++
	return i.increment
}

// Last returns the number (int) that was returned by the most recent call to this instance's Next()
func (i *IncrementalInt) Last() int {
	return i.increment
}

// Set changes the increment to given value, the succeeding call to Next() will return the given value+1
func (i *IncrementalInt) Set(value int) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.increment = value
}
