package bank

import (
	"github.com/felipefp2/boleto-api/citibank"
	"github.com/felipefp2/boleto-api/models"
)

func getIntegrationCitibank(boleto models.BoletoRequest) (Bank, error) {
	return citibank.New(), nil
}
