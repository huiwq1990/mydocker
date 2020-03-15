VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

############
# Building #
############

.PHONY: build
build:
	#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X $(REPOPATH)/pkg/version.version=$(VERSION)" \
	#  -a -o mydocker cmd/mydocker/main.go
	go build -o mydocker cmd/mydocker/main.go

alpine:
	rm -rf /root/alpine
	mkdir -p /root/alpine
	tar -xvf doc/alpine-minirootfs-3.11.3-x86_64.tar.gz -C /root/alpine

run:
	./mydocker run --tty  alpine /bin/sh
runv:
	./mydocker run --tty -v /vagrant:/test alpine /bin/sh


	./mydocker network delete net1
	./mydocker network create net1 --subnet 10.1.1.0/24

	./mydocker run --tty --net net1 alpine /bin/sh


ci: build run
