package models

import (
	"fmt"
	"regexp"

	"github.com/felipefp2/boleto-api/util"
)

// Agreement afiliação do cliente com o bano
type Agreement struct {
	AgreementNumber uint   `json:"agreementNumber,omitempty"`
	Wallet          uint16 `json:"wallet,omitempty"`
	WalletVariation uint16 `json:"walletVariation,omitempty"`
	Agency          string `json:"agency,omitempty"`
	AgencyDigit     string `json:"agencyDigit,omitempty"`
	Account         string `json:"account,omitempty"`
	AccountDigit    string `json:"accountDigit,omitempty"`
}

// IsAgencyValid retorna se é uma agência válida
func (a *Agreement) IsAgencyValid() error {
	re := regexp.MustCompile("(\\D+)")
	ag := re.ReplaceAllString(a.Agency, "")
	if len(ag) < 5 && len(ag) > 0 {
		a.Agency = util.PadLeft(ag, "0", 4)
		return nil
	}
	return NewErrorResponse("MPAgency", "Agência inválida, deve conter até 4 dígitos")
}

// CalculateAgencyDigit calcula dígito da agência
func (a *Agreement) CalculateAgencyDigit(digitCalculator func(agency string) string) {
	re := regexp.MustCompile("(\\D+)")
	ad := re.ReplaceAllString(a.AgencyDigit, "")
	if len(ad) == 1 {
		a.AgencyDigit = ad
	} else {
		a.AgencyDigit = digitCalculator(a.Agency)
	}
}

// IsAccountValid retorna se é uma conta válida
func (a *Agreement) IsAccountValid(accountLength int) error {
	re := regexp.MustCompile("(\\D+)")
	ac := re.ReplaceAllString(a.Account, "")
	if len(ac) < accountLength+1 && len(ac) > 0 {
		a.Account = util.PadLeft(ac, "0", uint(accountLength))
		return nil
	}
	return NewErrorResponse("MPAccount", fmt.Sprintf("Conta inválida, deve conter até %d dígitos", accountLength))
}

//CalculateAccountDigit calcula dígito da conta
func (a *Agreement) CalculateAccountDigit(digitCalculator func(agency, account string) string) {
	re := regexp.MustCompile("(\\D+)")
	ad := re.ReplaceAllString(a.AccountDigit, "")
	if len(ad) == 1 {
		a.AccountDigit = ad
	} else {
		a.AccountDigit = digitCalculator(a.Agency, a.Account)
	}
}
