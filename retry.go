package buaaclock

import (
	"crypto/rand"
	"math/big"
	"time"
)

// Config defines the config for addon.
type RetryConfig struct {
	// InitialInterval defines the initial time interval for backoff algorithm.
	//
	// Optional. Default: 1 * time.Second
	InitialInterval time.Duration

	// MaxBackoffTime defines maximum time duration for backoff algorithm. When
	// the algorithm is reached this time, rest of the retries will be maximum
	// 32 seconds.
	//
	// Optional. Default: 32 * time.Second
	MaxBackoffTime time.Duration

	// Multiplier defines multiplier number of the backoff algorithm.
	//
	// Optional. Default: 2.0
	Multiplier float64

	// MaxRetryCount defines maximum retry count for the backoff algorithm.
	//
	// Optional. Default: 10
	MaxRetryCount int

	// currentInterval tracks the current waiting time.
	//
	// Optional. Default: 1 * time.Second
	currentInterval time.Duration
}

// configDefault sets the config values if they are not set.
func retryConfigDefault(config ...RetryConfig) RetryConfig {
	// defaultConfig is the default config for retry.
	defaultConfig := RetryConfig{
		InitialInterval: 1 * time.Second,
		MaxBackoffTime:  32 * time.Second,
		Multiplier:      2.0,
		MaxRetryCount:   10,
		currentInterval: 1 * time.Second,
	}

	if len(config) == 0 {
		return defaultConfig
	}

	cfg := config[0]
	if cfg.InitialInterval == 0 {
		cfg.InitialInterval = defaultConfig.InitialInterval
	}

	if cfg.MaxBackoffTime == 0 {
		cfg.MaxBackoffTime = defaultConfig.MaxBackoffTime
	}

	if cfg.Multiplier <= 0 {
		cfg.Multiplier = defaultConfig.Multiplier
	}

	if cfg.MaxRetryCount <= 0 {
		cfg.MaxRetryCount = defaultConfig.MaxRetryCount
	}

	if cfg.currentInterval != cfg.InitialInterval {
		cfg.currentInterval = defaultConfig.currentInterval
	}

	return cfg
}

// ExponentialBackoff is a retry mechanism for retrying some calls.
type ExponentialBackoff struct {
	// InitialInterval is the initial time interval for backoff algorithm.
	InitialInterval time.Duration

	// MaxBackoffTime is the maximum time duration for backoff algorithm. It limits
	// the maximum sleep time.
	MaxBackoffTime time.Duration

	// Multiplier is a multiplier number of the backoff algorithm.
	Multiplier float64

	// MaxRetryCount is the maximum number of retry count.
	MaxRetryCount int

	// currentInterval tracks the current sleep time.
	currentInterval time.Duration
}

// NewExponentialBackoff creates a ExponentialBackoff with default values.
func NewExponentialBackoff(config ...RetryConfig) *ExponentialBackoff {
	cfg := retryConfigDefault(config...)

	return &ExponentialBackoff{
		InitialInterval: cfg.InitialInterval,
		MaxBackoffTime:  cfg.MaxBackoffTime,
		Multiplier:      cfg.Multiplier,
		MaxRetryCount:   cfg.MaxRetryCount,
		currentInterval: cfg.currentInterval,
	}
}

// Retry is the core logic of the retry mechanism. If the calling function returns
// nil as an error, then the Retry method is terminated with returning nil. Otherwise,
// if all function calls are returned error, then the method returns this error.
func (e *ExponentialBackoff) Retry(f func() error) error {
	if e.currentInterval <= 0 {
		e.currentInterval = e.InitialInterval
	}

	var err error
	for i := 0; i < e.MaxRetryCount; i++ {
		err = f()
		if err == nil {
			return nil
		}

		next := e.next()
		time.Sleep(next)
	}

	return err
}

// next calculates the next sleeping time interval.
func (e *ExponentialBackoff) next() time.Duration {
	// generate a random value between [0, 1000)
	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return e.MaxBackoffTime
	}

	t := e.currentInterval + (time.Duration(n.Int64()) * time.Millisecond)

	e.currentInterval = time.Duration(float64(e.currentInterval) * e.Multiplier)
	if t >= e.MaxBackoffTime {
		e.currentInterval = e.MaxBackoffTime
		return e.MaxBackoffTime
	}

	return t
}
