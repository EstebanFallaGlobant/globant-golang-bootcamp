package tests

import (
	"fmt"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part2/methods-types-structures/mts/ipaddr"
)

func TestIPAddrStringer(t *testing.T) {
	ipStr := "127.0.0.1"
	ip := ipaddr.IPAddr{127, 0, 0, 1}
	str := fmt.Sprint(ip)
	if str != ipStr {
		t.Errorf("IPAddr type string conversion failed. Got: %s, Want: %s", str, ipStr)
	}
}
