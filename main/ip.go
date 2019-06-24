package main

import (
	"fmt"
	"net"
)

var (
	localIpArray []string
)


func init() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		panic(fmt.Sprintf("get local ip failed,%v", err))
	}
	//获取本机的所有ipv4地址
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIpArray = append(localIpArray, ipnet.IP.String())
			}
		}
	}
	fmt.Println(localIpArray)
}
