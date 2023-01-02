package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	service "service-url-shortener/internal/errors"
	"service-url-shortener/internal/usecase/mocks"
)

const (
	tBlank    = "test_blank/"
	tShort    = "test_short"
	tShortURL = tBlank + tShort
	tOriginal = "https://google.com"
)

func TestNewShortener(t *testing.T) {
	t.Parallel()

	repo := mocks.MockUrlsRepo{}
	cache := mocks.MockUrlsCache{}
	digitiser := mocks.MockDigitiser{}
	blank := "testblank"

	expected := &ShortenerUseCase{
		repo:      repo,
		cache:     cache,
		digitiser: digitiser,
		blank:     blank,
	}

	result := NewShortener(repo, cache, digitiser, blank)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected shortener:%v\nGot: %v", expected, result)
	}
}

func TestShortenerUseCaseShortenExistErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expError error
		usecase  Shortener
	}{
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.exist: %w",
				fmt.Errorf("exist - uc.repo.GetID: %w",
					fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w",
						errors.New("sql error")))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(context.Context, string) (int, error) {
						return 0,
							fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w",
								errors.New("sql error"))
					},
				},
			},
		},
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.exist: %w",
				fmt.Errorf("exist - uc.digitiser.String: %w",
					errors.New("digitiser error"))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(context.Context, string) (int, error) { return 0, nil },
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(int) (string, error) {
						return "", errors.New("digitiser error")
					},
				},
			},
		},
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.exist: %w",
				fmt.Errorf("exist - uc.repo.Activate: %w",
					errors.New("UrlsRepo.Activate error"))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(context.Context, string) (int, error) { return 0, nil },
					ActivateFunc: func(context.Context, int, time.Duration) error {
						return errors.New("UrlsRepo.Activate error")
					},
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(int) (string, error) { return "test string", nil },
				},
			},
		},
	}

	for testNum, test := range tests {
		short, err := test.usecase.Shorten(context.Background(), tOriginal, time.Duration(0))
		if short != "" {
			t.Fatalf("Expected empty string(short), Got: %s", short)
		}

		if !reflect.DeepEqual(err, test.expError) {
			t.Fatalf("Test %d\nExpected err: %s\nGot: %s", testNum+1, test.expError, err)
		}
	}
}

func TestShortenerUseCaseShortenExistResult(t *testing.T) {
	t.Parallel()

	usecase := &ShortenerUseCase{
		blank: tBlank,
		repo: mocks.MockUrlsRepo{
			GetIDFunc:    func(context.Context, string) (int, error) { return 0, nil },
			ActivateFunc: func(context.Context, int, time.Duration) error { return nil },
		},
		digitiser: mocks.MockDigitiser{
			StringFunc: func(int) (string, error) { return tShort, nil },
		},
		cache: mocks.MockUrlsCache{
			SetFunc: func(s string, s2 string, duration time.Duration) error { return nil },
		},
	}

	result, err := usecase.Shorten(context.Background(), tOriginal, time.Duration(0))
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	if result != tShortURL {
		t.Fatalf("Expected result: %s\nGot: %s", tShortURL, result)
	}
}

func TestShortenerUseCaseShortenCreateErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expError error
		usecase  Shortener
	}{
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.create: %w",
				fmt.Errorf("create - uc.repo.Count: %w",
					errors.New("uc.repo.Count error"),
				),
			),
			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					GetIDFunc: func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
					CountFunc: func(context.Context) (int, error) { return 0, errors.New("uc.repo.Count error") },
				},
			},
		},
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.create: %w",
				fmt.Errorf("create - uc.repo.Create: %w",
					errors.New("uc.repo.Create error"),
				),
			),
			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					CreateFunc: func(context.Context, string, time.Duration) (int, error) {
						return 0, errors.New("uc.repo.Create error")
					},
					GetIDFunc: func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
					CountFunc: func(context.Context) (int, error) { return 0, nil },
				},
				digitiser: mocks.MockDigitiser{
					MaxFunc: func() int { return 1 },
				},
			},
		},
		{
			expError: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.create: %w",
				fmt.Errorf("create - uc.repo.Rewrite: %w",
					errors.New("uc.repo.Rewrite error"),
				),
			),
			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					RewriteFunc: func(context.Context, string, time.Duration) (int, error) {
						return 0, errors.New("uc.repo.Rewrite error")
					},
					GetIDFunc: func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
					CountFunc: func(context.Context) (int, error) { return 1, nil },
				},
				digitiser: mocks.MockDigitiser{
					MaxFunc: func() int { return 0 },
				},
			},
		},
		{
			expError: fmt.Errorf("ShortenerUseCase - Shorten - uc.create: %w",
				fmt.Errorf("create - uc.digitiser.String: %w",
					errors.New("uc.digitiser.String error"),
				),
			),
			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					RewriteFunc: func(context.Context, string, time.Duration) (int, error) { return 0, nil },
					GetIDFunc:   func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
					CountFunc:   func(context.Context) (int, error) { return 1, nil },
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(int) (string, error) {
						return "", errors.New("uc.digitiser.String error")
					},
					MaxFunc: func() int { return 0 },
				},
			},
		},
	}

	for testNum, test := range tests {
		short, err := test.usecase.Shorten(context.Background(), tOriginal, time.Duration(0))
		if short != "" {
			t.Fatalf("Expected empty string(short), Got: %s", short)
		}

		if !reflect.DeepEqual(err, test.expError) {
			t.Fatalf("Test %d\nExpected err: %s\nGot: %s", testNum+1, test.expError, err)
		}
	}
}

func TestShortenerUseCaseShortenCreateResult(t *testing.T) {
	t.Parallel()

	usecase := &ShortenerUseCase{
		blank: tBlank,
		repo: &mocks.MockUrlsRepo{
			RewriteFunc: func(context.Context, string, time.Duration) (int, error) { return 0, nil },
			GetIDFunc:   func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
			CountFunc:   func(context.Context) (int, error) { return 1, nil },
		},
		digitiser: mocks.MockDigitiser{
			StringFunc: func(int) (string, error) { return tShort, nil },
			MaxFunc:    func() int { return 0 },
		},
		cache: mocks.MockUrlsCache{
			SetFunc: func(string, string, time.Duration) error { return nil },
		},
	}

	result, err := usecase.Shorten(context.Background(), tOriginal, time.Duration(0))
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	if result != tShortURL {
		t.Fatalf("Expected result: %s, Got: %s", result, tShortURL)
	}
}

func TestShortenerUseCaseShortenCreateSetCacheError(t *testing.T) {
	t.Parallel()

	setFuncError := fmt.Errorf("UrlsCache - Set - c.Client.Set: %w",
		errors.New("failed to set cache because of something"),
	)

	usecase := &ShortenerUseCase{
		blank: tBlank,
		repo: &mocks.MockUrlsRepo{
			RewriteFunc: func(context.Context, string, time.Duration) (int, error) { return 0, nil },
			GetIDFunc:   func(context.Context, string) (int, error) { return 0, service.ErrNotFoundURL },
			CountFunc:   func(context.Context) (int, error) { return 1, nil },
		},
		digitiser: mocks.MockDigitiser{
			StringFunc: func(i int) (string, error) { return tShort, nil },
			MaxFunc:    func() int { return 0 },
		},
		cache: mocks.MockUrlsCache{
			SetFunc: func(string, string, time.Duration) error { return setFuncError },
		},
	}

	expError := fmt.Errorf("%w: %s", service.ErrCaching, setFuncError)

	result, err := usecase.Shorten(context.Background(), tOriginal, time.Duration(0))

	if result != tShortURL {
		t.Fatalf("Expected result: %s, Got: %s", result, tShortURL)
	}

	if !errors.Is(err, service.ErrCaching) ||
		!reflect.DeepEqual(err, expError) {
		t.Fatalf(
			"Expected error of type '%s'\nExpected error: %s\nGot: %s",
			service.ErrCaching,
			expError,
			err,
		)
	}
}

func TestShortenerUseCaseLengthenCacheGetResult(t *testing.T) {
	t.Parallel()

	o := tOriginal

	uc := NewShortener(
		&mocks.MockUrlsRepo{},
		&mocks.MockUrlsCache{GetFunc: func(string) (*string, error) { return &o, nil }},
		&mocks.MockDigitiser{},
		"",
	)

	original, err := uc.Lengthen(context.Background(), tShortURL)
	if err != nil {
		t.Fatalf("Expected error: nil\nGot:%s", err)
	}

	if original != tOriginal {
		t.Fatalf("Expected original url: %s\nGot%s", tOriginal, original)
	}

}

func TestShortenerUseCaseLengthenParseShortErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		shortURL string
		expError error
		uc       Shortener
	}{
		{
			shortURL: "https://:abc",
			expError: fmt.Errorf("ShortenerUseCase - Lengthen - uc.parseShort: %w",
				fmt.Errorf("parseShort - urlParse.Parse: %w",
					errors.New("parse \"https://:abc\": invalid port \":abc\" after host"),
				),
			),
			uc: &ShortenerUseCase{
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) { return nil, nil },
				},
				digitiser: mocks.MockDigitiser{},
			},
		},
		{
			shortURL: "s",
			expError: fmt.Errorf("ShortenerUseCase - Lengthen - uc.parseShort: %w",
				service.ErrImpossibleShortURL,
			),
			uc: &ShortenerUseCase{
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) { return nil, nil },
				},
				digitiser: mocks.MockDigitiser{},
			},
		},
	}

	for tNum, test := range tests {
		original, err := test.uc.Lengthen(context.Background(), test.shortURL)
		if original != "" {
			t.Fatalf("Test %d\nExpected empty string original\nGot:%s", tNum, original)
		}

		if err.Error() != test.expError.Error() {
			t.Fatalf("Test %d\nExpected error: %s\nGot:%s", tNum, test.expError, err)
		}
	}
}

func TestShortenerUseCaseLengthenGetErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		shortURL string
		expError error
		uc       Shortener
	}{
		{
			shortURL: tShortURL,
			expError: fmt.Errorf("ShortenerUseCase - Lengthen - uc.get: %w",
				fmt.Errorf("get - uc.digitiser.Digit: %w",
					fmt.Errorf("test digitiser error"),
				),
			),
			uc: &ShortenerUseCase{
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) { return nil, nil },
				},
				digitiser: mocks.MockDigitiser{
					DigitFunc: func(string) (int, error) {
						return 0, fmt.Errorf("test digitiser error")
					},
					LengthFunc: func() int { return len(tShortURL) + 1 },
				},
			},
		},
		{
			shortURL: tShortURL,
			expError: fmt.Errorf("ShortenerUseCase - Lengthen - uc.get: %w",
				service.ErrNotFoundURL,
			),
			uc: &ShortenerUseCase{
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) { return nil, nil },
				},
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(context.Context, int) (string, time.Time, error) {
						return "", time.Time{}, service.ErrNotFoundURL
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int { return len(tShortURL) + 1 },
				},
			},
		},
		{
			shortURL: tShortURL,
			expError: fmt.Errorf("ShortenerUseCase - Lengthen - uc.get: %w",
				fmt.Errorf("get - uc.repo.GetURL: %w",
					fmt.Errorf("UrlsRepo - GetURL - r.Pool.QueryRow.Scan: %w",
						fmt.Errorf("test sql error"),
					),
				),
			),
			uc: &ShortenerUseCase{
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) { return nil, nil },
				},
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(context.Context, int) (string, time.Time, error) {
						return "", time.Time{}, fmt.Errorf(
							"UrlsRepo - GetURL - r.Pool.QueryRow.Scan: %w",
							fmt.Errorf("test sql error"),
						)
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int { return len(tShortURL) + 1 },
				},
			},
		},
	}

	for tNum, test := range tests {
		tNum++
		original, err := test.uc.Lengthen(context.Background(), test.shortURL)
		if original != "" {
			t.Fatalf("Test %d\nExpected empty string original\nGot:%s", tNum, original)
		}

		if !reflect.DeepEqual(err, test.expError) {
			t.Fatalf("Test %d\nExpected error: %s\nGot: %s", tNum, test.expError, err)
		}
	}
}

func TestShortenerUseCaseLengthenResult(t *testing.T) {
	t.Parallel()

	uc := &ShortenerUseCase{
		cache: mocks.MockUrlsCache{
			GetFunc: func(string) (*string, error) { return nil, nil },
			SetFunc: func(string, string, time.Duration) error { return nil },
		},
		repo: mocks.MockUrlsRepo{
			GetURLFunc: func(context.Context, int) (string, time.Time, error) {
				return tOriginal, time.Time{}, nil
			},
		},
		digitiser: mocks.MockDigitiser{
			LengthFunc: func() int { return len(tShortURL) + 1 },
		},
	}

	original, err := uc.Lengthen(context.Background(), tShortURL)
	if err != nil {
		t.Fatalf("Expected error nil\nGot:%s", err)
	}

	if original != tOriginal {
		t.Fatalf("Expected original url: %s\nGot:%s", tOriginal, original)
	}
}

func TestShortenerUseCaseLengthenTTL(t *testing.T) {
	t.Parallel()

	pTime, _ := time.Parse(time.RFC3339, "2022-12-30T21:10:29.03Z")
	expDuration := pTime.Sub(time.Now().UTC())
	criterion := time.Second

	uc := &ShortenerUseCase{
		cache: mocks.MockUrlsCache{
			GetFunc: func(string) (*string, error) { return nil, nil },
			SetFunc: func(a string, b string, ttl time.Duration) error {
				diff := ttl - expDuration
				if diff > criterion {
					return fmt.Errorf(
						"\ndiff between TTL(%s) and expected Duration(%s) "+
							"is more than criterion(%s) diff: %s",
						ttl, expDuration,
						criterion,
						diff,
					)
				}
				return nil
			},
		},
		repo: mocks.MockUrlsRepo{
			GetURLFunc: func(context.Context, int) (string, time.Time, error) {
				return tOriginal, pTime, nil
			},
		},
		digitiser: mocks.MockDigitiser{
			LengthFunc: func() int { return len(tShortURL) + 1 },
		},
	}

	original, err := uc.Lengthen(context.Background(), tShortURL)
	if err != nil {
		if errors.Is(err, service.ErrCaching) {
			t.Fatalf("Time to live calculation failed: %s", err)
		}
		t.Fatalf("Expected error nil\nGot:%s", err)
	}

	if original != tOriginal {
		t.Fatalf("Expected original url: %s\nGot:%s", tOriginal, original)
	}
}

func TestShortenerUseCaseLengthenCacheError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expError error
		uc       Shortener
	}{
		{
			expError: fmt.Errorf("%w: %s", service.ErrCaching,
				fmt.Errorf("cache GET error"),
			),
			uc: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(context.Context, int) (string, time.Time, error) {
						return tOriginal, time.Time{}, nil
					},
				},
				cache: mocks.MockUrlsCache{
					GetFunc: func(string) (*string, error) {
						return nil, fmt.Errorf("cache GET error")
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int { return len(tShortURL) + 1 },
				},
			},
		},
		{
			expError: fmt.Errorf("%w: %s", service.ErrCaching,
				fmt.Errorf("cache SET error"),
			),
			uc: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(context.Context, int) (string, time.Time, error) {
						return tOriginal, time.Time{}, nil
					},
				},
				cache: mocks.MockUrlsCache{
					SetFunc: func(string, string, time.Duration) error {
						return fmt.Errorf("cache SET error")
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int { return len(tShortURL) + 1 },
				},
			},
		},
	}

	for tNum, test := range tests {
		tNum++
		original, err := test.uc.Lengthen(context.Background(), tShortURL)
		if original != tOriginal {
			t.Fatalf("Test %d\nExpected original url: %s\nGot:%s", tNum, tOriginal, original)
		}

		if !reflect.DeepEqual(err, test.expError) || !errors.Is(err, service.ErrCaching) {
			t.Fatalf("Test %d\nExpected error: %s\nGot: %s", tNum, test.expError, err)
		}
	}
}
