package itau

import (
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/validations"
)

func itauValidateAccount(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAccountValid(7)
		if err != nil {
			return err
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func itauValidateAgency(b interface{}) error {
	switch t := b.(type) {
	case *models.BoletoRequest:
		err := t.Agreement.IsAgencyValid()
		if err != nil {
			return models.NewErrorResponse("MP400", err.Error())
		}
		return nil
	default:
		return validations.InvalidType(t)
	}
}

func itauBoletoTypeValidate(b interface{}) error {
	bt := itauBoletoTypes()

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
