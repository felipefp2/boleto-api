package citibank

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PMoneda/flow"
	"github.com/felipefp2/boleto-api/config"
	"github.com/felipefp2/boleto-api/log"
	"github.com/felipefp2/boleto-api/metrics"
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/tmpl"
	"github.com/felipefp2/boleto-api/util"
	"github.com/felipefp2/boleto-api/validations"
)

var o = &sync.Once{}
var m map[string]string

type bankCiti struct {
	validate  *models.Validator
	log       *log.Log
	transport *http.Transport
}

func New() bankCiti {
	b := bankCiti{
		validate: models.NewValidator(),
		log:      log.CreateLog(),
	}
	b.validate.Push(validations.ValidateAmount)
	b.validate.Push(validations.ValidateExpireDate)
	b.validate.Push(validations.ValidateBuyerDocumentNumber)
	b.validate.Push(validations.ValidateRecipientDocumentNumber)
	b.validate.Push(citiValidateAgency)
	b.validate.Push(citiValidateAccount)
	b.validate.Push(citiValidateWallet)
	b.validate.Push(citiValidateBoletoType)
	transp, err := util.BuildTLSTransport(config.Get().CertBoletoPathCrt, config.Get().CertBoletoPathKey, config.Get().CertBoletoPathCa)
	if err != nil {
		//TODO
	}
	b.transport = transp
	return b
}

//Log retorna a referencia do log
func (b bankCiti) Log() *log.Log {
	return b.log
}

func (b bankCiti) RegisterBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {

	boleto.Title.BoletoType, boleto.Title.BoletoTypeCode = getBoletoType(boleto)

	boleto.Title.OurNumber = calculateOurNumber(boleto)
	r := flow.NewFlow()
	serviceURL := config.Get().URLCiti
	from := getResponseCiti()
	to := getAPIResponseCiti()
	bod := r.From("message://?source=inline", boleto, getRequestCiti(), tmpl.GetFuncMaps())
	bod.To("logseq://?type=request&url="+serviceURL, b.log)
	//TODO: change for tls flow connector (waiting for santander)
	var responseCiti string
	var status int
	var err error
	duration := util.Duration(func() {
		responseCiti, status, err = b.sendRequest(bod.GetBody().(string))
	})
	if err != nil {
		return models.BoletoResponse{}, err
	}
	metrics.PushTimingMetric("citibank-register-boleto-online", duration.Seconds())
	bod.To("set://?prop=header", map[string]string{"status": strconv.Itoa(status)})
	bod.To("set://?prop=body", responseCiti)
	bod.To("logseq://?type=response&url="+serviceURL, b.log)
	ch := bod.Choice()
	ch.When(flow.Header("status").IsEqualTo("200"))
	ch.To("transform://?format=xml", from, to, tmpl.GetFuncMaps())
	ch.Otherwise()
	ch.To("logseq://?type=response&url="+serviceURL, b.log).To("apierro://")

	switch t := bod.GetBody().(type) {
	case string:
		response := util.ParseJSON(t, new(models.BoletoResponse)).(*models.BoletoResponse)
		return *response, nil
	case models.BoletoResponse:
		return t, nil
	}
	return models.BoletoResponse{}, models.NewInternalServerError("MP500", "Internal error")
}

func (b bankCiti) ProcessBoleto(boleto *models.BoletoRequest) (models.BoletoResponse, error) {
	errs := b.ValidateBoleto(boleto)
	if len(errs) > 0 {
		return models.BoletoResponse{Errors: errs}, nil
	}
	return b.RegisterBoleto(boleto)
}

func (b bankCiti) ValidateBoleto(boleto *models.BoletoRequest) models.Errors {
	return models.Errors(b.validate.Assert(boleto))
}

func (b bankCiti) sendRequest(body string) (string, int, error) {
	serviceURL := config.Get().URLCiti
	if config.Get().MockMode {
		return util.Post(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"})
	} else {
		return util.PostTLS(serviceURL, body, config.Get().TimeoutDefault, map[string]string{"Soapaction": "RegisterBoleto"}, b.transport)
	}
}

//GetBankNumber retorna o codigo do banco
func (b bankCiti) GetBankNumber() models.BankNumber {
	return models.Citibank
}

func calculateOurNumber(boleto *models.BoletoRequest) uint {
	ourNumberWithDigit := strconv.Itoa(int(boleto.Title.OurNumber)) + util.OurNumberDv(strconv.Itoa(int(boleto.Title.OurNumber)), util.MOD11)
	value, _ := strconv.Atoi(ourNumberWithDigit)
	return uint(value)
}

func (b bankCiti) GetBankNameIntegration() string {
	return "Citibank"
}

func citiBoletoTypes() map[string]string {
	o.Do(func() {
		m = make(map[string]string)

		m["DMI"] = "03" //Duplicata Mercantil p/ Indicação
	})

	return m
}

func getBoletoType(boleto *models.BoletoRequest) (bt string, btc string) {
	if len(boleto.Title.BoletoType) < 1 {
		return "DMI", "03"
	}
	btm := citiBoletoTypes()

	if btm[strings.ToUpper(boleto.Title.BoletoType)] == "" {
		return "DMI", "03"
	}

	return boleto.Title.BoletoType, btm[strings.ToUpper(boleto.Title.BoletoType)]
}
