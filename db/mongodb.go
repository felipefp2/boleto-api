package db

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/felipefp2/boleto-api/config"
	"github.com/felipefp2/boleto-api/log"
	"github.com/felipefp2/boleto-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoDb Struct
type MongoDb struct {
	m sync.RWMutex
}

var dbName = "Boleto"

var (
	dbSession *mgo.Session
	err       error
)

//CreateMongo cria uma nova intancia de conexão com o mongodb
func CreateMongo(l *log.Log) (*MongoDb, error) {

	if dbSession == nil {
		dbSession, err = mgo.DialWithInfo(getInfo())

		if err != nil {
			l.Warn(err, fmt.Sprintf("Error create connection mongo %s", err.Error()))
			return nil, err
		}
	}

	db := new(MongoDb)

	return db, nil
}

func getInfo() *mgo.DialInfo {
	connMgo := strings.Split(config.Get().MongoURL, ",")
	return &mgo.DialInfo{
		Addrs:     connMgo,
		Timeout:   5 * time.Second,
		Database:  "Boleto",
		PoolLimit: 512,
		Username:  config.Get().MongoUser,
		Password:  config.Get().MongoPassword,
	}
}

//SaveBoleto salva um boleto no mongoDB
func (e *MongoDb) SaveBoleto(boleto models.BoletoView) error {

	e.m.Lock()
	defer e.m.Unlock()

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(dbName).C("boletos")
	err = c.Insert(boleto)
	return err
}

//GetBoletoByID busca um boleto pelo ID que vem na URL
func (e *MongoDb) GetBoletoByID(id, pk string) (models.BoletoView, error) {

	e.m.Lock()
	defer e.m.Unlock()
	result := models.BoletoView{}

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(dbName).C("boletos")

	if len(id) == 24 {
		d := bson.ObjectIdHex(id)
		err = c.Find(bson.M{"_id": d}).One(&result)
	} else {
		err = c.Find(bson.M{"id": id}).One(&result)
	}

	if err != nil || !hasValidKey(result, pk) {
		return models.BoletoView{}, models.NewHTTPNotFound("MP404", "Not Found")
	}

	return result, nil
}

//UpdateBoleto altera um boleto pelo ID
func (e *MongoDb) UpdateBoleto(id string, boleto models.BoletoView) error {

	e.m.Lock()
	defer e.m.Unlock()

	session := dbSession.Copy()

	defer session.Close()

	c := session.DB(dbName).C("boletos")

	selector := bson.M{"id": id}

	if len(id) == 24 {
		d := bson.ObjectIdHex(id)
		selector = bson.M{"_id": d}
	}

	updator := bson.M{
		"$set": bson.M{
			"boleto.Title.ExpireDate":        boleto.Boleto.Title.ExpireDate,
			"boleto.Title.ExpireDateTime":    boleto.Boleto.Title.ExpireDateTime,
			"boleto.Title.JuroDate":          boleto.Boleto.Title.JuroDate,
			"boleto.Title.JuroDateTime":      boleto.Boleto.Title.JuroDateTime,
			"boleto.Title.JuroInCents":       boleto.Boleto.Title.JuroInCents,
			"boleto.Title.JuroInPercentual":  boleto.Boleto.Title.JuroInPercentual,
			"boleto.Title.MultaDate":         boleto.Boleto.Title.MultaDate,
			"boleto.Title.MultaDateTime":     boleto.Boleto.Title.MultaDateTime,
			"boleto.Title.MultaInCents":      boleto.Boleto.Title.MultaInCents,
			"boleto.Title.MultaInPercentual": boleto.Boleto.Title.MultaInPercentual,
			"boleto.Title.AmountInCents":     boleto.Boleto.Title.AmountInCents,
		},
	}

	err = c.Update(selector, updator)

	if err != nil {
		return models.NewHTTPNotFound("MP404", "Not Found")
	}

	return nil
}

//Close Fecha a conexão
func (e *MongoDb) Close() {
	fmt.Println("Close Database Connection")
}

func hasValidKey(r models.BoletoView, pk string) bool {
	return (r.SecretKey == "" || r.PublicKey == pk)
}
