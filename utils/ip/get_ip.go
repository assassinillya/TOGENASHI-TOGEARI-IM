package ip

import (
	"fmt"
	"net"
)

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
