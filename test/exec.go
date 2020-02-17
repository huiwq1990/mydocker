package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

func main()  {

	ip,ipNet, _ := net.ParseCIDR("192.168.1.0/24")

c := 10
len := uint(len(ip))
	for t := len; t > 0; t-=1 {
		[]byte(ip)[len-t] += uint8(c >> ((t - 1) * 8))
	}
	ip[3]+=1
	log.Println(ip)

	log.Printf("alloc ip: %v", ip)
	log.Printf("alloc ip: %v", ipNet.IP)

	log.Printf("pid:%d\n", os.Getpid())


	_,err := net.InterfaceByName("bb")
	fmt.Println("%v",err)
	//  action := exec.Command("sh")
	//
	//  action.Stdin = os.Stdin
	//  action.Stderr = os.Stderr
	//  action.Stdout = os.Stdout
	//
	//  if err := action.Run(); err != nil {
	//	//      log.Printf("Init Run() function err : %v\n", err)
	//      log.Fatal(rr)e
	//  }
	command := "/bin/sh"
	if err := syscall.Exec(command, []string{command}, os.Environ()); err != nil {
		log.Printf("syscall.Exec err: %v\n", err)
		log.Fatal(err)
	}
}
