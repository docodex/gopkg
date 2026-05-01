package internal

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateIPv4(t *testing.T) {
	ip, err := PrivateIPv4()
	if err != nil {
		t.Skipf("skipping: no private IPv4 address available: %v", err)
		return
	}

	if len(ip) != 4 {
		t.Fatalf("expected IPv4 (4 bytes), got %d bytes: %v", len(ip), ip)
	}

	buf := strings.Builder{}
	for i := range ip {
		if buf.Len() != 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(strconv.Itoa(int(ip[i])))
	}
	fmt.Println(buf.String())

	lower8, err := Lower8BitPrivateIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lower8 != uint8(ip[3]) {
		t.Fatalf("expected %d, got %d", ip[3], lower8)
	}
	fmt.Println(Lower8BitPrivateIPv4())

	lower16, err := Lower16BitPrivateIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := uint16(ip[2])<<8 + uint16(ip[3])
	if lower16 != expected {
		t.Fatalf("expected %d, got %d", expected, lower16)
	}
	fmt.Println(Lower16BitPrivateIPv4())

	fmt.Println(uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
	fmt.Println(uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
}

func TestIsPrivateIPv4Table(t *testing.T) {
	cases := []struct {
		name string
		ip   net.IP
		want bool
	}{
		{"nil", nil, false},
		{"10.x.y.z", net.IPv4(10, 0, 0, 1).To4(), true},
		{"172.15 below range", net.IPv4(172, 15, 0, 1).To4(), false},
		{"172.16 low boundary", net.IPv4(172, 16, 0, 1).To4(), true},
		{"172.31 high boundary", net.IPv4(172, 31, 255, 254).To4(), true},
		{"172.32 above range", net.IPv4(172, 32, 0, 1).To4(), false},
		{"192.168.x.y", net.IPv4(192, 168, 1, 1).To4(), true},
		{"192.167.x.y", net.IPv4(192, 167, 1, 1).To4(), false},
		{"169.254.x.y", net.IPv4(169, 254, 1, 1).To4(), true},
		{"169.253.x.y", net.IPv4(169, 253, 1, 1).To4(), false},
		{"public 8.8.8.8", net.IPv4(8, 8, 8, 8).To4(), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, isPrivateIPv4(tc.ip))
		})
	}
}
