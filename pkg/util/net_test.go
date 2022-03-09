package util

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestGetInterfaceIpv4Addr(t *testing.T) {
	addr, err := GetInterfaceIpv4Addr("vEthernet (Internet)")
	if err != nil {
		t.Error(err)
	}
	require.Equal(t, "192.168.30.2", addr)
}

func TestInterfaces(t *testing.T) {
	netInterfaces, _ := net.Interfaces()
	for _, i := range netInterfaces {
		fmt.Println(i)
	}
}
