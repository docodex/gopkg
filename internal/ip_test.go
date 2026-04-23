package internal_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/docodex/gopkg/internal"
)

func TestPrivateIPv4(t *testing.T) {
	ip, err := internal.PrivateIPv4()
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

	lower8, err := internal.Lower8BitPrivateIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lower8 != uint8(ip[3]) {
		t.Fatalf("expected %d, got %d", ip[3], lower8)
	}
	fmt.Println(internal.Lower8BitPrivateIPv4())

	lower16, err := internal.Lower16BitPrivateIPv4()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := uint16(ip[2])<<8 + uint16(ip[3])
	if lower16 != expected {
		t.Fatalf("expected %d, got %d", expected, lower16)
	}
	fmt.Println(internal.Lower16BitPrivateIPv4())

	fmt.Println(uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
	fmt.Println(uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
}
