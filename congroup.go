package congroup

import (
	"sync"

	"github.com/pkg/errors"
)

// ConGroup excecutes functions concurrently, waits them done by sync.WaitGroup,
// return the first error occured, and is safe with panic.
type ConGroup struct {
	sync.WaitGroup
	errOnce sync.Once
	err     error
}

// New get a group without context control.
func New() *ConGroup {
	return &ConGroup{}
}

// Add run the f inputed concurrently.
func (g *ConGroup) Add(f func() error) {
	g.WaitGroup.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				g.setErrorOnce(errors.Errorf("%v", err))
			}
		}()
		defer g.WaitGroup.Done()
		err := f()
		if err != nil {
			g.setErrorOnce(err)
		}
	}()
}

// Wait is blocked until all functions from Add() done and returns a non-nil value if any error or panic occured.
func (g *ConGroup) Wait() error {
	g.WaitGroup.Wait()
	return g.err
}

func (g *ConGroup) setErrorOnce(err error) {
	g.errOnce.Do(func() {
		g.err = err
	})
}
