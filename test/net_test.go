package main

import (
	"net"
	"testing"
)

func TestAllocate(t *testing.T) {
	ip, _, _ := net.ParseCIDR("192.168.0.1/24")
	t.Logf("alloc ip: %v", ip)
}

func main()  {


}