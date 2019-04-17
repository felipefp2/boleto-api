package bank

import (
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/pefisa"
)

func getIntegrationPefisa(boleto models.BoletoRequest) (Bank, error) {
	return pefisa.New(), nil
}
