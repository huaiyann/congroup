package errs

import (
	"container/list"
	"fmt"
	"sync"
)

func New() *Errs {
	return &Errs{
		list: list.New(),
	}
}

type Errs struct {
	lock sync.RWMutex
	list *list.List
}

type ErrorInfo struct {
	Err  error
	From string
}

func (e *Errs) Add(err error, from string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.list.PushBack(&ErrorInfo{err, from})
}

func (e *Errs) Has() bool {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.list.Len() > 0
}

func (e *Errs) Error() string {
	e.lock.Lock()
	defer e.lock.Unlock()

	info := fmt.Sprintf("%d errors occurred", e.list.Len())
	cur := e.list.Front()
	for i := 1; cur != nil; i++ {
		v := cur.Value.(*ErrorInfo)
		info += fmt.Sprintf("\n\nNO.%d: %s\n\thandler added from %s", i, v.Err, v.From)
		cur = cur.Next()
	}
	return info
}
