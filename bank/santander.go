package bank

import (
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/santander"
)

func getIntegrationSantander(boleto models.BoletoRequest) (Bank, error) {
	return santander.New(), nil
}
