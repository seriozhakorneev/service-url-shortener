package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"service-url-shortener/internal/usecase/mocks"
)

func TestShortenerUseCase_Shorten_Errors(t *testing.T) {

	tests := []struct {
		expectedErr error
		usecase     Shortener
	}{
		// uc.exist errors
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
