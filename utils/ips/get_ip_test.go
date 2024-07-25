package ips

import (
	"fmt"
	"testing"
)

func TestGetIP(t *testing.T) {
	ip := GetIP()
	fmt.Println(ip)
}
