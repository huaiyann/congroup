package congroup

type NewOpt struct {
	concurrency *uint
}

func (o *NewOpt) merge(opts ...*NewOpt) *NewOpt {
	for _, opt := range opts {
		if opt.concurrency != nil {
			o.concurrency = opt.concurrency
		}
	}
	return o
}

// SetMaxConcurrency 设置执行handler时的最大并发数，大于0有效，默认DefaultConcurrency
func (o *NewOpt) SetMaxConcurrency(c uint) *NewOpt {
	o.concurrency = &c
	return o
}

func (o *NewOpt) GetMaxConcurrency() uint {
	if o.concurrency == nil || *o.concurrency == 0 {
		return DefaultConcurrency
	}
	return *o.concurrency
}

func NewOption() *NewOpt {
	return new(NewOpt)
}

type AddOpt struct {
	label *string
}

func (o *AddOpt) merge(opts ...*AddOpt) *AddOpt {
	for _, opt := range opts {
		if opt.label != nil {
			o.label = opt.label
		}
	}
	return o
}

// SetLabel 对任务设置label，用于tracing时的区分
func (o *AddOpt) SetLabel(l string) *AddOpt {
	o.label = &l
	return o
}

func (o *AddOpt) GetLabel() string {
	if o.label == nil {
		return ""
	}
	return *o.label
}

func AddOption() *AddOpt {
	return new(AddOpt)
}
