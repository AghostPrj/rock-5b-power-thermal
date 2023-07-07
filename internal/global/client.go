/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 3:37
 * @Desc:
 */

package global

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RabbitmqClient   *amqp.Connection
	RabbitmqChannels = make([]*amqp.Channel, 0)

	RabbitmqUploadChannel *amqp.Channel

	RabbitmqRoutingKey = ""

	MqttUploadTopic = ""

	MqttClient mqtt.Client
)
