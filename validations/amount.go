package validations

import "github.com/felipefp2/boleto-api/models"

//ValidateAmount valida o valor do titulo
func ValidateAmount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsAmountInCentsValid()
	default:
		return InvalidType(t)
	}
}
