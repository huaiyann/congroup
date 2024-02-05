package channel

import (
	"container/list"
	"sync"
	"time"
)

func New[T any]() *Channel[T] {
	c := new(Channel[T])
	c.in = make(chan T, 128)
	c.out = make(chan T, 128)
	c.list = list.New()

	go c.movePop()
	go c.movePush()

	return c
}

// 模拟无限缓冲的channel，与远程channel特性相同：close后in不能再写入，close并消费空后out返回nil,false
type Channel[T any] struct {
	sync.Mutex
	closed  bool
	in, out chan T
	list    *list.List
}

func (c *Channel[T]) In() chan<- T {
	return c.in
}

func (c *Channel[T]) Out() <-chan T {
	return c.out
}

func (c *Channel[T]) Close() {
	close(c.in)
}

func (c *Channel[T]) movePush() {
	for v := range c.in {
		c.Lock()
		c.list.PushBack(v)
		c.Unlock()
	}
	c.closed = true
}

func (c *Channel[T]) blockPop(t *time.Ticker) (data T, has bool) {
	for {
		c.Lock()
		switch {
		case c.closed && c.list.Len() == 0:
			c.Unlock()
			return
		case !c.closed && c.list.Len() == 0:
			c.Unlock()
			<-t.C
			continue
		default:
			front := c.list.Front()
			c.list.Remove(front)
			c.Unlock()
			return front.Value.(T), true
		}
	}
}

func (c *Channel[T]) movePop() {
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		data, has := c.blockPop(tick)
		if !has {
			break
		}
		c.out <- data
	}
	close(c.out)
}
