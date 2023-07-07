/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:17
 * @Desc:
 */

package constData

const (
	ApplicationName = "rock-5b-power-thermal"

	/* tpl
	ConfKey = ""
	EnvKey  = ""
	Default = ""
	*/

	ConfDebugFlagKey = "app.debug"
	EnvDebugFlagKey  = "app_debug"
	DefaultDebugFlag = false

	ConfTtyPathKey = "app.tty.path"
	EnvTtyPathKey  = "app_tty_path"
	DefaultTtyPath = "/dev/ttyS7"

	ConfServerListenPortKey = "app.server.listen.port"
	EnvServerListenPortKey  = "app_server_listen_port"
	DefaultServerListenPort = 11099

	ConfServerListenHostKey = "app.server.listen.host"
	EnvServerListenHostKey  = "app_server_listen_host"
	DefaultServerListenHost = ""

	ConfUploadIntervalKey = "app.upload.interval"
	EnvUploadIntervalKey  = "app_upload_interval"
	DefaultUploadInterval = int64(10)

	ConfAppAllowMqttUploadKey = "app.upload.mqtt.allow"
	EnvAppAllowMqttUploadKey  = "app_upload_mqtt_allow"
	DefaultAppAllowMqttUpload = false

	ConfAppAllowUploadRabbitmqKey = "app.upload.rabbitmq.allow"
	EnvAppAllowUploadRabbitmqKey  = "app_upload_rabbitmq_allow"
	DefaultAppAllowUploadRabbitmq = false

	ConfAllowGetNvmeKey = "app.nvme.allow"
	EnvAllowGetNvmeKey  = "app-nvme-allow"
	DefaultAllowGetNvme = true

	ConfNvmePathKey = "app.nvme.path"
	EnvNvmePathKey  = "app-nvme-path"
	DefaultNvmePath = "/dev/nvme0n1"
)
