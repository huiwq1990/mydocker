package main

import (
	log "github.com/sirupsen/logrus"

)

func main()  {
	arr := make([]string,0)
	arr = append(arr,"aa")
	log.Infof("%s",arr)
}
