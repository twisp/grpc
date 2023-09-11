package token

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type tokenRefresherAlways struct {
	generator TokenGenerator
}

func NewTokenRefresherAlways(generator TokenGenerator) TokenRefresher {
	return &tokenRefresherAlways{
		generator: generator,
	}
}

func (r *tokenRefresherAlways) Token() ([]byte, bool, error) {
	generated, err := r.generator.Generate()
	return generated, true, err
}

func (r *tokenRefresherAlways) Stop() {}

type tokenRefresherTTL struct {
	generator     TokenGenerator
	maxTokenAge   float64
	maxRefreshAge float64
	now           func() time.Time
	token         atomic.Pointer[tokenWithTTL]
	mux           sync.Mutex
	stopped       atomic.Bool
}

var _ TokenRefresher = (*tokenRefresherTTL)(nil)

type tokenWithTTL struct {
	token []byte
	exp   time.Time
	now   func() time.Time
}

func (t *tokenWithTTL) Valid() bool {
	if t == nil {
		return false
	}
	return t.exp.After(t.now())
}

func NewTokenRefresherTTL(
	generator TokenGenerator,
	maxTokenAge float64,
	maxRefreshAge float64,
	now func() time.Time,
) (TokenRefresher, error) {
	if maxTokenAge < 0.0 || maxTokenAge > 1.0 {
		return nil, errors.New("max percent age must be between 0.0 and 1.0")
	}

	r := &tokenRefresherTTL{
		generator:     generator,
		maxTokenAge:   maxTokenAge,
		maxRefreshAge: maxRefreshAge,
		now:           now,
	}

	go func() { _, _, _ = r.Token() }()

	return r, nil
}

func (r *tokenRefresherTTL) Token() ([]byte, bool, error) {
	if token := r.token.Load(); token.Valid() {
		return token.token, false, nil
	}
	return r.refresh()
}

func (r *tokenRefresherTTL) refresh() ([]byte, bool, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if token := r.token.Load(); token.Valid() {
		return token.token, false, nil
	}

	token, err := r.generator.Generate()
	if err != nil {
		return nil, false, err
	}

	now := r.now()
	ttl := r.generator.TTL()
	age := time.Duration(float64(ttl.Nanoseconds()) * r.maxTokenAge)
	exp := now.Add(age)

	r.token.Store(&tokenWithTTL{
		token: token,
		exp:   exp,
		now:   r.now,
	})

	if !r.stopped.Load() {
		when := time.Duration(float64(ttl.Nanoseconds()) * r.maxRefreshAge)
		timer := time.NewTimer(when)
		go func() {
			<-timer.C
			if !r.stopped.Load() {
				_, _, _ = r.refresh()
			}
		}()
	}

	return token, true, nil
}

func (r *tokenRefresherTTL) Stop() {
	r.stopped.Store(true)
}
