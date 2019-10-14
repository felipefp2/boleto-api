package validations

import "github.com/felipefp2/boleto-api/models"

//ValidateMultaDate valida se a data de multa do boleto é maior que a de vencimento
func ValidateDescontoDate(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsDescontoDateValid()
	default:
		return InvalidType(t)
	}
}
