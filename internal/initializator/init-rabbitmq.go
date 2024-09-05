/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2022/7/6 3:07
 * @Desc:
 */

package initializator

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"time"
)

func initRabbitmqClient() {
	dsn := "amqp://" + url.QueryEscape(viper.GetString(constData.ConfRabbitmqUserKey)) + ":" +
		url.QueryEscape(viper.GetString(constData.ConfRabbitmqPasswordKey)) +
		"@" + url.QueryEscape(viper.GetString(constData.ConfRabbitmqHostKey)) + ":" +
		viper.GetString(constData.ConfRabbitmqPortKey) + "/"
	config := amqp.Config{
		Vhost:      viper.GetString(constData.ConfRabbitmqVirtualHostKey),
		ChannelMax: viper.GetUint16(constData.ConfRabbitmqChannelMaxKey),
		FrameSize:  viper.GetInt(constData.ConfRabbitmqFrameSizeKey),
		Heartbeat:  viper.GetDuration(constData.ConfRabbitmqHeartBeatKey),
	}

	log.WithFields(log.Fields{
		constData.ConfRabbitmqUserKey:        viper.GetString(constData.ConfRabbitmqUserKey),
		constData.ConfRabbitmqPasswordKey:    viper.GetString(constData.ConfRabbitmqPasswordKey),
		constData.ConfRabbitmqHostKey:        viper.GetString(constData.ConfRabbitmqHostKey),
		constData.ConfRabbitmqPortKey:        viper.GetString(constData.ConfRabbitmqPortKey),
		constData.ConfRabbitmqVirtualHostKey: viper.GetString(constData.ConfRabbitmqVirtualHostKey),
		constData.ConfRabbitmqChannelMaxKey:  viper.GetInt(constData.ConfRabbitmqChannelMaxKey),
		constData.ConfRabbitmqFrameSizeKey:   viper.GetInt(constData.ConfRabbitmqFrameSizeKey),
		constData.ConfRabbitmqHeartBeatKey:   viper.GetDuration(constData.ConfRabbitmqHeartBeatKey),
	}).Trace("rabbit client info")

	var amqpClient *amqp.Connection

	for {
		count := 0
		tmpAmqpClient, err := amqp.DialConfig(dsn, config)
		if err != nil {
			count++
		} else {
			amqpClient = tmpAmqpClient
			break
		}

		if count > 10 {
			log.WithFields(log.Fields{
				"op":  "startup",
				"err": err,
			}).Panic("init rabbitmq client error")
		}

		time.Sleep(time.Millisecond * 50)
	}

	ch, err := amqpClient.Channel()
	if err != nil {
		log.WithFields(log.Fields{
			"op":  "startup",
			"err": err,
		}).Panic("get channel error")
	}
	defer func() {
		_ = ch.Close()
	}()

	global.RabbitmqClient = amqpClient

	initRabbitmqExchangeAndQueue()

	global.RabbitmqUploadChannel = initRabbitmqProviderChannel()

	go handlerRabbitmqReconnect()
}

func handlerRabbitmqReconnect() {
	errChan := make(chan *amqp.Error)
	_ = global.RabbitmqClient.NotifyClose(errChan)

	err, ok := <-errChan
	if ok {
		if err == nil {
			return
		} else {
			log.WithFields(log.Fields{
				"op":  "rabbitmq close handler",
				"err": err,
			}).Error()

			for {
				closed := 0
				for i := 0; i < len(global.RabbitmqChannels); i++ {
					_ = global.RabbitmqChannels[i].Close()
				}
				for i := 0; i < len(global.RabbitmqChannels); i++ {
					if global.RabbitmqChannels[i].IsClosed() {
						closed++
					}
				}
				if closed == len(global.RabbitmqChannels) {
					break
				}
			}

			_ = global.RabbitmqClient.Close()
			initRabbitmqClient()
			return
		}
	}

}

func initRabbitmqProviderChannel() *amqp.Channel {
	channel, err := global.RabbitmqClient.Channel()
	if err != nil {
		log.WithFields(log.Fields{
			"op":  "startup",
			"err": err,
		}).Panic("get Provider channel error")
	}

	global.RabbitmqChannels = append(global.RabbitmqChannels, channel)
	return channel

}

func initRabbitmqExchangeAndQueue() {
	ch, err := global.RabbitmqClient.Channel()
	if err != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init exchange and queue",
			"err":  err,
		}).Panic("get channel error")
	}
	defer func() {
		_ = ch.Close()
	}()

	err = ch.ExchangeDeclare(viper.GetString(constData.ConfAppRabbitmqUploadExchangeKey), "direct",
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init exchange and queue",
			"err":  err,
		}).Panic("create exchange error")
	}

	cpuSerialNum := getCpuSerialNum()

	queueName := viper.GetString(constData.ConfAppRabbitmqUploadQueuePrefixKey) + cpuSerialNum
	log.WithFields(log.Fields{
		"op":         "startup",
		"step":       "init exchange and queue",
		"queue_name": queueName,
	}).Info()

	_, err = ch.QueueDeclare(queueName,
		true, false, false, false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init exchange and queue",
			"err":  err,
		}).Panic("create queue " + queueName + " error")
	}

	routingKey := viper.GetString(constData.ConfAppRabbitmqUploadRoutingPrefixKey) + cpuSerialNum

	err = ch.QueueBind(queueName, routingKey,
		viper.GetString(constData.ConfAppRabbitmqUploadExchangeKey), false, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"op":   "startup",
			"step": "init exchange and queue",
			"err":  err,
		}).Panic("bind queue  to " + viper.GetString(constData.ConfAppRabbitmqUploadExchangeKey) + " error")
	}

	global.RabbitmqRoutingKey = routingKey

}
