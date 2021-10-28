/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 30/04/2021
 * @Desc: desc
 */

package enum

// redis锁时间
const (
	DefaultLockAcquireTimeout = 1 // 秒
	DefaultLockKeyTimeout     = 1
	LockPrefix                = "redis_lock:"
)
