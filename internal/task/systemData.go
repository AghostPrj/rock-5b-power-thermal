/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 2:07
 * @Desc:
 */

package task

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/constData"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/global"
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	"github.com/anatol/smart.go"
	linuxProc "github.com/c9s/goprocinfo/linux"
	"github.com/md14454/gosensors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	lastCpuIdle  uint64
	lastCpuTotal uint64
)

func getMemInfo() (memoryTotal, memoryFree, memoryAvailable uint64, memoryUsage float64, err error) {
	memoryTotal = 0
	memoryFree = 0
	memoryAvailable = 0
	memoryUsage = 0

	var b []byte
	b, err = os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	content := string(b)
	lines := strings.Split(content, "\n")

	varGets := 0
	for _, line := range lines {
		if varGets >= 3 {
			break
		}

		if strings.HasPrefix(line, "MemTotal") {
			varGets++

			tmp := strings.ReplaceAll(line, "MemTotal:", "")
			tmp = strings.ReplaceAll(tmp, "kB", "")
			tmp = strings.TrimSpace(tmp)
			parseUint, err := strconv.ParseUint(tmp, 10, 64)
			if err == nil {
				memoryTotal = parseUint
			}
			continue
		}
		if strings.HasPrefix(line, "MemFree") {
			varGets++

			tmp := strings.ReplaceAll(line, "MemFree:", "")
			tmp = strings.ReplaceAll(tmp, "kB", "")
			tmp = strings.TrimSpace(tmp)
			parseUint, err := strconv.ParseUint(tmp, 10, 64)
			if err == nil {
				memoryFree = parseUint
			}
			continue
		}
		if strings.HasPrefix(line, "MemAvailable") {
			varGets++

			tmp := strings.ReplaceAll(line, "MemAvailable:", "")
			tmp = strings.ReplaceAll(tmp, "kB", "")
			tmp = strings.TrimSpace(tmp)
			parseUint, err := strconv.ParseUint(tmp, 10, 64)
			if err == nil {
				memoryAvailable = parseUint
			}
			continue
		}
	}

	memoryUsage = math.Trunc((float64(memoryTotal-memoryFree)/float64(memoryTotal))*10000) / 100

	return
}

func getCpuInfo() (upTime int64, cpuUsage float64, err error) {
	upTime = 0
	cpuUsage = 0

	stat, err := linuxProc.ReadStat("/proc/stat")
	if err != nil || stat == nil {
		return
	}

	upTime = time.Now().Unix() - stat.BootTime.Unix()

	nowIdle := stat.CPUStatAll.Idle

	nowTotal := stat.CPUStatAll.User + stat.CPUStatAll.Nice + stat.CPUStatAll.System + stat.CPUStatAll.Idle +
		stat.CPUStatAll.IOWait + stat.CPUStatAll.IRQ + stat.CPUStatAll.SoftIRQ

	tmpTotal := nowTotal - lastCpuTotal
	tmpIdle := nowIdle - lastCpuIdle

	cpuUsage = math.Trunc((float64(tmpTotal-tmpIdle)/float64(tmpTotal))*10000) / 100

	lastCpuIdle = nowIdle
	lastCpuTotal = nowTotal

	return

}

func GetSocHotSpotTemperature() (hotSpotTemperature float64) {
	hotSpotTemperature = 0

	chips := gosensors.GetDetectedChips()

	for i := range chips {
		chip := chips[i]

		features := chip.GetFeatures()

		for j := range features {
			feature := features[j]

			if strings.HasPrefix(feature.GetLabel(), "temp") {
				value := feature.GetValue()
				if value > hotSpotTemperature {
					hotSpotTemperature = value
				}
			}
		}

	}
	return
}

func getNvmeInfo() (temperature int32, availSpare uint8, err error) {
	temperature = 0
	availSpare = 0

	device, err := smart.OpenNVMe(viper.GetString(constData.ConfNvmePathKey))
	if err != nil {
		return
	}

	defer func() {
		_ = device.Close()
	}()

	smartLog, err := device.ReadSMART()
	if err != nil {
		return
	}

	temperature = int32(smartLog.Temperature) - 273
	availSpare = smartLog.AvailSpare

	return

}

func GetSystemData() {
	for {
		time.Sleep(time.Millisecond * 1000)

		newSystemCacheData := object.SystemCacheData{}

		memoryTotal, memoryFree, memoryAvailable, memoryUsage, err := getMemInfo()
		if err != nil {
			log.WithFields(log.Fields{
				"op":   "get_system_data",
				"step": "get_memory_info",
				"err":  err,
			}).Warn()
		} else {
			newSystemCacheData.MemoryTotal = memoryTotal
			newSystemCacheData.MemoryFree = memoryFree
			newSystemCacheData.MemoryAvailable = memoryAvailable
			newSystemCacheData.MemoryUsage = memoryUsage
		}

		upTime, cpuUsage, err := getCpuInfo()
		if err != nil {
			log.WithFields(log.Fields{
				"op":   "get_system_data",
				"step": "get_cpu_info",
				"err":  err,
			}).Warn()
		} else {
			newSystemCacheData.UpTime = upTime
			newSystemCacheData.CpuUsage = cpuUsage
		}

		newSystemCacheData.SocHotSpotTemperature = GetSocHotSpotTemperature()

		if viper.GetBool(constData.ConfAllowGetNvmeKey) {
			nvmeTemperature, nvmeAvailSpare, err := getNvmeInfo()
			if err != nil {

			} else {
				newSystemCacheData.NvmeTemperature = nvmeTemperature
				newSystemCacheData.NvmeAvailSpare = nvmeAvailSpare

			}
		}

		newSystemCacheData.UpdateAt = time.Now().Unix()

		global.SystemDataCacheLock.Lock()
		global.SystemCachedData = &newSystemCacheData
		global.SystemDataCacheLock.Unlock()

	}

}
