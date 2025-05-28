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
		t.Logf("get PrivateIPv4 error: %v", ip)
		return
	}
	buf := strings.Builder{}
	for i := range ip {
		if buf.Len() != 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(strconv.Itoa(int(ip[i])))
	}
	fmt.Println(buf.String())
	fmt.Println(internal.Lower8BitPrivateIPv4())
	fmt.Println(internal.Lower16BitPrivateIPv4())
	fmt.Println(uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
	fmt.Println(uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3]))
}
