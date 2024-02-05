package panics

import (
	"container/list"
	"fmt"
	"sync"
)

func New() *Panics {
	return &Panics{
		list: list.New(),
	}
}

type Panics struct {
	lock sync.Mutex
	list *list.List
}

type PanicInfo struct {
	Reason interface{}
	Stack  []byte
}

func (p *Panics) Add(reason interface{}, stack []byte) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.list.PushBack(&PanicInfo{reason, stack})
}

func (p *Panics) Has() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.list.Len() > 0
}

func (p *Panics) Info() string {
	p.lock.Lock()
	defer p.lock.Unlock()

	info := fmt.Sprintf("%d panics occurred\n\n", p.list.Len())
	cur := p.list.Front()
	for i := 1; cur != nil; i++ {
		v := cur.Value.(*PanicInfo)
		info += fmt.Sprintf("NO.%d: %v\n%s\n", i, v.Reason, v.Stack)
		cur = cur.Next()
	}
	return info
}
