package mocks

import "context"

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
		return m.Max()
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
	CreateFunc  func(ctx context.Context, s string) (int, error)
	RewriteFunc func(ctx context.Context, s string) (int, error)
	GetURLFunc  func(ctx context.Context, i int) (string, error)
	GetIDFunc   func(ctx context.Context, s string) (int, error)
	TouchFunc   func(ctx context.Context, i int) error
	CountFunc   func(ctx context.Context) (int, error)
}

func (m MockUrlsRepo) Create(ctx context.Context, s string) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, s)
	}
	return 0, nil
}

func (m MockUrlsRepo) Rewrite(ctx context.Context, s string) (int, error) {
	if m.RewriteFunc != nil {
		return m.RewriteFunc(ctx, s)
	}
	return 0, nil
}

func (m MockUrlsRepo) GetURL(ctx context.Context, i int) (string, error) {
	if m.GetURLFunc != nil {
		return m.GetURLFunc(ctx, i)
	}
	return "", nil
}

func (m MockUrlsRepo) GetID(ctx context.Context, s string) (int, error) {
	if m.GetIDFunc != nil {
		return m.GetIDFunc(ctx, s)
	}
	return 0, nil
}

func (m MockUrlsRepo) Touch(ctx context.Context, i int) error {
	if m.TouchFunc != nil {
		return m.TouchFunc(ctx, i)
	}
	return nil
}

func (m MockUrlsRepo) Count(ctx context.Context) (int, error) {
	if m.CountFunc != nil {
		return m.Count(ctx)
	}
	return 0, nil
}
