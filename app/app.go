package app

import (
	"fmt"
	"os"

	"github.com/PMoneda/flow"
	"github.com/felipefp2/boleto-api/api"
	"github.com/felipefp2/boleto-api/config"
	"github.com/felipefp2/boleto-api/env"
	"github.com/felipefp2/boleto-api/log"
	"github.com/felipefp2/boleto-api/mock"
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/robot"
	"github.com/felipefp2/boleto-api/util"
)

//Params this struct contains all execution parameters to run application
type Params struct {
	DevMode    bool
	MockMode   bool
	DisableLog bool
}

//NewParams returns new Empty pointer to ExecutionParameters
func NewParams() *Params {
	return new(Params)
}

//Run starts boleto api Application
func Run(params *Params) {
	env.Config(params.DevMode, params.MockMode, params.DisableLog)

	if config.Get().MockMode {
		go mock.Run("9091")
	}

	installLog()

	installCertificates()

	go robot.RecoveryRobot(config.Get().RecoveryRobotExecutionEnabled)

	api.InstallRestAPI()

}

func installLog() {
	err := log.Install()
	if err != nil {
		fmt.Println("Log SEQ Fails")
		os.Exit(-1)
	}
}

func installCertificates() {
	if config.Get().DevMode == false {
		err := util.ListCert()
		if err != nil {
			fmt.Println("Copy Certificates Fails")
			os.Exit(-1)
		}
	}
}

func installflowConnectors() {
	flow.RegisterConnector("logseq", util.SeqLogConector)
	flow.RegisterConnector("apierro", models.BoletoErrorConector)
	flow.RegisterConnector("tls", util.TlsConector)
}
