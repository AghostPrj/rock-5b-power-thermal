/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 6:31
 * @Desc:
 */

package constData

const (
	ConfAppMqttHostKey = "app.mqtt.host"
	EnvAppMqttHostKey  = "app_mqtt_host"
	DefaultAppMqttHost = "127.0.0.1"

	ConfAppMqttPortKey = "app.mqtt.port"
	EnvAppMqttPortKey  = "app_mqtt_port"
	DefaultAppMqttPort = "1883"

	ConfAppMqttUserKey = "app.mqtt.user"
	EnvAppMqttUserKey  = "app_mqtt_user"
	DefaultAppMqttUser = "user"

	ConfAppMqttPasswordKey = "app.mqtt.password"
	EnvAppMqttPasswordKey  = "app_mqtt_password"
	DefaultAppMqttPassword = "password"

	ConfAppUploadMqttTopicPrefixKey = "app.upload.mqtt.topic_prefix"
	EnvAppUploadMqttTopicPrefixKey  = "app_upload_mqtt_topic_prefix"
	DefaultAppUploadMqttTopicPrefix = "topic/mqtt-ha/rock-5b/status-"
)
