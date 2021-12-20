package ipaddr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPAddrStringer(t *testing.T) {
	ipStr := "127.0.0.1"
	ip := IPAddr{127, 0, 0, 1}
	str := fmt.Sprint(ip)

	assert.Equal(t, ipStr, str)
}
