/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/7 21:21
 * @Desc:
 */

package main

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/initializator"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/router"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/task"
	"github.com/md14454/gosensors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"time"
)

func main() {

	gosensors.Init()
	initializator.InitApp()

	chanRcvSerialData := make(chan string)
	chanSingleData := make(chan string)

	go task.GetDataFromSerial(viper.GetString(constData.ConfTtyPathKey), chanRcvSerialData)
	go task.ProcessRawSerialData(chanRcvSerialData, chanSingleData)
	go task.CacheData(chanSingleData)
	go task.GetSystemData()

	time.Sleep(time.Second * 5)

	if viper.GetBool(constData.ConfAppAllowUploadRabbitmqKey) || viper.GetBool(constData.ConfAppAllowMqttUploadKey) {
		go task.UploadData()
	}

	go func() {
		for {
			time.Sleep(time.Second)
			log.WithField("op", "timer").
				WithField("runtimes", runtime.NumGoroutine()).
				WithField("cGoCall", runtime.NumCgoCall()).
				Trace("system running")
		}
	}()

	ginRouter := router.BuildGinRouter()
	log.WithFields(log.Fields{
		"op":   "startup",
		"host": viper.GetString(constData.ConfServerListenHostKey),
		"port": viper.GetString(constData.ConfServerListenPortKey),
	}).Info("http server start")

	router.StartGinServer(ginRouter)

}
