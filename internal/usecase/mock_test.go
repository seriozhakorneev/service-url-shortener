package usecase_test

import "context"

type mockDigitiser struct {
	DigitFunc  func(s string) (int, error)
	StringFunc func(i int) (string, error)
	MaxFunc    func() int
	LengthFunc func() int
}

func (m mockDigitiser) Digit(s string) (int, error) {
	if m.DigitFunc != nil {
		return m.DigitFunc(s)
	}
	return 0, nil
}

func (m mockDigitiser) String(i int) (string, error) {
	if m.StringFunc != nil {
		return m.StringFunc(i)
	}
	return "", nil
}

func (m mockDigitiser) Max() int {
	if m.MaxFunc != nil {
		return m.Max()
	}
	return 0
}

func (m mockDigitiser) Length() int {
	if m.LengthFunc != nil {
		return m.LengthFunc()
	}
	return 0
}

type mockUrlsRepo struct {
	CreateFunc  func(ctx context.Context, s string) (int, error)
	RewriteFunc func(ctx context.Context, s string) (int, error)
	GetURLFunc  func(ctx context.Context, i int) (string, error)
	GetIDFunc   func(ctx context.Context, s string) (int, error)
	TouchFunc   func(ctx context.Context, i int) error
	CountFunc   func(ctx context.Context) (int, error)
}

func (m mockUrlsRepo) Create(ctx context.Context, s string) (int, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, s)
	}
	return 0, nil
}

func (m mockUrlsRepo) Rewrite(ctx context.Context, s string) (int, error) {
	if m.RewriteFunc != nil {
		return m.RewriteFunc(ctx, s)
	}
	return 0, nil
}

func (m mockUrlsRepo) GetURL(ctx context.Context, i int) (string, error) {
	if m.GetURLFunc != nil {
		return m.GetURLFunc(ctx, i)
	}
	return "", nil
}

func (m mockUrlsRepo) GetID(ctx context.Context, s string) (int, error) {
	if m.GetIDFunc != nil {
		return m.GetIDFunc(ctx, s)
	}
	return 0, nil
}

func (m mockUrlsRepo) Touch(ctx context.Context, i int) error {
	if m.TouchFunc != nil {
		return m.TouchFunc(ctx, i)
	}
	return nil
}

func (m mockUrlsRepo) Count(ctx context.Context) (int, error) {
	if m.CountFunc != nil {
		return m.Count(ctx)
	}
	return 0, nil
}
