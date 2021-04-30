/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package enum

// redis锁时间
const (
	DEFAULT_LOCK_ACQUIRE_TIMEOUT = 1 // 秒
	DEFAULT_LOCK_KEY_TIMEOUT     = 1
	LOCK_PREFIX                  = "lock:"
)
