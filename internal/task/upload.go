/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 2:05
 * @Desc:
 */

package task

import (
	"encoding/json"
	"time"

	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func UploadData() {
	lastUpload := int64(0)
	lastUpdate := int64(0)
	uploadInterval := viper.GetInt64(constData.ConfUploadIntervalKey)

	if uploadInterval < 5 {
		uploadInterval = 5
	}

	for {
		time.Sleep(time.Millisecond * 100)

		global.McuDataCacheLock.Lock()
		tmpMcuCachedData := global.McuCachedData
		global.McuDataCacheLock.Unlock()

		if tmpMcuCachedData == nil {
			continue
		}

		global.SystemDataCacheLock.Lock()
		tmpSystemCachedData := global.SystemCachedData
		global.SystemDataCacheLock.Unlock()

		if tmpSystemCachedData == nil {
			continue
		}

		nowUnixTimestamp := time.Now().Unix()
		if (nowUnixTimestamp-lastUpload) > uploadInterval &&
			lastUpdate < tmpMcuCachedData.UpdateAt &&
			lastUpdate < tmpSystemCachedData.UpdateAt {

			lastUpload = nowUnixTimestamp

			lastUpdate = tmpMcuCachedData.UpdateAt
			if lastUpdate < tmpSystemCachedData.UpdateAt {
				lastUpdate = tmpSystemCachedData.UpdateAt
			}

			uploadData := object.UploadData{}
			uploadData.Parse(tmpSystemCachedData, tmpMcuCachedData)

			log.WithFields(log.Fields{
				"data": uploadData,
				"op":   "upload",
				"step": "log_upload_data",
			}).Debug()

			uploadDataJson, err := json.Marshal(uploadData)
			if err != nil {
				continue
			}

			if viper.GetBool(constData.ConfAppAllowMqttUploadKey) {
				go func(jsonString string) {
					publishToken := global.MqttClient.Publish(global.MqttUploadTopic, 0, false, jsonString)
					if !publishToken.WaitTimeout(5 * time.Second) {
						log.WithField("op", "upload").WithField("step", "mqtt_publish").Warn("mqtt publish timeout")
					}
				}(string(uploadDataJson))
			}

		}

	}
}
