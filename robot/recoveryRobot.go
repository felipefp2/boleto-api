package robot

import (
	"strconv"

	"github.com/jasonlvhit/gocron"
	"github.com/felipefp2/boleto-api/config"
	"github.com/felipefp2/boleto-api/db"
	"github.com/felipefp2/boleto-api/log"
	"github.com/felipefp2/boleto-api/util"
)

//RecoveryRobot robô que faz a resiliência de boletos
func RecoveryRobot(ex string) {

	if ex == "true" {
		go func() {
			e, _ := strconv.ParseUint(config.Get().RecoveryRobotExecutionInMinutes, 10, 64)
			gocron.Every(e).Minutes().Do(executionTask)
			<-gocron.Start()
		}()
	}

}

func executionTask() {

	lg := log.CreateLog()
	lg.Operation = "RecoveryRobot"

	lg.InitRobot()

	redis := db.CreateRedis()
	keys, _ := redis.GetAllJSON()

	mongo, errMongo := db.CreateMongo(lg)
	if util.CheckErrorRobot(errMongo) == false {
		for _, key := range keys {
			bol, errRedis := redis.GetBoletoJSONByKey(string(key), lg)
			if util.CheckErrorRobot(errRedis) == false {
				err := mongo.SaveBoleto(bol)

				if util.CheckErrorRobot(err) == false {
					lg.ResumeRobot(string(key))
					redis.DeleteBoletoJSONByKey(string(key), lg)
				}
			}
		}
	}

	lg.EndRobot()
}
