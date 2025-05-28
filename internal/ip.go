package internal

import (
	"errors"
	"net"
)

var ErrNoPrivateAddress = errors.New("no private ip address")

func PrivateIPv4() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}

	return nil, ErrNoPrivateAddress
}

func isPrivateIPv4(ip net.IP) bool {
	// Allow private IP addresses (RFC1918) and link-local addresses (RFC3927)
	return ip != nil &&
		(ip[0] == 10 || // 10.0.0.0 to 10.255.255.255
			(ip[0] == 172 && ip[1] >= 16 && ip[1] < 32) || // 172.16.0.0 to 172.31.255.255
			(ip[0] == 192 && ip[1] == 168) || // 192.168.0.0 to 192.168.255.255
			(ip[0] == 169 && ip[1] == 254)) // 169.254.0.0/16
}

func Lower8BitPrivateIPv4() uint8 {
	ip, err := PrivateIPv4()
	if err != nil || len(ip) != 4 {
		return 0
	}
	return uint8(ip[3])
}

func Lower16BitPrivateIPv4() uint16 {
	ip, err := PrivateIPv4()
	if err != nil || len(ip) != 4 {
		return 0
	}
	return uint16(ip[2])<<8 + uint16(ip[3])
}
