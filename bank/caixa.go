package bank

import (
	"github.com/felipefp2/boleto-api/caixa"
	"github.com/felipefp2/boleto-api/models"
)

func getIntegrationCaixa(boleto models.BoletoRequest) (Bank, error) {
	return caixa.New(), nil
}
