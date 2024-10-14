package hashcash

import "errors"

var (
	ErrSolutionFail  = errors.New("exceeded 2^20 iterations failed to find solution")
	ErrResourceEmpty = errors.New("empty hashcash resource")
	ErrInvalidHeader = errors.New("invalid hashcash header format")
	ErrNoCollision   = errors.New("no collision most significant bits are not zero")
	ErrTimestamp     = errors.New("time stamp is too far into the future or expired")
)
