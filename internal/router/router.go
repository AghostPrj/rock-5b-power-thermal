/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 4:23
 * @Desc:
 */

package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StartGinServer(router *gin.Engine) {
	addr := viper.GetString(constData.ConfServerListenHostKey) +
		":" + viper.GetString(constData.ConfServerListenPortKey)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("err", err.Error()).WithField("op", "startup").Fatal()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.WithField("op", "shutdown").Info("received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("err", err).WithField("op", "shutdown").Fatal()
	}
	log.WithField("op", "shutdown").Info("server stopped")
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
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	if viper.GetBool(constData.ConfDebugFlagKey) {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "x-transfer-token", "x-captcha-id", "uri"},
		}))
	}

	router.Any("/", handleAccess)
	router.Any("/info", handleAccess)

	return router
}
