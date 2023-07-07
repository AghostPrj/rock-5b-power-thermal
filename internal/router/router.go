/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 4:23
 * @Desc:
 */

package router

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func StartGinServer(router *gin.Engine) {
	err := router.Run(viper.GetString(constData.ConfServerListenHostKey) +
		":" + viper.GetString(constData.ConfServerListenPortKey))
	if err != nil {
		log.WithField("err", err.Error()).WithField("op", "startup").Fatal()
	}
}

func handleAccess(context *gin.Context) {
	global.McuDataCacheLock.Lock()
	tmpMcuCachedData := global.McuCachedData
	global.McuDataCacheLock.Unlock()

	if tmpMcuCachedData == nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	global.SystemDataCacheLock.Lock()
	tmpSystemCachedData := global.SystemCachedData
	global.SystemDataCacheLock.Unlock()

	if tmpSystemCachedData == nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, nil)
		return
	}

	uploadData := object.UploadData{}
	uploadData.Parse(tmpSystemCachedData, tmpMcuCachedData)

	context.JSON(http.StatusOK, &uploadData)
	return
}

func BuildGinRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	if viper.GetBool(constData.ConfDebugFlagKey) {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
			ExposeHeaders: []string{"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, " +
				"Cache-Control, Content-Language, Content-Type, x-transfer-token, x-captcha-id, uri"},
		}))
	}

	router.Any("/", handleAccess)
	router.Any("/info", handleAccess)

	return router
}
