package validations

import "github.com/felipefp2/boleto-api/models"

//ValidateMultaDate valida se a data de expiração do boleto não está no passado
func ValidateMultaDate(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsMultaDateValid()
	default:
		return InvalidType(t)
	}
}
