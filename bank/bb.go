package bank

import (
	"github.com/felipefp2/boleto-api/bb"
	"github.com/felipefp2/boleto-api/models"
)

func getIntegrationBB(boleto models.BoletoRequest) (Bank, error) {
	return bb.New(), nil
}