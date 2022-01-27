package tui

import (
	"github.com/kochavalabs/m8/internal/cfg"
	"github.com/kochavalabs/mazzaroth-go"
)

type options struct {
	Client mazzaroth.Client
	Config *cfg.Configuration
}

type Options interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fto *funcOption) apply(opt *options) {
	fto.f(opt)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithConfiguration(cfg *cfg.Configuration) Options {
	return newFuncOption(func(o *options) {
		o.Config = cfg
	})
}

func WithMazzarothClient(client mazzaroth.Client) Options {
	return newFuncOption(func(o *options) {
		o.Client = client
	})
}

func defaultOptions() *options {
	return nil
}
