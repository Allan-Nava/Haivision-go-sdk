package haivision

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validate è l'istanza condivisa del validatore. `validator/v10` è thread-safe e va creato una
// sola volta: mantiene una cache della reflection per tipo.
//
// Sostituisce `gopkg.in/validator.v2`, non manutenuto e con due problemi concreti:
//   - i tag `validate:"nonnil,min=1"` erano su TUTTI i campi, opzionali compresi, quindi una
//     route SRT senza `ttl`/`tos` veniva rifiutata lato client prima di partire;
//   - `min=1` su un `bool` produceva "unsupported type", cioè validazione sempre fallita.
var validate = validator.New(validator.WithRequiredStructEnabled())

// ValidationError segnala che il corpo della richiesta non è valido: la richiesta non è stata
// inviata al gateway.
type ValidationError struct {
	Op  string
	Err error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("haivision: %s: richiesta non valida: %v", e.Op, e.Err)
}

func (e *ValidationError) Unwrap() error { return e.Err }

func validateRequest(op string, v any) error {
	if err := validate.Struct(v); err != nil {
		return &ValidationError{Op: op, Err: err}
	}
	return nil
}
