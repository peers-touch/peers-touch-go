package source

import (
	"github.com/peers-touch/peers-touch/station/frame/core/option"
	"github.com/peers-touch/peers-touch/station/frame/core/pkg/cli"
	"github.com/peers-touch/peers-touch/station/frame/core/pkg/config/encoder"
	"github.com/peers-touch/peers-touch/station/frame/core/pkg/config/encoder/json"
)

type sourceOptionsKey struct{}
type contextKey struct{}

var wrapper = option.NewWrapper[Options](sourceOptionsKey{}, func(options *option.Options) *Options {
	return &Options{
		Options: options,
	}
})

type Options struct {
	*option.Options

	// Encoder
	Encoder encoder.Encoder
}

func GetOptions(root *option.Options) *Options {
	if root == nil || root.Ctx().Value(sourceOptionsKey{}) == nil {
		panic("rootOpts or sourceOptions is nil")
	}

	opts := root.Ctx().Value(sourceOptionsKey{}).(*Options)
	if opts.Encoder == nil {
		// todo check existed or set by more appropriate way
		opts.Encoder = json.NewEncoder()
	}

	return opts
}

// WithEncoder sets the source encoder
func WithEncoder(e encoder.Encoder) option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.Encoder = e
	})
}

// Context sets the cli context
func Context(c *cli.Context) option.Option {
	return wrapper.Wrap(func(opts *Options) {
		opts.AppendCtx(contextKey{}, c)
	})
}

func WrapOption(f func(*Options)) option.Option {
	return wrapper.Wrap(f)
}
