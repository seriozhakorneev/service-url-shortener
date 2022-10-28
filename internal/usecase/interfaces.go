// Package usecase implements application business logic. Each logic group in own file.
package usecase

type (
	// Echo -.
	Echo interface {
		Rewrite(m map[string]any)
	}

	// Rewriter -.
	Rewriter interface {
		Rewrite(m map[string]any)
	}
)
