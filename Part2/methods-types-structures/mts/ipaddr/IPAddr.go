package ipaddr

import "fmt"

type IPAddr [4]byte

func (p IPAddr) String() string {
	var result string

	for i, val := range p {
		if i > 0 {
			result += "."
		}
		result += fmt.Sprintf("%d", val)
	}

	return result
}
