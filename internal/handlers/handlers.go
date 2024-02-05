package handlers

import "context"

type Handler func(context.Context) error

type IHanlerInfo interface {
	GetFrom() string
	GetHandler() Handler
}

type HandlerInfo struct {
	H    Handler
	From string
}

func (h *HandlerInfo) GetFrom() string {
	return h.From
}

func (h *HandlerInfo) GetHandler() Handler {
	return h.H
}
