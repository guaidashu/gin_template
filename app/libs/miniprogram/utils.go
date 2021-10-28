/**
 * @Author: yy
 * @Author: 1023767856@qq.com
 * @Date: 28/09/2021
 * @Desc: 其他工具
 */

package miniprogram

import "net"

// LocalIP 获取机器的IP
func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}
