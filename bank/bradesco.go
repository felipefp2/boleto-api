package bank

import (
	"fmt"

	"github.com/felipefp2/boleto-api/bradescoNetEmpresa"
	"github.com/felipefp2/boleto-api/bradescoShopFacil"
	"github.com/felipefp2/boleto-api/models"
)

//Get retorna estrategia de acordo com a carteira ou erro caso o banco não exista
func getIntegrationBradesco(boleto models.BoletoRequest) (Bank, error) {
	switch boleto.Agreement.Wallet {
	case 4, 9, 19:
		return bradescoNetEmpresa.New(), nil
	case 25, 26:
		return bradescoShopFacil.New(), nil
	default:
		return nil, models.NewErrorResponse("MPWallet", fmt.Sprintf("Carteira %d não existe", boleto.Agreement.Wallet))
	}
}
