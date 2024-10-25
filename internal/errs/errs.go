package errs

import (
	"errors"
	"fmt"
	"sync"
)

func New() *Errs {
	return &Errs{}
}

type ErrorChain struct {
	cause *ErrorChain
	Err   error
	From  string
}

func (e *ErrorChain) Is(target error) bool {
	return errors.Is(e.Err, target)
}

func (e *ErrorChain) As(target any) bool {
	return errors.As(e.Err, target)
}

func (e *ErrorChain) Unwrap() error {
	if e.cause == nil {
		// 必需有这个，不能直接返回e.cause
		// 因为nil和=nil的*ErrorChain是不同的，后者虽然是空指针但是实现了Unwrap方法，导致errors.Is和As无法正常结束
		return nil
	}
	return e.cause
}

func (e *ErrorChain) Error() string {
	return e.Err.Error()
}

type Errs struct {
	lock sync.RWMutex
	err  *ErrorChain
}

func (e *Errs) Add(err error, from string) {
	if err == nil {
		return
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.err = &ErrorChain{Err: err, From: from, cause: e.err}
}

func (e *Errs) Unwrap() error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.err == nil {
		return nil
	}
	return e.err
}

func (e *Errs) Has() bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.err != nil
}

func (e *Errs) Error() string {
	e.lock.Lock()
	defer e.lock.Unlock()

	var info string
	cnt := 0
	for cur := e.err; cur != nil; cur = cur.cause {
		cnt++
		info += fmt.Sprintf("\n\nNO.%d: %s\n\thandler added from %s", cnt, cur.Err, cur.From)
	}
	info = fmt.Sprintf("%d errors occurred", cnt) + info

	return info
}
