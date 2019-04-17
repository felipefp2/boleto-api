package validations

import "github.com/felipefp2/boleto-api/models"

//ValidateExpireDate valida se a data de expiração do boleto não está no passado
func ValidateExpireDate(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		return t.Title.IsExpireDateValid()
	default:
		return InvalidType(t)
	}
}
