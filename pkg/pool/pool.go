package pool

import "sync"

// Resetter defines the interface for types that can reset their state
type Resetter interface {
	Reset()
}

// Pool is a generic object pool for types implementing Resetter
type Pool[T Resetter] struct {
	pool sync.Pool
	new  func() T
}

// New creates a new Pool with the provided object constructor
func New[T Resetter](newFunc func() T) *Pool[T] {
	return &Pool[T]{
		new: newFunc,
		pool: sync.Pool{
			New: func() any { return newFunc() },
		},
	}
}

// Get returns an object from the pool
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns the object to the pool after resetting its state
func (p *Pool[T]) Put(x T) {
	x.Reset()
	p.pool.Put(x)
}
