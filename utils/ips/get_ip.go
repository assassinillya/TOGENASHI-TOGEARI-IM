package ips

import (
	"fmt"
	"net"
)

// GetIP 此方法会获取外网的地址并非本机的ip, 不作使用只作为留档
func GetIP() (addr string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网卡信息出错: ", err)
		return ""
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("获取IP地址出错: ", err)
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}
	return
}
