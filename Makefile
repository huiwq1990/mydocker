VERSION_MAJOR ?= 0
VERSION_MINOR ?= 1
VERSION_BUILD ?= 0
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

$(shell mkdir -p ./out)

############
# Building #
############

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X $(REPOPATH)/pkg/version.version=$(VERSION)" \
	  -a -o mydocker cmd/mydocker/main.go

alpine:
	rm -rf /root/alpine
	mkdir -p /root/alpine
	tar -xvf doc/alpine-minirootfs-3.11.3-x86_64.tar.gz -C /root/alpine

busybox:
	cd /root/
	docker export `docker run -itd busybox:latest` > busybox.tar
	mkdir busybox && tar -xvf busybox.tar -C busybox

run:
	./mydocker run --tty  alpine /bin/sh

ci: build run
