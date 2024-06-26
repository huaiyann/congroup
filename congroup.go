package congroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/huaiyann/congroup/v2/internal/channel"
	"github.com/huaiyann/congroup/v2/internal/errs"
	"github.com/huaiyann/congroup/v2/internal/handlers"
	"github.com/huaiyann/congroup/v2/internal/panics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handler func(context.Context) error

type ConGroup struct {
	ctx        context.Context
	ctxErrOnce sync.Once
	wg         sync.WaitGroup
	errs       *errs.Errs
	panics     *panics.Panics
	c          *channel.Channel[handlers.IHanlerInfo]
	opt        *NewOpt
	span       trace.Span
}

func New(ctx context.Context, opts ...*NewOpt) *ConGroup {
	ctx, span := otel.GetTracerProvider().Tracer("congroup").Start(ctx, "congroup|New")

	cg := &ConGroup{
		ctx:    ctx,
		errs:   errs.New(),
		panics: panics.New(),
		c:      channel.New[handlers.IHanlerInfo](),
		opt:    NewOption().merge(opts...),
		span:   span,
	}

	for i := uint(0); i < cg.opt.GetMaxConcurrency(); i++ {
		cg.wg.Add(1)
		go cg.run()
	}

	return cg
}

func (g *ConGroup) Add(h Handler, opts ...*AddOpt) {
	opt := AddOption().merge(opts...)
	pc, file, line, _ := runtime.Caller(1)
	from := fmt.Sprintf("%s:%d", file, line)
	if callerFunc := runtime.FuncForPC(pc); callerFunc != nil {
		from += fmt.Sprintf("(%s)", callerFunc.Name())
	}
	g.c.In() <- &handlers.HandlerInfo{
		H:     handlers.Handler(h),
		From:  from,
		Label: opt.GetLabel(),
	}
}

func (g *ConGroup) Wait() error {
	g.c.Close()
	g.wg.Wait()
	g.span.End()
	if g.panics.Has() {
		panic(g.panics.Info())
	}
	if g.errs.Has() {
		return g.errs
	}
	return nil
}

func (g *ConGroup) run() {
	defer g.wg.Done()
	for {
		select {
		case v, ok := <-g.c.Out():
			if !ok {
				return
			}
			g.execHandler(v)
		case <-g.ctx.Done():
			g.ctxErrOnce.Do(func() {
				g.errs.Add(g.ctx.Err(), "context done when ConGroup running")
			})
			return
		}
	}

}

func (g *ConGroup) execHandler(h handlers.IHanlerInfo) {
	defer func() {
		reason := recover()
		if reason != nil {
			buf := make([]byte, 10240)
			length := runtime.Stack(buf, false)
			g.panics.Add(reason, buf[:length])
		}
	}()

	ctx := g.ctx
	name := "congroup|execHandler"
	if label := h.GetLabel(); label != "" {
		name += "|" + label
	}
	ctx, span := otel.GetTracerProvider().Tracer("congroup").Start(ctx, name)
	span.SetAttributes(attribute.String("added from", h.GetFrom()))
	defer span.End()

	err := h.GetHandler()(ctx)
	if err != nil {
		g.errs.Add(err, h.GetFrom())
	}
}
