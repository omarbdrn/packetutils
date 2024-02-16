package utils

import (
	"fmt"
	"net"
	"time"
)

var dialer = &net.Dialer{
	Timeout: 2 * time.Second,
}

func DialTCP(ip string, port uint16) (net.Conn, error) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	return dialer.Dial("tcp", addr)
}

func DialUDP(ip string, port uint16) (net.Conn, error) {
	addr := net.JoinHostPort(ip, fmt.Sprintf("%d", port))
	return dialer.Dial("udp", addr)
}

func getLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return net.ParseIP("127.0.0.1"), nil
			}
		}
	}

	return nil, fmt.Errorf("No suitable local IP address found")
}
