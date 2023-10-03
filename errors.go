package backoff

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid config: curve and func are required")
)
