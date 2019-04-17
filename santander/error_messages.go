package santander

import (
	"github.com/felipefp2/boleto-api/models"
)

/*
Este erros segue exatamente a documentação do banco, por isso as mensagens em português
*/
var ticketResponseErrors = map[string]string{
	"-1": "Certificado não encontrado",
	"1":  "Erro, dados de entrada inválidos",
	"2":  "Erro interno de criptografia",
	"3":  "Erro, Ticket já utilizado anteriormente",
	"4":  "Erro, Ticket gerado para outro sistema",
	"5":  "Erro, Ticket expirado",
	"6":  "Erro interno (dados)",
	"7":  "Erro interno (timestamp)",
}

func checkError(code string) error {
	msg, exist := ticketResponseErrors[code]
	if !exist {
		return nil
	}
	return models.NewFormatError(msg)
}
