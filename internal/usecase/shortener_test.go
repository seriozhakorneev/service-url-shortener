package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	service "service-url-shortener/internal/errors"
	"service-url-shortener/internal/usecase/mocks"
)

const (
	expectedResult = "success"
	testURL        = "https://google.com/test"
)

func TestNewShortener(t *testing.T) {
	t.Parallel()

	repo := mocks.MockUrlsRepo{}
	digitiser := mocks.MockDigitiser{}
	blank := "testblank"

	expected := &ShortenerUseCase{
		repo:      repo,
		digitiser: digitiser,
		blank:     blank,
	}

	result := NewShortener(repo, digitiser, blank)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Expected shortener:%v\nGot: %v", expected, result)
	}
}

func TestShortenerUseCaseShortenErrorsExist(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedErr error
		usecase     Shortener
	}{
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Shorten - %w",
				fmt.Errorf("exist - uc.repo.GetID: %w",
					fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w",
						errors.New("sql error")))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0,
							fmt.Errorf("UrlsRepo - GetID - r.Pool.QueryRow.Scan: %w",
								errors.New("sql error"))
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Shorten - %w",
				fmt.Errorf("exist - uc.digitiser.String: %w",
					errors.New("digitiser error"))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, nil
					},
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(i int) (string, error) {
						return "", errors.New("digitiser error")
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Shorten - %w",
				fmt.Errorf("exist - uc.repo.Activate: %w",
					errors.New("UrlsRepo.Activate error"))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, nil
					},
					TouchFunc: func(ctx context.Context, i int) error {
						return errors.New("UrlsRepo.Activate error")
					},
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(i int) (string, error) {
						return "test string", nil
					},
				},
			},
		},
	}

	for testNum, test := range tests {
		short, err := test.usecase.Shorten(context.Background(), "https://google.com")
		if short != "" {
			t.Fatalf("Expected empty string(short), Got: %s", short)
		}

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Fatalf("Test: %d\nExpected err: %s\nGot: %s", testNum, test.expectedErr, err)
		}
	}
}

func TestShortenerUseCaseShortenResultExist(t *testing.T) {
	t.Parallel()

	exResult := "test_short"

	usecase := &ShortenerUseCase{
		repo: mocks.MockUrlsRepo{
			GetIDFunc: func(ctx context.Context, s string) (int, error) {
				return 0, nil
			},
			TouchFunc: func(ctx context.Context, i int) error {
				return nil
			},
		},
		digitiser: mocks.MockDigitiser{
			StringFunc: func(i int) (string, error) {
				return exResult, nil
			},
		},
	}

	result, err := usecase.Shorten(context.Background(), "testurl")
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	if result != exResult {
		t.Fatalf("Expected result: %s\nGot: %s", exResult, result)
	}
}

func TestShortenerUseCaseShortenErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedErr error
		usecase     Shortener
	}{
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.repo.Count: %w",
				errors.New("uc.repo.Count error")),

			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, service.ErrNotFoundURL
					},
					CountFunc: func(ctx context.Context) (int, error) {
						return 0, errors.New("uc.repo.Count error")
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Shorten - uc.repo.Create: %w",
				errors.New("uc.repo.Create error")),

			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					CreateFunc: func(ctx context.Context, s string) (int, error) {
						return 0, errors.New("uc.repo.Create error")
					},
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, service.ErrNotFoundURL
					},
					CountFunc: func(ctx context.Context) (int, error) {
						return 0, nil
					},
				},
				digitiser: mocks.MockDigitiser{
					MaxFunc: func() int {
						return 1
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf("ShortenerUseCase - Shorten - uc.repo.Rewrite: %w",
				errors.New("uc.repo.Rewrite error")),

			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					RewriteFunc: func(ctx context.Context, s string) (int, error) {
						return 0, errors.New("uc.repo.Rewrite error")
					},
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, service.ErrNotFoundURL
					},
					CountFunc: func(ctx context.Context) (int, error) {
						return 1, nil
					},
				},
				digitiser: mocks.MockDigitiser{
					MaxFunc: func() int {
						return 0
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf("ShortenerUseCase - Shorten - uc.digitiser.String: %w",
				errors.New("uc.digitiser.String error")),

			usecase: &ShortenerUseCase{
				repo: &mocks.MockUrlsRepo{
					RewriteFunc: func(ctx context.Context, s string) (int, error) {
						return 0, nil
					},
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, service.ErrNotFoundURL
					},
					CountFunc: func(ctx context.Context) (int, error) {
						return 1, nil
					},
				},
				digitiser: mocks.MockDigitiser{
					StringFunc: func(i int) (string, error) {
						return "", errors.New("uc.digitiser.String error")
					},
					MaxFunc: func() int {
						return 0
					},
				},
			},
		},
	}

	for testNum, test := range tests {
		short, err := test.usecase.Shorten(context.Background(), "https://google.com")
		if short != "" {
			t.Fatalf("Expected empty string(short), Got: %s", short)
		}

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Fatalf("Test: %d\nExpected err: %s\nGot: %s", testNum, test.expectedErr, err)
		}
	}
}

func TestShortenerUseCaseShortenResult(t *testing.T) {
	t.Parallel()

	usecase := &ShortenerUseCase{
		repo: &mocks.MockUrlsRepo{
			RewriteFunc: func(ctx context.Context, s string) (int, error) {
				return 0, nil
			},
			GetIDFunc: func(ctx context.Context, s string) (int, error) {
				return 0, service.ErrNotFoundURL
			},
			CountFunc: func(ctx context.Context) (int, error) {
				return 1, nil
			},
		},
		digitiser: mocks.MockDigitiser{
			StringFunc: func(i int) (string, error) {
				return expectedResult, nil
			},
			MaxFunc: func() int {
				return 0
			},
		},
	}

	result, err := usecase.Shorten(context.Background(), "")
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	if result != expectedResult {
		t.Fatalf("Expected result: %s, Got: %s", result, expectedResult)
	}
}

func TestShortenerUseCaseLengthen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedErr error
		usecase     Shortener
	}{
		{
			expectedErr: service.ErrImpossibleShortURL,

			usecase: &ShortenerUseCase{
				digitiser: mocks.MockDigitiser{},
			},
		},
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Lengthen - uc.digitiser.Digit: %w",
				errors.New("uc.digitiser.Digit error")),

			usecase: &ShortenerUseCase{
				digitiser: mocks.MockDigitiser{
					DigitFunc: func(s string) (int, error) {
						return 0, errors.New("uc.digitiser.Digit error")
					},
					LengthFunc: func() int {
						return len(testURL) + 1
					},
				},
			},
		},
		{
			expectedErr: fmt.Errorf(
				"ShortenerUseCase - Lengthen - uc.repo.GetURL: %w",
				errors.New("uc.repo.GetURL error")),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(ctx context.Context, i int) (string, error) {
						return "", errors.New("uc.repo.GetURL error")
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int {
						return len(testURL) + 1
					},
				},
			},
		},
		{
			expectedErr: service.ErrNotFoundURL,

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetURLFunc: func(ctx context.Context, i int) (string, error) {
						return "", service.ErrNotFoundURL
					},
				},
				digitiser: mocks.MockDigitiser{
					LengthFunc: func() int {
						return len(testURL) + 1
					},
				},
			},
		},
	}

	for testNum, test := range tests {
		original, err := test.usecase.Lengthen(context.Background(), testURL)
		if original != "" {
			t.Fatalf("Expected empty string(original), Got: %s", original)
		}

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Fatalf("Test: %d\nExpected err: %s\nGot: %s", testNum, test.expectedErr, err)
		}
	}
}

func TestShortenerUseCaseLengthenResult(t *testing.T) {
	t.Parallel()

	testURL := "https://google.com/test"
	expectedResult := "success"

	usecase := &ShortenerUseCase{
		repo: mocks.MockUrlsRepo{
			GetURLFunc: func(ctx context.Context, i int) (string, error) {
				return expectedResult, nil
			},
		},
		digitiser: mocks.MockDigitiser{
			LengthFunc: func() int {
				return len(testURL) + 1
			},
		},
	}

	result, err := usecase.Lengthen(context.Background(), testURL)
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}

	if result != expectedResult {
		t.Fatalf("Expected result: %s\nGot: %s", expectedResult, result)
	}
}
