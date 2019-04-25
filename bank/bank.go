package bank

import (
	"fmt"

	"github.com/felipefp2/boleto-api/log"
	"github.com/felipefp2/boleto-api/models"
)

//Bank é a interface que vai oferecer os serviços em comum entre os bancos
type Bank interface {
	ProcessBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	RegisterBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
	ValidateBoleto(*models.BoletoRequest) models.Errors
	GetBankNumber() models.BankNumber
	GetBankNameIntegration() string
	Log() *log.Log
	ProcessBoletoForEdit(*models.BoletoRequest) (models.BoletoResponse, error)
	EditBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
}

// //EditOption é a interface que vai oferecer os serviços de alteração dos boletos
// type EditOption interface {
// 	ProcessBoletoForEdit(*models.BoletoRequest) (models.BoletoResponse, error)
// 	EditBoleto(*models.BoletoRequest) (models.BoletoResponse, error)
// }

// //ProcessBoletoForEdit processo o boleto para alteração de dados
// func ProcessBoletoForEdit(bank Bank) {
// 	if bankWithEdit, ok := bank.(EditOption); ok {
// 		bankWithEdit.ProcessBoletoForEdit()
// 	}
// }

//Get retorna estrategia de acordo com o banco ou erro caso o banco não exista
func Get(boleto models.BoletoRequest) (Bank, error) {
	switch boleto.BankNumber {
	case models.BancoDoBrasil:
		return getIntegrationBB(boleto)
	case models.Bradesco:
		return getIntegrationBradesco(boleto)
	case models.Caixa:
		return getIntegrationCaixa(boleto)
	case models.Citibank:
		return getIntegrationCitibank(boleto)
	case models.Santander:
		return getIntegrationSantander(boleto)
	case models.Itau:
		return getIntegrationItau(boleto)
	case models.Pefisa:
		return getIntegrationPefisa(boleto)
	default:
		return nil, models.NewErrorResponse("MPBankNumber", fmt.Sprintf("Banco %d não existe", boleto.BankNumber))
	}
}
