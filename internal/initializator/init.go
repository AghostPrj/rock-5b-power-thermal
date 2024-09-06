/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:17
 * @Desc:
 */

package initializator

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/ggg17226/aghost-go-base/pkg/utils/configUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"os/exec"
	"strings"
)

func InitApp() {
	configUtils.SetConfigFileName(constData.ApplicationName)
	bindApiAppConfigKey()
	bindApiAppConfigDefaultValue()
	configUtils.InitConfigAndLog()

	if viper.GetBool(constData.ConfAppAllowUploadRabbitmqKey) {
		initRabbitmqClient()
	}

	if viper.GetBool(constData.ConfAppAllowMqttUploadKey) {
		initMqtt()
	}

	log.WithFields(log.Fields{
		"op":   "startup",
		"step": "alert",
	}).Info()

}

func bindApiAppConfigKey() {
	configUtils.ConfigKeyList = append(configUtils.ConfigKeyList,
		[]string{constData.ConfDebugFlagKey, constData.EnvDebugFlagKey},

		[]string{constData.ConfServerListenPortKey, constData.EnvServerListenPortKey},
		[]string{constData.ConfServerListenHostKey, constData.EnvServerListenHostKey},

		[]string{constData.ConfAppMqttHostKey, constData.EnvAppMqttHostKey},
		[]string{constData.ConfAppMqttPortKey, constData.EnvAppMqttPortKey},
		[]string{constData.ConfAppMqttUserKey, constData.EnvAppMqttUserKey},
		[]string{constData.ConfAppMqttPasswordKey, constData.EnvAppMqttPasswordKey},
		[]string{constData.ConfAppUploadMqttTopicPrefixKey, constData.EnvAppUploadMqttTopicPrefixKey},

		[]string{constData.ConfTtyPathKey, constData.EnvTtyPathKey},
		[]string{constData.ConfUploadIntervalKey, constData.EnvUploadIntervalKey},
		[]string{constData.ConfAppAllowMqttUploadKey, constData.EnvAppAllowMqttUploadKey},
		[]string{constData.ConfAppAllowUploadRabbitmqKey, constData.EnvAppAllowUploadRabbitmqKey},
		[]string{constData.ConfAllowGetNvmeKey, constData.EnvAllowGetNvmeKey},
		[]string{constData.ConfNvmePathKey, constData.EnvNvmePathKey},

		[]string{constData.ConfRabbitmqHostKey, constData.EnvRabbitmqHostKey},
		[]string{constData.ConfRabbitmqPortKey, constData.EnvRabbitmqPortKey},
		[]string{constData.ConfRabbitmqUserKey, constData.EnvRabbitmqUserKey},
		[]string{constData.ConfRabbitmqPasswordKey, constData.EnvRabbitmqPasswordKey},
		[]string{constData.ConfRabbitmqVirtualHostKey, constData.EnvRabbitmqVirtualHostKey},
		[]string{constData.ConfRabbitmqChannelMaxKey, constData.EnvRabbitmqChannelMaxKey},
		[]string{constData.ConfRabbitmqFrameSizeKey, constData.EnvRabbitmqFrameSizeKey},
		[]string{constData.ConfRabbitmqHeartBeatKey, constData.EnvRabbitmqHeartBeatKey},
		[]string{constData.ConfAppRabbitmqUploadExchangeKey, constData.EnvAppRabbitmqUploadExchangeKey},
		[]string{constData.ConfAppRabbitmqUploadQueuePrefixKey, constData.EnvAppRabbitmqUploadQueuePrefixKey},
		[]string{constData.ConfAppRabbitmqUploadRoutingPrefixKey, constData.EnvAppRabbitmqUploadRoutingPrefixKey},
	)
}
func bindApiAppConfigDefaultValue() {
	viper.SetDefault(constData.ConfDebugFlagKey, constData.DefaultDebugFlag)

	viper.SetDefault(constData.ConfServerListenPortKey, constData.DefaultServerListenPort)
	viper.SetDefault(constData.ConfServerListenHostKey, constData.DefaultServerListenHost)

	viper.SetDefault(constData.ConfAppMqttHostKey, constData.DefaultAppMqttHost)
	viper.SetDefault(constData.ConfAppMqttPortKey, constData.DefaultAppMqttPort)
	viper.SetDefault(constData.ConfAppMqttUserKey, constData.DefaultAppMqttUser)
	viper.SetDefault(constData.ConfAppMqttPasswordKey, constData.DefaultAppMqttPassword)
	viper.SetDefault(constData.ConfAppUploadMqttTopicPrefixKey, constData.DefaultAppUploadMqttTopicPrefix)

	viper.SetDefault(constData.ConfTtyPathKey, constData.DefaultTtyPath)
	viper.SetDefault(constData.ConfUploadIntervalKey, constData.DefaultUploadInterval)
	viper.SetDefault(constData.ConfAppAllowMqttUploadKey, constData.DefaultAppAllowMqttUpload)
	viper.SetDefault(constData.ConfAppAllowUploadRabbitmqKey, constData.DefaultAppAllowUploadRabbitmq)
	viper.SetDefault(constData.ConfAllowGetNvmeKey, constData.DefaultAllowGetNvme)
	viper.SetDefault(constData.ConfNvmePathKey, constData.DefaultNvmePath)

	viper.SetDefault(constData.ConfRabbitmqHostKey, constData.DefaultRabbitmqHost)
	viper.SetDefault(constData.ConfRabbitmqPortKey, constData.DefaultRabbitmqPort)
	viper.SetDefault(constData.ConfRabbitmqUserKey, constData.DefaultRabbitmqUser)
	viper.SetDefault(constData.ConfRabbitmqPasswordKey, constData.DefaultRabbitmqPassword)
	viper.SetDefault(constData.ConfRabbitmqVirtualHostKey, constData.DefaultRabbitmqVirtualHost)
	viper.SetDefault(constData.ConfRabbitmqChannelMaxKey, constData.DefaultRabbitmqChannelMax)
	viper.SetDefault(constData.ConfRabbitmqFrameSizeKey, constData.DefaultRabbitmqFrameSize)
	viper.SetDefault(constData.ConfRabbitmqHeartBeatKey, constData.DefaultRabbitmqHeartBeat)
	viper.SetDefault(constData.ConfAppRabbitmqUploadExchangeKey, constData.DefaultAppRabbitmqUploadExchange)
	viper.SetDefault(constData.ConfAppRabbitmqUploadQueuePrefixKey, constData.DefaultAppRabbitmqUploadQueuePrefix)
	viper.SetDefault(constData.ConfAppRabbitmqUploadRoutingPrefixKey, constData.DefaultAppRabbitmqUploadRoutingPrefix)
}

func getCpuSerialNum() (result string) {
	result = ""

	cmd := exec.Command("bash", "-c", "grep Serial /proc/cpuinfo | awk '{print $3}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		result = "---"
	} else {
		result = strings.TrimSpace(strings.ReplaceAll(string(out), "\n", ""))
	}

	return
}

func getMacAddress() (result string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, netInterface := range interfaces {
		mac := netInterface.HardwareAddr.String()
		if mac != "" && len(mac) == 17 && strings.Trim(strings.ReplaceAll(mac, ":", ""), "0") != "" {
			return strings.ReplaceAll(mac, ":", "")
		}
	}
	return ""
}

func getIdString() (result string) {
	cpuSerialNum := getCpuSerialNum()
	if cpuSerialNum != "" && len(cpuSerialNum) > 0 && strings.Trim(cpuSerialNum, "0") != "" {
		return cpuSerialNum
	} else {
		return getMacAddress()
	}
}
