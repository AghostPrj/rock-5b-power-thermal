/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:17
 * @Desc:
 */

package initializator

import (
	"net"
	"os/exec"
	"strings"

	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/ggg17226/aghost-go-base/pkg/utils/configUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitApp() {
	configUtils.SetConfigFileName(constData.ApplicationName)
	bindApiAppConfigKey()
	bindApiAppConfigDefaultValue()
	configUtils.InitConfigAndLog()

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
		[]string{constData.ConfAllowGetNvmeKey, constData.EnvAllowGetNvmeKey},
		[]string{constData.ConfNvmePathKey, constData.EnvNvmePathKey},
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
	viper.SetDefault(constData.ConfAllowGetNvmeKey, constData.DefaultAllowGetNvme)
	viper.SetDefault(constData.ConfNvmePathKey, constData.DefaultNvmePath)
}

func getCpuSerialNum() (result string) {
	result = ""

	cmd := exec.Command("bash", "-c", "grep Serial /proc/cpuinfo | awk '{print $3}'")
	out, err := cmd.CombinedOutput()
	if err != nil {
		result = ""
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

	isValidMac := func(mac string) bool {
		return mac != "" && len(mac) == 17 && strings.Trim(strings.ReplaceAll(mac, ":", ""), "0") != ""
	}

	for _, netInterface := range interfaces {
		mac := netInterface.HardwareAddr.String()
		if isValidMac(mac) {
			if strings.HasPrefix(netInterface.Name, "eth") ||
				strings.HasPrefix(netInterface.Name, "en") ||
				strings.HasPrefix(netInterface.Name, "wlan") ||
				strings.HasPrefix(netInterface.Name, "wl") {
				return strings.ReplaceAll(mac, ":", "")
			}
		}
	}

	for _, netInterface := range interfaces {
		mac := netInterface.HardwareAddr.String()
		if isValidMac(mac) {
			return strings.ReplaceAll(mac, ":", "")
		}
	}
	return ""
}

func getIdString() (result string) {
	cpuSerialNum := getCpuSerialNum()
	if cpuSerialNum != "" &&
		cpuSerialNum != "---" &&
		strings.Trim(cpuSerialNum, "0") != "" {
		return cpuSerialNum
	} else {
		return getMacAddress()
	}
}
