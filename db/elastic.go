package db

import (
	"net/http"
	"strings"

	"fmt"

	"io/ioutil"

	"encoding/json"

	"github.com/felipefp2/boleto-api/config"
	"github.com/felipefp2/boleto-api/models"
	"github.com/felipefp2/boleto-api/util"
)

type elasticDb struct{}

//SaveBoleto salva um boleto no elasticsearch
func (e *elasticDb) SaveBoleto(boleto models.BoletoView) error {
	client := util.DefaultHTTPClient()
	url := fmt.Sprintf("%s/boletoapi/boleto/%s", config.Get().ElasticURL, boleto.ID)
	req, err := http.NewRequest("POST", url, strings.NewReader(boleto.ToJSON()))
	if err != nil {
		return err
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return errResp
	}
	defer req.Body.Close()
	defer resp.Body.Close()
	return nil
}

//GetBoletoById busca um boleto pelo ID que vem na URL
func (e *elasticDb) GetBoletoByID(id string) (models.BoletoView, error) {
	client := util.DefaultHTTPClient()
	url := fmt.Sprintf("%s/boletoapi/boleto/%s", config.Get().ElasticURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.BoletoView{}, err
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return models.BoletoView{}, errResp
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	elasticData := struct {
		Source models.BoletoView `json:"_source"`
	}{}
	json.Unmarshal(data, &elasticData)
	return elasticData.Source, nil
}

func (e *elasticDb) Close() {}
