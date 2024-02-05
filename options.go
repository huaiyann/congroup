package congroup

type Option struct {
	concurrency *uint
}

func (o *Option) merge(opts ...*Option) {
	for _, opt := range opts {
		if opt.concurrency != nil {
			o.concurrency = opt.concurrency
		}
	}
}

func (o *Option) GetMaxConcurrency() uint {
	if o.concurrency == nil || *o.concurrency == 0 {
		return DefaultConcurrency
	}
	return *o.concurrency
}

// MaxConcurrency 设置执行handler时的最大并发数，大于0有效，默认DefaultConcurrency
func MaxConcurrency(c uint) *Option {
	return &Option{
		concurrency: &c,
	}
}
