package delivery

import "github.com/pkg/errors"

var (
	// Handler errors
	InvalidBody       = errors.New("invalid body in request")
	InvalidParameters = errors.New("invalid parameters")
)
