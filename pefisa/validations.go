package pefisa

import (
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/validations"
)

func pefisaBoletoTypeValidate(b interface{}) error {
	bt := pefisaBoletoTypes()

	switch t := b.(type) {

	case *models.BoletoRequest:
		if len(t.Title.BoletoType) > 0 && bt[t.Title.BoletoType] == "" {
			return models.NewErrorResponse("MP400", "espécie de boleto informada não existente")
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}
