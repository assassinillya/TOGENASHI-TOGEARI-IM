package addr

import (
	"fmt"
	"github.com/cc14514/go-geoip2"
	geoip2db "github.com/cc14514/go-geoip2-db"
	"net"
)

var addrDB *geoip2.DBReader

func init() {
	addrDB1, err := geoip2db.NewGeoipDbByStatik()
	if err != nil {
		panic(err)
	}
	addrDB = addrDB1
}

func IsIntranetIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return true
	}
	// 192.168
	// 172.16 - 172.31
	// 10
	// 169.254
	return (ip4[0] == 192 && ip4[1] == 168) ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 32) ||
		(ip4[0] == 10) ||
		(ip4[0] == 169 && ip4[1] == 254)
}

func GetAddr(ip string) string {
	parseIP := net.ParseIP(ip)
	if IsIntranetIP(parseIP) {
		return "内部地址"
	}
	record, err := addrDB.City(net.ParseIP(ip))
	if err != nil {
		return "错误的地址"
	}
	var province string
	if len(record.Subdivisions) > 0 {
		province = record.Subdivisions[0].Names["zh-CN"]
	}
	city := record.City.Names["zh-CN"]
	return fmt.Sprintf("%s-%s", province, city)
}
