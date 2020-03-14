package main

import (
	"net"
	"testing"
)

func TestAllocate(t *testing.T) {
	ip, ipnet, _ := net.ParseCIDR("192.168.0.1/24")

	a,b := ipnet.Mask.Size()
	t.Logf("alloc ip: %v, mask:%v %v", ip, a,b )
}