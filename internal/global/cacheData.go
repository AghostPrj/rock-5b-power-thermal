/**
 * @Author: aghost<ggg17226@gmail.com>
 * @Date: 2023/7/8 1:25
 * @Desc:
 */

package global

import (
	"github.com/AghostPrj/rock-5b-power-thermal/internal/object"
	"sync"
)

var (
	McuDataCacheLock = sync.Mutex{}

	SystemDataCacheLock = sync.Mutex{}

	McuCachedData *object.McuCacheData

	SystemCachedData *object.SystemCacheData
)
