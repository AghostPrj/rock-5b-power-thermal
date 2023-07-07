/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 2:05
 * @Desc:
 */

package task

import (
	"context"
	"encoding/json"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
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
					_ = publishToken.Wait()
				}(string(uploadDataJson))
			}

			if viper.GetBool(constData.ConfAppAllowUploadRabbitmqKey) {
				go func(jsonByte []byte) {
					testCount := 0
					for {
						if global.RabbitmqUploadChannel.IsClosed() {
							testCount++
							time.Sleep(time.Millisecond * 5)
						} else {
							break
						}

						if testCount > 10 {
							break
						}
					}

					if testCount >= 10 {
						return
					}

					err = global.RabbitmqUploadChannel.PublishWithContext(context.Background(),
						viper.GetString(constData.ConfAppRabbitmqUploadExchangeKey), global.RabbitmqRoutingKey, false, false, amqp.Publishing{
							ContentType: "application/json",
							Body:        jsonByte,
						})
					if err != nil {
						log.WithFields(log.Fields{
							"op":   "upload",
							"step": "upload_rabbitmq",
							"err":  err,
						}).Warn()
					}
				}(uploadDataJson)

			}

		}

	}
}
