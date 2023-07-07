/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:07
 * @Desc:
 */

package task

import (
	"encoding/json"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	"strings"
	"time"
)

func ProcessRawSerialData(chanRcvSerialData <-chan string, chanSingleData chan<- string) {
	var residue string
	for {
		input, ok := <-chanRcvSerialData
		if !ok {
			break
		}

		// 拼接新接收的字符串
		combined := residue + input

		// 以换行符分割
		parts := strings.Split(combined, "\n")

		// 取出最后一个元素作为下一次拼接的开头
		residue = parts[len(parts)-1]

		// 将除最后一个之外的其他切分出来的数据一个个的输出到chan2
		for _, part := range parts[:len(parts)-1] {
			chanSingleData <- strings.ReplaceAll(part, "\n", "")
		}
	}
}

func CacheData(chanSingleData <-chan string) {

	transferData := &object.McuTransferData{}

	outputData := &object.McuCacheData{}

	for {
		input, ok := <-chanSingleData
		if !ok {
			break
		}

		if strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}") {
			transferData = &object.McuTransferData{}

			err := json.Unmarshal([]byte(input), transferData)
			if err != nil {
				continue
			}

			if transferData != nil && transferData.Data != nil && transferData.Operation == "data_transfer" {

				outputData = &object.McuCacheData{
					UpdateAt:          time.Now().Unix(),
					McuDeviceId:       transferData.Data.DeviceId,
					RawInputVoltage:   [3]float32{},
					McuSystemTick:     transferData.Data.SystemTick,
					VddaValue:         float32(transferData.Data.VddaValue) / 1000,
					McuTemperature:    transferData.Data.McuTemperature,
					InputCurrent:      float32(transferData.Data.AdcValueChannel4) / 500,
					SensorTemperature: transferData.Data.SensorTemperature,
					FanPwmDuty:        transferData.Data.FanPwmDuty,
					FanSpeed:          transferData.Data.FanSpeed,
				}

				divNum := 0

				if transferData.Data.AdcValueChannel1 <= (transferData.Data.VddaValue - 50) {
					outputData.RawInputVoltage[0] = float32(transferData.Data.AdcValueChannel1) * 5.7 / 1000
					divNum++
				} else {
					outputData.RawInputVoltage[0] = 0
				}

				if transferData.Data.AdcValueChannel2 <= (transferData.Data.VddaValue - 50) {
					outputData.RawInputVoltage[1] = float32(transferData.Data.AdcValueChannel2) * 3.7 / 1000
					divNum++
				} else {
					outputData.RawInputVoltage[1] = 0
				}

				if transferData.Data.AdcValueChannel3 <= (transferData.Data.VddaValue - 50) {
					outputData.RawInputVoltage[2] = float32(transferData.Data.AdcValueChannel3) * 7.8 / 1000
					divNum++
				} else {
					outputData.RawInputVoltage[2] = 0
				}

				outputData.InputVoltage = (outputData.RawInputVoltage[0] + outputData.RawInputVoltage[1] + outputData.RawInputVoltage[2]) / float32(divNum)
				outputData.InputPower = outputData.InputCurrent * outputData.InputVoltage

				global.McuDataCacheLock.Lock()
				global.McuCachedData = outputData
				global.McuDataCacheLock.Unlock()

			}

		}
	}
}
