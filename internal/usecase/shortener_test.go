package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	internal "service-url-shortener/internal/errors"
	"service-url-shortener/internal/usecase/mocks"
)

func TestNewShortener(t *testing.T) {
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

func TestShortenerUseCase_Shorten_Errors_exist(t *testing.T) {
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
				fmt.Errorf("exist - uc.repo.Touch: %w",
					errors.New("UrlsRepo.Touch error"))),

			usecase: &ShortenerUseCase{
				repo: mocks.MockUrlsRepo{
					GetIDFunc: func(ctx context.Context, s string) (int, error) {
						return 0, nil
					},
					TouchFunc: func(ctx context.Context, i int) error {
						return errors.New("UrlsRepo.Touch error")
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

func TestShortenerUseCase_Shorten_Result_exist(t *testing.T) {
	expectedResult := "test_short"

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
				return expectedResult, nil
			},
		},
	}

	result, err := usecase.Shorten(context.Background(), "testurl")
	if err != nil {
		t.Fatalf("Unexpected error in test: %s", err)
	}
	if result != expectedResult {
		t.Fatalf("Expected result: %s\nGot: %s", expectedResult, result)
	}
}

func TestShortenerUseCase_Shorten_Errors(t *testing.T) {
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
						return 0, internal.ErrNotFoundURL
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
						return 0, internal.ErrNotFoundURL
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
						return 0, internal.ErrNotFoundURL
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
						return 0, internal.ErrNotFoundURL
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

func TestShortenerUseCase_Shorten_Result(t *testing.T) {
	expectedResult := "success"

	usecase := &ShortenerUseCase{
		repo: &mocks.MockUrlsRepo{
			RewriteFunc: func(ctx context.Context, s string) (int, error) {
				return 0, nil
			},
			GetIDFunc: func(ctx context.Context, s string) (int, error) {
				return 0, internal.ErrNotFoundURL
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

func TestShortenerUseCase_Lengthen(t *testing.T) {
	tests := []struct {
		expectedErr error
		usecase     Shortener
	}{
		{
			expectedErr: internal.ErrImpossibleShortURL,
			usecase: &ShortenerUseCase{
				digitiser: mocks.MockDigitiser{},
			},
		},
	}

	for testNum, test := range tests {

		original, err := test.usecase.Lengthen(context.Background(), "https://google.com/test")
		if original != "" {
			t.Fatalf("Expected empty string(original), Got: %s", original)
		}
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Fatalf("Test: %d\nExpected err: %s\nGot: %s", testNum, test.expectedErr, err)
		}
	}
}
