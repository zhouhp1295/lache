package lache

import "time"

type ItemMode int

const (
	Interval ItemMode = 0 // 自动更新模式
	Expire   ItemMode = 1 // 过期模式
)

const DefaultGroup = "Default"

type IntervalHandler func() interface{}

type ItemOptions struct {
	Mode            ItemMode        //模式 interval 循环更新模式, expire 过期模式
	Group           string          // 分组
	Interval        time.Duration   // 循环更新时间
	IntervalHandler IntervalHandler //循环更新
	Expires         time.Duration   // 有效期限
}

type Option func(opt *ItemOptions)

func WithMode(mode ItemMode) Option {
	return func(opt *ItemOptions) {
		opt.Mode = mode
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opt *ItemOptions) {
		opt.Interval = interval
	}
}

func WithIntervalHandler(handler IntervalHandler) Option {
	return func(opt *ItemOptions) {
		opt.IntervalHandler = handler
	}
}

func WithGroup(group string) Option {
	return func(opt *ItemOptions) {
		opt.Group = group
	}
}
