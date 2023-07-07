/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2022/7/6 2:55
 * @Desc:
 */

package constData

import "time"

const (
	ConfRabbitmqHostKey = "app.rabbitmq.host"
	EnvRabbitmqHostKey  = "app_rabbitmq_host"
	DefaultRabbitmqHost = "127.0.0.1"

	ConfRabbitmqPortKey = "app.rabbitmq.port"
	EnvRabbitmqPortKey  = "app_rabbitmq_port"
	DefaultRabbitmqPort = "5672"

	ConfRabbitmqUserKey = "app.rabbitmq.user"
	EnvRabbitmqUserKey  = "app_rabbitmq_user"
	DefaultRabbitmqUser = "guest"

	ConfRabbitmqPasswordKey = "app.rabbitmq.password"
	EnvRabbitmqPasswordKey  = "app_rabbitmq_password"
	DefaultRabbitmqPassword = "guest"

	ConfRabbitmqVirtualHostKey = "app.rabbitmq.virtual_host"
	EnvRabbitmqVirtualHostKey  = "app_rabbitmq_virtual_host"
	DefaultRabbitmqVirtualHost = "/"

	ConfRabbitmqChannelMaxKey = "app.rabbitmq.channel_max"
	EnvRabbitmqChannelMaxKey  = "app_rabbitmq_channel_max"
	DefaultRabbitmqChannelMax = 0

	ConfRabbitmqFrameSizeKey = "app.rabbitmq.frame_size"
	EnvRabbitmqFrameSizeKey  = "app_rabbitmq_frame_size"
	DefaultRabbitmqFrameSize = 0

	ConfRabbitmqHeartBeatKey = "app.rabbitmq.heart_beat"
	EnvRabbitmqHeartBeatKey  = "app_rabbitmq_heart_beat"
	DefaultRabbitmqHeartBeat = time.Second

	ConfAppRabbitmqUploadExchangeKey = "app.upload.rabbitmq.exchange"
	EnvAppRabbitmqUploadExchangeKey  = "app_upload_rabbitmq_exchange"
	DefaultAppRabbitmqUploadExchange = "rock_5b_status_exchange"

	ConfAppRabbitmqUploadQueuePrefixKey = "app.upload.rabbitmq.queue_prefix"
	EnvAppRabbitmqUploadQueuePrefixKey  = "app_upload_rabbitmq_queue_prefix"
	DefaultAppRabbitmqUploadQueuePrefix = "rock_5b_status_queue-"

	ConfAppRabbitmqUploadRoutingPrefixKey = "app.upload.rabbitmq.routing_prefix"
	EnvAppRabbitmqUploadRoutingPrefixKey  = "app_upload_rabbitmq_routing_prefix"
	DefaultAppRabbitmqUploadRoutingPrefix = "rock_5b_status_routing-"
)
