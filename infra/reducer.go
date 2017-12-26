package infra

import (
	"github.com/pkg/errors"
)

// Reduce transitions the local dnote state by consuming the action returned
// from the server
func Reduce(ctx DnoteCtx, a Action) error {
	switch a.Type {
	case ActionAddNote:
		// TODO: assert data type and consume
	default:
		return errors.Errorf("Unsupported action %s", a.Type)
	}

	return nil
}
