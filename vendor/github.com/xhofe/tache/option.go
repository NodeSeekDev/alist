package tache

import (
	"log/slog"
	"time"
)

// Options is the options for manager
type Options struct {
	Works                int
	MaxRetry             int
	Timeout              *time.Duration
	PersistPath          string
	PersistDebounce      *time.Duration
	Running              bool
	Logger               *slog.Logger
	PersistReadFunction  func() ([]byte, error)
	PersistWriteFunction func([]byte) error
}

// DefaultOptions returns default options
func DefaultOptions() *Options {
	persistDebounce := 3 * time.Second
	return &Options{
		Works: 5,
		//MaxRetry: 1,
		PersistDebounce: &persistDebounce,
		Running:         true,
		Logger:          slog.Default(),
	}
}

// Option is the option for manager
type Option func(*Options)

// WithOptions set options
func WithOptions(opts Options) Option {
	return func(o *Options) {
		*o = opts
	}
}

// WithWorks set works
func WithWorks(works int) Option {
	return func(o *Options) {
		o.Works = works
	}
}

// WithMaxRetry set retry
func WithMaxRetry(maxRetry int) Option {
	return func(o *Options) {
		o.MaxRetry = maxRetry
	}
}

// WithTimeout set timeout
func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.Timeout = &timeout
	}
}

// WithPersistPath set persist path
func WithPersistPath(path string) Option {
	return func(o *Options) {
		o.PersistPath = path
	}
}

func WithPersistFunction(r func() ([]byte, error), w func([]byte) error) Option {
	return func(o *Options) {
		o.PersistReadFunction = r
		o.PersistWriteFunction = w
	}
}

// WithPersistDebounce set persist debounce
func WithPersistDebounce(debounce time.Duration) Option {
	return func(o *Options) {
		o.PersistDebounce = &debounce
	}
}

// WithRunning set running
func WithRunning(running bool) Option {
	return func(o *Options) {
		o.Running = running
	}
}

// WithLogger set logger
func WithLogger(logger *slog.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}
