package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/felipefp2/boleto-api/util"
)

// Title título de cobrança de entrada
type Title struct {
	CreateDate           time.Time `json:"createDate,omitempty"`
	ExpireDateTime       time.Time `json:"expireDateTime,omitempty"`
	ExpireDate           string    `json:"expireDate,omitempty"`
	AmountInCents        uint64    `json:"amountInCents,omitempty"`
	OurNumber            uint      `json:"ourNumber,omitempty"`
	Instructions         string    `json:"instructions,omitempty"`
	DocumentNumber       string    `json:"documentNumber,omitempty"`
	NSU                  string    `json:"nsu,omitempty"`
	JuroDate             string    `json:"juroDate,omitempty"`
	JuroDateTime         time.Time `json:"JuroDateTime,omitempty"`
	JuroInCents          uint64    `json:"juroInCents,omitempty"`
	JuroInPercentual     float64   `json:"juroInPercentual,omitempty"`
	MultaDate            string    `json:"multaDate,omitempty"`
	MultaDateTime        time.Time `json:"MultaDateTime,omitempty"`
	MultaInCents         uint64    `json:"multaInCents,omitempty"`
	MultaInPercentual    float64   `json:"multaInPercentual,omitempty"`
	DescontoDate         string    `json:"descontoDate,omitempty"`
	DescontoDateTime     time.Time `json:"descontoDateTime,omitempty"`
	DescontoInCents      uint64    `json:"descontoInCents,omitempty"`
	DescontoInPercentual float64   `json:"descontoInPercentual,omitempty"`
	BoletoType           string    `json:"boletoType,omitempty"`
	BoletoTypeCode       string
}

//ValidateInstructionsLength valida se texto das instruções possui quantidade de caracteres corretos
func (t Title) ValidateInstructionsLength(max int) error {
	if len(t.Instructions) > max {
		return NewErrorResponse("MPInstructions", fmt.Sprintf("Instruções não podem passar de %d caracteres", max))
	}
	return nil
}

//ValidateDocumentNumber número do documento
func (t *Title) ValidateDocumentNumber() error {
	re := regexp.MustCompile("(\\D+)")
	ad := re.ReplaceAllString(t.DocumentNumber, "")
	if ad == "" {
		t.DocumentNumber = ad
	} else if len(ad) < 10 {
		t.DocumentNumber = util.PadLeft(ad, "0", 10)
	} else {
		t.DocumentNumber = ad[:10]
	}
	return nil
}

//IsExpireDateValid retorna um erro se a data de expiração for inválida
func (t *Title) IsExpireDateValid() error {
	d, err := parseDate(t.ExpireDate)
	if err != nil {
		return NewErrorResponse("MPExpireDate", fmt.Sprintf("Data em um formato inválido, esperamos AAAA-MM-DD e recebemos %s", t.ExpireDate))
	}
	n, _ := parseDate(util.BrNow().Format("2006-01-02"))
	t.CreateDate = n
	t.ExpireDateTime = d
	if t.CreateDate.After(t.ExpireDateTime) {
		return NewErrorResponse("MPExpireDate", "Data de expiração não pode ser menor que a data de hoje")
	}
	return nil
}

//IsMultaDateValid retorna um erro se a data de multa for inválida
func (t *Title) IsMultaDateValid() error {
	if t.MultaDate == "" {
		return nil
	}
	d, err := parseDate(t.MultaDate)
	if err != nil {
		return NewErrorResponse("MPMultaDate", fmt.Sprintf("Data em um formato inválido, esperamos AAAA-MM-DD e recebemos %s", t.ExpireDate))
	}
	t.MultaDateTime = d
	if t.ExpireDateTime == t.MultaDateTime || t.ExpireDateTime.After(t.MultaDateTime) {
		return NewErrorResponse("MPMultaDate", "Data de multa não pode ser menor/igual que a data de vencimento")
	}
	return nil
}

//IsJuroDateValid retorna um erro se a data de juro for inválida
func (t *Title) IsJuroDateValid() error {
	if t.JuroDate == "" {
		return nil
	}
	d, err := parseDate(t.JuroDate)
	if err != nil {
		return NewErrorResponse("MPJuroDate", fmt.Sprintf("Data em um formato inválido, esperamos AAAA-MM-DD e recebemos %s", t.ExpireDate))
	}
	t.JuroDateTime = d
	if t.ExpireDateTime == t.JuroDateTime || t.ExpireDateTime.After(t.JuroDateTime) {
		return NewErrorResponse("MPJuroDate", "Data de juro não pode ser menor/igual que a data de vencimento")
	}
	return nil
}

//IsDescontoDateValid retorna um erro se a data de desconto for inválida
func (t *Title) IsDescontoDateValid() error {
	if t.DescontoDate == "" {
		return nil
	}
	d, err := parseDate(t.DescontoDate)
	if err != nil {
		return NewErrorResponse("MPDescontoDate", fmt.Sprintf("Data em um formato inválido, esperamos AAAA-MM-DD e recebemos %s", t.ExpireDate))
	}
	t.DescontoDateTime = d
	if t.ExpireDateTime == t.DescontoDateTime || t.DescontoDateTime.After(t.ExpireDateTime) {
		return NewErrorResponse("MPDescontoDate", "Data de desconto não pode ser maior/igual que a data de vencimento")
	}
	return nil
}

//IsAmountInCentsValid retorna um erro se o valor em centavos for inválido
func (t *Title) IsAmountInCentsValid() error {
	if t.AmountInCents < 1 {
		return NewErrorResponse("MPAmountInCents", "Valor não pode ser menor do que 1 centavo")
	}
	return nil
}

func parseDate(t string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", t)
	if err != nil {
		return time.Now(), err
	}
	return date, nil
}
