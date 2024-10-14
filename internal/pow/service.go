package pow

import (
	"github.com/maxrasky/faraway/internal/pow/hashcash"
)

const (
	defaultTargetBits = 20
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Challenge() []byte {
	return hashcash.NewResourceToken(defaultTargetBits)
}

func (s *Service) Verify(challenge, solution []byte) error {
	hc, err := hashcash.New(challenge)
	if err != nil {
		return err
	}

	ok, err := hc.Verify(string(solution))
	if err != nil {
		return err
	}

	if !ok {
		return hashcash.ErrSolutionFail
	}

	return nil
}

func (s *Service) Solve(challenge []byte) ([]byte, error) {
	hc, err := hashcash.New(challenge)
	if err != nil {
		return nil, err
	}

	solution, err := hc.Compute()
	if err != nil {
		return nil, err
	}

	return []byte(solution), nil
}
