package bank

import (
	"github.com/felipefp2/boleto-api/itau"
	"github.com/felipefp2/boleto-api/models"
)

func getIntegrationItau(boleto models.BoletoRequest) (Bank, error) {
	return itau.New(), nil
}
