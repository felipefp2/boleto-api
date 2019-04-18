package validations

import "github.com/felipefp2/boleto-api/models"

//ValidateJuroDate valida se a data de juro do boleto  é maior que a de vencimento
func ValidateJuroDate(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsJuroDateValid()
	default:
		return InvalidType(t)
	}
}
