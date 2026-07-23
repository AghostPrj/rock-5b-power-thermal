/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 3:37
 * @Desc:
 */

package global

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	MqttUploadTopic = ""

	MqttClient mqtt.Client
)
