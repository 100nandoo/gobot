package reddit

import "sync"

type singleflightCall[T any] struct {
	done chan struct{}
	val  T
	err  error
}

type singleflight[T any] struct {
	fn      func() (T, error)
	mu      sync.Mutex
	current *singleflightCall[T]
}

func (s *singleflight[T]) do() (T, error) {
	s.mu.Lock()
	if s.current != nil {
		c := s.current
		s.mu.Unlock()
		<-c.done
		return c.val, c.err
	}

	c := &singleflightCall[T]{done: make(chan struct{})}
	s.current = c
	s.mu.Unlock()

	c.val, c.err = s.fn()

	s.mu.Lock()
	s.current = nil
	s.mu.Unlock()
	close(c.done)

	return c.val, c.err
}

func newSingleflight[T any](fn func() (T, error)) func() (T, error) {
	sf := &singleflight[T]{fn: fn}
	return sf.do
}
