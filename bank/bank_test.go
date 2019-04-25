package bank

import (
	"testing"

	"github.com/felipefp2/boleto-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldExecuteBankStrategy(t *testing.T) {
	assert := func(n models.BoletoRequest) {
		bank, err := Get(n)
		number := bank.GetBankNumber()
		So(err, ShouldBeNil)
		So(number.IsBankNumberValid(), ShouldBeTrue)
		So(number, ShouldEqual, n.BankNumber)
	}
	Convey("deve-se verificar o retorno da estrategia de cada banco", t, func() {
		assert(models.BoletoRequest{BankNumber: models.Bradesco, Agreement: models.Agreement{Wallet: 9}})
		assert(models.BoletoRequest{BankNumber: models.Bradesco, Agreement: models.Agreement{Wallet: 25}})
		assert(models.BoletoRequest{BankNumber: models.BancoDoBrasil})
		assert(models.BoletoRequest{BankNumber: models.Citibank})
		assert(models.BoletoRequest{BankNumber: models.Santander})
		assert(models.BoletoRequest{BankNumber: models.Itau})
		assert(models.BoletoRequest{BankNumber: models.Caixa})

		_, err := Get(models.BoletoRequest{BankNumber: 88})
		So(err, ShouldNotBeNil)
	})

}
