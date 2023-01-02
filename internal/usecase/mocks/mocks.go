package mocks

import (
	"context"
	"time"
)

type MockDigitiser struct {
	DigitFunc  func(s string) (int, error)
	StringFunc func(i int) (string, error)
	MaxFunc    func() int
	LengthFunc func() int
}

func (m MockDigitiser) Digit(s string) (int, error) {
	if m.DigitFunc != nil {
		return m.DigitFunc(s)
	}
	return 0, nil
}

func (m MockDigitiser) String(i int) (string, error) {
	if m.StringFunc != nil {
		return m.StringFunc(i)
	}
	return "", nil
}

func (m MockDigitiser) Max() int {
	if m.MaxFunc != nil {
		return m.MaxFunc()
	}
	return 0
}

func (m MockDigitiser) Length() int {
	if m.LengthFunc != nil {
		return m.LengthFunc()
	}
	return 0
}

type MockUrlsRepo struct {
	CreateFunc   func(context.Context, string, time.Duration) (int, error)
	RewriteFunc  func(context.Context, string, time.Duration) (int, error)
	GetURLFunc   func(context.Context, int) (string, time.Time, error)
	GetIDFunc    func(context.Context, string) (int, error)
	ActivateFunc func(context.Context, int, time.Duration) error
	CountFunc    func(context.Context) (int, error)
}

func (m MockUrlsRepo) Create(ctx context.Context, s string, d time.Duration) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, s, d)
	}
	return 0, nil
}

func (m MockUrlsRepo) Rewrite(ctx context.Context, s string, d time.Duration) (int, error) {
	if m.RewriteFunc != nil {
		return m.RewriteFunc(ctx, s, d)
	}
	return 0, nil
}

func (m MockUrlsRepo) GetURL(ctx context.Context, i int) (string, time.Time, error) {
	if m.GetURLFunc != nil {
		return m.GetURLFunc(ctx, i)
	}
	return "", time.Time{}, nil
}

func (m MockUrlsRepo) GetID(ctx context.Context, s string) (int, error) {
	if m.GetIDFunc != nil {
		return m.GetIDFunc(ctx, s)
	}
	return 0, nil
}

func (m MockUrlsRepo) Activate(ctx context.Context, i int, d time.Duration) error {
	if m.ActivateFunc != nil {
		return m.ActivateFunc(ctx, i, d)
	}
	return nil
}

func (m MockUrlsRepo) Last(ctx context.Context) (int, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

type MockUrlsCache struct {
	SetFunc func(string, string, time.Duration) error
	GetFunc func(string) (*string, error)
}

func (m MockUrlsCache) Set(k string, v string, d time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(k, v, d)
	}
	return nil
}

func (m MockUrlsCache) Get(k string) (*string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(k)
	}
	return nil, nil
}
