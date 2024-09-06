/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 6:30
 * @Desc:
 */

package initializator

import (
	"fmt"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func initMqtt() {
	global.MqttUploadTopic = viper.GetString(constData.ConfAppUploadMqttTopicPrefixKey) + getIdString()

	log.WithFields(log.Fields{
		"op":            "startup",
		"step":          "init_mqtt_client",
		"publish_topic": global.MqttUploadTopic,
	}).Info()

	opts := mqtt.NewClientOptions()

	random, err := uuid.NewRandom()
	if err != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init_mqtt_client",
			"err":  err,
		}).Panic()
		return
	}

	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",
		viper.GetString(constData.ConfAppMqttHostKey), viper.GetUint(constData.ConfAppMqttPortKey)))
	opts.SetClientID(random.String())
	opts.SetUsername(viper.GetString(constData.ConfAppMqttUserKey))
	opts.SetPassword(viper.GetString(constData.ConfAppMqttPasswordKey))
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(time.Millisecond * 500)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init_mqtt_client",
			"err":  token.Error(),
		}).Panic()
	}

	global.MqttClient = client
}
