/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:08
 * @Desc:
 */

package object

type PowerAndThermalData struct {
	DeviceId          string  `json:"dev_id"`
	SystemTick        uint64  `json:"sys_tick"`
	VddaValue         uint32  `json:"vdda"`
	McuTemperature    uint32  `json:"mcu_temp"`
	AdcValueChannel1  uint32  `json:"adc_ch_1"`
	AdcValueChannel2  uint32  `json:"adc_ch_2"`
	AdcValueChannel3  uint32  `json:"adc_ch_3"`
	AdcValueChannel4  uint32  `json:"adc_ch_4"`
	SensorTemperature float32 `json:"sensor_temp"`
	FanPwmDuty        float32 `json:"duty_pwm"`
	FanSpeed          uint64  `json:"fan_speed"`
}

type McuTransferData struct {
	Operation string               `json:"op"`
	Version   string               `json:"ver"`
	Data      *PowerAndThermalData `json:"data"`
}

type McuCacheData struct {
	UpdateAt          int64      `json:"update_at"`
	McuDeviceId       string     `json:"dev_id"`
	McuSystemTick     uint64     `json:"mcu_sys_tick"`
	VddaValue         float32    `json:"vdda"`
	McuTemperature    uint32     `json:"mcu_temp"`
	InputVoltage      float32    `json:"input_voltage"`
	RawInputVoltage   [3]float32 `json:"-"`
	InputCurrent      float32    `json:"input_current"`
	InputPower        float32    `json:"input_power"`
	SensorTemperature float32    `json:"sensor_temp"`
	FanPwmDuty        float32    `json:"duty_pwm"`
	FanSpeed          uint64     `json:"fan_speed"`
}

type SystemCacheData struct {
	UpdateAt              int64   `json:"update_at"`
	MemoryTotal           uint64  `json:"memory_total"`
	MemoryFree            uint64  `json:"memory_free"`
	MemoryAvailable       uint64  `json:"memory_available"`
	MemoryUsage           float64 `json:"memory_usage"`
	CpuUsage              float64 `json:"cpu_usage"`
	UpTime                int64   `json:"up_time"`
	SocHotSpotTemperature float64 `json:"soc_hot_spot_temperature"`
	NvmeTemperature       int32   `json:"nvme_temperature"`
	NvmeAvailSpare        uint8   `json:"nvme_avail_spare"`
}

type UploadData struct {
	UpdateAt              int64   `json:"update_at"`
	MemoryTotal           uint64  `json:"memory_total"`
	MemoryFree            uint64  `json:"memory_free"`
	MemoryAvailable       uint64  `json:"memory_available"`
	MemoryUsage           float64 `json:"memory_usage"`
	CpuUsage              float64 `json:"cpu_usage"`
	UpTime                int64   `json:"up_time"`
	SocHotSpotTemperature float64 `json:"soc_hot_spot_temperature"`
	NvmeTemperature       int32   `json:"nvme_temperature"`
	NvmeAvailSpare        uint8   `json:"nvme_avail_spare"`
	McuTemperature        uint32  `json:"mcu_temp"`
	InputVoltage          float32 `json:"input_voltage"`
	InputCurrent          float32 `json:"input_current"`
	InputPower            float32 `json:"input_power"`
	SensorTemperature     float32 `json:"sensor_temp"`
	FanPwmDuty            float32 `json:"duty_pwm"`
	FanSpeed              uint64  `json:"fan_speed"`
}

func (u *UploadData) Parse(sysData *SystemCacheData, mcuData *McuCacheData) {
	u.UpdateAt = sysData.UpdateAt
	if mcuData.UpdateAt > u.UpdateAt {
		u.UpdateAt = mcuData.UpdateAt
	}

	u.MemoryTotal = sysData.MemoryTotal
	u.MemoryFree = sysData.MemoryFree
	u.MemoryAvailable = sysData.MemoryAvailable
	u.MemoryUsage = sysData.MemoryUsage
	u.CpuUsage = sysData.CpuUsage
	u.UpTime = sysData.UpTime
	u.SocHotSpotTemperature = sysData.SocHotSpotTemperature
	u.NvmeTemperature = sysData.NvmeTemperature
	u.NvmeAvailSpare = sysData.NvmeAvailSpare

	u.McuTemperature = mcuData.McuTemperature
	u.InputVoltage = mcuData.InputVoltage
	u.InputCurrent = mcuData.InputCurrent
	u.InputPower = mcuData.InputPower
	u.SensorTemperature = mcuData.SensorTemperature
	u.FanPwmDuty = mcuData.FanPwmDuty
	u.FanSpeed = mcuData.FanSpeed

}
