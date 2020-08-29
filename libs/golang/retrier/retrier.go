package retrier

import "time"

// Retrier retries attempts a function until it succeeds or until it's
// performed a set number of retries.
type Retrier struct {
	f       func() error
	retries uint
}

// New instantiates a new Retrier.
func New(retries uint, f func() error) *Retrier {
	return &Retrier{f: f, retries: retries}
}

// Start begins the retry schedule.
func (r *Retrier) Start() error {
	var err error
	for i := 0; i < int(r.retries); i++ {
		err = r.f()
		if err == nil {
			return nil
		}
		if i == (int(r.retries) - 1) {
			return err
		}
		time.Sleep(time.Second << uint(i))
	}
	return err
}
