package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func Test001(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatalf("os.pipe error:%v\n", err)
	}
	_, err = writer.Write([]byte("pipe content"))
	if err != nil {
		log.Fatalf("writer.Write error:%v\n", err)
	}

	buf := make([]byte, 20)
	n, err := reader.Read(buf)
	if err != nil {
		log.Fatalf("reader.Read(buf) error:%v\n", err)
	}
	log.Printf("Read Content:%q\n", string(buf[:n]))
}


func Test002(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatalf("os.pipe error:%v\n", err)
	}
	go func() {
		for i := 0; i < 10; i++ {
			content := fmt.Sprintf("%s-%d\n", "pipe content", i)
			_, err = writer.Write([]byte(content))
			if err != nil {
				log.Fatalf("writer.Write error:%v\n", err)
			}
		}
		writer.Close()
	}()

	go func() {
		n, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalf("reader.Read(buf) error:%v\n", err)
		}
		log.Printf("Read Content:%q\n", n)
	}()

	for i := 0; i <= 100; i++{
		time.Sleep(1 * time.Second)
	}
}


