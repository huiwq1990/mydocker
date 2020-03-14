package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
)

// 文件内容 {"192.168.0.0/24":"0100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"}
const ipamDefaultAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"

type IPAM struct {
	SubnetAllocatorPath string
	Subnets *map[string]string
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: ipamDefaultAllocatorPath,
}

func (ipam *IPAM) load() error {
	if _, err := os.Stat(ipam.SubnetAllocatorPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}
	subnetConfigFile, err := os.Open(ipam.SubnetAllocatorPath)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}
	subnetJson := make([]byte, 2000)
	n, err := subnetConfigFile.Read(subnetJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(subnetJson[:n], ipam.Subnets)
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
		return err
	}
	return nil
}

func (ipam *IPAM) dump() error {
	ipamConfigFileDir, _ := path.Split(ipam.SubnetAllocatorPath)
	if _, err := os.Stat(ipamConfigFileDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(ipamConfigFileDir, 0644)
		} else {
			return err
		}
	}
	subnetConfigFile, err := os.OpenFile(ipam.SubnetAllocatorPath, os.O_TRUNC | os.O_WRONLY | os.O_CREATE, 0644)
	defer subnetConfigFile.Close()
	if err != nil {
		return err
	}

	ipamConfigJson, err := json.Marshal(ipam.Subnets)
	if err != nil {
		return err
	}

	_, err = subnetConfigFile.Write(ipamConfigJson)
	if err != nil {
		return err
	}

	return nil
}

func (ipam *IPAM) Allocate(subnet *net.IPNet) (net.IP, error) {
	// 存放网段中地址分配信息的数组
	ipam.Subnets = &map[string]string{}

	// 从文件中加载已经分配的网段信息
	err := ipam.load()
	if err != nil {
		log.Errorf("load allocation info, %v", err)
		return nil,err
	}

	subnetAllocBitMap,exist := (*ipam.Subnets)[subnet.String()]
	// 如果不存在，初始化
	if !exist {
		ones,bits := subnet.Mask.Size()
		num := 1 << uint(bits - ones )
		log.Debugf("subnet, mask: %v %v %v",ones,bits,num)
		for i:=0;i < num;i++{
			subnetAllocBitMap += "0"
		}
	}

	var allocIP net.IP
	for c := range subnetAllocBitMap {
		if subnetAllocBitMap[c] == '0' {
			ipalloc := []byte(subnetAllocBitMap)
			log.Debugf("before alloc: %s",subnetAllocBitMap)
			ipalloc[c] = '1'
			(*ipam.Subnets)[subnet.String()] = string(ipalloc)
			log.Debugf("after alloc: %s, val: %d",(*ipam.Subnets)[subnet.String()],c)
			// 即使是ipv4, 长度也可能为16
			allocIP = subnet.IP.To4()
			ipLen := uint(len(allocIP))
			log.Debugf("alloc start ip: %s, len: %v",allocIP,ipLen)
			for t := uint(ipLen); t > 0; t-=1 {
				[]byte(allocIP)[ipLen-t] += uint8(c >> ((t - 1) * 8))
			}
			allocIP[3]+=1
			break
		}
	}

	ipam.dump()
	log.Debugf("alloc ip result: %s",allocIP.String())
	return allocIP,nil
}

func (ipam *IPAM) Release(subnet *net.IPNet, ipaddr *net.IP) error {
	ipam.Subnets = &map[string]string{}

	_, subnet, _ = net.ParseCIDR(subnet.String())

	err := ipam.load()
	if err != nil {
		log.Errorf("Error dump allocation info, %v", err)
	}

	c := 0
	releaseIP := ipaddr.To4()
	releaseIP[3]-=1
	for t := uint(4); t > 0; t-=1 {
		c += int(releaseIP[t-1] - subnet.IP[t-1]) << ((4-t) * 8)
	}

	ipalloc := []byte((*ipam.Subnets)[subnet.String()])
	ipalloc[c] = '0'
	(*ipam.Subnets)[subnet.String()] = string(ipalloc)

	ipam.dump()
	return nil
}