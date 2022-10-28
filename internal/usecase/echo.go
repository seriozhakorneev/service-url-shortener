package usecase

// EchoUseCase -.
type EchoUseCase struct {
	rewriter Rewriter
}

// New -.
func New(r Rewriter) *EchoUseCase {
	return &EchoUseCase{rewriter: r}
}

// Rewrite - runs rewriter if rewrite rules are active.
func (uc *EchoUseCase) Rewrite(data map[string]any) {
	if uc.rewriter != nil {
		uc.rewriter.Rewrite(data)
	}
}
